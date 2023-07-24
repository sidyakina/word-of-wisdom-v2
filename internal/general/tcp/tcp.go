package tcp

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	readDataBuffer = 1024
	TCP            = "tcp"
)

func CloseConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		log.Printf("failed to close connection with server: %v", err)
	}
}

func ReadWithDeadline(clientConn net.Conn, connReadDeadline time.Duration) ([]byte, error) {
	// set connection timeout for cases when connection open but message isn't send
	err := clientConn.SetReadDeadline(time.Now().UTC().Add(connReadDeadline))
	if err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	data := make([]byte, readDataBuffer)
	n, err := clientConn.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	// reset deadline because we don't wait more messages at this time
	err = clientConn.SetReadDeadline(time.Time{})
	if err != nil {
		return nil, fmt.Errorf("failed to reset read deadline: %w", err)
	}

	return data[:n], nil
}
