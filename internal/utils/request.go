package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateRequestID generates a unique request ID
// Format: timestamp-randomhex
func GenerateRequestID() string {
	// Use current timestamp as part of the ID
	timestamp := time.Now().UnixNano()

	// Generate a random component
	// Note: In Go 1.20+, rand.Intn is automatically seeded
	randomPart := rand.Intn(0xFFFFFF) // 6 hex digits

	// Combine them into a string
	return fmt.Sprintf("%d-%06x", timestamp, randomPart)
}
