package quotesrepo

import (
	_ "embed"
	"errors"
	"log"
	"math/rand"
	"strings"
)

//go:embed data.txt
var quotesFile string

type Repo struct {
	quotes []string
}

func New() (*Repo, error) {
	quotes := strings.Split(quotesFile, "\n")

	log.Printf("loaded %v quotes", len(quotes))

	if len(quotes) == 0 {
		return nil, errors.New("empty quotes list")
	}

	return &Repo{quotes: quotes}, nil
}

func (r *Repo) GetQuote() (string, error) {
	n := rand.Intn(len(r.quotes))

	quote := r.quotes[n]

	log.Printf("quotes[%v] = %v", n, quote)

	return quote, nil
}
