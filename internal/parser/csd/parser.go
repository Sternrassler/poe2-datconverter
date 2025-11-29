package csd

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"unicode/utf16"
)

// CsdFile represents a parsed .csd file.
type CsdFile struct {
	Includes []string
	Content  string // Full decoded text
}

// Read parses a .csd file from the given reader.
func Read(r io.Reader) (*CsdFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Decode UTF-16LE
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid utf-16le data: odd length")
	}

	u16s := make([]uint16, len(data)/2)
	for i := 0; i < len(u16s); i++ {
		u16s[i] = binary.LittleEndian.Uint16(data[i*2 : i*2+2])
	}

	// Remove BOM if present
	if len(u16s) > 0 && u16s[0] == 0xFEFF {
		u16s = u16s[1:]
	}

	text := string(utf16.Decode(u16s))

	csdFile := &CsdFile{
		Content: text,
	}

	// Simple parsing
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "include") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				csdFile.Includes = append(csdFile.Includes, strings.Trim(parts[1], "\""))
			}
		}
	}

	return csdFile, nil
}
