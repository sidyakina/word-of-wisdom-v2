package main

import (
	"log"
	clientpkg "word-of-wisdom-v2/internal/client"
	"word-of-wisdom-v2/internal/general/challenger"
)

func main() {
	config, err := parseConfig()
	if err != nil {
		log.Panicf("failed to parse config: %v", err)
	}

	log.Printf("create client for server %v and request quote", config.ServerURL)

	solver := challenger.NewSolver(config.SolutionTimeout)
	client := clientpkg.New(config.ServerURL, config.ReadTimeout, solver)

	quote, err := client.GetQuote()
	if err != nil {
		log.Printf("failed to request quote: %v", err)
	} else {
		log.Printf("---------------------")
		log.Printf("wisdom of the day: %v", quote)
		log.Printf("---------------------")
	}
}
