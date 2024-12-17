package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Initiate the websocket server
// Room logic creation (room communication through channels)
// Handle player connections to a room (default to creating one room, wait for connections -> if two players join, create a new room)

type RoomCollection map[string]*Room

type RoomHandler struct {
	latestRoomId int16
	rooms RoomCollection

	mu sync.Mutex

	Inch chan ServerMsg

	cfg *Config
}

const (
	defaultReadBufferSize = 1024
	defaultWriteBufferSize = 1024
	defaultClientAddr = "http://localhost:5173"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  defaultReadBufferSize,
    WriteBufferSize: defaultWriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == defaultClientAddr
	},
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

	// Send a handshake message to the client
	msg := NewWSMsg(ServerMsg{
		Typ: MessageHandshake,
		PlayerId: newPlayer.Id,
	})
	if err := SendToClient(newPlayer, &msg); err != nil {
		log.Printf("Error sending message to client: %v\n", err)
	}
}

// Destroying unused rooms -> set a timeout context when players leave?
func (rh *RoomHandler) acceptLoop() {
	for {
		select {
		case msg := <-rh.Inch:
			switch msg.Typ {
			case MessageRoomDestroyed:
				// Cleanup
				delete(rh.rooms, fmt.Sprintf("%d", msg.RoomId))
			}
		}
	}
}

func (rh *RoomHandler) resolveRoomConnection(p *Player) {
	rh.mu.Lock()
	defer rh.mu.Unlock()

	for _, v := range rh.rooms {
		if v.Player1 == nil {
			v.Player1 = p
			p.Id = 1
			p.outch = v.Inch
			p.Player1 = true
			p.rId = v.Id
			v.Inch <- ServerMsg{
				Typ: MessagePlayerJoined,
				PlayerId: p.Id,
			}
			return
		}
		if v.Player2 == nil {
			v.Player2 = p
			p.Id = 2
			p.outch = v.Inch
			p.rId = v.Id
			v.Inch <- ServerMsg{
				Typ: MessagePlayerJoined,
				PlayerId: p.Id,
			}
			return
		}
	}

	newId := rh.latestRoomId+1
	sId := fmt.Sprintf("%d", newId)

	rh.rooms[sId] = NewRoom(newId, rh.cfg)
	r := rh.rooms[sId]

	p.outch = r.Inch
	p.Id = 1
	p.Player1 = true
	p.rId = r.Id
	r.Player1 = p
	rh.latestRoomId = newId
	
	go r.Start()
	
	r.Inch <- NewServerMsg(MessagePlayerJoined, r.Id, p.Id, Card{}, false, 0, false)

}

func main() {
	var (
		addr = flag.String("addr", ":3000", "game server address")
	)

	flag.Parse()
	cfg := NewConfig(*addr) // Will be passed to the player through handshake

	rh := &RoomHandler{
		rooms: make(RoomCollection),
		Inch: make(chan ServerMsg),
		cfg: cfg,
	}

	http.HandleFunc("/", rh.ServeWebSocket)
	http.ListenAndServe(cfg.Addr, nil)
}