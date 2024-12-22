package aes

import (
	"fmt"
	"math"
)

const KeyEntropyThreshold = 4.0  // Example entropy threshold, adjust as needed
const FileEntropyThreshold = 4.5 // Example entropy threshold, adjust as needed

// IsHighEntropy checks if the given bytes has high entropy
func IsHighEntropy(key []byte, entropyThreshold float64) bool {
	entropyValue := CalculateEntropy(key)
	fmt.Println("The entropy of given bytes is", entropyValue)
	return entropyValue > entropyThreshold
}

// CalculateEntropy calculates the entropy of the given data
func CalculateEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0.0
	}

	freq := make(map[byte]float64)
	for _, b := range data {
		freq[b]++
	}

	var entropy float64
	for _, count := range freq {
		p := count / float64(len(data))
		entropy -= p * math.Log2(p)
	}

	return entropy
}
