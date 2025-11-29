# PoE2 Data Converter

Ein umfassender Datenkonverter für Path of Exile 2, der proprietäre Spieldatenformate in JSON umwandelt.

## Übersicht

Dieses Projekt konvertiert die binären Spieldatendateien von Path of Exile 2 in ein menschenlesbares JSON-Format:

- **DATC64-Dateien**: Binäre Datentabellen mit Spielinformationen
- **CSD-Dateien**: Client String Description-Dateien mit Textbeschreibungen und IDs
- **IT-Dateien**: Metadaten-Dateien (Item Templates) mit Vererbungshierarchien

**Status:** Work in Progress - Wird später als NuGet-Paket veröffentlicht.

## Features

- **Binäres Daten-Parsing**: Direktes Memory-Marshalling des DATC64-Binärformats
- **Cross-Reference-Auflösung**: Automatische Auflösung von Referenzen zwischen Tabellen (Strings, Arrays, Tabellenreferenzen)
- **Parallele Verarbeitung**: Optimierte Multi-Thread-Serialisierung für Performance
- **Multi-Language-Support**: Verarbeitet mehrere Sprachordner
- **Code-Generierung**: Generiert C#-Structs für externe Nutzung
- **Zwei Output-Modi**: Mit/ohne RowIndex für verschiedene Anwendungsfälle

## Technologie-Stack

- **Sprache:** C# .NET 8.0
- **Framework:** Console Application (x64)
- **Abhängigkeiten:** Newtonsoft.Json 13.0.3
- **Features:** Unsafe Code für direkte Memory-Manipulation, parallele Verarbeitung

## Projektstruktur

```text
src/
├── Generated/              # Auto-generierte Strukturen (814 Structs + 144 Enums)
│   ├── Enums/             # 144 Enum-Definitionen
│   └── Structs/           # 814 Struct-Definitionen
├── Parsers/               # Binary Parser für DATC64-Format
│   └── DatReader.cs       # Haupt-DATC64-Binary-Parser
├── Serialization/         # JSON-Serialisierung mit Referenzauflösung
│   └── DatStructSerializer.cs  # Haupt-Serialisierer
├── Generators/            # C#-Code-Generator
│   └── CSharpClassGenerator.cs
├── Models/                # Basis-Datenmodelle
│   ├── ArrayReference.cs
│   ├── StringReference.cs
│   └── TableReference.cs
├── Program.cs             # Haupteinstiegspunkt
├── Config.cs              # Pfadkonfiguration
└── Constants.cs           # Projekt-Konstanten
```

## Funktionsweise

### DATC64-Dateiformat

DATC64-Dateien sind Path of Exiles binäres Datentabellenformat:

1. **Header**: 4 Bytes mit Zeilenanzahl
2. **Zeilen-Daten**: Feste Länge pro Zeile mit gepackten Binary Structs
3. **Separator**: `0xBBBBBBBBBBBBBBBB` (8 Bytes)
4. **Datensektion**: Variable Daten (Strings, Arrays)

**Spezielle Werte:**

- Null-Wert: `0xFEFEFEFEFEFEFEFE` (8 Bytes)
- String-Terminator: `[0, 0, 0, 0]` (Unicode)

### Verarbeitungs-Pipeline

```text
DATC64-Datei → DatReader → Unsafe Memory Marshal → C# Struct → DatStructSerializer → JSON
```

1. Datei als Byte-Array einlesen
2. Zeilenanzahl aus Header extrahieren
3. Datensektion nach Separator identifizieren
4. Jede Zeile in C#-Struct konvertieren mit `Marshal.PtrToStructure`
5. Referenzen (Strings, Arrays, Tabellen) aus Datensektion auflösen
6. Mit paralleler Verarbeitung zu JSON serialisieren

### Datenstrukturen

Das Projekt enthält umfangreiche auto-generierte Strukturen:

- **814 Structs**: Repräsentieren verschiedene Spieltabellen (Characters, Skills, Items, Monsters, etc.)
- **144 Enums**: Spielkonstanten und Flags

Beispiel-Strukturen:

- `Characters.cs`: Charakterdaten (Klassen, Stats, Skills, Visuals)
- `SkillGems.cs`: Skill-Gem-Eigenschaften
- `BaseItemTypes.cs`: Item-Basis-Definitionen
- `MonsterVarieties.cs`: Monster-Definitionen
- `WorldAreas.cs`: Zonen-/Gebietsdaten

### Referenzsystem

Der Konverter verarbeitet drei Arten von Referenzen:

1. **StringReference**: Löst String-Offsets zu echtem Text auf
2. **ArrayReference**: Extrahiert Listen von Werten aus Datensektion
3. **TableReference**: Cross-Table-Referenzen mit Id/RowIndex-Auflösung

Attribute werden verwendet, um Referenztypen zu markieren:

```csharp
[ReferenceTable("TableName")]     // Markiert Fremdschlüssel
[ElementType(typeof(Type))]       // Array-Elementtyp
[EnumName("EnumName")]            // Enum-Zuordnung
```

## Einrichtung

1. **Repository klonen**

   ```bash
   git clone https://github.com/Sternrassler/poe2-datconverter.git
   cd poe2-datconverter
   ```

2. **`Config.cs` mit deinen Pfaden bearbeiten:**

   - `ExtractedFilesPath`: Pfad zu extrahierten Spieldateien
   - `DataOutputPath`: Ausgabepfad für konvertierte JSON-Dateien
   - `RepoDataOutputPath`: Ausgabe für Repository
   - `ModelsOutputPath`: Generierte C#-Modelle (optional)

3. **Build und Run**

   ```bash
   dotnet build
   dotnet run
   ```

## Output-Modi

Der Konverter unterstützt zwei Serialisierungsmodi:

- **Mit RowIndex**: Für Entwicklung, enthält alle Referenzdetails
- **Ohne RowIndex**: Für Repository, kompaktere Ausgabe

## Aktueller Status

### Letzte Updates

- ✅ Parallele Serialisierung implementiert (Performance-Optimierung)
- ✅ Charaktersystem-Mapping (Characters, Ascendancy, Chanceable Items)
- ✅ Skill-Crafting-Datenstrukturen
- ✅ 814 Structs und 144 Enums definiert

### Work in Progress

- 🔧 Viele "Unk" (Unknown) Felder werden noch reverse-engineered
- 🔧 Zusätzliche Spieltabellen-Mappings
- 🔧 NuGet-Paket-Vorbereitung

## Statistiken

- **Gesamt-Dateien**: 980 C#-Dateien
- **Generierte Structs**: 814
- **Generierte Enums**: 144
- **Verarbeitung**: Parallel (Multi-Threaded)

## Lizenz

MIT License - Copyright 2024 gaming.tools

_Dieses Produkt ist nicht mit Grinding Gear Games verbunden oder von ihnen unterstützt._

## Beiträge

Dies ist ein Community-Reverse-Engineering-Projekt. Beiträge sind willkommen, insbesondere für:

- Identifizierung unbekannter Struct-Felder
- Hinzufügen neuer Tabellen-Mappings
- Performance-Verbesserungen
- Dokumentation

## Danksagungen

Entwickelt für die Path of Exile Community zur Datenanalyse und Tool-Entwicklung.
