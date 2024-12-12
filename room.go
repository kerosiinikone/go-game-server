package main

import (
	"fmt"
	"log"
	"sync"
)

type RoomState interface {
	acceptLoop()
	name() string
}

type Room struct {
	Id int16 
	
	State RoomState
	
	Player1 *Player
	Player2 *Player
	
	Inch chan ServerMsg
	
	cfg *Config
	mu sync.Mutex
}

type RoomWaitingForPlayers struct {
    r *Room 
}

func NewRoom(id int16, cfg *Config) *Room {
	r := &Room{
		Id: id,
		cfg: cfg,
		Inch: make(chan ServerMsg),
		mu: sync.Mutex{},
	}
	r.State = &RoomWaitingForPlayers{
		r: r,
	}
	return r
}

func (r *Room) Start() {
	log.Printf("State %s activated in room %d\n", r.State.name(), r.Id)
	r.State.acceptLoop()
}

func (r *Room) setState(state RoomState) {
	log.Printf("Room %d state changed to %T\n", r.Id, state)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.State = state
}


// Wait for players to connect
func (r *RoomWaitingForPlayers) acceptLoop() {
	for {
		select {
		case msg := <- r.r.Inch:
			switch msg.typ {
			case MessagePlayerJoined:
				fmt.Printf("New player %d joined to %d\n", msg.playerId, r.r.Id)
				if r.r.Player1 != nil && r.r.Player2 != nil {
					r.r.setState(&Player1Turn{
						r: r.r,
					})

					go r.r.Start()

					r.r.Player1.Inch <- ServerMsg{
						typ: MessageGameStarted,
					}
					r.r.Player2.Inch <- ServerMsg{
						typ: MessageGameStarted,
					}
					
					return
				}
			}
		}
	}
}

func (r *RoomWaitingForPlayers) name() string {
	return "RoomWaitingForPlayers"
}

// A few considerations:
// A queue for all the incoming requests -> not needed since
// only one player goes at a time with one input
func (r *Player1Turn) acceptLoop() {
	for {
		select {
		case msg := <- r.r.Inch:
			switch msg.typ {
			case MessagePlayer1Turn:
				// Handle
			case MessagePlayer2Turn:
				// Put on a queue?
				// or illegal
			case MessagePlayerJoined:
				panic("Illegal")
			}

		}
	}
}

func (r *Player1Turn) name() string {
	return "Player1Turn"
}