package main

type Config struct {
	Addr         string
	MaxRoomCount int
	RoomTimeout  int
}

func NewConfig(addr string) *Config {
	return &Config{
		Addr:        addr,
		RoomTimeout: 5,
	}
}