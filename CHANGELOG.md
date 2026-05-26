# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Removed
- Standalone `VERSION` file. SemVer single source of truth is now `CHANGELOG.md` (top `## [X.Y.Z]` entry) plus git tags (`vX.Y.Z`).

### Changed
- `.agent/rules.md`: versioning rule now references `CHANGELOG.md` + git tags instead of the `VERSION` file.

## [0.1.0] - 2025-11-29

### Added
- Initial project structure and documentation.
- `VERSION` file for tracking project version.
- `CHANGELOG.md` for tracking changes.
- **Rules**: Added Development Methodology section (TDD, Clean Code) and derived project rules to `.agent/rules.md`.

### Changed
- **Docker Build**: Switched to using the public `Sternrassler/ooz` repository for the `poe-export` image build.
- **Docker Image**: Reverted image name from `poe-export-custom` to `poe-export` in `Makefile` and `Dockerfile`.
- **Documentation**: Updated `docs/EXTRACTION.md` to reflect the current Docker image name and extraction process.
- **Parser**: Updated string parsing logic in `internal/parser/dat/reader.go` to use 4-byte null terminator and enforce 2-byte alignment (matching `dat2json`).
