# PoE2 Bundle Extraktions-Spezifikation

## Übersicht

Path of Exile 2 verwendet ein Bundle-System für Spieldaten. Der Einstiegspunkt ist `_.index.bin`, welche mehrere `.bundle.bin` Dateien referenziert.

## Dateiformate

### Bundle Index (`_.index.bin`)

- Enthält eine Liste aller Bundles und Dateien.
- Struktur:
  - Header
  - Bundle-Einträge (Records)
  - Datei-Einträge (Records)
  - Pfad-Generierungsdaten

### Bundle Datei (`.bundle.bin`)

- Enthält komprimierte Datenblöcke.
- **Kompression**: Oodle (Kraken, Leviathan, Mermaid).
- **Header**: 2 Bytes, die den Kompressionstyp anzeigen? Oder komplexer?

## Implementierungs-Herausforderungen

### 1. Oodle Dekompression

Oodle ist ein proprietärer Kompressionsalgorithmus von RAD Game Tools.

- **Problem**: Es gibt keine Standard-Go-Implementierung für Oodle.
- **Lösungen**:
  1. **CGO + Oodle DLL/SO**: Nutzung von `cgo` zur Einbindung der offiziellen `oo2core` Bibliothek (aus dem Spiel extrahiert).
  2. **Externes Tool**: Aufruf eines existierenden Tools (wie `oo2core` CLI oder ein eigener C++ Wrapper) via `os/exec`.
  3. **Reverse Engineered Port**: Unwahrscheinlich stabil/verfügbar.

> [!IMPORTANT]
> Wir benötigen eine Strategie für die Oodle-Dekompression. Die Nutzung der Shared Library des Spiels (z.B. `liboo2corelinux.so` unter Linux) via CGO ist der technisch sinnvollste Ansatz, erfordert aber, dass der Nutzer den Pfad zur Library bereitstellt.

## Vorgeschlagene Architektur

### `internal/bundles`

- `IndexReader`: Parst `_.index.bin`.
- `BundleReader`: Liest `.bundle.bin` Chunks.
- `Decompressor`: Interface für Dekompression.
  - `OodleDecompressor`: Implementierung mittels CGO/DLL.

### Workflow

1. Lese `_.index.bin`, um eine Datei-Map zu erstellen (Pfad -> BundleID + Offset + Größe).
2. Für eine angeforderte Datei:
   - Lokalisiere die `.bundle.bin`.
   - Lese den komprimierten Chunk.
   - Dekomprimiere ihn.
   - Extrahiere die spezifischen Dateidaten.
