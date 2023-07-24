package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"word-of-wisdom-v2/internal/general/challenger"
	"word-of-wisdom-v2/internal/general/tcp"
	"word-of-wisdom-v2/pkg/api"
)

func (s *Server) HandleNewConnect(clientConn net.Conn) {
	defer func() {
		defer func() {
			log.Printf("closing connection to client")
			tcp.CloseConnection(clientConn)
		}()
	}()

	challenge, err := s.sendChallenge(clientConn)
	if err != nil {
		log.Printf("failed to send challenge: %v", err)

		return
	}

	solution, err := s.waitSolution(clientConn)
	if err != nil {
		log.Printf("failed to receive solution: %v", err)

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

func (s *Server) sendChallenge(clientConn net.Conn) (*challenger.ChallengeInfo, error) {
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
	data, err := tcp.ReadWithDeadline(clientConn, s.connReadDeadline)
	if err != nil {
		return "", err
	}

	message := api.Solution{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal data %s to solution: %w", data, err)
	}

	return message.Solution, nil
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
