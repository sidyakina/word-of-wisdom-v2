package challenger

import (
	"io"
	"log"
	"testing"
	"time"
)

func BenchmarkSolveChallenge_5zeros_5symbols(b *testing.B) {
	benchmarkSolveChallenge(5, 5, b)
}

func BenchmarkSolveChallenge_5zeros_20symbols(b *testing.B) {
	benchmarkSolveChallenge(5, 20, b)
}

func BenchmarkSolveChallenge_10zeros_5symbols(b *testing.B) {
	benchmarkSolveChallenge(10, 5, b)
}

func BenchmarkSolveChallenge_10zeros_20symbols(b *testing.B) {
	benchmarkSolveChallenge(10, 20, b)
}

func BenchmarkSolveChallenge_15zeros_5symbols(b *testing.B) {
	benchmarkSolveChallenge(15, 5, b)
}

func BenchmarkSolveChallenge_15zeros_20symbols(b *testing.B) {
	benchmarkSolveChallenge(15, 20, b)
}

func BenchmarkSolveChallenge_20zeros_5symbols(b *testing.B) {
	benchmarkSolveChallenge(20, 5, b)
}

func BenchmarkSolveChallenge_20zeros_20symbols(b *testing.B) {
	benchmarkSolveChallenge(20, 20, b)
}

func benchmarkSolveChallenge(numberLeadingZeros, numberSymbols int32, b *testing.B) {
	log.SetOutput(io.Discard)

	solver := NewSolver(time.Hour)
	challenge := ChallengeInfo{
		RandomString:       "randomstring",
		NumberLeadingZeros: numberLeadingZeros,
		NumberSymbols:      numberSymbols,
	}

	for i := 0; i < b.N; i++ {
		_, err := solver.SolveChallenge(challenge)
		if err != nil {
			panic(err) // instead of benchmark will be wrong
		}
	}
}
