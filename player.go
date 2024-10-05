package main

import "github.com/gorilla/websocket"

type Player struct {
	Inch  chan ServerMessage
	// Room inch
	outch chan ServerMessage

	conn  *websocket.Conn
}

func NewPlayer(conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
		Inch: make(chan ServerMessage),
	}
}

// Handle player and room communication on the server
func (p *Player) acceptLoop() {
	for {
		select {
		case msg := <- p.Inch:
			_ = msg
		}
	}
}