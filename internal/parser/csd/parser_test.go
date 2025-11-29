package csd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Sternrassler/poe2-datconverter/internal/parser/csd"
)

func TestParseCsdFile(t *testing.T) {
	// Path to extracted file
	path := "/home/ix/PoE2/Extracted/metadata/statdescriptions/gem_stat_descriptions.csd"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	csdFile, err := csd.Read(f)
	if err != nil {
		t.Fatalf("failed to read csd file: %v", err)
	}

	fmt.Printf("Includes: %v\n", csdFile.Includes)
	fmt.Printf("Content Length: %d\n", len(csdFile.Content))

	if len(csdFile.Includes) == 0 {
		t.Logf("Warning: no includes found (might be valid for this file)")
	}
}
