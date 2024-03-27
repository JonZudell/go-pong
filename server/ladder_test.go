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

func Test1000Connections(t *testing.T) {
	ladder := NewLadder()
	go ladder.run()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server with twenty websocket connections
	for i := 0; i < 1000; i++ {
		ws, _, wsErr := websocket.DefaultDialer.Dial(u, nil)
		if wsErr != nil {
			t.Fatalf("ws%d: %v", i+1, wsErr)
		}
		ws.WriteJSON(map[string]interface{}{
			"type": "ready",
		})
		defer ws.Close()
		require.NoError(t, wsErr)
	}
	time.Sleep(time.Second / 2)
	require.Equal(t, 1000, len(ladder.clients))
	require.Equal(t, 500, len(ladder.games))
}
func Test500Connections(t *testing.T) {
	ladder := NewLadder()
	go ladder.run()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server with twenty websocket connections
	for i := 0; i < 500; i++ {
		ws, _, wsErr := websocket.DefaultDialer.Dial(u, nil)
		if wsErr != nil {
			t.Fatalf("ws%d: %v", i+1, wsErr)
		}
		if i < 250 {
			ws.WriteJSON(map[string]interface{}{
				"type": "ready",
			})
		}
		defer ws.Close()
		require.NoError(t, wsErr)
	}
	time.Sleep(time.Second / 2)
	require.Equal(t, 500, len(ladder.clients))
	require.Equal(t, 125, len(ladder.games))
}
