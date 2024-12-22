package aes

import (
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	key, err := LoadKey("aeskey.txt") // Load the key from the file

	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	plaintext := []byte("Hello, World!")
	ciphertext, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted text does not match plaintext. Got %s, want %s", decrypted, plaintext)
	}
}
