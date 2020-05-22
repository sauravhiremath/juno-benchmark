# juno-benchmark

Perform load testing on GO to measure request-response latency for [juno](https://github.com/bytesonus/juno). Supports concurrent tcp connections and requests <br>
[Juno Communication Protocols here](https://github.com/bytesonus/juno/blob/develop/docs/COMMUNICATION-PROTOCOL.md)

## Get started

* Start Juno on a TCPv4 server

* Run the following juno test module, with flag `-s` with juno socket-address as `<host>:<port>`. <br>
This servers as a module, juno-benchmark will call the function from.
_[ This will be changed on next fix]_

```javascript
const minimist = require('minimist');
const JunoModule = require('juno-node');

const argv = minimist(process.argv.slice(2));

async function main() {
    let module = await JunoModule.default(argv['s'] || argv['socket-location']);
    await module.initialize('tester', '1.0.0');

    const resp = await module.declareFunction('test', (arg) => { return '200 OK' });
    console.log('Tester module test declared.');
    console.log(JSON.stringify(resp));
}

main();
```

* Start `juno-benchmark` CLI. <br>
Ex command: <br>
`go run main.go -s 127.0.0.1:4000 -c 10 -r 10 -t 1m`

## Help

```
NAME:
   juno-benchmark - Specify socket address and other parameters to start stress

USAGE:
   main [global options] command [command options] [arguments...]

DESCRIPTION:
   TCP load generator for juno. Written in Go

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --socket-address value, -s value  TCP v4 address to establish tcp connection
   --connections value, -c value     Connections to keep open to the destination (default: 1)
   --rate value, -r value            Messages per second to send in a connection (default: 1)
   --time value, -t value            Exit after the specified amount of time. Valid time units are ns, us (or µs), ms, s, m, h (default: 10s)
   --help, -h                        show help (default: false)
```