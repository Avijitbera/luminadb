package main

import (
	"net"
)

type ConnectionPool struct {
	jobs    chan net.Conn
	workers int
}

func NewConnectionPool(workers int) *ConnectionPool {
	return &ConnectionPool{
		jobs:    make(chan net.Conn, workers),
		workers: workers,
	}
}
