package server

import (
	"fmt"
	"log"
	"net"
	"time"
	"word-of-wisdom-v2/internal/general/challenger"
	"word-of-wisdom-v2/internal/general/tcp"
)

type Challenger interface {
	GetChallenge() challenger.ChallengeInfo
	ValidateSolution(challenge challenger.ChallengeInfo, solution string) bool
}

type QuotesRepo interface {
	GetQuote() (string, error)
}

type Server struct {
	listener   net.Listener
	challenger Challenger
	quotesRepo QuotesRepo

	connReadDeadline time.Duration
}

func New(port uint32, connReadDeadline time.Duration, challenger Challenger, quotesRepo QuotesRepo) (*Server, error) {
	listener, err := net.Listen(tcp.TCP, fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &Server{
		listener:         listener,
		challenger:       challenger,
		quotesRepo:       quotesRepo,
		connReadDeadline: connReadDeadline,
	}

	return server, nil
}

func (s *Server) Listen() {
	for {
		clientConn, err := s.listener.Accept()
		if err != nil {
			log.Printf("failed to accept connect: %v", err)
			continue
		}

		log.Printf("new client connection accepted")

		go s.HandleNewConnect(clientConn)
	}
}
