package utils

import (
	"fmt"
	"time"
)

// TimestampMsg() prepends the original message with
// an RFC3339 compliant timestamp and prints to STDOUT.
// RFC3339 timestamp example: 2023-01-24T09:59:58+08:00
func TimestampMsg(message string) {
	fmt.Println(timestampMsg(message, time.Time{}))
}

// timestampMsg() is a helper function to facilitate unit
// testing by injecting a predefined time.Time argument.
func timestampMsg(message string, timestamp time.Time) string {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	return fmt.Sprintf("[ %s ] %s", timestamp.Format(time.RFC3339), message)
	// fmt.Printf("[ %s ] %s\n", timestamp.Format(time.RFC3339), message)
}
