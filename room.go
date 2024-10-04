package main

import (
	"sync"
)

type Room struct {
	Id int16 

	// WaitingForPlayers, Player1Turn, Player2Turn, GameOver
	State interface{}

	Player1 *Player
	Player2 *Player

	// For incoming messages
	Inch chan struct{}

	cfg *Config
	mu sync.Mutex
}


func NewRoom(id int16, cfg *Config) *Room {
	return &Room{
		Id: id,
		cfg: cfg,
		Inch: make(chan struct{}),
		mu: sync.Mutex{},
	}
}

func (r *Room) Start() {
	// Init stuff and start accepting messages
	r.acceptLoop()
}

// Wait for players to connect
func (r *Room) acceptLoop() {
	for {
		select {
		case msg := <- r.Inch:
			_ = msg
		}
	}
}