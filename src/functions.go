// Package src is a collection of functions for juno-benchmark
package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
)

var conns chan net.Conn
var stats chan time.Duration

/*

RegisterModule initializes the given connection.

Sends the Marshaled register module request to Juno

*/
func RegisterModule(c *net.Conn, sub *sync.WaitGroup) {
	defer sub.Done()

	var i int = 0
	uuid := uuid.NewV4().String()
	iMsg, _ := json.Marshal(InitMsg{
		RequestID: uuid,
		Type:      1,
		ModuleID:  "benchmark" + uuid,
		Version:   "1.0.0",
	})
	fmt.Fprintf(*c, string(iMsg)+"\n")
	scanner := bufio.NewScanner(*c)
	for scanner.Scan() {
		// log.Println("Initial message " + strconv.Itoa(i) + scanner.Text())
		i++
		if i > 1 {
			break
		}
	}
	// log.Println("[INFO] Finished registering module...")
}

/*

SendMessage handles message transmission over given connection.

Creates new request ID (UUID V4) for each transmission.

Calculates latency between request and response

*/
func SendMessage(c *net.Conn, sub *sync.WaitGroup) {
	// Create unique requestID
	defer sub.Done()

	msg, _ := json.Marshal(Job{
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
	// log.Println(time.Since(start))
	stats <- time.Since(start)
	return
}

/*

QueueConns handles registering N connections.

All connections are spawned concurrently and the function
exists on successfully registration of spawned connections

*/
func QueueConns(ADDR string, n int) {
	var sub sync.WaitGroup

	for i := 0; i < n; i++ {
		conn, _ := net.Dial("tcp4", ADDR)
		conns <- conn
		sub.Add(1)
		go RegisterModule(&conn, &sub)
	}
	sub.Wait()
}

/*

QueueJobs handles transmitting N messages over 1 connection.

All messages are spawned concurrently and the function
exists on getting successfull responses from all transmissions sent

*/
func QueueJobs(c *net.Conn, rate int64, TTL time.Duration, main *sync.WaitGroup) {
	defer main.Done()

	var sub sync.WaitGroup
	done := make(chan bool)

	ticker := time.NewTicker(time.Duration(int64(time.Second) / rate))
	defer ticker.Stop()

	go func() {
		time.Sleep(TTL)
		done <- true
	}()
	for {
		select {
		case <-done:
			ticker.Stop()
			sub.Wait()
			(*c).Close()
			log.Println("[DEBUG] Closing connection...")
			return
		case _ = <-ticker.C:
			sub.Add(1)
			go SendMessage(c, &sub)
		}
	}
}

/*

RuntimeStats displays info on latency and other memStats

Returns:

Current latency as time.Duration()
Average latency in microseconds (µ)

*/
func RuntimeStats(main *sync.WaitGroup) {
	defer main.Done()

	var avgLatency int64 = 0
	var count int64 = 0

	for latency := range stats {
		fmt.Printf("Current Latency: %v  \t Average Latency: %v µs   \r", latency, avgLatency)
		avgLatency = (avgLatency*count + latency.Microseconds()) / (count + 1)
		count++
	}
}

/*

Start runs the task specified in the arguments. Exits upon successfull completion of the task.

Closes all open connections

*/
func Start(ADDR string, CONN, RATE int64, TTL time.Duration) {
	var main sync.WaitGroup
	conns = make(chan net.Conn, CONN)
	stats = make(chan time.Duration)

	main.Add(1)
	go RuntimeStats(&main)

	QueueConns(ADDR, cap(conns))

	for i := 0; i < cap(conns); i++ {
		conn, more := <-conns
		if more {
			main.Add(1)
			go QueueJobs(&conn, RATE, TTL, &main)
		} else {
			log.Println("[INFO] Spawned all connections.")
			close(conns)
			return
		}
	}
	main.Wait()
	log.Printf("[DEBUG] Program ended gracefully...\n")
}

/*

GracefulAbort closes all the open connections if any and returns back end channel

Listens to SIGTERM and SIGNINT to trigger process abortion

*/
func GracefulAbort() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s
		fmt.Println("Process aborting... [Press Ctrl+C to force exit]")
		close(conns)
		os.Exit(0)
	}()
}
