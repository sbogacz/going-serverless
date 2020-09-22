package toy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gofrs/uuid"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	s    *Server
	addr string
)

func TestMain(m *testing.M) {
	if err := setupServer(); err != nil {
		log.WithError(err).Fatal("failed to set up server for tests")
	}

	go s.Start()
	time.Sleep(500 * time.Millisecond)
	status := m.Run()
	s.Stop()
	os.Exit(status)
}

func TestHappyPath(t *testing.T) {
	testBlob := "this is a test blob"
	var key string
	t.Run("create a blob", func(t *testing.T) {
		resp, err := http.Post(addr, "application/text", strings.NewReader(testBlob)) // HL
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

func setupServer() error {
	port, err := freeport.GetFreePort()
	if err != nil {
		return errors.Wrap(err, "failed to get free port")
	}

	c := &Config{
		BucketName: "test-bucket",
		Port:       port,
	}
	addr = fmt.Sprintf("http://127.0.0.1:%d/blobs", port)
	log.Infof("starting server at %s", addr)
	s = New(c, NewLocalStore())
	return nil
}
