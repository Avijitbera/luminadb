package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	// address := ":6379"
	// certFile := "server.crt"
	// keyFile := "server.key"
	// poolSize := 10

	mode := flag.String("mode", "server", "Mode: server or repl")
	address := flag.String("address", ":6379", "Server address")
	certFile := flag.String("cert", "server.crt", "TLS certificate file")
	keyFile := flag.String("key", "server.key", "TLS key file")
	poolSize := flag.Int("pool", 10, "Connection pool size")

	flag.Parse()
	if *mode == "server" {

		db := NewDatabase()
		db.AddUser("user1", "password1")
		db.AddUser("user2", "password2")
		db.AddUser("user3", "password3")

		go func() {
			for {
				db.CleanupExpiredKeys()
				time.Sleep(1 * time.Second)
			}
		}()

		if err := StartServer(*address, *certFile, *keyFile, db, *poolSize); err != nil {
			log.Fatal(err)
		}
	} else if *mode == "repl" {
		log.Println("Starting REPL client...")
		StartREPLClient(*address, *certFile)
	} else {
		log.Println("Unknown mode: ", *mode)
	}

}
