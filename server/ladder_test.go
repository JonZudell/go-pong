package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

func TestNConnections(t *testing.T) {
	var connections int = 25000
	ladder := NewLadder()
	go ladder.run()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	var wg sync.WaitGroup
	wg.Add(connections)

	var websockets []*websocket.Conn
	// Connect to the server with twenty websocket connections concurrently
	for i := 0; i < connections; i++ {
		func() {
			ws, _, wsErr := websocket.DefaultDialer.Dial(u, nil)
			if wsErr != nil {
				t.Fatalf("ws%d: %v", i+1, wsErr)
			}
			require.NoError(t, wsErr)
			websockets = append(websockets, ws)
			defer wg.Done()
		}()
	}

	wg.Wait()

	time.Sleep(time.Second)
	// Check the number of connections
	require.Equal(t, connections, len(ladder.conns))

	// Close all the connections
	for _, ws := range websockets {
		err := ws.Close()
		if err != nil {
			t.Fatalf("ws close: %v", err)
		}
	}
	time.Sleep(time.Second / 4)
	require.Equal(t, 0, len(ladder.conns))

}
