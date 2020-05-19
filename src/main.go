package main

import (
	"bufio"
	"fmt"
	"net"
)

var initialize = "{\"requestId\": \"jb-1\", \"type\": 1, \"moduleId\": \"juno-benchmark\", \"version\": \"1.0.0\"}\n"
var message = "{\"requestId\": \"test\", \"type\": 3, \"function\": \"tester.test\", \"arguments\": { \"req\": \"sent\" } }\n"

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:4000")
	// Send and receive initialize message here
	fmt.Fprintf(conn, initialize)
	initializeRes, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf(initializeRes)

	// Send data here
	fmt.Fprintf(conn, message)

	// wait for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf(message)

}
