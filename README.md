# word-of-wisdom-v2

It's a new version of project https://github.com/sidyakina/word-of-wisdom

## Description: 
Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge

## Commands

### Linter 
 ```bash
$ make linter
 ```

### Unit tests
 ```bash
$ make tests
 ```


### Benchmark
 ```bash
$ make benchmark
 ```

### Build all
 ```bash
$ make build
 ```

### Start
Port for server can be changed in `build/*.env`
 ```bash
$ make start
 ```
or 
 ```bash
$ make start-server 
$ make start-client
 ```

### Clean
 ```bash
$ make stop-server
 ```

## Benchmark results (./internal/general/challenger/solver_test.go)
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


## Algorithm

Algorithm based on HashCash algorithm (https://en.wikipedia.org/wiki/Hashcash) but have several differences:

• SHA2 was used instead of SHA1, because SHA1 is a cryptographically broken. SHA2 take more time for generate hash and in our case it is benefit, because it adds more time difficulty for solve challenge.
• HashCash uses email info for challenge, but we haven't such info for our client. Random string is used instead.

Algorithm steps:

1. Server generates random string with `challengeNumberSymbols` symbols `a-z, A-z, 0-9`
2. Client in loop finds random string with `solutionNumberSymbols` which must be such:
```
    numberLeadingZeros(SHA2(server_random_string + solution)) == numberLeadingZeros
```
3. Server validate `solution` against this condition and check length of string. Types of symbols aren't checked because it isn't matter which symbols client used (it doesn't affect generation time) 

Current values: `challengeNumberSymbols`=20, `solutionNumberSymbols`=10, `numberLeadingZeros`=20

Value `challengeNumberSymbols` affects number possible server strings. 
If hackers would to prepare hashes for all possible server strings they need prepare 64^20 hashes. 
For 20 leading zeros it will take 64^20 * 0.9s (if made without concurrency). If it isn't enough we can randomize number of symbols in range.

We can also make it more difficult but randomizing `solutionNumberSymbols` in some range instead of predefined number of symbols.
It won't significantly change solution generation time (for example: 5 symbols takes approximately 0.8s, 20 symbols takes 0.9s), 
but will increase number of server string variants in several times.

Also, we can randomize `numberLeadingZeros`, but it will affect generation time, 
so if we request too many zeros, client couldn't make it in time (`READ_SOLUTION_TIMEOUT` can be changed to wait less or more seconds for solution).

## Credentials 
Quotes about space are taken from https://www.thefactsite.com/100-space-facts/