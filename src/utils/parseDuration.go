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

func ParsePlanDuration(durationStr string) int64 {
	var durationMap = map[string]int64{
		"1M": 29 * 24 * 60 * 60,     // 30 days
		"6M": 6 * 29 * 24 * 60 * 60, // 6 months
		"1Y": 365 * 24 * 60 * 60,    // 1 year
	}

	return durationMap[durationStr]

}
