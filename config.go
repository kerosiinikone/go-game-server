package main

type Config struct {
	Addr         string
	MaxRoomCount int
}

func NewConfig(addr string) *Config {
	return &Config{
		Addr: addr,
	}
}