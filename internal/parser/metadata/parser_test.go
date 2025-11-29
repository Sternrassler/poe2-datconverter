package metadata_test

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unicode/utf16"

	"github.com/Sternrassler/poe2-datconverter/internal/parser/metadata"
)

func TestParseItFile(t *testing.T) {
	// Mock content
	content := `version 2
extends "metadata/items/weapons/weapon_base"

base_item_type = "Bow"
{
    range = 120
    attack_speed = 1.5
}`

	// Encode to UTF-16LE
	u16s := utf16.Encode([]rune(content))
	buf := new(bytes.Buffer)
	for _, u := range u16s {
		binary.Write(buf, binary.LittleEndian, u)
	}

	itFile, err := metadata.Read(buf)
	if err != nil {
		t.Fatalf("failed to read it file: %v", err)
	}

	if itFile.Version != "2" {
		t.Errorf("expected version 2, got %s", itFile.Version)
	}
	if itFile.Extends != "metadata/items/weapons/weapon_base" {
		t.Errorf("expected extends 'metadata/items/weapons/weapon_base', got '%s'", itFile.Extends)
	}

	// Verify structure
	if len(itFile.Root.Children) != 2 { // base_item_type and the block
		t.Errorf("expected 2 children in root, got %d", len(itFile.Root.Children))
	}

	// Check first child: base_item_type = "Bow"
	child1 := itFile.Root.Children[0]
	if child1.Key != "base_item_type" || child1.Value != "Bow" {
		t.Errorf("unexpected first child: %s = %s", child1.Key, child1.Value)
	}

	// Check second child: block
	child2 := itFile.Root.Children[1]
	if len(child2.Children) != 2 {
		t.Errorf("expected 2 children in block, got %d", len(child2.Children))
	}

	blockChild1 := child2.Children[0]
	if blockChild1.Key != "range" || blockChild1.Value != "120" {
		t.Errorf("unexpected block child 1: %s = %s", blockChild1.Key, blockChild1.Value)
	}
}
