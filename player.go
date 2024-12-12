package main

import "github.com/gorilla/websocket"

type Player struct {
	Id   	int16
	rId 	int16
	conn  	*websocket.Conn
	Player1 bool
	
	Inch  	chan ServerMsg
	outch 	chan ServerMsg
}

func NewPlayer(conn *websocket.Conn) *Player {
	p := &Player{
		conn: conn,
		Inch: make(chan ServerMsg),
	}
	go p.listenForClient()

	return p
}

// Two types of messages from client -> turnCmd and leaveCmd 
func (p *Player) listenForClient() {
	for {
		_, p, err := p.conn.ReadMessage()
		// p is the message from the client

		// If the client disconnects, remove them from the room
		// If the client sends a leave message, remove them from the room
		// If the client sends a turn message, send it to the room

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
			switch msg.typ {
			case MessageGameStarted:
				// Use a helper to send a response to client (GameStarted)
				// Wait for client to acknowledge
			}
		}
	}
}