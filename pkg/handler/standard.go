package handler

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
)

// Allow all origins
func allowAllOrigins(wsHandler websocket.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Upgrade, Connection, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions")
		w.Header().Set("Connection", "Upgrade")
		w.Header().Set("Upgrade", "websocket")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Forcefully set the Origin header (standard library rejects requests without it)
		if r.Header.Get("Origin") == "" {
			r.Header.Set("Origin", "http://localhost")
		}

		// Upgrade connection
		wsHandler.ServeHTTP(w, r)
	}
}

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client:", ws.RemoteAddr())
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected or read error:", err)
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[0:n]
		fmt.Println("Received:", string(msg))
		s.broadcast(msg)

		//for single client write
		//_, err = ws.Write([]byte("Thank you for the message"))
		//if err != nil {
		//	fmt.Println("Write error:", err)
		//}
	}
}

func (s *Server) broadcast(msg []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("Write error:", err)
			}
		}(ws)
	}
}

func testStandard() {
	server := NewServer()

	// Apply CORS fix
	http.Handle("/ws", allowAllOrigins(websocket.Handler(server.handleWS)))

	fmt.Println("WebSocket server started on ws://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
