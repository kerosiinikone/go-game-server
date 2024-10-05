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

type RoomHandler struct {
	latestRoomId int16
	Inch chan ServerMessage
	rooms RoomCollection
	cfg *Config
}

const (
	defaultReadBufferSize = 1024
	defaultWriteBufferSize = 1024
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  defaultReadBufferSize,
    WriteBufferSize: defaultWriteBufferSize,
}

// Wrap for errors?
func (rh *RoomHandler) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        // Logger
        return
    }

	// Assign the request to a new player
	newPlayer := NewPlayer(conn)

	// Connect the player to an available room -> if no available rooms, spawn a new one
	rh.resolveRoomConnection(newPlayer)
	
	for {
		// Convert to a new message type in proto
		_, p, err := conn.ReadMessage()
		
		if err != nil {
			// Signal the room that the player disconnected -> clear everything
			return
		}
		_ = p
	}
}

func (rh *RoomHandler) acceptLoop() {
	for {
		select {
		case msg := <-rh.Inch:
			_ = msg
		}
	}
}

func (rh *RoomHandler) resolveRoomConnection(p *Player) {
	for _, v := range rh.rooms {
		if v.Player1 == nil {
			v.Player1 = p
			v.Inch <- ServerMessage{
				typ: MessagePlayerJoined,
			}
			return
		}
		if v.Player2 == nil {
			v.Player2 = p
			v.Inch <- ServerMessage{
				typ: MessagePlayerJoined,
			}
			return
		}
	}

	// Simplify
	newId := rh.latestRoomId+1
	rh.rooms[fmt.Sprintf("%d", newId)] = NewRoom(newId, rh.cfg)
	r := rh.rooms[fmt.Sprintf("%d", newId)]
	r.Player1 = p
	p.outch = r.Inch
	rh.latestRoomId = newId
	
	go r.Start()
	
	// Message wrapper func -> enough time for the server to start accepting comms?
	r.Inch <- ServerMessage{
		typ: MessagePlayerJoined,
	}

}

func main() {
	var (
		addr = flag.String("addr", ":3000", "game server address")
	)
	flag.Parse()
	cfg := NewConfig(*addr)

	rh := &RoomHandler{
		rooms: make(RoomCollection),
		Inch: make(chan ServerMessage),
		cfg: cfg,
	}
	// For every individual connection player
	http.HandleFunc("/", rh.ServeWebSocket)

	http.ListenAndServe(cfg.Addr, nil)
}