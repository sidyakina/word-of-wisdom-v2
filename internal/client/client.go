package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"word-of-wisdom-v2/internal/general/challenger"
	"word-of-wisdom-v2/internal/general/tcp"
	"word-of-wisdom-v2/pkg/api"
)

type Solver interface {
	SolveChallenge(info challenger.ChallengeInfo) (string, error)
}

type Client struct {
	solver           Solver
	url              string
	connReadDeadline time.Duration
}

func New(url string, connReadDeadline time.Duration, solver Solver) *Client {
	return &Client{
		solver:           solver,
		url:              url,
		connReadDeadline: connReadDeadline,
	}
}

func (c *Client) GetQuote() (string, error) {
	serverConn, err := net.Dial(tcp.TCP, c.url)
	if err != nil {
		return "", fmt.Errorf("failed to connect with server on %v: %w", c.url, err)
	}

	defer func() {
		log.Printf("closing connection to server")
		tcp.CloseConnection(serverConn)
	}()

	log.Printf("connected with server")

	challenge, err := c.waitChallenge(serverConn)
	if err != nil {
		return "", fmt.Errorf("failed to received challenge: %w", err)
	}

	solution, err := c.solver.SolveChallenge(*challenge)
	if err != nil {
		return "", fmt.Errorf("failed to solve challenge: %w", err)
	}

	err = sendSolution(serverConn, solution)
	if err != nil {
		return "", fmt.Errorf("failed to send solution to server: %w", err)
	}

	quote, err := c.waitQuote(serverConn)
	if err != nil {
		return "", fmt.Errorf("failed to receive quote: %w", err)
	}

	return quote, nil
}

func (c *Client) waitChallenge(serverConn net.Conn) (*challenger.ChallengeInfo, error) {
	data, err := tcp.ReadWithDeadline(serverConn, c.connReadDeadline)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	message := api.Challenge{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal challenge %s: %w", data, err)
	}

	challenge := challenger.ChallengeInfo{
		RandomString:          message.RandomString,
		NumberLeadingZeros:    message.NumberLeadingZeros,
		SolutionNumberSymbols: message.SolutionNumberSymbols,
	}

	log.Printf("challenge received: %+v", challenge)

	return &challenge, nil
}

func sendSolution(serverConn net.Conn, solution string) error {
	message := api.Solution{
		Solution: solution,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal api.Solution: %w", err)
	}

	log.Printf("send solution %s to server", data)

	_, err = serverConn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send solution to server: %w", err)
	}

	return nil
}

func (c *Client) waitQuote(serverConn net.Conn) (string, error) {
	data, err := tcp.ReadWithDeadline(serverConn, c.connReadDeadline)
	if err != nil {
		return "", fmt.Errorf("failed to read message: %w", err)
	}

	message := api.Quote{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal quote %s: %w", data, err)
	}

	log.Printf("quote received: %+v", message)

	return message.Quote, nil
}
