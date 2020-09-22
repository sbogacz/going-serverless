package toy

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/sbogacz/gophercon18-kickoff-talk/third/internal/httperrs"
)

// Server represents all of the config and clients
// needed to run our app
type Server struct {
	store Store
	cfg   *Config

	router *chi.Mux      // HL
	cancel chan struct{} // OMIT
}

// New tries to cerate a new instance of Service
func New(c *Config, s Store) *Server {
	return &Server{
		cfg:    c,
		store:  s,
		router: chi.NewRouter(),
		cancel: make(chan struct{}),
	}
}

// Start starts the server
func (s *Server) Start() {
	s.router.Route("/blobs", func(r chi.Router) { // HL
		r.Post("/", s.storeBlob)
		r.Route("/{key}", func(r chi.Router) {
			r.Get("/", s.getBlob)
			r.Delete("/", s.deleteBlob)
		})
	})
	h := &http.Server{ // OMIT
		Addr:         fmt.Sprintf(":%d", s.cfg.Port), // OMIT
		ReadTimeout:  5 * time.Second,                // OMIT
		WriteTimeout: 5 * time.Second,                // OMIT
		Handler:      s.router,                       // OMIT
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

// Router exposes our chi Route externally
func (s *Server) Router() *chi.Mux {
	return s.router
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

	if err := s.store.Set(context.TODO(), key, string(b)); err != nil {
		http.Error(w, "failed to store", httperrs.StatusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
}

func (s *Server) getBlob(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")

	_, err := uuid.FromString(key)
	if err != nil {
		http.Error(w, "invalid key", http.StatusBadRequest)
		return
	}

	data, err := s.store.Get(context.TODO(), key)
	if err != nil {
		http.Error(w, "failed to retrieve object", httperrs.StatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func (s *Server) deleteBlob(w http.ResponseWriter, req *http.Request) {
	key := chi.URLParam(req, "key")

	_, err := uuid.FromString(key)
	if err != nil {
		http.Error(w, "invalid key", http.StatusBadRequest)
		return
	}
	if err := s.store.Del(context.TODO(), key); err != nil {
		http.Error(w, "failed to delete object", httperrs.StatusCode(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
