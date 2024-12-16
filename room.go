package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

type RoomState interface {
	acceptLoop()
	Name() string
}

type RoomWaitingForPlayers struct {
	r *Room 
}

type Room struct {
	Id int16 
	
	// Room state
	State RoomState

	// Game state
	Deck []Card
	Player1Drawn []Card
	Player2Drawn []Card
	Player1Cards []Card
	Player2Cards []Card
	
	Player1 *Player
	Player2 *Player
	
	Inch chan ServerMsg
	
	cfg *Config
	mu sync.Mutex
}

func NewRoom(id int16, cfg *Config) *Room {
	r := &Room{
		Id: id,
		cfg: cfg,
		Inch: make(chan ServerMsg),
		mu: sync.Mutex{},
		// Two decks of cards
		Deck: append(deck, deck...),
	}
	r.State = &RoomWaitingForPlayers{
		r: r,
	}
	return r
}

func (r *Room) Start() {
	log.Printf("State %s activated in room %d\n", r.State.Name(), r.Id)
	r.State.acceptLoop()
}

func (r *Room) setState(state RoomState) {
	log.Printf("Room %d state changed to %T\n", r.Id, state)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.State = state
}


// Wait for players to connect
func (r *RoomWaitingForPlayers) acceptLoop() {
	for {
		select {
		case msg := <- r.r.Inch:
			switch msg.Typ {
			case MessagePlayerJoined:
				fmt.Printf("New player %d joined to %d\n", msg.PlayerId, r.r.Id)

				if r.r.Player1 != nil && msg.PlayerId == 2 {
					log.Printf("Player 2 joined room %d\n", r.r.Id)
					r.r.Player1.Inch <- ServerMsg{
						Typ: MessagePlayerJoined,
						PlayerId: msg.PlayerId,
					}
				}

				if r.r.Player1 != nil && r.r.Player2 != nil {
					r.r.setState(&Player1Turn{
						r: r.r,
					})

					go r.r.Start()

					r.r.Player1.Inch <- ServerMsg{
						Typ: MessagePlayer1Turn,
					}
					r.r.Player2.Inch <- ServerMsg{
						Typ: MessagePlayer1Turn,
					}	

					return
				}
			}
		}
	}
}

func (r *RoomWaitingForPlayers) Name() string {
	return "RoomWaitingForPlayers"
}

// A few considerations:
// A queue for all the incoming requests -> not needed since
// only one player goes at a time with one input
func (r *Player1Turn) acceptLoop() {
	for {
		select {
		case msg := <- r.r.Inch:
			switch msg.Typ {
			case MessagePlayer1Played:
				// Update state
				// Notify player2

				rc := r.r.RandomCard()
				if rc.Suit == "" && rc.Value == "" {
					// No cards left
					// Game over, check for winner
					var winner int16
					
					// Set state to game over
					r.r.setState(&GameOver{
						r: r.r,
					})

					go r.r.Start()

					if len(r.r.Player1Cards) > len(r.r.Player2Cards) {
						winner = r.r.Player1.Id
					} else {
						winner = r.r.Player2.Id
					}
					r.r.Player1.Inch <- ServerMsg{
						Typ: MessageGameOver,
						Winner: winner,
					}
					r.r.Player2.Inch <- ServerMsg{
						Typ: MessageGameOver,
						Winner: winner,
					}
					return
				}

				r.r.setState(&Player2Turn{
					r: r.r,
				})
				
				go r.r.Start()

				r.r.Player1Drawn = append(r.r.Player1Drawn, rc)

				r.r.Player1.Inch <- ServerMsg{
					Typ: MessagePlayer2Turn,
					Card: rc,
				}
				r.r.Player2.Inch <- ServerMsg{
					Typ: MessagePlayer2Turn,
					Card: rc,
				}
				return
			case MessagePlayer2Played:
				// Put on a queue?
				// or illegal
			case MessagePlayer2Turn:
				// Illegal
			case MessagePlayerJoined:
				panic("Illegal")
			}
		}
	}
}

func (r *Player1Turn) Name() string {
	return "Player1Turn"
}

func (r *Player2Turn) acceptLoop() {
	for {
		select {
		case msg := <- r.r.Inch:
			switch msg.Typ {
			case MessagePlayer2Played:
				// Update state
				// Notify players
				rc := r.r.RandomCard()
				if rc.Suit == "" && rc.Value == "" {
					// No cards left
					// Game over, check for winner
					var winner int16
					
					// Set state to game over
					r.r.setState(&GameOver{
						r: r.r,
					})

					go r.r.Start()

					if len(r.r.Player1Cards) > len(r.r.Player2Cards) {
						winner = r.r.Player1.Id
					} else {
						winner = r.r.Player2.Id
					}
					r.r.Player1.Inch <- ServerMsg{
						Typ: MessageGameOver,
						Winner: winner,
					}
					r.r.Player2.Inch <- ServerMsg{
						Typ: MessageGameOver,
						Winner: winner,
					}
					return
				}

				r.r.setState(&Player1Turn{
					r: r.r,
				})

				go r.r.Start()

				r.r.Player2Drawn = append(r.r.Player2Drawn, rc)

				if len(r.r.Player1Drawn) == len(r.r.Player2Drawn) {
					// Check for winner
					player2Card := r.r.Player2Drawn[len(r.r.Player2Drawn)-1]
					player1Card := r.r.Player1Drawn[len(r.r.Player1Drawn)-1]

					// Parse the values to integers
					p1Score, err := strconv.Atoi(player1Card.Value)
					if err != nil {
						panic(err)
					}
					p2Score, err := strconv.Atoi(player2Card.Value)
					if err != nil {
						panic(err)
					}

					if p1Score > p2Score {
						// Player 1 wins
						won := []Card{player1Card, player2Card}
						r.r.Player1Cards = append(r.r.Player1Cards, won...)

						r.r.Player1.Inch <- ServerMsg{
							Typ: MessagePlayer1Turn,
							Won: true,
							Card: rc,
						}
						r.r.Player2.Inch <- ServerMsg{
							Typ: MessagePlayer1Turn,
							Card: rc,
						}	

					} else if p1Score < p2Score {
						// Player 2 wins
						won := []Card{player1Card, player2Card}
						r.r.Player2Cards = append(r.r.Player2Cards, won...)

						r.r.Player1.Inch <- ServerMsg{
							Typ: MessagePlayer1Turn,
							Card: rc,
						}
						r.r.Player2.Inch <- ServerMsg{
							Typ: MessagePlayer1Turn,
							Won: true,
							Card: rc,
						}	
					} else {
						// Draw -> this is where the war occurs
						r.r.Player1.Inch <- ServerMsg{
							Typ: MessagePlayer1Turn,
							Card: rc,
						}
						r.r.Player2.Inch <- ServerMsg{
							Typ: MessagePlayer1Turn,
							Card: rc,
						}	
					}
				}

				return
			case MessagePlayer1Played:
				// Put on a queue?
				// or illegal
			case MessagePlayer1Turn:
				// Illegal
			case MessagePlayerJoined:
				panic("Illegal")
			}
		}
	}
}

func (r *Player2Turn) Name() string {
	return "Player2Turn"
}

func (r *GameOver) acceptLoop() {}

func (r *GameOver) Name() string {
	return "GameOver"
}