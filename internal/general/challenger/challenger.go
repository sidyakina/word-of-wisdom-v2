package challenger

import (
	"log"
)

type Challenger struct {
	numberLeadingZeros     int32
	challengeNumberSymbols int32
	solutionNumberSymbols  int32
}

func NewChallenger(numberLeadingZeros, challengeNumberSymbols int32, solutionNumberSymbols int32) *Challenger {
	return &Challenger{
		numberLeadingZeros:     numberLeadingZeros,
		challengeNumberSymbols: challengeNumberSymbols,
		solutionNumberSymbols:  solutionNumberSymbols,
	}
}

func (c *Challenger) GetChallenge() ChallengeInfo {
	return ChallengeInfo{
		RandomString:       generateCryptoRandomString(c.challengeNumberSymbols),
		NumberLeadingZeros: c.numberLeadingZeros,
		NumberSymbols:      c.solutionNumberSymbols,
	}
}

func (c *Challenger) ValidateSolution(challenge ChallengeInfo, solution string) bool {
	if len(solution) != int(challenge.NumberSymbols) {
		log.Printf("wrong number symbols in solution: %v, want: %v", len(solution), challenge.NumberSymbols)

		return false
	}

	return isValid(challenge.RandomString+solution, challenge.NumberLeadingZeros)
}
