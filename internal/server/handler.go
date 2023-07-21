package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"word-of-wisdom-v2/internal/general"
	"word-of-wisdom-v2/pkg/api"
)

const readDataBuffer = 512

func (s *Server) HandleNewConnect(clientConn net.Conn) {
	defer func() {
		log.Printf("closing connect to client")

		err := clientConn.Close()
		if err != nil {
			log.Printf("failed to close connect to client: %v", err)
		}
	}()

	challenge, err := s.sendChallenge(clientConn)
	if err != nil {
		log.Printf("failed to send challenge: %v", err)

		return
	}

	solution, err := s.waitSolution(clientConn)
	if err != nil {
		log.Printf("failed to send challenge: %v", err)

		return
	}

	if !s.challenger.ValidateSolution(*challenge, solution) {
		log.Printf("solution %v is wrong", solution)

		return
	}

	err = s.sendQuote(clientConn)
	if err != nil {
		log.Printf("failed to send quote: %v", err)

		return
	}

	log.Printf("client conn successfully handled")
}

func (s *Server) sendChallenge(clientConn net.Conn) (*general.ChallengeInfo, error) {
	challenge := s.challenger.GetChallenge()
	log.Printf("challenge for new connection: %+v", challenge)

	message := api.Challenge{
		RandomString:       challenge.RandomString,
		NumberLeadingZeros: challenge.NumberLeadingZeros,
		NumberSymbols:      challenge.NumberSymbols,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal api.Challenge: %w", err)
	}

	_, err = clientConn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send challenge to client: %w", err)
	}

	return &challenge, nil
}

func (s *Server) waitSolution(clientConn net.Conn) (string, error) {
	data, err := readWithDeadline(clientConn, s.connReadDeadline)
	if err != nil {
		return "", err
	}

	message := api.Solution{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal data %s to solution: %w", data, err)
	}

	return message.RandomString, nil
}

func (s *Server) sendQuote(clientConn net.Conn) error {
	quote, err := s.quotesRepo.GetQuote()
	if err != nil {
		return fmt.Errorf("failed to get quote: %w", err)
	}

	message := api.Quote{
		Quote: quote,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal api.Quote: %w", err)
	}

	log.Printf("send quote %s to client", data)

	_, err = clientConn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send quote to client: %w", err)
	}

	return nil
}

func readWithDeadline(clientConn net.Conn, connReadDeadline time.Duration) ([]byte, error) {
	// set connection timeout for cases when connection open but solution isn't send
	err := clientConn.SetReadDeadline(time.Now().UTC().Add(connReadDeadline))
	if err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	data := make([]byte, readDataBuffer)
	_, err = clientConn.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	// reset deadline because client will never write again
	err = clientConn.SetReadDeadline(time.Time{})
	if err != nil {
		return nil, fmt.Errorf("failed to reset read deadline: %w", err)
	}

	return data, nil
}
