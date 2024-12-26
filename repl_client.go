package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
)

func StartREPLClient(serverAddress, certFile string) {
	certPool := tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", serverAddress, &certPool)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to server. Enter commands (type 'exit' to quit):")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)

		}
		input = strings.TrimSpace(input)

		if input == "exit" {
			log.Println("Exiting...")
			break
		}
	}

}
