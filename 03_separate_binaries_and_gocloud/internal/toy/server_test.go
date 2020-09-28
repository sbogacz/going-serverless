package toy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/phayes/freeport"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob/fileblob"
)

var (
	s    *Server
	addr string
)

func TestMain(m *testing.M) {
	dir, cleanup := newTempDir()
	defer cleanup()

	if err := setupServer(dir); err != nil {
		log.WithError(err).Fatal("failed to set up server for tests")
	}

	go s.Start()
	time.Sleep(500 * time.Millisecond) // OMIT
	status := m.Run()
	s.Stop()
	os.Exit(status)
}

func TestHappyPath(t *testing.T) {
	testBlob := "this is a test blob"
	var key string
	t.Run("create a blob", func(t *testing.T) {
		resp, err := http.Post(addr, "application/text", strings.NewReader(testBlob))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Body)

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		key = string(b)
	})

	itemAddr := fmt.Sprintf("%s/%s", addr, key)
	t.Run("fetch the blob", func(t *testing.T) {
		resp, err := http.Get(itemAddr)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Body)

		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, testBlob, string(b))
	})

	t.Run("delete the blob", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", itemAddr, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("fetch the deleted blob", func(t *testing.T) {
		resp, err := http.Get(itemAddr)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestNotFounds(t *testing.T) {
	u, err := uuid.NewV4()
	require.NoError(t, err)
	itemAddr := fmt.Sprintf("%s/%s", addr, u.String())
	t.Run("fetch the blob", func(t *testing.T) {
		resp, err := http.Get(itemAddr)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("delete the blob", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", itemAddr, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func setupServer(dir string) error {
	port, err := freeport.GetFreePort()
	if err != nil {
		return fmt.Errorf("failed to get free port: %w", err)
	}

	c := &Config{
		BucketName: "test-bucket",
		Port:       port,
	}
	addr = fmt.Sprintf("http://127.0.0.1:%d/blobs", port)
	log.Infof("starting server at %s", addr)

	store, err := fileblob.OpenBucket(dir, nil) // HL
	if err != nil {
		return fmt.Errorf("failed to set up local store: %w", err)
	}
	s = New(c, store) // HL
	return nil
}

func newTempDir() (string, func()) {
	dir, err := ioutil.TempDir("", "toy-test-files")
	if err != nil {
		panic(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}
