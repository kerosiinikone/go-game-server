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
			log.Printf("Player %d played\n", p.Id)
			p.outch <- ServerMsg{
				Typ: MessagePlayer1Played,
				PlayerId: p.Id,
			}
		case MessagePlayer2Played:
			log.Printf("Player %d played\n", p.Id)
			p.outch <- ServerMsg{
				Typ: MessagePlayer2Played,
				PlayerId: p.Id,
			}
		}
		// If the client disconnects, remove them from the room
		// If the client sends a leave message, remove them from the room
		// If the client sends a turn message, send it to the room
	}
}

// Handle player and room communication on the server
// For example -> room notifies the player1 that it is their turn 
// Receive the channel comm here and pass it onto the client
func (p *Player) acceptLoop() {
	for {
		select {
		case msg := <- p.Inch:
			switch msg.Typ {
			case MessagePlayer1Turn:
				var (
					bytes []byte
					clientMsg WSMsg
				)
				if p.Player1 {
					clientMsg = WSMsg{
						Typ: MessagePlayer1Turn,
						PlayerId: p.Id,
						Card: msg.Card,
						ScoreCards: msg.ScoreCards,
					}
				} else {
					clientMsg = WSMsg{
						Typ: MessagePlayer1Turn,
						PlayerId: p.Id,
						ScoreCards: msg.ScoreCards,
					}
				}
				bytes, err := json.Marshal(&clientMsg)
				if err != nil {
					log.Printf("Error unmarshalling message: %v\n", err)
				}
				p.conn.WriteMessage(websocket.TextMessage, bytes)
			case MessagePlayer2Turn:
				var (
					bytes []byte
					clientMsg WSMsg
				)
				if !p.Player1 {
					clientMsg = WSMsg{
						Typ: MessagePlayer2Turn,
						PlayerId: p.Id,
						Card: msg.Card,
					}
				} else {
					clientMsg = WSMsg{
						Typ: MessagePlayer2Turn,
						PlayerId: p.Id,
					}
				}
				bytes, err := json.Marshal(&clientMsg)
				if err != nil {
					log.Printf("Error unmarshalling message: %v\n", err)
				}
				p.conn.WriteMessage(websocket.TextMessage, bytes)
			}
		}
	}
}