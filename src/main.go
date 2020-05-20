package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type initMsg struct {
	RequestID string `json:"requestId"`
	Type      int    `json:"type"`
	ModuleID  string `json:"moduleId"`
	Version   string `json:"version"`
}

type job struct {
	RequestID string    `json:"requestId"`
	Type      int       `json:"type"`
	Function  string    `json:"function"`
	Arguments arguments `json:"arguments"`
}

type arguments struct {
	Req string `json:"req"`
}

func registerModule(c *net.Conn, sub *sync.WaitGroup) {
	// Send module initialization request for each conn (register juno module)
	// Create unique requestID and moduleID
	// If not debug, dont print message, use io.Copy(ioutil.Discard) instead
	defer sub.Done()

	var i int = 0
	uuid := uuid.NewV4().String()
	iMsg, _ := json.Marshal(initMsg{
		RequestID: uuid,
		Type:      1,
		ModuleID:  "benchmark" + uuid,
		Version:   "1.0.0",
	})
	fmt.Fprintf(*c, string(iMsg)+"\n")
	scanner := bufio.NewScanner(*c)
	for scanner.Scan() {
		log.Println("Initial message " + strconv.Itoa(i) + scanner.Text())
		i++
		if i > 1 {
			break
		}
	}
	// log.Println("[INFO] Finished registering module...")
}

func callFunction(c *net.Conn, sub *sync.WaitGroup) {
	// Send call Function request to juno
	// Create unique requestID
	// UnMarshal the response to check if it was successfull, (data == "200 OK")
	// If not debug then do not store message use io.Copy(ioutil.Discard) instead
	defer sub.Done()

	msg, _ := json.Marshal(job{
		RequestID: uuid.NewV4().String(),
		Type:      3,
		Function:  "tester.test",
		Arguments: arguments{
			Req: "sent",
		},
	})
	start := time.Now()
	fmt.Fprintf(*c, string(msg)+"\n")
	message, err := bufio.NewReader(*c).ReadString('\n')
	if err != nil {
		log.Panic("[ERROR] Message recieve failed: ", err)
		return
	}
	log.Println(time.Since(start), message)
	return
}

func initialize(c chan net.Conn, n int) {
	var sub sync.WaitGroup

	for i := 0; i < n; i++ {
		conn, _ := net.Dial("tcp4", "127.0.0.1:4000")
		c <- conn
		sub.Add(1)
		go registerModule(&conn, &sub)
	}
	sub.Wait()
}

func runJobs(c *net.Conn, n int, main *sync.WaitGroup) {
	defer main.Done()

	var sub sync.WaitGroup
	for i := 0; i < n; i++ {
		sub.Add(1)
		go callFunction(c, &sub)
	}
	sub.Wait()

	// log.Println("[DEBUG] Closing connection...")
	(*c).Close()
}

func runtimeStats() {
	// Display info on throughput and other memStats
}

func main() {
	// ################################# CHANGE CONNECTIONS AND JOBS HERE #########################################
	var CONN int = 100
	var JOBS int = 100
	// ############################################################################################################

	var main sync.WaitGroup
	conns := make(chan net.Conn, CONN)

	initialize(conns, cap(conns))

	for i := 0; i < cap(conns); i++ {
		conn, more := <-conns
		if more {
			main.Add(1)
			go runJobs(&conn, JOBS, &main)
		} else {
			log.Println("[INFO] Spawned all connections.")
			close(conns)
			return
		}
	}
	main.Wait()

	log.Printf("[DEBUG] Program ended gracefully...\n")
}
