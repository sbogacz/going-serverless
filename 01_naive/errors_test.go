package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrors(t *testing.T) {
	t.Run("bad request error gives 400s", func(t *testing.T) {
		err := newBadRequestErr(fmt.Errorf("your oops"), "")
		require.Error(t, err)
		resp := errorResponse(err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// works with a more specific message
		err = newBadRequestErr(fmt.Errorf("your oops"), "seriously, you stink")
		require.Error(t, err)
		resp = errorResponse(err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("internal server error gives 500s", func(t *testing.T) {
		err := newInternalServerErr(fmt.Errorf("our oops"), "")
		require.Error(t, err)
		resp := errorResponse(err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// works with a more specific message
		err = newInternalServerErr(fmt.Errorf("our oops"), "seriously, I stink")
		require.Error(t, err)
		resp = errorResponse(err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
	t.Run("not found error gives 404s", func(t *testing.T) {
		err := newNotFoundErr(fmt.Errorf("it's not there"), "")
		require.Error(t, err)
		resp := errorResponse(err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		// works with a more specific message
		err = newNotFoundErr(fmt.Errorf("it's not there"), "seriously, I can't find it anywhere")
		require.Error(t, err)
		resp = errorResponse(err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
