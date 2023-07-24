package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	challengerpkg "word-of-wisdom-v2/internal/general/challenger"
	serverpkg "word-of-wisdom-v2/internal/server"
	quotesrepo "word-of-wisdom-v2/internal/server/quotes-repo"
)

func main() {
	config, err := parseConfig()
	if err != nil {
		log.Panicf("failed to parse config: %v", err)
	}

	challenger := challengerpkg.NewChallenger(
		config.NumberZeroBits, config.LenChallengeString, config.LenSolutionString,
	)
	quotes, err := quotesrepo.New()
	if err != nil {
		log.Panicf("failed to create quotes repo: %v", err)
	}

	server, err := serverpkg.New(config.Port, config.ReadSolutionTimeout, challenger, quotes)
	if err != nil {
		log.Panicf("failed to create server %v", err)
	}

	go server.Listen()

	log.Printf("quotes server started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch

	log.Println("got signal, quotes server stopped")
}
