package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	gcm           cipher.AEAD
	cipherOnce    sync.Once
	cipherInitErr error
)

func ensureCipher() error {
	cipherOnce.Do(func() {
		key, err := getEncryptionKey()
		if err != nil {
			cipherInitErr = err
			return
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			cipherInitErr = fmt.Errorf("failed to create cipher: %w", err)
			return
		}

		g, err := cipher.NewGCM(block)
		if err != nil {
			cipherInitErr = fmt.Errorf("failed to create GCM: %w", err)
			return
		}
		gcm = g
	})
	return cipherInitErr
}

func getEncryptionKey() ([]byte, error) {
	keyBase64 := os.Getenv("APP_ENC_KEY")
	if keyBase64 == "" {
		return nil, fmt.Errorf("APP_ENC_KEY environment variable is required (32-byte base64 key)")
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid APP_ENC_KEY format: %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("APP_ENC_KEY must be 32 bytes when decoded")
	}

	return key, nil
}

func EncryptPassword(password string) (string, error) {
	if err := ensureCipher(); err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptPassword(encrypted string) (string, error) {
	if err := ensureCipher(); err != nil {
		return "", err
	}

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
