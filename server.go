package main

import (
	"crypto/tls"
	"log"
)

func StartServer(address, certFile, keyFile string, db *Database, poolSize int) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		return err
	}

	defer listener.Close()

	pool := NewConnectionPool(poolSize)
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
