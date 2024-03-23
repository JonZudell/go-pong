package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	s := NewServer("")
	s.ServeHTTP(res, req)
	require.Equal(t, 200, res.Code)
}
