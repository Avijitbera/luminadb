package main

import (
	"fmt"
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
	db := &Database{
		users:      make(map[string]string),
		data:       make(map[string]map[string]string),
		expiration: make(map[string]map[string]time.Time),
	}
	return db
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
	fmt.Println(db.users)
	fmt.Println(username, password)
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

func (db *Database) Get(username, key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if exireTime, exists := db.expiration[username][key]; exists && exireTime.Before(time.Now()) {
		return "", false
	}
	value, exists := db.data[username][key]
	return value, exists
}

func (db *Database) Delete(username, key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, exists := db.data[username][key]; exists {
		delete(db.data[username], key)
		delete(db.expiration[username], key)
		return true
	}
	return false
}

func (db *Database) CleanupExpiredKeys() {
	db.mu.Lock()
	defer db.mu.Unlock()
	now := time.Now()
	for user, keys := range db.expiration {
		for key, expireTime := range keys {
			if now.After(expireTime) {
				delete(db.data[user], key)
				delete(db.expiration[user], key)
			}
		}
	}
}
