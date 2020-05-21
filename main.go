// Package Juno Benchmark is a TCP load generator for Juno. Written in Go
package main

import (
	"juno-benchmark/src"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var (
		ADDR string
		CONN int64
		JOBS int64
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
	}
	app.Action = func(c *cli.Context) error {
		src.Start(ADDR, CONN, JOBS)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
