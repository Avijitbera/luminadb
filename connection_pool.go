package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
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
	for conn := range p.jobs {
		handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *Database) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.Write([]byte("Enter credentials (username password): "))
	authLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	authParts := strings.Split(strings.TrimSpace(authLine), " ")
	if len(authParts) != 2 || !db.Authenticate(authParts[0], authParts[1]) {
		conn.Write([]byte("Invalid credentials\n"))
		return
	}
	username := authParts[0]
	conn.Write([]byte("Authenticated\n"))
	for {

		conn.Write([]byte("> "))
		commandLine, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		command := strings.Fields(strings.TrimSpace(commandLine))
		if len(command) == 0 {
			conn.Write([]byte("Enter credentials (username password): "))
			continue
		}
		switch strings.ToUpper(command[0]) {
		case "SET":
			fmt.Println(command)
			if len(command) < 3 {
				conn.Write([]byte("Invalid command\n"))

			} else {
				key, value := command[1], command[2]
				ttlSeconds, err := strconv.Atoi(command[3])
				if err != nil {
					conn.Write([]byte("Invalid TTL\n"))
					continue
				}
				ttl := time.Duration(ttlSeconds) * time.Second
				db.Set(username, key, value, ttl)
				conn.Write([]byte("OK\n"))
			}

		case "GET":
			if len(command) != 2 {
				conn.Write([]byte("Invalid command\n"))
			} else {
				key := command[1]
				value, exists := db.Get(username, key)
				if exists {
					conn.Write([]byte(value + "\n"))
				} else {
					conn.Write([]byte("Key not found\n"))
				}
			}
		case "DEL":
			if len(command) != 2 {
				conn.Write([]byte("Invalid command\n"))
			} else {
				key := command[1]
				if db.Delete(username, key) {
					conn.Write([]byte("OK\n"))
				} else {
					conn.Write([]byte("Key not found\n"))
				}
			}
		case "EXIT":
			conn.Write([]byte("Goodbye\n"))
			return
		default:
			conn.Write([]byte("Invalid command\n"))
		}
	}
}
