# iperfweb
A simple webserver with html files that renders the iperf3 output on a guage using google visualization library


# Simple Example
- On a HOSTA machine run `iperf3 -s` - Runs iperf3 in server mode (default port)
- On the same HOSTA or different HOSTB 
```
go build server.go
./server HOSTA
```

- Open a Browser
open http://HOSTA:8080


