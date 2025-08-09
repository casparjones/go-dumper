package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

var (
	gcm cipher.AEAD
)

func init() {
	key := getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("failed to create cipher: %v", err))
	}

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Sprintf("failed to create GCM: %v", err))
	}
}

func getEncryptionKey() []byte {
	keyBase64 := os.Getenv("APP_ENC_KEY")
	if keyBase64 == "" {
		panic("APP_ENC_KEY environment variable is required (32-byte base64 key)")
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		panic(fmt.Sprintf("invalid APP_ENC_KEY format: %v", err))
	}

	if len(key) != 32 {
		panic("APP_ENC_KEY must be 32 bytes when decoded")
	}

	return key
}

func EncryptPassword(password string) (string, error) {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptPassword(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted password: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("encrypted data too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	return string(plaintext), nil
}

func GenerateEncryptionKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}