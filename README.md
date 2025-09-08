# File Extractor

A Go package for determining file types and extracting text content from files.

## Usage

```go
package main

import (
    "fmt"
    "log"
    
    file_extractor "file-extractor"
)

func main() {
    // Extract text from a file
    success, text, err := file_extractor.ExtractText("/path/to/document.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    if success {
        fmt.Printf("Extracted text: %s\n", text)
    } else {
        fmt.Println("File is not a text file or couldn't extract text")
    }
}
```

## API

### `ExtractText(filePath string) (bool, string, error)`

Extracts text content from a file if possible.

**Parameters:**
- `filePath`: Path to the file to process

**Returns:**
- `success`: `true` if text was successfully extracted
- `text`: The extracted text content (empty if success is false)
- `error`: Any error that occurred during processing

**Detection Methods:**

The package uses a 4-layer fallback detection system:

1. **File Extension Check**: Fast check against known text file extensions (.txt, .md, .json, .py, etc.)
2. **HTTP Content Detection**: Uses Go's `net/http.DetectContentType()` on first 512 bytes
3. **Content Type Analysis**: Checks if detected MIME type is text-based
4. **Binary Heuristic**: Analyzes content for UTF-8 validity, null bytes, and printable character ratio

## Supported File Types

**Documents:** `.txt`, `.md`, `.rst`  
**Data:** `.json`, `.xml`, `.yaml`, `.yml`, `.csv`, `.tsv`  
**Code:** `.py`, `.go`, `.js`, `.ts`, `.java`, `.c`, `.cpp`, `.h`, `.hpp`, `.php`, `.rb`, `.pl`, `.r`, `.sql`  
**Web:** `.html`, `.htm`, `.css`  
**Scripts:** `.sh`, `.bash`, `.zsh`, `.fish`, `.ps1`  
**Config:** `.conf`, `.cfg`, `.ini`  
**Academic:** `.tex`, `.bib`  
**Logs:** `.log`  
**No Extension:** Files without extensions (like `README`)

The package will reject binary files (images, videos, executables, PDFs, etc.) and only process actual text content.