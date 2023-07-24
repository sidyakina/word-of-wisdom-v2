package challenger

import (
	"fmt"
	"log"
	"time"
)

type Solver struct {
	timeout time.Duration
}

func NewSolver(timeout time.Duration) *Solver {
	return &Solver{
		timeout: timeout,
	}
}

func (s *Solver) SolveChallenge(challenge ChallengeInfo) (string, error) {
	timer := time.NewTimer(s.timeout)
	defer timer.Stop()

	var (
		numberAttempts int
		solution       string
	)
	for {
		select {
		case <-timer.C:
			return "", fmt.Errorf("timeout: (%v)", s.timeout)
		default:
			solution = generateMathRandomString(challenge.NumberSymbols)
			numberAttempts++
		}

		if isValid(challenge.RandomString+solution, challenge.NumberLeadingZeros) {
			log.Printf("solution %v found after %v attempts", solution, numberAttempts)

			break
		}
	}

	return solution, nil
}
