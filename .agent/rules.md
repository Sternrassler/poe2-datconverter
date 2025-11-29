# Project Rules

## 1. Technology Stack

- **Language**: Go (Golang)
- **Version**: Latest stable (1.23+)
- **Type**: CLI Application

## 2. Architecture & Patterns

- **Performance**:
  - Use efficient binary reading (e.g., `encoding/binary`, buffered I/O).
  - Use Goroutines for parallel processing of multiple files.
- **Data Flow**:
  - **Input**: Extracted game bundles.
  - **Processing**: Binary -> Go Structs -> JSON.
  - **Output**: JSON files.
- **Project Structure (Standard Go Layout)**:
  - `cmd/poe2-extractor/`: Main entry point.
  - `internal/`: Private application code.
    - `parser/`: Parsing logic (`dat`, `metadata`, `csd`).
    - `model/`: Data models.
    - `serializer/`: JSON conversion.
  - `pkg/`: Library code (if applicable).

## 3. Naming Conventions

- **Go Style**: `camelCase` for internal, `PascalCase` for exported.
- **Files**: `snake_case.go`.

## 4. Documentation

- Maintain `DATENFLUSS.md` updates if logic diverges significantly.

## 5. Kommunikation & Dokumentation

- **Sprache**: Immer auf Deutsch kommunizieren.
- **Dokumentation**: Alle Markdown-Dateien (`.md`) müssen in Deutsch verfasst sein.

## 6. Agent Configuration

- **Permissions**: WebSearch is allowed.
