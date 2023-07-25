# word-of-wisdom-v2

It's a new version of project https://github.com/sidyakina/word-of-wisdom

## Description: 
Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge

## go bench results (./internal/general/challenger/solver_test.go)
goos: linux
goarch: amd64
pkg: word-of-wisdom-v2/internal/general/challenger
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz

| description                               | iterations | nanoseconds per operaion | 
| ------------------------------------------ | --- | ------------------- |  
|  BenchmarkSolveChallenge_5zeros_5symbols  |       	   65101 | 	     18755 ns/op | 
|  BenchmarkSolveChallenge_5zeros_20symbols |      	   41742 | 	     28401 ns/op | 
|  BenchmarkSolveChallenge_10zeros_5symbols |      	    1942 | 	    624274 ns/op | 
|  BenchmarkSolveChallenge_10zeros_20symbols |     	    1305 | 	    952450 ns/op | 
|  BenchmarkSolveChallenge_15zeros_5symbols |      	      66 | 	  16534175 ns/op | 
|  BenchmarkSolveChallenge_15zeros_20symbols |     	      37 | 	  34304761 ns/op | 
|  BenchmarkSolveChallenge_20zeros_5symbols |      	       3 | 	 761173731 ns/op | 
|  BenchmarkSolveChallenge_20zeros_20symbols |     	       2 | 	 896642020 ns/op | 

