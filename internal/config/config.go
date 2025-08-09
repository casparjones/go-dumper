package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnvFiles loads environment variables from .env files
// Priority order: .env.local > .env > system environment variables
func LoadEnvFiles() error {
	envFiles := []string{
		".env.local",
		".env",
	}

	for _, file := range envFiles {
		if err := loadEnvFile(file); err != nil {
			// Continue if file doesn't exist, but return error for other issues
			if !os.IsNotExist(err) {
				return fmt.Errorf("error loading %s: %w", file, err)
			}
		}
	}

	return nil
}

// loadEnvFile loads environment variables from a single .env file
func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format in %s at line %d: %s", filename, lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes from value if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
				(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if not already set (allows .env.local to override .env)
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("failed to set environment variable %s: %w", key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", filename, err)
	}

	return nil
}

// GetEnv returns an environment variable with a fallback value
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// RequireEnv returns an environment variable or panics if not set
func RequireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

// Config holds application configuration
type Config struct {
	Port       string
	EncKey     string
	SQLitePath string
	BackupDir  string
	AdminUser  string
	AdminPass  string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:       GetEnv("APP_PORT", "8080"),
		EncKey:     RequireEnv("APP_ENC_KEY"),
		SQLitePath: GetEnv("SQLITE_PATH", "/data/app/app.db"),
		BackupDir:  GetEnv("BACKUP_DIR", "/data/backups"),
		AdminUser:  GetEnv("ADMIN_USER", ""),
		AdminPass:  GetEnv("ADMIN_PASS", ""),
	}
}
