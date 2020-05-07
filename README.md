# juno-benchmark

Benchmarking tool to measure request-response for juno. Supports multi-threading

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
