package challenger

import (
	cryptorand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	mathrand "math/rand"
)

type ChallengeInfo struct {
	RandomString          string
	NumberLeadingZeros    int32
	SolutionNumberSymbols int32
}

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCryptoRandomString(numberSymbols int32) string {
	bytes := make([]byte, numberSymbols)
	_, err := cryptorand.Read(bytes)
	if err != nil {
		log.Printf("failed to cryptorand read: %v, use mathrand string", err)

		return generateMathRandomString(numberSymbols)
	}

	for i := 0; i < int(numberSymbols); i++ {
		bytes[i] = symbols[int(bytes[i])%len(symbols)]
	}

	return string(bytes)
}

// I didn't call mathrand.Seed because it was deprecated since Go 1.20
func generateMathRandomString(lenString int32) string {
	str := make([]rune, lenString)

	for i := 0; i < int(lenString); i++ {
		k := mathrand.Intn(len(symbols))
		str[i] = rune(symbols[k])
	}

	return string(str)
}

func isValid(data string, numberLeadingZeros int32) bool {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
	zeros := calculateLeadingZeros(hash)

	return zeros >= numberLeadingZeros
}

func calculateLeadingZeros(hash string) int32 {
	var zeros int32

	for _, v := range hash {
		switch v {
		case '0':
			zeros += 4 // 0 -> 0000

			continue
		case '1':
			zeros += 3 // 1 -> 0001
		case '2', '3':
			zeros += 2
		case '4', '5', '6', '7':
			zeros += 1
		}

		break
	}

	return zeros
}
