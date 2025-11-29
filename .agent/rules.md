# Project Rules

## 1. Technology Stack

- **Language**: Go (Golang)
- **Version**: 1.24.0
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


## 7. Versionierung & Changelog

- **Changelog**: Alle Änderungen müssen im `CHANGELOG.md` dokumentiert werden.
- **Version**: Die Versionsnummer wird ausschließlich in der Datei `VERSION` gepflegt.

## 8. Build & Extraction

- **Tool**: `make` (Targets: `extract-data`, `build-extractor`).
- **Docker**: Image `poe-export` (based on `Sternrassler/ooz`).
- **Input**: `extracted/` directory.

## 9. Error Handling

- **Startup**: Fail fast (`log.Fatalf`) for critical config/path errors.
- **Processing**: Log & Continue (`log.Printf`) for individual file errors.

## 10. Project Structure

- `extracted/`: Input game data (Gitignored).
- `internal/`: Private application logic.
- `cmd/`: Application entry points.


## 11. Development Methodology

- **TDD (Test Driven Development)**: Wir entwickeln Test-Driven. Zuerst der Test, dann der Code.
- **Clean Code**: Wir befolgen Clean Code Prinzipien (Sprechende Namen, kleine Funktionen, Single Responsibility).
