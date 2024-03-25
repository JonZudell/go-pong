package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestLadder(t *testing.T) {
	ladder := NewLadder()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
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

func TestMultipleConnections(t *testing.T) {
	ladder := NewLadder()
	go ladder.run()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	require.Equal(t, 0, len(ladder.clients))
	require.Equal(t, 0, len(ladder.games))
	// Connect to the server with two websocket connections
	ws1, _, wsErr1 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr1 != nil {
		t.Fatalf("ws1: %v", wsErr1)
	}
	defer ws1.Close()

	ws2, _, wsErr2 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr2 != nil {
		t.Fatalf("ws2: %v", wsErr2)
	}
	defer ws2.Close()
	ladder.ladderTick()
	require.NoError(t, wsErr1)
	require.NoError(t, wsErr2)
	require.Equal(t, 2, len(ladder.clients))
	ws3, _, wsErr3 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr3 != nil {
		t.Fatalf("ws3: %v", wsErr3)
	}
	defer ws3.Close()

	require.NoError(t, wsErr3)
	require.Equal(t, 3, len(ladder.clients))
	time.Sleep(time.Second)
	require.Equal(t, 1, len(ladder.games))

}
