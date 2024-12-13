package main

import (
	"fmt"
	"log"
	"sync"
)

type RoomState interface {
	acceptLoop()
	Name() string
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
	log.Printf("State %s activated in room %d\n", r.State.Name(), r.Id)
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
						typ: MessagePlayer1Turn,
					}
					r.r.Player2.Inch <- ServerMsg{
						typ: MessagePlayer1Turn,
					}	

					return
				}
			}
		}
	}
}

func (r *RoomWaitingForPlayers) Name() string {
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
			case MessagePlayer1Played:
				// Update state
				// Notify player2
				r.r.setState(&Player2Turn{
					r: r.r,
				})
				
				go r.r.Start()

				r.r.Player1.Inch <- ServerMsg{
					typ: MessagePlayer2Turn,
				}
				r.r.Player2.Inch <- ServerMsg{
					typ: MessagePlayer2Turn,
				}	

				return
			case MessagePlayer2Played:
				// Put on a queue?
				// or illegal
			case MessagePlayer2Turn:
				// Illegal
			case MessagePlayerJoined:
				panic("Illegal")
			}

		}
	}
}

func (r *Player1Turn) Name() string {
	return "Player1Turn"
}

func (r *Player2Turn) acceptLoop() {
	for {
		select {
		case msg := <- r.r.Inch:
			switch msg.typ {
			case MessagePlayer2Played:
				// Update state
				// Notify players
				r.r.setState(&Player1Turn{
					r: r.r,
				})

				go r.r.Start()

				r.r.Player1.Inch <- ServerMsg{
					typ: MessagePlayer1Turn,
				}
				r.r.Player2.Inch <- ServerMsg{
					typ: MessagePlayer1Turn,
				}

				return
			case MessagePlayer1Played:
				// Put on a queue?
				// or illegal
			case MessagePlayer1Turn:
				// Illegal
			case MessagePlayerJoined:
				panic("Illegal")
			}
		}
	}
}

func (r *Player2Turn) Name() string {
	return "Player2Turn"
}