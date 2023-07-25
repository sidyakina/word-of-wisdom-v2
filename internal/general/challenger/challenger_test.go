package challenger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChallenger_ValidateSolution(t *testing.T) {
	tests := []struct {
		name      string
		challenge ChallengeInfo
		solution  string
		want      bool
	}{
		{
			name: "ok, hash = 00000bc... with 20 zeros = 20 wanted",
			challenge: ChallengeInfo{
				RandomString:       "randomstring",
				NumberLeadingZeros: 20,
				NumberSymbols:      7,
			},
			solution: "VuA8fgf",
			want:     true,
		},
		{
			name: "ok, hash = 00000bc... with 20 zeros > 19 wanted",
			challenge: ChallengeInfo{
				RandomString:       "randomstring",
				NumberLeadingZeros: 19,
				NumberSymbols:      7,
			},
			solution: "VuA8fgf",
			want:     true,
		},
		{
			name: "not ok, hash = 00000bc... with 20 zeros, but wanted 21 zeros",
			challenge: ChallengeInfo{
				RandomString:       "randomstring",
				NumberLeadingZeros: 21,
				NumberSymbols:      8,
			},
			solution: "VuA8fgf",
			want:     false,
		},
		{
			name: "not ok, hash = 00000bc... with 20 zeros, but wanted 8 symbols",
			challenge: ChallengeInfo{
				RandomString:       "randomstring",
				NumberLeadingZeros: 20,
				NumberSymbols:      8,
			},
			solution: "VuA8fgf",
			want:     false,
		},
		{
			name: "not ok, hash = 00000bc... with 20 zeros, but wanted 6 symbols",
			challenge: ChallengeInfo{
				RandomString:       "randomstring",
				NumberLeadingZeros: 20,
				NumberSymbols:      6,
			},
			solution: "VuA8fgf",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Challenger{}
			result := c.ValidateSolution(tt.challenge, tt.solution)

			assert.Equal(t, tt.want, result)
		})
	}
}
