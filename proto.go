package main

// Message types

const (
	MessageRoomClosed = iota
	MessagePlayerJoined
	MessageGameStarted
	// ...
)

// Channel comms
type ServerMessage struct {
	typ int
	// ...
}

// JSON
type WSMessage struct{}