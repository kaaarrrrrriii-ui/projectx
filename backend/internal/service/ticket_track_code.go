package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

const trackCodeAlphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

func GenerateTrackCode() (string, error) {
	const length = 8

	result := make([]byte, length)
	max := big.NewInt(int64(len(trackCodeAlphabet)))

	for i := range result {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		result[i] = trackCodeAlphabet[n.Int64()]
	}

	return fmt.Sprintf(
		"ОТК-%s-%s",
		result[:4],
		result[4:],
	), nil
}

func HashTrackCode(code string, secret []byte) []byte {
	code = strings.ToUpper(strings.TrimSpace(code))

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(code))

	return mac.Sum(nil)
}