package logging

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

var (
	codechars = []string{
		"A", "B", "C", "D", "E", "F", "G", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
)

// Generate an error code for the current timestamp
func GenErrCode() string {
	timestamp := time.Now().UTC().UnixMilli() / 1000
	timePart := fmt.Sprintf("%X", timestamp)

	randomPart := ""
	for i := 0; i < 4; i++ {
		randomPart += codechars[rand.Intn(len(codechars))]
	}

	return timePart + randomPart
}
