package client

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
	"time"
	"word-of-wisdom-v2/internal/general/challenger"
	"word-of-wisdom-v2/internal/general/tcp"
)

type solverMock struct {
	info *challenger.ChallengeInfo

	solution string
	err      error
}

func (s *solverMock) SolveChallenge(info challenger.ChallengeInfo) (string, error) {
	s.info = &info

	return s.solution, s.err
}

type serverMock struct {
	listener net.Listener

	solution []byte

	closeBeforeSendChallenge bool
	closeBeforeSolution      bool
	closeBeforeQuote         bool
}

func (s *serverMock) handleConnect() {
	conn, err := s.listener.Accept()
	if err != nil {
		log.Printf("err while accept: %v", err)
		return
	}

	defer func() { _ = conn.Close() }()

	if s.closeBeforeSendChallenge {
		return
	}

	_, _ = conn.Write([]byte(`{"random_string":"rstring1","number_leading_zeros":5,"number_symbols":3}`))

	if s.closeBeforeSolution {
		return
	}

	data := make([]byte, 1024)
	n, _ := conn.Read(data)

	s.solution = data[:n]

	if s.closeBeforeQuote {
		return
	}

	_, _ = conn.Write([]byte(`{"quote":"quote1"}`))
}

func TestClient_GetQuote(t *testing.T) {
	listener, err := net.Listen(tcp.TCP, ":8080")
	if err != nil {
		log.Printf("err while listen: %v", err)
		return
	}

	defer func() { _ = listener.Close() }()

	tests := []struct {
		name string

		// solver params
		solution  string
		solverErr error
		// server mock params
		closeBeforeSendChallenge bool
		closeBeforeSolution      bool
		closeBeforeQuote         bool

		// to check
		wantQuote     string
		wantErr       bool
		wantChallenge *challenger.ChallengeInfo
		wantSolution  []byte
	}{
		{
			name:                     "all ok",
			solution:                 "solution1",
			solverErr:                nil,
			closeBeforeSendChallenge: false,
			closeBeforeSolution:      false,
			closeBeforeQuote:         false,
			wantQuote:                "quote1",
			wantErr:                  false,
			wantChallenge: &challenger.ChallengeInfo{
				RandomString:       "rstring1",
				NumberLeadingZeros: 5,
				NumberSymbols:      3,
			},
			wantSolution: []byte(`{"solution":"solution1"}`),
		},
		{
			name:                     "quote wasn't sent",
			solution:                 "solution1",
			solverErr:                nil,
			closeBeforeSendChallenge: false,
			closeBeforeSolution:      false,
			closeBeforeQuote:         true,
			wantQuote:                "",
			wantErr:                  true,
			wantChallenge: &challenger.ChallengeInfo{
				RandomString:       "rstring1",
				NumberLeadingZeros: 5,
				NumberSymbols:      3,
			},
			wantSolution: []byte(`{"solution":"solution1"}`),
		},
		{
			name:                     "server didn't wait solution",
			solution:                 "solution1",
			solverErr:                nil,
			closeBeforeSendChallenge: false,
			closeBeforeSolution:      true,
			closeBeforeQuote:         false,
			wantQuote:                "",
			wantErr:                  true,
			wantChallenge: &challenger.ChallengeInfo{
				RandomString:       "rstring1",
				NumberLeadingZeros: 5,
				NumberSymbols:      3,
			},
			wantSolution: nil,
		},
		{
			name:                     "solver didn't solve",
			solution:                 "",
			solverErr:                fmt.Errorf("solver error"),
			closeBeforeSendChallenge: false,
			closeBeforeSolution:      false,
			closeBeforeQuote:         false,
			wantQuote:                "",
			wantErr:                  true,
			wantChallenge: &challenger.ChallengeInfo{
				RandomString:       "rstring1",
				NumberLeadingZeros: 5,
				NumberSymbols:      3,
			},
			wantSolution: []byte{}, // server tried to read, but there was nothing
		},
		{
			name:                     "no challenge from server",
			solution:                 "",
			solverErr:                nil,
			closeBeforeSendChallenge: true,
			closeBeforeSolution:      false,
			closeBeforeQuote:         false,
			wantQuote:                "",
			wantErr:                  true,
			wantChallenge:            nil,
			wantSolution:             nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solver := &solverMock{err: tt.solverErr, solution: tt.solution}
			server := &serverMock{
				listener:                 listener,
				closeBeforeSendChallenge: tt.closeBeforeSendChallenge,
				closeBeforeSolution:      tt.closeBeforeSolution,
				closeBeforeQuote:         tt.closeBeforeQuote,
			}

			go server.handleConnect()

			c := New("localhost:8080", time.Minute, solver)

			quote, err := c.GetQuote()
			assert.Equalf(t, tt.wantErr, err != nil, "err = %v", err)
			assert.Equal(t, tt.wantQuote, quote)
			assert.Equal(t, tt.wantChallenge, solver.info)
			assert.Equal(t, tt.wantSolution, server.solution)
		})
	}
}
