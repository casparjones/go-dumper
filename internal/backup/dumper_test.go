package backup

import (
	"strings"
	"testing"
	"time"

	"github.com/user/go-dumper/internal/store"
)

func TestFormatValue(t *testing.T) {
	d := &Dumper{}

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "NULL",
		},
		{
			name:     "string value",
			input:    "hello world",
			expected: "'hello world'",
		},
		{
			name:     "string with quotes",
			input:    "hello 'world'",
			expected: "'hello \\'world\\''",
		},
		{
			name:     "string with backslashes",
			input:    "hello\\world",
			expected: "'hello\\\\world'",
		},
		{
			name:     "string with newlines",
			input:    "hello\nworld",
			expected: "'hello\\nworld'",
		},
		{
			name:     "bytes value",
			input:    []byte("hello"),
			expected: "'hello'",
		},
		{
			name:     "integer value",
			input:    42,
			expected: "42",
		},
		{
			name:     "float value",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "boolean true",
			input:    true,
			expected: "1",
		},
		{
			name:     "boolean false",
			input:    false,
			expected: "0",
		},
		{
			name:     "time value",
			input:    time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
			expected: "'2023-12-25 15:30:45'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.formatValue(tt.input)
			if result != tt.expected {
				t.Errorf("formatValue(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapeString(t *testing.T) {
	d := &Dumper{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string with single quotes",
			input:    "don't",
			expected: "don\\'t",
		},
		{
			name:     "string with double quotes",
			input:    "say \"hello\"",
			expected: "say \\\"hello\\\"",
		},
		{
			name:     "string with backslashes",
			input:    "path\\to\\file",
			expected: "path\\\\to\\\\file",
		},
		{
			name:     "string with newlines",
			input:    "line1\nline2",
			expected: "line1\\nline2",
		},
		{
			name:     "string with carriage returns",
			input:    "line1\rline2",
			expected: "line1\\rline2",
		},
		{
			name:     "string with tabs",
			input:    "col1\tcol2",
			expected: "col1\\tcol2",
		},
		{
			name:     "complex string",
			input:    "It's a \"complex\" string\nwith\\backslashes\tand\ttabs\r\n",
			expected: "It\\'s a \\\"complex\\\" string\\nwith\\\\backslashes\\tand\\ttabs\\r\\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.escapeString(tt.input)
			if result != tt.expected {
				t.Errorf("escapeString(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWriteHeader(t *testing.T) {
	d := &Dumper{}
	var buf strings.Builder
	
	target := &store.Target{
		Host:   "localhost",
		DBName: "testdb",
	}

	err := d.writeHeader(&buf, target)
	if err != nil {
		t.Fatalf("writeHeader failed: %v", err)
	}

	output := buf.String()

	// Check for required elements
	if !strings.Contains(output, "MySQL dump created by go-dumper") {
		t.Error("Header missing application signature")
	}

	if !strings.Contains(output, "Host: localhost") {
		t.Error("Header missing host information")
	}

	if !strings.Contains(output, "Database: testdb") {
		t.Error("Header missing database information")
	}

	if !strings.Contains(output, "SET NAMES utf8mb4") {
		t.Error("Header missing charset setting")
	}

	if !strings.Contains(output, "SET FOREIGN_KEY_CHECKS=0") {
		t.Error("Header missing foreign key checks disable")
	}
}

func TestWriteInsert(t *testing.T) {
	d := &Dumper{}
	var buf strings.Builder

	columns := []string{"`id`", "`name`", "`email`"}
	values := []string{"(1, 'John', 'john@example.com')", "(2, 'Jane', 'jane@example.com')"}

	err := d.writeInsert(&buf, "users", columns, values)
	if err != nil {
		t.Fatalf("writeInsert failed: %v", err)
	}

	output := buf.String()
	expected := "INSERT INTO `users` (`id`, `name`, `email`) VALUES\n(1, 'John', 'john@example.com'),\n(2, 'Jane', 'jane@example.com');\n"

	if output != expected {
		t.Errorf("writeInsert output mismatch:\nGot:\n%q\nExpected:\n%q", output, expected)
	}
}

func TestParseScheduleTime(t *testing.T) {
	tests := []struct {
		name          string
		scheduleTime  string
		expectedHour  int
		expectedMin   int
		expectedError bool
	}{
		{
			name:         "valid time",
			scheduleTime: "14:30",
			expectedHour: 14,
			expectedMin:  30,
		},
		{
			name:         "midnight",
			scheduleTime: "00:00",
			expectedHour: 0,
			expectedMin:  0,
		},
		{
			name:         "end of day",
			scheduleTime: "23:59",
			expectedHour: 23,
			expectedMin:  59,
		},
		{
			name:          "invalid format - no colon",
			scheduleTime:  "1430",
			expectedError: true,
		},
		{
			name:          "invalid format - too many parts",
			scheduleTime:  "14:30:00",
			expectedError: true,
		},
		{
			name:          "invalid hour - too high",
			scheduleTime:  "25:30",
			expectedError: true,
		},
		{
			name:          "invalid hour - negative",
			scheduleTime:  "-1:30",
			expectedError: true,
		},
		{
			name:          "invalid minute - too high",
			scheduleTime:  "14:60",
			expectedError: true,
		},
		{
			name:          "invalid minute - negative",
			scheduleTime:  "14:-1",
			expectedError: true,
		},
		{
			name:          "non-numeric hour",
			scheduleTime:  "ab:30",
			expectedError: true,
		},
		{
			name:          "non-numeric minute",
			scheduleTime:  "14:cd",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, minute, err := parseScheduleTime(tt.scheduleTime)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.scheduleTime)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.scheduleTime, err)
				return
			}

			if hour != tt.expectedHour {
				t.Errorf("Hour mismatch for %q: expected %d, got %d", tt.scheduleTime, tt.expectedHour, hour)
			}

			if minute != tt.expectedMin {
				t.Errorf("Minute mismatch for %q: expected %d, got %d", tt.scheduleTime, tt.expectedMin, minute)
			}
		})
	}
}