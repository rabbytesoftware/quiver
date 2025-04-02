package websocket

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func Init() {
	server := socketio.NewServer(nil)

	// Handle connection event
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Client connected:", s.ID())
		return nil
	})

	// Handle disconnect event
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Client disconnected:", s.ID(), reason)
	})

	// Handle a custom event called "message"
	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		fmt.Println("Message received:", msg)
		// Broadcast the message to all clients
		s.Emit("response", "Server received: "+msg)
		server.BroadcastToNamespace("/", "broadcast", s.ID()+": "+msg)
	})

	// Handle error
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Error:", e)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()
	defer server.Close()

	http.Handle("/", server)
	
	fmt.Println("Server started on http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("HTTP server error: ", err)
	}
}