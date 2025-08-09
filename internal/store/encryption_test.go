package store

import (
	"os"
	"testing"
)

func TestEncryptDecryptPassword(t *testing.T) {
	// Set up test encryption key
	testKey := "dGVzdGtleTEyMzQ1Njc4OTBhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ej0="
	os.Setenv("APP_ENC_KEY", testKey)
	defer os.Unsetenv("APP_ENC_KEY")

	// Reinitialize encryption with test key
	initEncryption()

	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!@#$%^&*()",
		},
		{
			name:     "empty password",
			password: "",
		},
		{
			name:     "unicode password",
			password: "пароль123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := EncryptPassword(tt.password)
			if err != nil {
				t.Fatalf("EncryptPassword failed: %v", err)
			}

			if encrypted == "" {
				t.Fatal("Encrypted password is empty")
			}

			// Decrypt
			decrypted, err := DecryptPassword(encrypted)
			if err != nil {
				t.Fatalf("DecryptPassword failed: %v", err)
			}

			if decrypted != tt.password {
				t.Errorf("Password mismatch: expected %q, got %q", tt.password, decrypted)
			}

			// Each encryption should produce different results
			encrypted2, err := EncryptPassword(tt.password)
			if err != nil {
				t.Fatalf("Second EncryptPassword failed: %v", err)
			}

			if encrypted == encrypted2 {
				t.Error("Encryption should produce different results each time")
			}

			// But both should decrypt to the same password
			decrypted2, err := DecryptPassword(encrypted2)
			if err != nil {
				t.Fatalf("Second DecryptPassword failed: %v", err)
			}

			if decrypted2 != tt.password {
				t.Errorf("Second decryption mismatch: expected %q, got %q", tt.password, decrypted2)
			}
		})
	}
}

func TestInvalidDecryption(t *testing.T) {
	// Set up test encryption key
	testKey := "dGVzdGtleTEyMzQ1Njc4OTBhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ej0="
	os.Setenv("APP_ENC_KEY", testKey)
	defer os.Unsetenv("APP_ENC_KEY")

	initEncryption()

	tests := []struct {
		name      string
		encrypted string
		wantError bool
	}{
		{
			name:      "invalid base64",
			encrypted: "invalid-base64!",
			wantError: true,
		},
		{
			name:      "too short",
			encrypted: "dGVzdA==",
			wantError: true,
		},
		{
			name:      "wrong data",
			encrypted: "aW52YWxpZGRhdGFpbnZhbGlkZGF0YWludmFsaWRkYXRh",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptPassword(tt.encrypted)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	key := GenerateEncryptionKey()
	
	if len(key) == 0 {
		t.Fatal("Generated key is empty")
	}

	// Generate another key and ensure they're different
	key2 := GenerateEncryptionKey()
	if key == key2 {
		t.Error("Generated keys should be different")
	}

	// Test that generated key can be used for encryption
	os.Setenv("APP_ENC_KEY", key)
	defer os.Unsetenv("APP_ENC_KEY")

	initEncryption()

	testPassword := "testpass123"
	encrypted, err := EncryptPassword(testPassword)
	if err != nil {
		t.Fatalf("Failed to encrypt with generated key: %v", err)
	}

	decrypted, err := DecryptPassword(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt with generated key: %v", err)
	}

	if decrypted != testPassword {
		t.Errorf("Password mismatch: expected %q, got %q", testPassword, decrypted)
	}
}

// Helper function to reinitialize encryption for testing
func initEncryption() {
	key := getEncryptionKey()
	// Reinitialize the global gcm variable
	// This is a bit hacky but necessary for testing
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
}