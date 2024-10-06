package main

import "github.com/gorilla/websocket"

type Player struct {
	rId 	int16
	conn  	*websocket.Conn
	
	Inch  	chan ServerMessage
	outch 	chan ServerMessage
}

func NewPlayer(conn *websocket.Conn) *Player {
	p := &Player{
		conn: conn,
		Inch: make(chan ServerMessage),
	}
	go p.listenForClient()

	return p
}

// Two types of messages from client -> turnCmd and leaveCmd 
func (p *Player) listenForClient() {
	for {
		_, p, err := p.conn.ReadMessage()
		_ = p
		if err != nil {
			// Cleanup
			return
		}
	}
}

// Handle player and room communication on the server
// For example -> room notifies the player1 that it is their turn 
// Receive the channel comm here and pass it onto the client
func (p *Player) acceptLoop() {
	for {
		select {
		case msg := <- p.Inch:
			_ = msg
		}
	}
}