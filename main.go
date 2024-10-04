package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Initiate the websocket server
// Room logic creation (room communication through channels)
// Handle player connections to a room (default to creating one room, wait for connections -> if two players join, create a new room)

type RoomCollection map[string]*Room

const (
	defaultReadBufferSize = 1024
	defaultWriteBufferSize = 1024
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  defaultReadBufferSize,
    WriteBufferSize: defaultWriteBufferSize,
}

// Wrap for errors?
func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        // Logger
        return
    }

	// Assign the request to a new player
	// Connect the player to a available room -> if no available rooms, spawn a new one 
	
	for {
		// Convert to a new message type in proto
		_, p, err := conn.ReadMessage()
		
		if err != nil {
			// Logger
			return
		}
		_ = p
	}
}

func resolveRoomConnection(p *Player) {
	// Find empty rooms
	// If none -> create a new one
}

func main() {
	var (
		addr = flag.String("addr", ":3000", "game server address")
		rooms RoomCollection = make(RoomCollection)
	)
	flag.Parse()
	cfg := NewConfig(*addr)

	// Always create one room at the start
	firstRoom := NewRoom(1, cfg)
	rooms[fmt.Sprintf("%d", firstRoom.Id)] = firstRoom

	go firstRoom.Start()

	// For every individual connection player
	http.HandleFunc("/", ServeWebSocket)

	http.ListenAndServe(cfg.Addr, nil)
}