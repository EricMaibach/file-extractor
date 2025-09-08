package file_extractor

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// ExtractText extracts text content from a file if possible
// Returns (success, text, error)
// - success: true if text was successfully extracted
// - text: the extracted text content (empty if success is false)
// - error: any error that occurred during processing
func ExtractText(filePath string) (bool, string, error) {
	// Check if file is a supported text type
	isText, _, err := isTextFile(filePath)
	if err != nil {
		return false, "", fmt.Errorf("failed to analyze file type: %v", err)
	}

	if !isText {
		return false, "", nil // Not an error, just not a text file
	}

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	content := string(data)

	// Validate that content is valid UTF-8 text
	if !utf8.ValidString(content) {
		return false, "", nil // Not valid UTF-8, can't extract as text
	}

	return true, content, nil
}

// isTextFile determines if a file is a text file using multiple detection methods
func isTextFile(filePath string) (bool, string, error) {
	// Method 1: Check by file extension first (fast)
	if isTextByExtension(filePath) {
		return true, "text/plain", nil
	}

	// Method 2: Use HTTP content detection with file sample
	file, err := os.Open(filePath)
	if err != nil {
		return false, "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read first 512 bytes for content type detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, "", fmt.Errorf("failed to read file sample: %v", err)
	}

	// Detect content type using HTTP package
	contentType := http.DetectContentType(buffer[:n])
	
	// Method 3: Check if detected type is text-based
	if isTextContentType(contentType) {
		return true, contentType, nil
	}

	// Method 4: Binary heuristic - check if content is mostly printable UTF-8
	if n > 0 && isLikelyTextContent(buffer[:n]) {
		return true, "text/plain", nil
	}

	return false, contentType, nil
}

// isTextByExtension checks if file extension indicates text content
func isTextByExtension(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	textExtensions := map[string]bool{
		".txt":      true,
		".md":       true,
		".markdown": true,
		".rst":      true,
		".csv":      true,
		".tsv":      true,
		".log":      true,
		".conf":     true,
		".cfg":      true,
		".ini":      true,
		".yaml":     true,
		".yml":      true,
		".json":     true,
		".xml":      true,
		".html":     true,
		".htm":      true,
		".css":      true,
		".js":       true,
		".ts":       true,
		".py":       true,
		".go":       true,
		".java":     true,
		".c":        true,
		".cpp":      true,
		".h":        true,
		".hpp":      true,
		".sh":       true,
		".bash":     true,
		".zsh":      true,
		".fish":     true,
		".ps1":      true,
		".sql":      true,
		".r":        true,
		".rb":       true,
		".php":      true,
		".pl":       true,
		".tex":      true,
		".bib":      true,
		"":          true, // files without extension might be text
	}
	
	return textExtensions[ext]
}

// isTextContentType checks if HTTP-detected content type is text-based
func isTextContentType(contentType string) bool {
	// Split off charset if present
	mainType := strings.Split(contentType, ";")[0]
	mainType = strings.TrimSpace(strings.ToLower(mainType))
	
	textTypes := map[string]bool{
		"text/plain":             true,
		"text/html":              true,
		"text/css":               true,
		"text/javascript":        true,
		"text/csv":               true,
		"text/xml":               true,
		"application/json":       true,
		"application/xml":        true,
		"application/javascript": true,
		"application/x-sh":       true,
		"application/x-python":   true,
		"application/x-perl":     true,
		"application/x-ruby":     true,
		"application/x-php":      true,
		"application/sql":        true,
		"application/yaml":       true,
		"application/x-yaml":     true,
	}
	
	// Also check if it starts with "text/"
	if strings.HasPrefix(mainType, "text/") {
		return true
	}
	
	return textTypes[mainType]
}

// isLikelyTextContent uses heuristics to determine if binary data is likely text
func isLikelyTextContent(data []byte) bool {
	if len(data) == 0 {
		return true
	}

	// Check if content is valid UTF-8
	if !utf8.Valid(data) {
		return false
	}

	// Reject files with null bytes (common in binary files)
	for _, b := range data {
		if b == 0 {
			return false
		}
	}

	// Count printable vs non-printable characters
	printableCount := 0
	controlCount := 0
	
	for _, b := range data {
		switch {
		case b >= 32 && b <= 126: // ASCII printable
			printableCount++
		case b == '\t' || b == '\n' || b == '\r': // Common whitespace
			printableCount++
		case b < 32: // Control characters
			controlCount++
		}
	}
	
	// If more than 85% of characters are printable, consider it text
	totalChars := len(data)
	if totalChars == 0 {
		return true
	}
	
	printableRatio := float64(printableCount) / float64(totalChars)
	return printableRatio > 0.85
}