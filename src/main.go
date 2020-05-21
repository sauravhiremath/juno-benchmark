package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
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
	_, err := bufio.NewReader(*c).ReadString('\n')
	if err != nil {
		log.Panic("[ERROR] Message recieve failed: ", err)
		return
	}
	log.Println(time.Since(start))
	return
}

func initialize(c chan net.Conn, ADDR string, n int) {
	var sub sync.WaitGroup

	for i := 0; i < n; i++ {
		conn, _ := net.Dial("tcp4", ADDR)
		c <- conn
		sub.Add(1)
		go registerModule(&conn, &sub)
	}
	sub.Wait()
}

func runJobs(c *net.Conn, n int64, main *sync.WaitGroup) {
	defer main.Done()

	var sub sync.WaitGroup
	var i int64 = 0
	for ; i < n; i++ {
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

func start(ADDR string, CONN int64, JOBS int64) {
	var main sync.WaitGroup
	conns := make(chan net.Conn, CONN)

	initialize(conns, ADDR, cap(conns))

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

func main() {
	// ################################# CHANGE CONNECTIONS AND JOBS NUMBER HERE ##################################

	var (
		ADDR string
		CONN int64
		JOBS int64
	)

	// ############################################################################################################

	app := cli.NewApp()
	app.Name = "juno-benchmark"
	app.Usage = "Benchmark throughputs for juno. Written in Go"
	app.Description = "Benchmark throughputs for juno written in Go"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "socket-address",
			Aliases:     []string{"s"},
			Usage:       "Tcp v4 address to establish tcp connection",
			Required:    true,
			Destination: &ADDR,
		},
		&cli.Int64Flag{
			Name:        "connections",
			Aliases:     []string{"c"},
			Value:       1,
			Usage:       "Connections to keep open to the destination",
			Required:    true,
			Destination: &CONN,
		},
		&cli.Int64Flag{
			Name:        "rate",
			Aliases:     []string{"r"},
			Value:       1,
			Usage:       "Messages per second to send in a connection",
			Required:    true,
			Destination: &JOBS,
		},
	}
	app.Action = func(c *cli.Context) error {
		start(ADDR, CONN, JOBS)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
