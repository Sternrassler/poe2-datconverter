package serializer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sternrassler/poe2-datconverter/internal/parser/csd"
	"github.com/Sternrassler/poe2-datconverter/internal/parser/dat"
	"github.com/Sternrassler/poe2-datconverter/internal/parser/metadata"
)

// ConvertToJSON converts the given parsed file to JSON and writes it to the output path.
func ConvertToJSON(data interface{}, outputPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Encoder
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	// Convert specific types to a map for better JSON structure if needed
	// For now, we just dump the struct.
	// In the future, we might want to transform DatFile into a list of rows/objects.

	var output interface{} = data

	switch v := data.(type) {
	case *dat.DatFile:
		// For DatFile, we might want to export the raw data as base64 or similar if we don't parse rows yet.
		// Since we only have basic parsing (header, data sections), we can't export meaningful rows yet without a schema.
		// So we export metadata about the file.
		output = map[string]interface{}{
			"rowCount":      v.RowCount,
			"rowLength":     v.RowLength,
			"fixedDataSize": len(v.DataFixed),
			"varDataSize":   len(v.DataVariable),
			// "rows": ... // TODO: Implement row parsing when schema is available
		}
	case *metadata.ItFile:
		output = map[string]interface{}{
			"version": v.Version,
			"extends": v.Extends,
			"root":    v.Root,
		}
	case *csd.CsdFile:
		output = map[string]interface{}{
			"includes": v.Includes,
			"content":  v.Content,
		}
	}

	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}
