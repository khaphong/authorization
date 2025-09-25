package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateUUIDv7 generates a UUID v7 (time-ordered UUID)
func GenerateUUIDv7() string {
	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixMilli()
	
	// Create 16 bytes for the UUID
	var uuidBytes [16]byte
	
	// First 6 bytes: timestamp (48 bits)
	binary.BigEndian.PutUint64(uuidBytes[:8], uint64(timestamp))
	// Shift to use only first 6 bytes
	copy(uuidBytes[:6], uuidBytes[:6])
	
	// Next 2 bytes: version and variant
	// 12 bits of randomness for sub-millisecond ordering
	randBytes := make([]byte, 2)
	rand.Read(randBytes)
	
	// Set version (7) in the most significant 4 bits of the 7th byte
	uuidBytes[6] = (randBytes[0] & 0x0f) | 0x70
	
	// Set variant (10) in the most significant 2 bits of the 9th byte
	uuidBytes[8] = (randBytes[1] & 0x3f) | 0x80
	
	// Fill remaining bytes with random data
	rand.Read(uuidBytes[9:])
	
	// Convert to UUID format
	return formatUUID(uuidBytes[:])
}

// GenerateUUIDv4 generates a standard UUID v4 (fallback)
func GenerateUUIDv4() string {
	return uuid.New().String()
}

// formatUUID formats 16 bytes as a UUID string
func formatUUID(bytes []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}
