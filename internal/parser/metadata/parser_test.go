package metadata_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Sternrassler/poe2-datconverter/internal/parser/metadata"
)

func TestParseItFile(t *testing.T) {
	// Path to extracted file
	path := "/home/ix/PoE2/Extracted/metadata/items/mapfragments/maven/mavenmapinsidebottomright5.it"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	itFile, err := metadata.Read(f)
	if err != nil {
		t.Fatalf("failed to read it file: %v", err)
	}

	fmt.Printf("Version: %s\n", itFile.Version)
	fmt.Printf("Extends: %s\n", itFile.Extends)
	fmt.Printf("Content Length: %d\n", len(itFile.Content))

	if itFile.Version != "2" {
		t.Errorf("expected version 2, got %s", itFile.Version)
	}
	if itFile.Extends == "" {
		t.Errorf("expected extends to be set")
	}
}
