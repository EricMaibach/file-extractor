package file_extractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractText(t *testing.T) {
	tempDir := os.TempDir()

	// Test text file
	textFile := filepath.Join(tempDir, "test.txt")
	textContent := "This is a test text file.\nWith multiple lines."
	err := os.WriteFile(textFile, []byte(textContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test text file: %v", err)
	}
	defer os.Remove(textFile)

	// Test binary file
	binaryFile := filepath.Join(tempDir, "test.bin")
	binaryData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	err = os.WriteFile(binaryFile, binaryData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test binary file: %v", err)
	}
	defer os.Remove(binaryFile)

	tests := []struct {
		name           string
		filePath       string
		expectedSuccess bool
		expectedText   string
	}{
		{
			name:           "text file",
			filePath:       textFile,
			expectedSuccess: true,
			expectedText:   textContent,
		},
		{
			name:           "binary file",
			filePath:       binaryFile,
			expectedSuccess: false,
			expectedText:   "",
		},
		{
			name:           "non-existent file",
			filePath:       "/path/that/does/not/exist.txt",
			expectedSuccess: false,
			expectedText:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, text, err := ExtractText(tt.filePath)

			if tt.name == "non-existent file" {
				if err == nil {
					t.Error("Expected error for non-existent file")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if success != tt.expectedSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectedSuccess, success)
			}

			if success && text != tt.expectedText {
				t.Errorf("Expected text %q, got %q", tt.expectedText, text)
			}
		})
	}
}