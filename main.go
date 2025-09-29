package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"
)

func main() {
	port := 3000

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: port,
	})
	if err != nil {
		slog.Error("failed listening on port 3005",
			slog.Int("port", port),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer listener.Close()

	slog.Info("server listening", slog.Int("port", port))

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			slog.Error("failed receiving tcp connection",
				slog.String("error", err.Error()),
			)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	slog.Info("client connected")

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 10)

	for {
		buff := make([]byte, 1_000)
		readedBytes, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("client disconnected")
				return
			}

			slog.Error("error reading from the client",
				slog.String("error", err.Error()),
			)
			continue
		}

		response := buff[:readedBytes]

		fmt.Println(response)
	}
}
