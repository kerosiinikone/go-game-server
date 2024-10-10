package main

import (
	"fmt"
	"sync"
)

type RoomState interface {
	acceptLoop()
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
	// Init stuff and start accepting messages
	r.State.acceptLoop()
}

func (r *Room) setState(state RoomState) {
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
				fmt.Printf("New player joined to %d", r.r.Id)

				if r.r.Player1 != nil && r.r.Player2 != nil {
					// Start the game
					r.r.Player1.Inch <- ServerMsg{
						typ: MessageGameStarted,
					}
					r.r.Player2.Inch <- ServerMsg{
						typ: MessageGameStarted,
					}
				}

				// A better way to achieve this?
				r.r.setState(&Player1Turn{
					r: r.r,
				})

				defer func ()  {
					go r.r.Start()
				}()

				return
			}
		}
	}
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