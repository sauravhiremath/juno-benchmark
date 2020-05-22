// Package Juno Benchmark is a TCP load generator for Juno. Written in Go
package main

import (
	"juno-benchmark/src"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	var (
		ADDR   string
		CONN   int64
		JOBS   int64
		TTL, _ = time.ParseDuration("10s")
	)

	app := cli.NewApp()
	app.Name = "juno-benchmark"
	app.Usage = "TCP load generator for juno. Written in Go"
	app.Description = "Benchmark throughputs for juno written in Go"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "socket-address",
			Aliases:     []string{"s"},
			Usage:       "TCP v4 address to establish tcp connection",
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
		&cli.DurationFlag{
			Name:        "time",
			Aliases:     []string{"t"},
			DefaultText: "10s",
			Usage:       "Exit after the specified amount of time. Valid time units are ns, us (or µs), ms, s, m, h",
			Required:    true,
			Destination: &TTL,
		},
	}
	app.Action = func(c *cli.Context) error {
		src.Start(ADDR, CONN, JOBS, TTL)
		go src.GracefulAbort()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	go src.GracefulAbort()
	forever := make(chan int)
	<-forever
}
