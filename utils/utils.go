package utils

import (
	"fmt"
	"time"
)

// TimestampMsg prepends the original message with
// an RFC3339 compliant timestamp and prints to STDOUT.
// RFC3339 timestamp example: 2023-01-24T09:59:58+08:00
func TimestampMsg(message string) {
	fmt.Printf("[ %s ] %s\n", time.Now().Format(time.RFC3339), message)
}
