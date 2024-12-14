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
)

type ServerMsg struct {
	Typ 		int
	RoomId 		int16
	PlayerId 	int16
	Card 		Card
	ScoreCards	[]Card
}

type WSMsg struct {
	Typ 		int
	PlayerId 	int16
	Card 		Card
	ScoreCards	[]Card
}