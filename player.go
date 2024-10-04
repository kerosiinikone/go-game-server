package main

type Player struct {
	Inch  chan struct{}
	outch chan struct{}
}