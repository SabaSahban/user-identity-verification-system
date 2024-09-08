package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	// Create a new SHA-256 hash object.
	hash := sha256.New()

	// Write the input string to the hash object.
	hash.Write([]byte(input))

	// Get the hashed bytes.
	hashedBytes := hash.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string.
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
