package main

import (
	"crypto/tls"
	"log"
)

func StartServer(address, certFile, keyFile string, db *Database) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	pool := NewConnectionPool(10)
	pool.Start(db)
	log.Printf("Server listening on %s with TLS", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		pool.AddConnection(conn)
	}

}
