package server

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
	"word-of-wisdom-v2/internal/general/challenger"
)

type clientMock struct {
	readErr  error
	solution []byte

	messages           []map[string]interface{}
	currentWriteNumber int32
	writeErrors        map[int32]error
}

func (c *clientMock) Read(b []byte) (n int, err error) {
	if c.readErr != nil {
		return 0, c.readErr
	}

	copy(b, c.solution)

	return len(c.solution), nil
}

func (c *clientMock) Write(b []byte) (n int, err error) {
	c.currentWriteNumber++

	m := map[string]interface{}{}
	_ = json.Unmarshal(b, &m)
	c.messages = append(c.messages, m)

	if c.writeErrors[c.currentWriteNumber] != nil {
		return 0, c.writeErrors[c.currentWriteNumber]
	}

	return len(b), nil
}

func (c *clientMock) Close() error { return nil }

func (c *clientMock) SetReadDeadline(_ time.Time) error { return nil }

// only for implementing interface
func (c *clientMock) LocalAddr() net.Addr { panic("implement me") }

func (c *clientMock) RemoteAddr() net.Addr { panic("implement me") }

func (c *clientMock) SetDeadline(_ time.Time) error { panic("implement me") }

func (c *clientMock) SetWriteDeadline(_ time.Time) error { panic("implement me") }

type quotesRepoMock struct {
	quote string
	err   error
}

func (q *quotesRepoMock) GetQuote() (string, error) {
	return q.quote, q.err
}

type challengerMock struct {
	challenge *challenger.ChallengeInfo
	solution  string
	result    bool
}

func (c *challengerMock) GetChallenge() challenger.ChallengeInfo {
	return challenger.ChallengeInfo{
		RandomString:          "randomstring",
		NumberLeadingZeros:    7,
		SolutionNumberSymbols: 8,
	}
}

func (c *challengerMock) ValidateSolution(challenge challenger.ChallengeInfo, solution string) bool {
	c.challenge = &challenge
	c.solution = solution

	return c.result
}

func TestServer_HandleNewConnect(t *testing.T) {
	tests := []struct {
		name string
		// parameters
		challengerResult      bool
		quotesErr             error
		quote                 string
		clientReadErr         error
		clientSolutionMessage []byte
		clientWriteErrors     map[int32]error

		// result
		wantChallenge *challenger.ChallengeInfo
		wantSolution  string
		wantMessages  []map[string]interface{}
	}{
		{
			name:                  "ok",
			challengerResult:      true,
			quotesErr:             nil,
			quote:                 "wise words",
			clientReadErr:         nil,
			clientSolutionMessage: []byte(`{"solution":"solution1"}`),
			clientWriteErrors:     map[int32]error{1: nil, 2: nil},
			wantChallenge: &challenger.ChallengeInfo{
				RandomString:          "randomstring",
				NumberLeadingZeros:    7,
				SolutionNumberSymbols: 8,
			},
			wantSolution: "solution1",
			wantMessages: []map[string]interface{}{
				{"random_string": "randomstring", "number_leading_zeros": 7., "solution_number_symbols": 8.},
				{"quote": "wise words"},
			},
		},
		{
			name:                  "quotes repo err, no quote",
			challengerResult:      true,
			quotesErr:             fmt.Errorf("quotes repo err"),
			quote:                 "",
			clientReadErr:         nil,
			clientSolutionMessage: []byte(`{"solution":"solution1"}`),
			clientWriteErrors:     map[int32]error{1: nil, 2: nil},
			wantChallenge: &challenger.ChallengeInfo{
				RandomString:          "randomstring",
				NumberLeadingZeros:    7,
				SolutionNumberSymbols: 8,
			},
			wantSolution: "solution1",
			wantMessages: []map[string]interface{}{
				{"random_string": "randomstring", "number_leading_zeros": 7., "solution_number_symbols": 8.},
			},
		},
		{
			name:                  "solution was wrong - no quote",
			challengerResult:      false,
			quotesErr:             nil,
			quote:                 "wise words",
			clientReadErr:         nil,
			clientSolutionMessage: []byte(`{"solution":"solution1"}`),
			clientWriteErrors:     map[int32]error{1: nil, 2: nil},
			wantChallenge: &challenger.ChallengeInfo{
				RandomString:          "randomstring",
				NumberLeadingZeros:    7,
				SolutionNumberSymbols: 8,
			},
			wantSolution: "solution1",
			wantMessages: []map[string]interface{}{
				{"random_string": "randomstring", "number_leading_zeros": 7., "solution_number_symbols": 8.},
			},
		},
		{
			name:                  "solution wasn't received - read from client error",
			challengerResult:      true,
			quotesErr:             nil,
			quote:                 "wise words",
			clientReadErr:         fmt.Errorf("client error"),
			clientSolutionMessage: []byte{},
			clientWriteErrors:     map[int32]error{1: nil, 2: nil},
			wantChallenge:         nil,
			wantSolution:          "",
			wantMessages: []map[string]interface{}{
				{"random_string": "randomstring", "number_leading_zeros": 7., "solution_number_symbols": 8.},
			},
		},
		{
			name:                  "can't write challenge",
			challengerResult:      true,
			quotesErr:             nil,
			quote:                 "wise words",
			clientReadErr:         nil,
			clientSolutionMessage: []byte(`{"solution":"solution1"}`),
			clientWriteErrors:     map[int32]error{1: fmt.Errorf("client error"), 2: nil},
			wantChallenge:         nil,
			wantSolution:          "",
			wantMessages: []map[string]interface{}{
				{"random_string": "randomstring", "number_leading_zeros": 7., "solution_number_symbols": 8.},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &challengerMock{result: tt.challengerResult}
			quotes := &quotesRepoMock{err: tt.quotesErr, quote: tt.quote}

			clientConn := &clientMock{
				readErr:     tt.clientReadErr,
				solution:    tt.clientSolutionMessage,
				messages:    nil,
				writeErrors: tt.clientWriteErrors,
			}

			s := &Server{
				challenger: ch,
				quotesRepo: quotes,
			}
			s.HandleNewConnect(clientConn)

			assert.Equal(t, tt.wantChallenge, ch.challenge)
			assert.Equal(t, tt.wantSolution, ch.solution)
			assert.Equal(t, tt.wantMessages, clientConn.messages)
		})
	}
}
