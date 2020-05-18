package main

import (
	"bufio"
	"fmt"
	"net"
)

func connect() {
	conn, _ := net.Dial("tcp", "127.0.0.1:4000")

	// Send data here
	fmt.Fprintf(conn, "hello\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		// handle error
	}
	fmt.Printf(status)
}

func main() {
	connect()
}
