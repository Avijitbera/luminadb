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
func (p *ConnectionPool) AddConnection(conn net.Conn) {
	p.jobs <- conn
}

func (p *ConnectionPool) Start(db *Database) {
	for i := 0; i < p.workers; i++ {
		go p.worker(db)
	}
}

func (p *ConnectionPool) worker(db *Database) {

}
