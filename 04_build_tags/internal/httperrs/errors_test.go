package httperrs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	t.Run("bad request error gives 400s", func(t *testing.T) {
		err := BadRequest(fmt.Errorf("your oops"), "")
		require.Error(t, err)
		code := StatusCode(err)
		require.Equal(t, http.StatusBadRequest, code)

		// works with a more specific message
		err = BadRequest(fmt.Errorf("your oops"), "seriously, you stink")
		require.Error(t, err)
		code = StatusCode(err)
		require.Equal(t, http.StatusBadRequest, code)
	})
	t.Run("internal server error gives 500s", func(t *testing.T) {
		err := InternalServer(fmt.Errorf("our oops"), "")
		require.Error(t, err)
		code := StatusCode(err)
		require.Equal(t, http.StatusInternalServerError, code)

		// works with a more specific message
		err = InternalServer(fmt.Errorf("our oops"), "seriously, I stink")
		require.Error(t, err)
		code = StatusCode(err)
		require.Equal(t, http.StatusInternalServerError, code)
	})
	t.Run("not found error gives 404s", func(t *testing.T) {
		err := NotFound(fmt.Errorf("it's not there"), "")
		require.Error(t, err)
		code := StatusCode(err)
		require.Equal(t, http.StatusNotFound, code)

		// works with a more specific message
		err = NotFound(fmt.Errorf("it's not there"), "seriously, I can't find it anywhere")
		require.Error(t, err)
		code = StatusCode(err)
		require.Equal(t, http.StatusNotFound, code)
	})
}
