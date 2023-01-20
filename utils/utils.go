package utils

import (
	"fmt"
	"time"
)

func FormatMessage(message string) {
	fmt.Printf("[ %s ] %s\n", time.Now().Format(time.RFC3339), message)
}
