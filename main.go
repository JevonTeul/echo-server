package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	port    = flag.Int("port", 4000, "TCP port to listen on")
	timeout = flag.Int("timeout", 30, "Client inactivity timeout in seconds")
)

const (
	goodbyeMsg = "Goodbye! Closing connection...\n\n"
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server started on port %d with %d second timeout", *port, *timeout)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	clientAddr := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	log.Printf("Client connected: %s", clientAddr)
	defer func() {
		conn.Close()
		log.Printf("Client disconnected: %s", clientAddr)
	}()

	logFile, err := os.Create(fmt.Sprintf("%s.log", strings.ReplaceAll(clientAddr, ":", "_")))
	if err != nil {
		log.Printf("Error creating log file: %v", err)
		return
	}
	defer logFile.Close()

	reader := bufio.NewReaderSize(conn, 1024)
	writer := bufio.NewWriter(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(*timeout) * time.Second))

		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("%s: Connection closed by client", clientAddr)
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("%s: Connection timed out", clientAddr)
				writer.WriteString("Connection timed out due to inactivity\n")
				writer.Flush()
			} else {
				log.Printf("%s: Read error: %v", clientAddr, err)
			}
			return
		}

		message = strings.TrimSpace(message)
		log.Printf("%s >> %s", clientAddr, message)
		fmt.Fprintf(logFile, "[%s] %s\n", time.Now().Format(time.RFC3339), message)

		response, shouldClose := processMessage(message)

		if _, err := writer.WriteString(response); err != nil {
			log.Printf("%s: Write error: %v", clientAddr, err)
			return
		}

		// Critical fix: Flush before closing connection
		if shouldClose {
			writer.Flush()
			time.Sleep(100 * time.Millisecond)
		}

		if err := writer.Flush(); err != nil {
			log.Printf("%s: Flush error: %v", clientAddr, err)
			return
		}
	}
}

func processMessage(msg string) (string, bool) {
	msg = strings.TrimSpace(strings.ToLower(msg))
	switch {
	case msg == "":
		return "Say something...\n\n", false
	case msg == "hello":
		return "Hi there!\n\n", false
	case msg == "bye" || strings.HasPrefix(msg, "/quit"):
		return goodbyeMsg, true
	case strings.HasPrefix(msg, "/time"):
		return time.Now().Format("15:04:05") + "\n\n", false
	case strings.HasPrefix(msg, "/echo "):
		return strings.TrimPrefix(msg, "/echo ") + "\n\n", false
	default:
		return msg + "\n\n", false
	}
}
