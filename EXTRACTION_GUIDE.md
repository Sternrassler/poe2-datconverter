# PoE2 Daten-Extraktion mit Docker

Diese Anleitung beschreibt, wie man die Spieldaten von Path of Exile 2 mithilfe des Docker-Containers extrahiert.

## Voraussetzungen

1. **Docker**: Muss installiert und lauffähig sein.
2. **Path of Exile 2 Installation**: Der Pfad zum Spielverzeichnis (enthält `Bundles2`).
    - Beispiel Linux (Steam): `~/.steam/debian-installation/steamapps/common/Path of Exile 2`
3. **Basis-Image**: Das Docker-Image `poe-export` muss lokal verfügbar sein.

## 1. Docker-Image bauen

Das Projekt enthält ein `Dockerfile`, das auf `poe-export` basiert und zusätzliche Tools installiert.

```bash
# Im Projektverzeichnis ausführen
docker build -t poe-export-custom .
```

## 2. Extraktion durchführen

Das Tool `bun_extract_file` im Container wird für die Extraktion verwendet. Es benötigt Zugriff auf das Spielverzeichnis und einen Ausgabeordner.

### Syntax

```bash
docker run --rm \
  -v "<PFAD_ZUM_SPIEL>":/game \
  -v "<AUSGABE_PFAD>":/output \
  poe-export-custom \
  bun_extract_file extract-files /game /output [DATEI_PFADE...]
```

### Beispiel: Einzelne Datei extrahieren

```bash
# Erstelle Ausgabeordner
mkdir -p extracted_data

# Extrahiere eine Datei
docker run --rm \
  -v "/home/ix/.steam/debian-installation/steamapps/common/Path of Exile 2":/game \
  -v "$(pwd)/extracted_data":/output \
  poe-export-custom \
  bun_extract_file extract-files /game /output "data/cursetypes.datc64"
```

### Beispiel: Massen-Extraktion (Dateiliste)

Sie können eine Liste von Dateien über die Standardeingabe (stdin) übergeben. Dies ist effizienter für viele Dateien.

```bash
# Verwende die existierende Liste files_to_extract.txt
cat files_to_extract.txt | docker run --rm -i \
  -v "/home/ix/.steam/debian-installation/steamapps/common/Path of Exile 2":/game \
  -v "$(pwd)/extracted_data":/output \
  poe-export-custom \
  bun_extract_file extract-files /game /output
```

## 3. Verzeichnisstruktur nach Extraktion

Die Dateien werden unter Beibehaltung der Verzeichnisstruktur in den Ausgabeordner extrahiert:

```text
extracted_data/
├── data/
│   ├── cursetypes.datc64
│   └── ...
├── metadata/
│   └── ...
└── ...
```

## Fehlerbehebung

- **"image poe-export not found"**: Stellen Sie sicher, dass das Basis-Image `poe-export` existiert (`docker images`).
- **Berechtigungen**: Die extrahierten Dateien gehören möglicherweise dem `root`-User (da Docker als root läuft). Passen Sie die Rechte ggf. an: `sudo chown -R $USER:$USER extracted_data`.
