package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type initMsg struct {
	RequestID string `json:"requestId"`
	Type      int    `json:"type"`
	ModuleID  string `json:"moduleId"`
	Version   string `json:"version"`
}

type funcCall struct {
	RequestID string    `json:"requestId"`
	Type      int       `json:"type"`
	Function  string    `json:"function"`
	Arguments arguments `json:"arguments"`
}

type arguments struct {
	Req string `json:"req"`
}

// var initialMsg = "{\"requestId\": \"jb-1\", \"type\": 1, \"moduleId\": \"juno-benchmark\", \"version\": \"1.0.0\"}\n"
// var message = "{\"requestId\": \"test\", \"type\": 3, \"function\": \"tester.test\", \"arguments\": { \"req\": \"sent\" } }\n"

func registerModule(c *net.Conn) {
	// Send module initialization request for each conn (register juno module)
	// Create unique requestID and moduleID
	// If not debug, dont print message, use io.Copy(ioutil.Discard) instead
	var i int = 0
	iMsg, _ := json.Marshal(initMsg{
		RequestID: uuid.NewV4().String(),
		Type:      1,
		ModuleID:  "juno-benchmark",
		Version:   "1.0.0",
	})
	fmt.Println(string(iMsg))
	fmt.Fprintf(*c, string(iMsg)+"\n")
	scanner := bufio.NewScanner(*c)
	for scanner.Scan() {
		fmt.Println("Initial message " + strconv.Itoa(i) + scanner.Text())
		i++
		if i > 1 {
			break
		}
	}
	return
}

func callFunction(c *net.Conn) {
	// Send callFunction payload to juno
	// Create unique requestID
	// Parse response to check if it was successfull, (data == "200 OK")
	// If not debug then do not store message use io.Copy(ioutil.Discard) instead
	msg, _ := json.Marshal(funcCall{
		RequestID: uuid.NewV4().String(),
		Type:      3,
		Function:  "tester.test",
		Arguments: arguments{
			Req: "sent",
		},
	})
	start := time.Now()
	fmt.Fprintf(*c, string(msg)+"\n")
	message, _ := bufio.NewReader(*c).ReadString('\n')
	fmt.Println(time.Since(start))
	fmt.Printf(message)
	(*c).Close()
	return
}

func spawnConn(c chan net.Conn, n int) {
	for i := 0; i < n; i++ {
		conn, _ := net.Dial("tcp", "127.0.0.1:4000")
		c <- conn
	}
}

func runtimeStats() {
	// Display info on throughput and other memStats
}

// func doWork(id int, j job) {
// 	fmt.Printf("worker%d: started %s, working for %f seconds\n", id, j.requestID, j.throughput.Seconds())
// }

// func reqHandler(jobs <-chan job) {
// 	conn, _ := net.Dial("tcp", "127.0.0.1:4000")
// 	// Send and receive initialize message here
// 	fmt.Fprintf(conn, initialMsg)
// 	initializeRes, _ := bufio.NewReader(conn).ReadString('\n')
// 	fmt.Printf(initializeRes)
// }

func main() {
	// var wg sync.WaitGroup
	conns := make(chan net.Conn, 10)
	spawnConn(conns, cap(conns))
	for conn := range conns {
		fmt.Printf("In for...\n")
		registerModule(&conn)
		callFunction(&conn)
		fmt.Printf("finished all...\n")
	}
	fmt.Printf("finished all outside...\n")
	close(conns)
	// for conn := range conns {
	// 	defer conn.Close()
	// }
}
