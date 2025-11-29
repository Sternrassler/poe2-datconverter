package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	GamePath           string
	ExtractedFilesPath string
	RawDataOutputPath  string
	ModelsOutputPath   string
	Limit              int
	ListOnly           bool
	Filter             string
}

func Load() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	homeDir, _ := os.UserHomeDir()
	defaultGamePath := filepath.Join(homeDir, ".steam/debian-installation/steamapps/common/Path of Exile 2")
	defaultBase := filepath.Join(homeDir, "PoE2")

	// Override defaults with env vars if present
	if envGame := os.Getenv("POE2_GAME_PATH"); envGame != "" {
		defaultGamePath = envGame
	}
	if envBase := os.Getenv("POE2_BASE_PATH"); envBase != "" {
		defaultBase = envBase
	}

	cfg := &Config{}

	flag.StringVar(&cfg.GamePath, "game", defaultGamePath, "Path to Path of Exile 2 installation")
	flag.StringVar(&cfg.ExtractedFilesPath, "out", filepath.Join(defaultBase, "Extracted"), "Output path for extracted files")
	flag.StringVar(&cfg.ModelsOutputPath, "models", filepath.Join(defaultBase, "GameModels"), "Path for generated models")
	flag.IntVar(&cfg.Limit, "limit", 0, "Limit number of files to extract (0 = no limit)")
	flag.BoolVar(&cfg.ListOnly, "list", false, "List files without extracting")
	flag.StringVar(&cfg.Filter, "filter", "", "Filter files by hash (hex) or name (if available)")

	flag.Parse()

	return cfg
}
