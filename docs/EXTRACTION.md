# PoE2 Daten-Extraktion

Dieses Dokument behandelt die Extraktion der Spieldaten von Path of Exile 2. Es besteht aus zwei Teilen:

1. **Praktische Anleitung**: Wie man die Daten mittels Docker extrahiert.
2. **Technische Spezifikation**: Hintergründe zum Bundle-Format und den Herausforderungen (Oodle).

---

## Teil 1: Extraktions-Anleitung (Docker)

Diese Anleitung beschreibt, wie man die Spieldaten von Path of Exile 2 mithilfe des Docker-Containers extrahiert.

### Voraussetzungen

1. **Docker**: Muss installiert und lauffähig sein.
2. **Make**: Zum Ausführen der Automatisierungsskripte.
3. **Path of Exile 2 Installation**: Der Pfad zum Spielverzeichnis (Standard: `~/.steam/debian-installation/steamapps/common/Path of Exile 2`).
   - Kann im `Makefile` über die Variable `GAME_PATH` angepasst werden.
4. **Basis-Image**: Das Docker-Image `poe-export` muss lokal verfügbar sein.

### 1. Docker-Image bauen

Das Projekt enthält ein `Makefile`, das den Build-Prozess automatisiert.

```bash
make build-extractor
```

### 2. Extraktion durchführen

Die Extraktion erfolgt vollautomatisch über das `make extract-data` Target. Dies führt folgende Schritte aus:

1. Baut das Docker-Image (falls nötig).
2. Listet alle Dateien im Spielverzeichnis auf (`all_files.txt`).
3. Filtert dynamisch alle relevanten Dateien (`data/*.datc64` und `metadata/`) in `files_to_extract_dynamic.txt`.
4. Extrahiert die gefilterten Dateien.
5. Korrigiert die Dateiberechtigungen (Owner auf aktuellen User setzen).

```bash
make extract-data
```

### 3. Verzeichnisstruktur nach Extraktion

Die Dateien werden in den Ordner `extracted/` im Projektverzeichnis extrahiert:

```text
extracted/
├── data/
│   ├── cursetypes.datc64
│   └── ...
├── metadata/
│   └── ...
└── ...
```

### Fehlerbehebung

- **"image poe-export not found"**: Stellen Sie sicher, dass das Basis-Image `poe-export` existiert (`docker images`).
- **Falscher Spielpfad**: Wenn Ihr Spiel an einem anderen Ort installiert ist, können Sie den Pfad beim Aufruf überschreiben:

  ```bash
  make extract-data GAME_PATH="/pfad/zu/ihrem/spiel"
  ```

---

## Teil 2: Technische Spezifikation (Bundle-Format)

### Übersicht

Path of Exile 2 verwendet ein Bundle-System für Spieldaten. Der Einstiegspunkt ist `_.index.bin`, welche mehrere `.bundle.bin` Dateien referenziert.

### Dateiformate

#### Bundle Index (`_.index.bin`)

- Enthält eine Liste aller Bundles und Dateien.
- Struktur:
  - Header
  - Bundle-Einträge (Records)
  - Datei-Einträge (Records)
  - Pfad-Generierungsdaten

#### Bundle Datei (`.bundle.bin`)

- Enthält komprimierte Datenblöcke.
- **Kompression**: Oodle (Kraken, Leviathan, Mermaid).
- **Header**: 2 Bytes, die den Kompressionstyp anzeigen? Oder komplexer?

### Implementierungs-Herausforderungen

#### 1. Oodle Dekompression

Oodle ist ein proprietärer Kompressionsalgorithmus von RAD Game Tools.

- **Problem**: Es gibt keine Standard-Go-Implementierung für Oodle.
- **Lösungen**:
  1. **CGO + Oodle DLL/SO**: Nutzung von `cgo` zur Einbindung der offiziellen `oo2core` Bibliothek (aus dem Spiel extrahiert).
  2. **Externes Tool**: Aufruf eines existierenden Tools (wie `oo2core` CLI oder ein eigener C++ Wrapper) via `os/exec`.
  3. **Reverse Engineered Port**: Unwahrscheinlich stabil/verfügbar.

> [!IMPORTANT]
> Wir benötigen eine Strategie für die Oodle-Dekompression. Die Nutzung der Shared Library des Spiels (z.B. `liboo2corelinux.so` unter Linux) via CGO ist der technisch sinnvollste Ansatz, erfordert aber, dass der Nutzer den Pfad zur Library bereitstellt.

### Vorgeschlagene Architektur (Internal)

#### `internal/bundles`

- `IndexReader`: Parst `_.index.bin`.
- `BundleReader`: Liest `.bundle.bin` Chunks.
- `Decompressor`: Interface für Dekompression.
  - `OodleDecompressor`: Implementierung mittels CGO/DLL.

#### Workflow

1. Lese `_.index.bin`, um eine Datei-Map zu erstellen (Pfad -> BundleID + Offset + Größe).
2. Für eine angeforderte Datei:
   - Lokalisiere die `.bundle.bin`.
   - Lese den komprimierten Chunk.
   - Dekomprimiere ihn.
   - Extrahiere die spezifischen Dateidaten.
