package main

import (
	"fmt"
	"sync"
)

// Own functions, own loops for different situations
type RoomState interface {
	acceptLoop()
}

type Room struct {
	Id int16 
	
	State RoomState
	
	Player1 *Player
	Player2 *Player
	
	// For incoming messages
	Inch chan ServerMessage
	
	cfg *Config
	mu sync.Mutex
}

type RoomWaitingForPlayers struct {
    r *Room 
}

type RoomGameStarted struct {
	r *Room
}

func NewRoom(id int16, cfg *Config) *Room {
	r := &Room{
		Id: id,
		cfg: cfg,
		Inch: make(chan ServerMessage),
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
					r.r.Player1.Inch <- ServerMessage{
						typ: MessageGameStarted,
					}
					r.r.Player2.Inch <- ServerMessage{
						typ: MessageGameStarted,
					}
				}

				// A better way to achieve this?
				r.r.setState(&RoomGameStarted{
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

func (r *RoomGameStarted) acceptLoop() {
	for {
		select {
		case msg := <- r.r.Inch:
			switch msg.typ {}
		}
	}
}