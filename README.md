# juno-benchmark

Perform load testing to measure request-response latency for [juno](https://github.com/bytesonus/juno). Supports multi-threading <br>
[Juno Communication Protocols here](https://github.com/bytesonus/juno/blob/develop/docs/COMMUNICATION-PROTOCOL.md)

### Tested on Linux
```bash
cc socket.c -o client -lpthread
```

Run the binary with juno running in background. Uhh
```
./client
```

Returns
* Request-response throughput
* Moving average of the throughputs

### For windows
Replace `\n` with `\r\n`. Feel free to contribute more changes
