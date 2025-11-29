package metadata

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"unicode/utf16"
)

// ItFile represents a parsed .it file.
type ItFile struct {
	Version string
	Extends string
	Root    *Node
}

// Node represents a node in the metadata tree.
type Node struct {
	Key      string  `json:"key,omitempty"`
	Value    string  `json:"value,omitempty"`
	Children []*Node `json:"children,omitempty"`
	Parent   *Node   `json:"-"`
}

// Read parses a .it file from the given reader.
func Read(r io.Reader) (*ItFile, error) {
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

	return ParseText(text)
}

// ParseText parses the decoded text content of a .it file.
func ParseText(text string) (*ItFile, error) {
	itFile := &ItFile{
		Root: &Node{},
	}

	scanner := bufio.NewScanner(strings.NewReader(text))
	currentNode := itFile.Root

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "version") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				itFile.Version = strings.Trim(parts[1], "\"")
			}
			continue
		}

		if strings.HasPrefix(line, "extends") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				itFile.Extends = strings.Trim(parts[1], "\"")
			}
			continue
		}

		if line == "{" {
			// Start new block
			newNode := &Node{Parent: currentNode}
			currentNode.Children = append(currentNode.Children, newNode)
			currentNode = newNode
			continue
		}

		if line == "}" {
			// End block
			if currentNode.Parent != nil {
				currentNode = currentNode.Parent
			}
			continue
		}

		// Key = Value
		if idx := strings.Index(line, "="); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])
			value = strings.Trim(value, "\"")

			child := &Node{
				Key:    key,
				Value:  value,
				Parent: currentNode,
			}
			currentNode.Children = append(currentNode.Children, child)
			continue
		}

		// Just a key or directive? Treat as key with empty value
		child := &Node{
			Key:    line,
			Parent: currentNode,
		}
		currentNode.Children = append(currentNode.Children, child)
	}

	return itFile, nil
}
