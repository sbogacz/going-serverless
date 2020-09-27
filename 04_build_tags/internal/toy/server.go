package toy

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/sbogacz/going-serverless/04_build_tags/internal/httperrs"
	"gocloud.dev/blob"
	"gocloud.dev/gcerrors"
)

// Server represents all of the config and clients
// needed to run our app
type Server struct {
	store *blob.Bucket
	cfg   *Config

	// Router is exposed to allow us to call ServeHTTP directly when needed
	Router *chi.Mux
	cancel chan struct{}
}

// New tries to cerate a new instance of Service
func New(c *Config, bucket *blob.Bucket) *Server {
	s := &Server{
		cfg:    c,
		store:  bucket,
		Router: chi.NewRouter(),
		cancel: make(chan struct{}),
	}
	s.Router.Route("/blobs", func(r chi.Router) { // HL
		r.Post("/", s.storeBlob)
		r.Route("/{key}", func(r chi.Router) {
			r.Get("/", s.getBlob)
			r.Delete("/", s.deleteBlob)
		})
	})
	return s
}

// Start starts the server
func (s *Server) Start() {
	h := &http.Server{ // OMIT
		Addr:         fmt.Sprintf(":%d", s.cfg.Port), // OMIT
		ReadTimeout:  5 * time.Second,                // OMIT
		WriteTimeout: 5 * time.Second,                // OMIT
		Handler:      s.Router,                       // OMIT
	} // OMIT
	// OMIT
	go func() { // OMIT
		<-s.cancel                           // OMIT
		_ = h.Shutdown(context.Background()) // OMIT
	}() // OMIT
	// OMIT
	if err := h.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}

// Stop stops the server gracefully
func (s *Server) Stop() {
	s.cancel <- struct{}{}
}

func (s *Server) storeBlob(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "couldn't parse the request body", http.StatusBadRequest)
		return
	}

	// create key
	u, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "failed to generate key", http.StatusInternalServerError)
		return
	}
	key := u.String()

	bw, err := s.store.NewWriter(context.TODO(), key, nil)
	if err != nil {
		http.Error(w, "failed to store", http.StatusInternalServerError)
		return
	}

	if _, err = bw.Write(b); err != nil {
		http.Error(w, "failed to store", http.StatusInternalServerError)
		return
	}

	if err = bw.Close(); err != nil {
		http.Error(w, "failed to store", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, key) //nolint: errcheck,gosec
}

func (s *Server) getBlob(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")

	_, err := uuid.FromString(key)
	if err != nil {
		http.Error(w, "invalid key", http.StatusBadRequest)
		return
	}

	rdr, err := s.store.NewReader(context.TODO(), key, nil)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			http.Error(w, "no such object", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to retrieve object", http.StatusInternalServerError)
		return
	}

	if _, err = io.Copy(w, rdr); err != nil {
		http.Error(w, "failed to retrieve object", http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteBlob(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")

	_, err := uuid.FromString(key)
	if err != nil {
		http.Error(w, "invalid key", http.StatusBadRequest)
		return
	}
	if err := s.store.Delete(context.TODO(), key); err != nil {
		if gcerrors.Code(err) == gcerrors.NotFound {
			http.Error(w, "no such object to delete", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete object", httperrs.StatusCode(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
