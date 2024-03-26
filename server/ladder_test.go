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
	time.Sleep(time.Second / 4)
	require.Equal(t, 1, len(ladder.games))
	ws3.Close()
	time.Sleep(time.Second / 4)
	require.Equal(t, 2, len(ladder.clients))
	// Connect to the server with another websocket connection
	ws4, _, wsErr4 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr4 != nil {
		t.Fatalf("ws4: %v", wsErr4)
	}
	defer ws4.Close()

	require.NoError(t, wsErr4)
	time.Sleep(time.Second / 4)
	require.Equal(t, 3, len(ladder.clients))
	ws2.Close()
	time.Sleep(time.Second / 4)
	require.Equal(t, 2, len(ladder.clients))
}
func TestTwoGames(t *testing.T) {
	ladder := NewLadder()
	go ladder.run()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	defer s.Close()
	// Connect to the server with four websocket connections
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws5, _, wsErr5 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr5 != nil {
		t.Fatalf("ws5: %v", wsErr5)
	}
	defer ws5.Close()

	require.NoError(t, wsErr5)
	ws6, _, wsErr6 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr6 != nil {
		t.Fatalf("ws6: %v", wsErr6)
	}
	defer ws6.Close()

	require.NoError(t, wsErr6)
	ws7, _, wsErr7 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr7 != nil {
		t.Fatalf("ws7: %v", wsErr7)
	}
	defer ws7.Close()

	require.NoError(t, wsErr7)
	ws8, _, wsErr8 := websocket.DefaultDialer.Dial(u, nil)
	if wsErr8 != nil {
		t.Fatalf("ws8: %v", wsErr8)
	}
	defer ws8.Close()

	require.NoError(t, wsErr8)
	time.Sleep(time.Second / 4)
	require.Equal(t, 4, len(ladder.clients))
	require.Equal(t, 2, len(ladder.games))
}
func TestTwentyConnections(t *testing.T) {
	ladder := NewLadder()
	go ladder.run()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ladder.ServeHTTP(w, r)
	}))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server with twenty websocket connections
	for i := 0; i < 20; i++ {
		ws, _, wsErr := websocket.DefaultDialer.Dial(u, nil)
		if wsErr != nil {
			t.Fatalf("ws%d: %v", i+1, wsErr)
		}
		defer ws.Close()
		require.NoError(t, wsErr)
	}

	time.Sleep(time.Second / 4)
	require.Equal(t, 20, len(ladder.clients))
}
