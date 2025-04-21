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

	// Create log file
	logFile, err := os.Create(fmt.Sprintf("%s.log", strings.ReplaceAll(clientAddr, ":", "_")))
	if err != nil {
		log.Printf("Error creating log file: %v", err)
		return
	}
	defer logFile.Close()

	reader := bufio.NewReaderSize(conn, 1024) // Set max message size
	writer := bufio.NewWriter(conn)

	for {
		// Set read timeout
		conn.SetReadDeadline(time.Now().Add(time.Duration(*timeout) * time.Second))

		// Read until newline
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

		// Trim and process message
		message = strings.TrimSpace(message)
		log.Printf("%s >> %s", clientAddr, message)
		fmt.Fprintf(logFile, "[%s] %s\n", time.Now().Format(time.RFC3339), message)

		// Process message and get response
		response, shouldClose := processMessage(message)

		// Write response
		if _, err := writer.WriteString(response + "\n"); err != nil {
			log.Printf("%s: Write error: %v", clientAddr, err)
			return
		}
		if err := writer.Flush(); err != nil {
			log.Printf("%s: Flush error: %v", clientAddr, err)
			return
		}

		if shouldClose {
			return
		}
	}
}

func processMessage(msg string) (string, bool) {
	msg = strings.TrimSpace(msg)
	switch {
	case msg == "":
		return "Say something...\n\n", false
	case msg == "hello":
		return "Hi there!\n\n", false
	case msg == "bye":
		return "Goodbye!\n\n", true
	case strings.HasPrefix(msg, "/time"):
		return time.Now().Format("15:04:05") + "\n\n", false
	case strings.HasPrefix(msg, "/quit"):
		return "Closing connection\n\n", true
	case strings.HasPrefix(msg, "/echo "):
		return strings.TrimPrefix(msg, "/echo ") + "\n\n", false
	default:
		return msg + "\n\n", false
	}
}

func handleCommand(cmd string) (string, bool) {
	parts := strings.SplitN(cmd, " ", 2)
	switch parts[0] {
	case "/time":
		return time.Now().Format("2006-01-02 15:04:05 MST"), false
	case "/quit":
		return "Closing connection", true
	case "/echo":
		if len(parts) > 1 {
			return parts[1], false
		}
		return "Echo what?", false
	default:
		return "Unknown command", false
	}
}
