package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Sternrassler/poe2-datconverter/internal/config"
	"github.com/Sternrassler/poe2-datconverter/internal/parser/csd"
	"github.com/Sternrassler/poe2-datconverter/internal/parser/dat"
	"github.com/Sternrassler/poe2-datconverter/internal/parser/metadata"
)

func main() {
	cfg := config.Load()

	fmt.Printf("PoE2 Data Converter\n")
	fmt.Printf("Extracted Path: %s\n", cfg.ExtractedFilesPath)

	// Check if extracted path exists
	if _, err := os.Stat(cfg.ExtractedFilesPath); os.IsNotExist(err) {
		log.Fatalf("Error: Extracted files path does not exist: %s\nPlease run the Docker extraction first.", cfg.ExtractedFilesPath)
	}

	fmt.Println("Scanning for extracted files...")

	// Walk through the extracted directory
	err := filepath.Walk(cfg.ExtractedFilesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		switch ext {
		case ".datc64":
			processDatFile(path)
		case ".it":
			processItFile(path)
		case ".csd":
			processCsdFile(path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking extracted directory: %v", err)
	}

	fmt.Println("Processing complete.")
}

func processDatFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open %s: %v", path, err)
		return
	}
	defer f.Close()

	datFile, err := dat.Read(f)
	if err != nil {
		log.Printf("Failed to parse %s: %v", path, err)
		return
	}
	fmt.Printf("Parsed .datc64: %s (Rows: %d)\n", filepath.Base(path), datFile.RowCount)
}

func processItFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open %s: %v", path, err)
		return
	}
	defer f.Close()

	itFile, err := metadata.Read(f)
	if err != nil {
		log.Printf("Failed to parse %s: %v", path, err)
		return
	}
	fmt.Printf("Parsed .it: %s (Version: %s)\n", filepath.Base(path), itFile.Version)
}

func processCsdFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open %s: %v", path, err)
		return
	}
	defer f.Close()

	csdFile, err := csd.Read(f)
	if err != nil {
		log.Printf("Failed to parse %s: %v", path, err)
		return
	}
	fmt.Printf("Parsed .csd: %s (Includes: %d)\n", filepath.Base(path), len(csdFile.Includes))
}
