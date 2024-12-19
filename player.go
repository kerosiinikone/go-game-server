package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

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
	go p.acceptLoop()

	return p
}

// Two types of messages from client -> turnCmd and leaveCmd 
func (p *Player) listenForClient() {
	defer func ()  {
		log.Printf("Player %d left\n", p.Id)
		p.outch <- NewServerMsg(MessagePlayerLeft, p.rId, p.Id, Card{}, false, 0, false)
		p.conn.Close()	
	}()
	for {
		var msg WSMsg
		_, bytes, err := p.conn.ReadMessage()
		if err != nil {
			return
		}
		if err := json.Unmarshal(bytes, &msg); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
		}
		// Also check for state of the room
		switch msg.Typ {
		case MessagePlayer1Played:
			p.outch <- NewServerMsg(MessagePlayer1Played, p.rId, p.Id, msg.Card, msg.Won, 0, msg.War)
		case MessagePlayer2Played:
			p.outch <- NewServerMsg(MessagePlayer2Played, p.rId, p.Id, msg.Card, msg.Won, 0, msg.War)
		}
		// If the client disconnects, remove them from the room
		// If the client sends a leave message, remove them from the room
		// If the client sends a turn message, send it to the room
	}
}

// ERRORS

// Handle player and room communication on the server
// For example -> room notifies the player1 that it is their turn 
// Receive the channel comm here and pass it onto the client
func (p *Player) acceptLoop() {
	for {
		select {
		case msg := <- p.Inch:
			switch msg.Typ {
			case MessagePlayerJoined:
				clientMsg := NewWSMsg(msg)
				SendToClient(p, clientMsg)
			case MessagePlayer1Turn:
				clientMsg := NewWSMsg(msg)
				SendToClient(p, clientMsg)
			case MessagePlayer2Turn:
				clientMsg := NewWSMsg(msg)
				SendToClient(p, clientMsg)
			case MessageGameOver:
				clientMsg := NewWSMsg(msg)
				SendToClient(p, clientMsg)
				p.close()
				return
			}
		}
	}
}

func (p *Player) close() {
	log.Printf("Player %d disconnected\n", p.Id)
	p.conn.Close()
}