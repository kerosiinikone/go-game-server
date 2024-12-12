package main

// Message types

const (
	MessageRoomClosed = iota
	MessagePlayerJoined
	MessageGameStarted
	MessageRoomDestroyed
	MessagePlayer1Turn
	MessagePlayer2Turn
)

type ServerMsg struct {
	typ int
	roomId int16
	playerId int16
}

type WSMsg struct{}