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
	typ int
	roomId int16
	playerId int16
}

type WSMsg struct {
	Typ int
	PlayerId int16
	Data []byte
}