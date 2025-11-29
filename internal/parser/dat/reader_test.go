package dat_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Sternrassler/poe2-datconverter/internal/parser/dat"
)

func TestParseActiveSkills(t *testing.T) {
	// Path to extracted file
	path := "/home/ix/PoE2/Extracted/data/activeskills.datc64"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	datFile, err := dat.Read(f)
	if err != nil {
		t.Fatalf("failed to read dat file: %v", err)
	}

	fmt.Printf("RowCount: %d\n", datFile.RowCount)
	fmt.Printf("RowLength: %d\n", datFile.RowLength)
	fmt.Printf("FixedData Size: %d\n", len(datFile.DataFixed))
	fmt.Printf("VariableData Size: %d\n", len(datFile.DataVariable))

	// Basic validation
	if datFile.RowCount == 0 {
		t.Errorf("expected row count > 0")
	}
	if datFile.RowLength == 0 {
		t.Errorf("expected row length > 0")
	}
}
