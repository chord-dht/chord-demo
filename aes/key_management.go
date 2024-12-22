package aes

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

const keySize = 32 // AES-256, 256 bits = 32 bytes = 64 hex characters

// LoadKey loads the AES key from a file and validates its length and entropy
func LoadKey(keyFileName string) ([]byte, error) {
	keyHex, err := os.ReadFile(keyFileName)
	if err != nil {
		return nil, err
	}
	keyHexStr := strings.TrimSpace(string(keyHex))
	key, err := hex.DecodeString(string(keyHexStr))
	if err != nil {
		return nil, err
	}
	// Validate key length
	if len(key) != keySize {
		return nil, fmt.Errorf("invalid key length: expected %d bytes, got %d bytes", keySize, len(key))
	}
	// Validate key entropy
	keyEntropy := CalculateEntropy(key)
	if keyEntropy < KeyEntropyThreshold {
		fmt.Printf(
			"The entropy of the key is %f < KeyEntropyThreshold %f, "+
				"key has insufficient entropy!\n", keyEntropy, KeyEntropyThreshold,
		)
		return nil, fmt.Errorf("key has insufficient entropy")
	} else {
		fmt.Printf(
			"The entropy of the key is %f > KeyEntropyThreshold %f, "+
				"good enough!\n", keyEntropy, KeyEntropyThreshold,
		)
		return key, nil
	}
}
