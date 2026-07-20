package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	// make an empty slice of bytes of a total size of 32 bytes
	token := make([]byte, 32)
	// fill empty byte slice with cryptographically random bytes
	rand.Read(token)
	// This takes the 32 raw bytes and converts them into a 64-character string of hex digits,
	// which is safe to store in a database column and send in JSON.
	return hex.EncodeToString(token)
}
