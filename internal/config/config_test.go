package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	
	// Create test .env file
	envContent := `# Test environment file
APP_PORT=9999
APP_ENC_KEY="test_key_123"
SQLITE_PATH=/tmp/test.db
# Comment line
BACKUP_DIR=/tmp/backups
ADMIN_USER=testuser
ADMIN_PASS='test password with spaces'
`
	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Clear existing env vars
	envVars := []string{"APP_PORT", "APP_ENC_KEY", "SQLITE_PATH", "BACKUP_DIR", "ADMIN_USER", "ADMIN_PASS"}
	for _, key := range envVars {
		os.Unsetenv(key)
	}

	// Load env file
	if err := LoadEnvFiles(); err != nil {
		t.Fatalf("LoadEnvFiles failed: %v", err)
	}

	// Test values
	tests := []struct {
		key      string
		expected string
	}{
		{"APP_PORT", "9999"},
		{"APP_ENC_KEY", "test_key_123"},
		{"SQLITE_PATH", "/tmp/test.db"},
		{"BACKUP_DIR", "/tmp/backups"},
		{"ADMIN_USER", "testuser"},
		{"ADMIN_PASS", "test password with spaces"},
	}

	for _, tt := range tests {
		if got := os.Getenv(tt.key); got != tt.expected {
			t.Errorf("Expected %s=%s, got %s=%s", tt.key, tt.expected, tt.key, got)
		}
	}
}

func TestLoadEnvFilesPriority(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	
	// Create .env file
	envContent := `APP_PORT=8080
SHARED_VAR=from_env
ENV_ONLY=env_value`
	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .env.local file (should override .env)
	envLocalContent := `APP_PORT=3000
SHARED_VAR=from_local
LOCAL_ONLY=local_value`
	envLocalPath := filepath.Join(tempDir, ".env.local")
	if err := os.WriteFile(envLocalPath, []byte(envLocalContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Clear existing env vars
	envVars := []string{"APP_PORT", "SHARED_VAR", "ENV_ONLY", "LOCAL_ONLY"}
	for _, key := range envVars {
		os.Unsetenv(key)
	}

	// Load env files
	if err := LoadEnvFiles(); err != nil {
		t.Fatalf("LoadEnvFiles failed: %v", err)
	}

	// Test priority: .env.local should override .env
	tests := []struct {
		key      string
		expected string
		desc     string
	}{
		{"APP_PORT", "3000", ".env.local should override .env"},
		{"SHARED_VAR", "from_local", ".env.local should override .env for shared vars"},
		{"ENV_ONLY", "env_value", ".env values should be loaded when not in .env.local"},
		{"LOCAL_ONLY", "local_value", ".env.local only values should be loaded"},
	}

	for _, tt := range tests {
		if got := os.Getenv(tt.key); got != tt.expected {
			t.Errorf("%s: Expected %s=%s, got %s=%s", tt.desc, tt.key, tt.expected, tt.key, got)
		}
	}
}

func TestSystemEnvPriority(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	
	// Create .env file
	envContent := `APP_PORT=8080`
	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Set system environment variable (should have highest priority)
	os.Setenv("APP_PORT", "7777")
	defer os.Unsetenv("APP_PORT")

	// Load env files
	if err := LoadEnvFiles(); err != nil {
		t.Fatalf("LoadEnvFiles failed: %v", err)
	}

	// System env should override .env file
	if got := os.Getenv("APP_PORT"); got != "7777" {
		t.Errorf("System env should have priority. Expected APP_PORT=7777, got APP_PORT=%s", got)
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	if got := GetEnv("TEST_VAR", "fallback"); got != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", got)
	}

	if got := GetEnv("NONEXISTENT_VAR", "fallback"); got != "fallback" {
		t.Errorf("Expected 'fallback', got '%s'", got)
	}
}

func TestRequireEnv(t *testing.T) {
	os.Setenv("REQUIRED_VAR", "required_value")
	defer os.Unsetenv("REQUIRED_VAR")

	if got := RequireEnv("REQUIRED_VAR"); got != "required_value" {
		t.Errorf("Expected 'required_value', got '%s'", got)
	}

	// Test panic for missing required var
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for missing required variable")
		}
	}()
	RequireEnv("MISSING_REQUIRED_VAR")
}

func TestConfigLoad(t *testing.T) {
	// Set test environment variables
	testEnvs := map[string]string{
		"APP_PORT":      "9999",
		"APP_ENC_KEY":   "test_key_123",
		"SQLITE_PATH":   "/test/app.db",
		"BACKUP_DIR":    "/test/backups",
		"ADMIN_USER":    "testadmin",
		"ADMIN_PASS":    "testpass",
	}

	for key, value := range testEnvs {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	cfg := Load()

	if cfg.Port != "9999" {
		t.Errorf("Expected Port='9999', got Port='%s'", cfg.Port)
	}
	if cfg.EncKey != "test_key_123" {
		t.Errorf("Expected EncKey='test_key_123', got EncKey='%s'", cfg.EncKey)
	}
	if cfg.SQLitePath != "/test/app.db" {
		t.Errorf("Expected SQLitePath='/test/app.db', got SQLitePath='%s'", cfg.SQLitePath)
	}
	if cfg.BackupDir != "/test/backups" {
		t.Errorf("Expected BackupDir='/test/backups', got BackupDir='%s'", cfg.BackupDir)
	}
	if cfg.AdminUser != "testadmin" {
		t.Errorf("Expected AdminUser='testadmin', got AdminUser='%s'", cfg.AdminUser)
	}
	if cfg.AdminPass != "testpass" {
		t.Errorf("Expected AdminPass='testpass', got AdminPass='%s'", cfg.AdminPass)
	}
}