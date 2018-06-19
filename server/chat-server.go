package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func serveClient(conn net.Conn) {

	message, _ := bufio.NewReader(conn).ReadString('\n')
	newmessage := strings.ToUpper(message)
	conn.Write([]byte(newmessage + "\n"))

}

func main() {

	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Errorf("Launching server error:", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf("Error handling connection!", err)
		} else {
			go serveClient(conn)
		}
	}
}
