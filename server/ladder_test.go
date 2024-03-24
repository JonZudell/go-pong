package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"github.com/theckman/go-securerandom"
)

func TestLadder(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-Websocket-Version", "13")
	rStr, baseErr := securerandom.Base64OfBytes(16)
	req.Header.Set("Sec-WebSocket-Key", rStr)
	require.NoError(t, err)
	require.NoError(t, baseErr)
	ladder := NewLadder()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	// Be sure to clean up the server or you might run out of file descriptors!
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, wsErr := websocket.DefaultDialer.Dial(u, nil)
	if wsErr != nil {
		t.Fatalf("%v", wsErr)
	}
	defer ws.Close()
	require.NoError(t, wsErr)
}
