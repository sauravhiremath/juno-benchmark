// Package src is a collection of functions for juno-benchmark
package src

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	conns     chan net.Conn
	moduleIDs chan string
	stats     chan time.Duration
)

/*

QueueConnections handles registering `CONN` number of connections.

All connections are spawned concurrently and the function
exists on successfully registration of spawned connections

*/
func QueueConnections(ADDR string, n int) {
	var sub sync.WaitGroup

	fmt.Printf("Initializing connections...   \r")
	for i := 0; i < n; i++ {
		conn, _ := net.Dial("tcp4", ADDR)
		conns <- conn
		sub.Add(1)
		go func() {
			RegisterModule(&conn)
			DeclareFunction(&conn)
			sub.Done()
		}()
	}
	sub.Wait()
}

/*

QueueFunctionCalls handles transmitting N messages over 1 connection.

All messages are spawned concurrently and the function
exists on getting successfull responses from all transmissions sent

*/
func QueueFunctionCalls(c *net.Conn, rate int64, TTL time.Duration, main *sync.WaitGroup) {
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
			log.Println("[DEBUG] Closing connection...")
			return
		case <-ticker.C:
			select {
			case moduleID := <-moduleIDs:
				sub.Add(1)
				go FunctionCall(c, &sub, moduleID)
				go func() {
					moduleIDs <- moduleID
				}()
				// fmt.Println("finished one fn call")
			default:
				// fmt.Println("end of first loop")
				break
			}
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
		fmt.Printf("Current Latency: %v    \t    Average Latency: %v µs       \r", latency, avgLatency)
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
	moduleIDs = make(chan string, CONN)

	main.Add(1)
	go RuntimeStats(&main)

	QueueConnections(ADDR, cap(conns))

	for i := 0; i < cap(conns); i++ {
		conn, more := <-conns
		if more {
			main.Add(1)
			go QueueFunctionCalls(&conn, RATE, TTL, &main)
		} else {
			close(conns)
			break
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
		fmt.Println("Process aborting... \n[Press Ctrl+C to force exit]")
		close(conns)
		os.Exit(0)
	}()
}
