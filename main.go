package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("/usr/src/app/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v\n", err)
	}
	defer file.Close()

	log.SetOutput(file)

	port := 3000

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		log.Fatalf("failed listening on port 3000: %v\n", err)
	}
	defer listener.Close()

	log.Printf("server listening on port %d\n", port)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("failed receiving tcp connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	log.Println("client connected")

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 10)

	for {
		buff := make([]byte, 1_000)
		readedBytes, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("client disconnected")
				return
			}

			log.Printf("error reading from the client: %v\n", err)

			continue
		}

		response := buff[:readedBytes]

		go processMsg(response)
	}
}

func processMsg(b []byte) {
	msg := string(b)

	log.Printf("received: %s\n", msg)
}
