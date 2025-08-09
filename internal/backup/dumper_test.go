package backup

import (
	"strings"
	"testing"
	"time"

	"github.com/casparjones/go-dumper/internal/store"
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

}

func TestDisableAndEnableForeignKeyChecks(t *testing.T) {
	d := &Dumper{}
	var buf strings.Builder

	if err := d.disableForeignKeyChecks(&buf); err != nil {
		t.Fatalf("disableForeignKeyChecks failed: %v", err)
	}
	if buf.String() != "SET FOREIGN_KEY_CHECKS=0;\n\n" {
		t.Errorf("unexpected disable output: %q", buf.String())
	}

	buf.Reset()

	if err := d.enableForeignKeyChecks(&buf); err != nil {
		t.Fatalf("enableForeignKeyChecks failed: %v", err)
	}
	if buf.String() != "\nSET FOREIGN_KEY_CHECKS=1;\n" {
		t.Errorf("unexpected enable output: %q", buf.String())
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
