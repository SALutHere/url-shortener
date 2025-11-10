package random

import (
	"math/rand"
	"strings"
	"time"
)

func NewRandomString(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"1234567890")

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var buf strings.Builder
	for range length {
		buf.WriteRune(chars[rnd.Intn(len(chars))])
	}

	return buf.String()
}
