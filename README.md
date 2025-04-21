# Enhanced TCP Echo Server

A concurrent TCP server with logging, command processing, and client management features implemented in Go.

## How to Run

### Prerequisites
- Go 1.16+
- Linux/Unix environment
- Netcat (`nc`) for testing

### Installation & Execution

# Clone repository (if applicable)
`git clone https://github.com/yourusername/tcp-echo-server.git`
`cd tcp-echo-server`

# Build and run with defaults (port 4000)
`go build -o echoserver && ./echoserver`

# Custom configuration
`./echoserver -port 5000 -timeout 45`


## Educational Highlights

### Most Enriching Feature  
**Concurrent Connection Handling**  
Implementing goroutines to manage multiple clients simultaneously provided deep insight into:  
- Race condition prevention  
- Goroutine lifecycle management  
- Shared resource handling (log files)  
- Channel-free concurrency patterns  

### Most Research-Intensive Feature  
**Graceful Connection Termination**  
Required significant research into:  
- TCP/IP connection states  
- Buffered I/O flushing guarantees  
- Network timeout best practices  
- Error handling edge cases (EOF vs. timeout vs. reset)  

**Key resources:**  
- [Go net package documentation](https://pkg.go.dev/net)  
- [TCP/IP Illustrated Vol.1](https://en.wikipedia.org/wiki/TCP/IP_Illustrated)  
- [Go Blog: Concurrency Patterns](https://go.dev/blog/pipelines)  

# Demo Video
https://youtu.be/kUOTWxD52D0