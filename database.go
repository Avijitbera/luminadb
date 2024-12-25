package main

import (
	"sync"
	"time"
)

type Database struct {
	mu         sync.RWMutex
	users      map[string]string
	data       map[string]map[string]string
	expiration map[string]map[string]time.Time
}

func NewDatabase() *Database {
	return &Database{
		users:      make(map[string]string),
		data:       make(map[string]map[string]string),
		expiration: make(map[string]map[string]time.Time),
	}
}

func (db *Database) AddUser(username, password string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.users[username] = password
	db.data[username] = make(map[string]string)
	db.expiration[username] = make(map[string]time.Time)
}

func (db *Database) Authenticate(username, password string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	storedPassword, exists := db.users[username]
	return exists && storedPassword == password
}

func (db *Database) Set(username, key, value string, ttl time.Duration) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[username][key] = value
	if ttl > 0 {
		db.expiration[username][key] = time.Now().Add(ttl)
	} else {
		delete(db.expiration[username], key)
	}

}
