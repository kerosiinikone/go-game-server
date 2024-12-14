package main

import "math/rand"

type Player1Turn struct {
	r *Room
}

type Player2Turn struct {
	r *Room
}

type Card struct {
	Suit  string
	Value string
}

var deck = []Card{
	{Suit: "Spades", Value: "2"},
	{Suit: "Spades", Value: "3"},
	{Suit: "Spades", Value: "4"},
	{Suit: "Spades", Value: "5"},
	{Suit: "Spades", Value: "6"},
	{Suit: "Spades", Value: "7"},
	{Suit: "Spades", Value: "8"},
	{Suit: "Spades", Value: "9"},
	{Suit: "Spades", Value: "10"},
	{Suit: "Spades", Value: "11"},
	{Suit: "Spades", Value: "12"},
	{Suit: "Spades", Value: "13"},
	{Suit: "Spades", Value: "14"},
	{Suit: "Hearts", Value: "2"},
	{Suit: "Hearts", Value: "3"},
	{Suit: "Hearts", Value: "4"},
	{Suit: "Hearts", Value: "5"},
	{Suit: "Hearts", Value: "6"},
	{Suit: "Hearts", Value: "7"},
	{Suit: "Hearts", Value: "8"},
	{Suit: "Hearts", Value: "9"},
	{Suit: "Hearts", Value: "10"},
	{Suit: "Hearts", Value: "11"},
	{Suit: "Hearts", Value: "12"},
	{Suit: "Hearts", Value: "13"},
	{Suit: "Hearts", Value: "14"},
	{Suit: "Clubs", Value: "2"},
	{Suit: "Clubs", Value: "3"},
	{Suit: "Clubs", Value: "4"},
	{Suit: "Clubs", Value: "5"},
	{Suit: "Clubs", Value: "6"},
	{Suit: "Clubs", Value: "7"},
	{Suit: "Clubs", Value: "8"},
	{Suit: "Clubs", Value: "9"},
	{Suit: "Clubs", Value: "10"},
	{Suit: "Clubs", Value: "11"},
	{Suit: "Clubs", Value: "12"},
	{Suit: "Clubs", Value: "13"},
	{Suit: "Clubs", Value: "14"},
	{Suit: "Diamonds", Value: "2"},
	{Suit: "Diamonds", Value: "3"},
	{Suit: "Diamonds", Value: "4"},
	{Suit: "Diamonds", Value: "5"},
	{Suit: "Diamonds", Value: "6"},
	{Suit: "Diamonds", Value: "7"},
	{Suit: "Diamonds", Value: "8"},
	{Suit: "Diamonds", Value: "9"},
	{Suit: "Diamonds", Value: "10"},
	{Suit: "Diamonds", Value: "11"},
	{Suit: "Diamonds", Value: "12"},
	{Suit: "Diamonds", Value: "13"},
	{Suit: "Diamonds", Value: "14"},
}

func (r *Room) RandomCard() Card {
	if len(r.Deck) == 0 {
		r.Deck = deck
	}
	randIndex := rand.Intn(len(r.Deck))
	c := r.Deck[randIndex]
	r.Deck = append(r.Deck[:randIndex], r.Deck[randIndex+1:]...)
	return c
}