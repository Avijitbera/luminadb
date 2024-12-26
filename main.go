package main

import (
	"log"
	"time"
)

func main() {
	address := ":6379"
	certFile := "server.crt"
	keyFile := "server.key"
	poolSize := 10

	db := NewDatabase()
	db.AddUser("user1", "password1")
	db.AddUser("user2", "password2")
	db.AddUser("user3", "password3")

	go func() {
		for {
			time.Sleep(1 * time.Second)
			db.CleanupExpiredKeys()
		}
	}()

	if err := StartServer(address, certFile, keyFile, db, poolSize); err != nil {
		log.Fatal(err)
	}

}
