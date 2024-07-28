package server

import "time"

type Config struct {
	Host    string
	Port    int // generally int
	Timeout time.Duration
	MaxConn int
}
