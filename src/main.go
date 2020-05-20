package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var initialize = "{\"requestId\": \"jb-1\", \"type\": 1, \"moduleId\": \"juno-benchmark\", \"version\": \"1.0.0\"}\n"
var message = "{\"requestId\": \"test\", \"type\": 3, \"function\": \"tester.test\", \"arguments\": { \"req\": \"sent\" } }\n"

type job struct {
	conn       net.Conn
	requestID  string
	throughput time.Duration
}

func spawnConn() {
	// Create a go routine to spawn connections to server
}

func initReq() {
	// Send module initialization requests for each conn
}

func tCPConnWrite() {

}

func tCPConnRead() {

}

func runtimeStats() {
	// Display info on throughput and other memStats
}
func doWork(id int, j job) {
	fmt.Printf("worker%d: started %s, working for %f seconds\n", id, j.requestID, j.throughput.Seconds())
}

func reqHandler(jobs <-chan job) {
	conn, _ := net.Dial("tcp", "127.0.0.1:4000")
	// Send and receive initialize message here
	fmt.Fprintf(conn, initialize)
	initializeRes, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf(initializeRes)
}

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:4000")
	// Send and receive initialize message here
	fmt.Fprintf(conn, initialize)
	initializeRes, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf(initializeRes)

	//Queue Jobs here

	// Send data here
	fmt.Fprintf(conn, message)

	// wait for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf(message)

}
