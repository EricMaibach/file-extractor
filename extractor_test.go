package file_extractor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractText_TextFiles(t *testing.T) {
	tests := []struct {
		name            string
		fileName        string
		expectedSuccess bool
		shouldContain   []string
	}{
		{
			name:            "plain text file",
			fileName:        "testtext.txt",
			expectedSuccess: true,
			shouldContain:   []string{"This is a test text file", "With multiple lines"},
		},
		{
			name:            "markdown file",
			fileName:        "testmarkdown.md",
			expectedSuccess: true,
			shouldContain:   []string{"Header", "Subheader", "List", "Values"},
		},
		{
			name:            "go source file",
			fileName:        "main.go",
			expectedSuccess: true,
			shouldContain:   []string{"package main", "func helloworld()", "Hello World!!"},
		},
		{
			name:            "json file",
			fileName:        "test.json",
			expectedSuccess: true,
			shouldContain:   []string{"name", "version", "1.0.0", "nested"},
		},
		{
			name:            "yaml file",
			fileName:        "test.yaml",
			expectedSuccess: true,
			shouldContain:   []string{"name: test", "version: 1.0.0", "items:", "key: value"},
		},
		{
			name:            "csv file",
			fileName:        "test.csv",
			expectedSuccess: true,
			shouldContain:   []string{"Name,Age,City", "John,30,New York", "Jane,25,Los Angeles"},
		},
		{
			name:            "log file",
			fileName:        "test.log",
			expectedSuccess: true,
			shouldContain:   []string{"INFO Starting application", "ERROR Failed to load configuration"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join("testdata", tt.fileName)
			success, text, err := ExtractText(filePath)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if success != tt.expectedSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectedSuccess, success)
			}

			if success {
				for _, expected := range tt.shouldContain {
					if !strings.Contains(text, expected) {
						t.Errorf("Expected text to contain %q, but it didn't", expected)
					}
				}
			}
		})
	}
}

func TestExtractText_BinaryFiles(t *testing.T) {
	tests := []struct {
		name            string
		fileName        string
		expectedSuccess bool
	}{
		{
			name:            "binary file with PNG header",
			fileName:        "test.bin",
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join("testdata", tt.fileName)
			success, text, err := ExtractText(filePath)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if success != tt.expectedSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectedSuccess, success)
			}

			if !success && text != "" {
				t.Errorf("Expected empty text for binary file, got %q", text)
			}
		})
	}
}

func TestExtractText_PDFFiles(t *testing.T) {
	// Test PDF extraction capability
	filePath := filepath.Join("testdata", "test.pdf")
	success, text, err := ExtractText(filePath)

	if err != nil {
		t.Errorf("Unexpected error during PDF extraction: %v", err)
	}

	// The go.pdf might be extractable or not depending on its format
	// We'll just verify the function handles it without crashing
	if success {
		if len(strings.TrimSpace(text)) == 0 {
			t.Error("ExtractText returned success but no text was extracted")
		} else {
			t.Logf("Successfully extracted %d characters from PDF", len(text))
			// Show a sample of extracted text
			sample := text
			if len(sample) > 200 {
				sample = sample[:200] + "..."
			}
			t.Logf("Sample of extracted text: %q", sample)
		}
	} else {
		t.Log("PDF was not text-extractable (may be image-based or encrypted)")
	}
}

func TestExtractText_ErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		shouldError   bool
		errorContains string
	}{
		{
			name:          "non-existent file",
			filePath:      "testdata/does_not_exist.txt",
			shouldError:   true,
			errorContains: "failed to read file",
		},
		{
			name:          "invalid path",
			filePath:      "/invalid/path/that/does/not/exist.txt",
			shouldError:   true,
			errorContains: "failed to read file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, text, err := ExtractText(tt.filePath)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got %v", tt.errorContains, err)
				}
				if success {
					t.Error("Expected success to be false for error case")
				}
				if text != "" {
					t.Errorf("Expected empty text for error case, got %q", text)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestIsTextByExtension(t *testing.T) {
	tests := []struct {
		filePath string
		expected bool
	}{
		{"test.txt", true},
		{"test.md", true},
		{"test.go", true},
		{"test.json", true},
		{"test.yaml", true},
		{"test.yml", true},
		{"test.csv", true},
		{"test.log", true},
		{"test.py", true},
		{"test.js", true},
		{"test.html", true},
		{"test.xml", true},
		{"test.sh", true},
		{"test.sql", true},
		{"test.rb", true},
		{"test.php", true},
		{"test.bin", false},
		{"test.exe", false},
		{"test.jpg", false},
		{"test.png", false},
		{"test.pdf", false},
		{"test.docx", false},
		{"test.zip", false},
		{"README", true}, // files without extension
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			result := isTextByExtension(tt.filePath)
			if result != tt.expected {
				t.Errorf("For %s: expected %v, got %v", tt.filePath, tt.expected, result)
			}
		})
	}
}

func TestIsTextContentType(t *testing.T) {
	tests := []struct {
		contentType string
		expected    bool
	}{
		{"text/plain", true},
		{"text/plain; charset=utf-8", true},
		{"text/html", true},
		{"text/css", true},
		{"text/javascript", true},
		{"application/json", true},
		{"application/xml", true},
		{"application/javascript", true},
		{"application/x-sh", true},
		{"text/anything", true}, // any text/* should be true
		{"image/png", false},
		{"image/jpeg", false},
		{"application/pdf", false},
		{"application/octet-stream", false},
		{"video/mp4", false},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			result := isTextContentType(tt.contentType)
			if result != tt.expected {
				t.Errorf("For %s: expected %v, got %v", tt.contentType, tt.expected, result)
			}
		})
	}
}

func TestIsLikelyTextContent(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "valid ASCII text",
			data:     []byte("This is valid ASCII text with numbers 123 and symbols !@#"),
			expected: true,
		},
		{
			name:     "text with newlines and tabs",
			data:     []byte("Line 1\nLine 2\tTabbed\rCarriage return"),
			expected: true,
		},
		{
			name:     "binary with null bytes",
			data:     []byte{0x00, 0x01, 0x02, 0x03},
			expected: false,
		},
		{
			name:     "PNG header",
			data:     []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expected: false,
		},
		{
			name:     "mostly control characters",
			data:     []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			expected: false,
		},
		{
			name:     "empty data",
			data:     []byte{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLikelyTextContent(tt.data)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestExtractText_UTF8Validation(t *testing.T) {
	tempDir := os.TempDir()

	// Create a file with invalid UTF-8
	invalidUTF8File := filepath.Join(tempDir, "invalid_utf8.txt")
	invalidData := []byte{0xFF, 0xFE, 0xFD} // Invalid UTF-8 sequence
	err := os.WriteFile(invalidUTF8File, invalidData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(invalidUTF8File)

	success, text, err := ExtractText(invalidUTF8File)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if success {
		t.Error("Expected extraction to fail for invalid UTF-8")
	}

	if text != "" {
		t.Errorf("Expected empty text for invalid UTF-8, got %q", text)
	}
}
