package main

// Message types

const (
	MessageRoomClosed = iota
	MessageHandshake
	MessagePlayerJoined
	MessageGameStarted
	MessageRoomDestroyed
	MessagePlayer1Turn
	MessagePlayer2Turn
	MessagePlayer1Played
	MessagePlayer2Played
	MessageGameOver
)

type ServerMsg struct {
	Typ 		int
	RoomId 		int16
	PlayerId 	int16
	Card 		Card
	Won			bool
	Winner 		int16
	War 		bool
}

type WSMsg struct {
	Typ 		int
	PlayerId 	int16
	Card 		Card
	Won			bool
	Winner 		int16
	War 		bool
}

func NewServerMsg(t int, r int16, p int16, c Card, w bool, wn int16, war bool) ServerMsg {
	return ServerMsg{
		Typ: t,
		RoomId: r,
		PlayerId: p,
		Card: c,
		Won: w,
		Winner: wn,
		War: war,
	}
}

func NewWSMsg(s ServerMsg) WSMsg {
	return WSMsg{
		Typ: s.Typ,
		PlayerId: s.PlayerId,
		Card: s.Card,
		Won: s.Won,
		Winner: s.Winner,
		War: s.War,
	}
}

func SendToClient(p *Player, w *WSMsg) error {
	return p.conn.WriteJSON(w)
}