package utils

import (
	"log"
	"time"
)

// ParseDuration converts a duration string (e.g., "1h", "30m", "15s") to seconds.
// If parsing fails, it falls back to the provided default duration in seconds.
func ParseDuration(durationStr string, defaultSeconds int) int64 {
	if durationStr == "" {
		return int64(defaultSeconds)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Printf("Invalid duration format '%s', using default %d seconds", durationStr, defaultSeconds)
		return int64(defaultSeconds)
	}

	return int64(duration.Seconds())
}
