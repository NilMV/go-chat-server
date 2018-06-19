package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", message)
}

func main() {

	clientCount := 0

	allClients := make(map[net.Conn]int)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	messages := make(chan string)

	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", ":6000")
	check(err, "Server is ready.")

	// Tell the server to accept connections forever
	// and push new connections into the newConnections channel.
	//
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			newConnections <- conn
		}
	}()

	for {

		select {
		case conn := <-newConnections:

			log.Printf("Accepted new client, #%d", clientCount)

			allClients[conn] = clientCount
			clientCount += 1
			go func(conn net.Conn, clientId int) {
				reader := bufio.NewReader(conn)
				for {
					incoming, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					messages <- fmt.Sprintf("Client %d > %s", clientId, incoming)
				}
				deadConnections <- conn

			}(conn, allClients[conn])

		case message := <-messages:

			for conn, _ := range allClients {

				go func(conn net.Conn, message string) {
					_, err := conn.Write([]byte(message + "\n"))

					if err != nil {
						deadConnections <- conn
					}
				}(conn, message)
			}
			log.Printf("New message: %s", message)
			log.Printf("Broadcast to %d clients", len(allClients))

		case conn := <-deadConnections:
			log.Printf("Client %d disconnected", allClients[conn])
			delete(allClients, conn)
		}
	}
}
