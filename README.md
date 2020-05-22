# juno-benchmark

Perform load testing on GO to measure request-response latency for [juno](https://github.com/bytesonus/juno). Supports concurrent tcp connections and requests <br>
[Juno Communication Protocols here](https://github.com/bytesonus/juno/blob/develop/docs/COMMUNICATION-PROTOCOL.md)

## Get started

* Start Juno as TCP4 server

* Start `juno-benchmark` CLI.

## Example

```sh
go run . -s 127.0.0.1:4000 -c 10 -r 10 -t 1m
```

## Help

```
NAME:
   juno-benchmark - Specify socket address and other parameters to start stress

USAGE:
   juno-benchmark [global options] command [command options] [arguments...]

DESCRIPTION:
   TCP load generator for juno. Written in Go

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --socket-address value, -s value  TCP v4 address to establish tcp connection (default: "127.0.0.1:4000")
   --connections value, -c value     Connections to keep open to the destination (default: 10)
   --rate value, -r value            Messages per second to send in a connection (default: 5)
   --time value, -t value            Exit after the specified amount of time. Valid time units are ns, us (or µs), ms, s, m, h (default: 10s)
   --help, -h                        show help (default: false)
```
