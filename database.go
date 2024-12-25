package main

import "sync"

type Database struct {
	mu    sync.RWMutex
	users map[string]string
	data  map[string]map[string]string
}

func NewDatabase() *Database {
	return &Database{
		users: make(map[string]string),
		data:  make(map[string]map[string]string),
	}
}
