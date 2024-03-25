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
func TestServerError(t *testing.T) {
	req, err := http.NewRequest("GET", "/404", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	s := NewServer("")
	s.ServeHTTP(res, req)
	require.Equal(t, 404, res.Code)
}
func TestServerInvalidRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/upgrade", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	s := NewServer("")
	s.ServeHTTP(res, req)
	require.Equal(t, 400, res.Code)
}
