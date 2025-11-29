# Path of Exile 2 Data Converter - Datenfluss-Dokumentation

Diese Dokumentation beschreibt den kompletten Datenfluss des PoE2 Data Converters von den Eingabedateien bis zur JSON-Ausgabe.

## Übersicht: Haupt-Workflow

```mermaid
flowchart TD
    Start([Programmstart]) --> Gen[Code-Generierung<br/>CSharpClassGenerator.Generate]
    Gen --> Meta[Metadata-Dateien konvertieren<br/>.it Dateien]
    Meta --> CSD[CSD-Dateien konvertieren<br/>.csd Dateien]
    CSD --> Data[DATC64-Dateien konvertieren<br/>.datc64 Dateien]
    Data --> Check{UpdateRepo?}
    Check -->|Nein| End([Ende])
    Check -->|Ja| CleanRepo[Repository-Ordner leeren]
    CleanRepo --> RepoMeta[Metadata für Repo konvertieren<br/>ohne RowIndex]
    RepoMeta --> RepoCSD[CSD für Repo konvertieren<br/>ohne RowIndex]
    RepoCSD --> RepoData[DATC64 für Repo konvertieren<br/>ohne RowIndex]
    RepoData --> Lang[Alle Sprachen durchlaufen]
    Lang --> LangLoop[Für jede Sprache:<br/>DATC64-Dateien konvertieren]
    LangLoop --> End
```

## Quelldaten: Verzeichnisstruktur und Dateiarten

### ⚠️ Wichtig: Bundle-Extraktion erforderlich

Die Spieldaten von Path of Exile 2 liegen **nicht direkt lesbar** vor, sondern sind in **komprimierten Bundle-Archiven** (.bundle.bin) gespeichert.

**Detaillierte Anleitung zur Extraktion:** Siehe [docs/EXTRACTION.md](EXTRACTION.md).

**Nach erfolgreicher Extraktion liegt folgende Struktur vor:**

```text
ExtractedFilesPath/
├── metadata/
│   ├── items/**/*.it         # Item Templates
│   ├── *.csd                 # Client String Descriptions
│   └── ...
└── data/
    ├── *.datc64              # Haupt-Tabellen
    ├── English/*.datc64      # Sprachspezifisch
    ├── German/*.datc64
    └── ...
```

### Konfigurierte Pfade (internal/config/config.go)

Das Projekt verwendet konfigurierbare Pfade:

```go
ExtractedFilesPath   = "./extracted"          // Eingabe: Extrahierte Spieldateien
// ... weitere Pfade
```

### Detaillierte Eingabe-Verzeichnisstruktur

```mermaid
graph TB
    Root["ExtractedFilesPath<br/>(./extracted)"]

    Root --> Meta["metadata/<br/>Spiel-Metadaten"]
    Root --> Data["data/<br/>Haupt-Datentabellen"]

    Meta --> IT["*.it Dateien<br/>Item Templates"]
    Meta --> CSD["*.csd Dateien<br/>Client String Descriptions"]
    Meta --> MetaSub["Unterverzeichnisse<br/>(rekursiv durchsucht)"]

    Data --> DATC64["*.datc64 Dateien<br/>Sprach-neutrale Tabellen"]
    Data --> LangDE["German/<br/>Deutsche Übersetzungen"]
    Data --> LangEN["English/<br/>Englische Übersetzungen"]
    Data --> LangFR["French/<br/>Französische Übersetzungen"]
    Data --> LangES["Spanish/<br/>Spanische Übersetzungen"]
    Data --> LangMore["Weitere Sprachen..."]

    LangDE --> DATDE["*.datc64 Dateien<br/>Deutsche Strings"]
    LangEN --> DATEN["*.datc64 Dateien<br/>Englische Strings"]
    LangFR --> DATFR["*.datc64 Dateien<br/>Französische Strings"]

    IT --> ITEx["Beispiele:<br/>- items/weapons/bow.it<br/>- monsters/boss.it<br/>- characters/warrior.it"]
    CSD --> CSDEx["Beispiele:<br/>- stat_descriptions.csd<br/>- skill_stat_descriptions.csd"]
    DATC64 --> DATEx["Beispiele:<br/>- Characters.datc64<br/>- SkillGems.datc64<br/>- BaseItemTypes.datc64<br/>- MonsterVarieties.datc64"]
```

### Konkrete Dateibeispiele nach Typ

#### 1. Metadata-Dateien (.it)

**Speicherort:** `ExtractedFilesPath/metadata/**/*.it` (rekursiv)

**Charakteristik:**

- Textformat mit Key-Value-Paaren
- Vererbungshierarchie via `extends`
- Verschachtelte Strukturen mit `{ }` Blöcken

**Beispiel-Pfade:**

```text
metadata/items/weapons/bows/bow_base.it
metadata/monsters/act1/zombie.it
metadata/characters/warrior.it
metadata/terrains/acts/act1_town.it
```

**Beispiel-Inhalt:**

```text
extends "metadata/items/weapons/weapon_base"

base_item_type = "Bow"
{
    range = 120
    attack_speed = 1.5
}
```

**Verarbeitung:**

- Alle `.it` Dateien werden rekursiv gefunden: `Directory.GetFiles(..., "*.it", SearchOption.AllDirectories)`
- Von `MetadataParser.Parse()` geparst
- Vererbung wird durch `MetadataParser.Merge()` aufgelöst
- Ausgabe: `DataOutputPath/xt/` und `RepoDataOutputPath/xt/`

#### 2. CSD-Dateien (.csd)

**Speicherort:** `ExtractedFilesPath/metadata/**/*.csd` (rekursiv)

**Charakteristik:**

- Unicode-Textformat
- Beschreibungen für Stat-IDs
- Operator-basierte Formatierung
- Parameter mit Werten

**Beispiel-Pfade:**

```text
metadata/stat_descriptions.csd
metadata/skill_stat_descriptions.csd
metadata/passive_skill_stat_descriptions.csd
```

**Beispiel-Inhalt:**

```text
description
2 increased_damage fire_damage
3
+ "Erhöht Feuerschaden um {0}%" canonical_line 0 1
+ "Adds {0} to {1} Fire Damage" 0 2
# "Feuerschaden um {0}% erhöht" 0 1

no_description some_internal_stat
```

**Verarbeitung:**

- Alle `.csd` Dateien werden rekursiv gefunden
- Von `CsdParser.Parse()` geparst (Unicode Encoding)
- IDs, Descriptions, Operators und Parameter werden extrahiert
- Ausgabe: `DataOutputPath/csd/` und `RepoDataOutputPath/csd/`

#### 3. DATC64-Dateien (Haupt-Tabellen)

**Speicherort:** `ExtractedFilesPath/data/*.datc64` (nur Top-Level)

**Charakteristik:**

- Binärformat
- Sprach-neutrale Spielmechanik-Daten
- Feste Zeilen-Struktur + variable Datensektion
- Enthält numerische IDs, Flags, Referenzen

**Beispiel-Dateien:**

```text
Characters.datc64              # Charakter-Klassen (Warrior, Ranger, etc.)
SkillGems.datc64              # Skill-Gems und ihre Eigenschaften
BaseItemTypes.datc64          # Basis-Item-Definitionen
MonsterVarieties.datc64       # Monster-Typen und Stats
WorldAreas.datc64             # Zonen und Gebiete
PassiveSkills.datc64          # Passiv-Skill-Baum
Mods.datc64                   # Item-Modifikatoren
Stats.datc64                  # Stat-Definitionen
AscendancyClasses.datc64      # Ascendancy-Klassen
ChanceableItemClasses.datc64  # Chance-bare Items
```

**Verarbeitung:**

- Nur Top-Level-Dateien: `Directory.GetFiles(dataFilesPath, "*.datc64", SearchOption.TopDirectoryOnly)`
- Von `ReaderFactory.GetReader()` und `DatReader.Read()` geparst
- Memory-Marshal zu C# Structs
- Parallele Serialisierung mit `DatStructSerializer`
- Ausgabe: `DataOutputPath/data/` und `RepoDataOutputPath/data/`

#### 4. DATC64-Dateien (Sprachspezifisch)

**Speicherort:** `ExtractedFilesPath/data/{language}/*.datc64`

**Charakteristik:**

- Gleiche Binärstruktur wie Haupt-Tabellen
- Enthalten übersetzte Strings
- Gleiche Tabellen-Namen wie Hauptdaten
- RowIndex korrespondiert zu Hauptdaten

**Sprach-Ordner:**

```text
data/English/      # Englische Übersetzungen
data/German/       # Deutsche Übersetzungen
data/French/       # Französische Übersetzungen
data/Spanish/      # Spanische Übersetzungen
data/Portuguese/   # Portugiesische Übersetzungen
data/Russian/      # Russische Übersetzungen
data/Thai/         # Thailändische Übersetzungen
data/Japanese/     # Japanische Übersetzungen
data/Korean/       # Koreanische Übersetzungen
data/SimplifiedChinese/   # Vereinfachtes Chinesisch
data/TraditionalChinese/  # Traditionelles Chinesisch
```

**Beispiel-Dateien pro Sprache:**

```text
data/English/Characters.datc64        # Charakter-Namen in Englisch
data/German/Characters.datc64         # Charakter-Namen in Deutsch
data/English/SkillGems.datc64        # Skill-Namen/-Beschreibungen (EN)
data/German/SkillGems.datc64         # Skill-Namen/-Beschreibungen (DE)
```

**Verarbeitung:**

- Sprach-Ordner werden dynamisch erkannt: `Directory.GetDirectories(Path.Combine(Config.ExtractedFilesPath, "data"))`
- Nur beim Repository-Update: Jede Sprache wird separat verarbeitet
- Ausgabe: `RepoDataOutputPath/data/{language}/`

### Datenfluss: Von Quelle zu Ausgabe

```mermaid
flowchart LR
    subgraph "Eingabe-Quellen"
        A1["ExtractedFilesPath/metadata/**/*.it"]
        A2["ExtractedFilesPath/metadata/**/*.csd"]
        A3["ExtractedFilesPath/data/*.datc64"]
        A4["ExtractedFilesPath/data/{lang}/*.datc64"]
    end

    subgraph "Parser"
        A1 --> P1[MetadataParser]
        A2 --> P2[CsdParser]
        A3 --> P3[DatReader]
        A4 --> P4[DatReader pro Sprache]
    end

    subgraph "Serialisierer"
        P1 --> S1[JSON Serializer]
        P2 --> S2[JSON Serializer]
        P3 --> S3[DatStructSerializer]
        P4 --> S4[DatStructSerializer]
    end

    subgraph "Entwicklungs-Ausgabe (mit RowIndex)"
        S1 --> O1["DataOutputPath/xt/*.json"]
        S2 --> O2["DataOutputPath/csd/*.json"]
        S3 --> O3["DataOutputPath/data/*.json"]
    end

    subgraph "Repository-Ausgabe (ohne RowIndex)"
        S1 --> R1["RepoDataOutputPath/xt/*.json"]
        S2 --> R2["RepoDataOutputPath/csd/*.json"]
        S3 --> R3["RepoDataOutputPath/data/*.json"]
        S4 --> R4["RepoDataOutputPath/data/{lang}/*.json"]
    end
```

### Dateimengen (Beispiel)

Basierend auf dem Projekt-Status:

| Kategorie | Anzahl Eingabe-Dateien | Anzahl Ausgabe-Dateien |
|-----------|------------------------|------------------------|
| *.it Dateien | ~2.000+ (geschätzt) | ~2.000+ JSON |
| *.csd Dateien | ~10-20 (geschätzt) | ~10-20 JSON |
| *.datc64 (Haupt) | Unbekannt | Verarbeitet durch 814 Structs |
| *.datc64 (pro Sprache) | Unbekannt × 11 Sprachen | JSON pro Sprache |

### Wichtige Code-Stellen für Datei-Zugriff

| Operation | Datei | Zeile | Code |
|-----------|-------|-------|------|
| Metadata-Dateien finden | `src/Program.cs` | 36 | `Directory.GetFiles(Config.GetExtractedFilePath("metadata"), "*.it", SearchOption.AllDirectories)` |
| CSD-Dateien finden | `src/Program.cs` | 62 | `Directory.GetFiles(Config.GetExtractedFilePath("metadata"), "*.csd", SearchOption.AllDirectories)` |
| DATC64-Dateien finden | `src/Program.cs` | 96 | `Directory.GetFiles(dataFilesPath, "*.datc64", SearchOption.TopDirectoryOnly)` |
| Sprachen ermitteln | `src/Program.cs` | 24-25 | `Directory.GetDirectories(Path.Combine(Config.ExtractedFilesPath, "data"))` |
| Dateien lesen | `src/Program.cs` | 102 | `File.ReadAllBytes(file)` |

## 1. Dateiquellen und Input-Formate (Übersicht)

```mermaid
graph LR
    subgraph "Extrahierte Spieldateien"
        A[Config.ExtractedFilesPath] --> B[metadata/]
        A --> C[data/]

        B --> D[*.it Dateien<br/>Item Templates<br/>Vererbungshierarchie]
        B --> E[*.csd Dateien<br/>Client String Descriptions<br/>Text + IDs]

        C --> F[*.datc64 Dateien<br/>Binäre Datentabellen<br/>Spielinformationen]
        C --> G[Sprach-Ordner<br/>English, German, etc.]

        G --> H[*.datc64 Dateien<br/>Sprachspezifische Daten]
    end
```

## 2. DATC64 Dateiformat-Struktur

```mermaid
graph TB
    subgraph "DATC64 Binärdatei"
        A[Header<br/>4 Bytes] --> B[Zeilenanzahl<br/>Int32]

        C[Zeilen-Daten<br/>Variable Länge] --> D[Zeile 1<br/>Feste Länge pro Zeile]
        C --> E[Zeile 2<br/>Feste Länge pro Zeile]
        C --> F[...]
        C --> G[Zeile N<br/>Feste Länge pro Zeile]

        H[Separator<br/>8 Bytes] --> I[0xBBBBBBBBBBBBBBBB]

        J[Datensektion<br/>Variable Länge] --> K[Strings<br/>Unicode, null-terminiert]
        J --> L[Arrays<br/>Variable Elemente]
        J --> M[Referenz-Daten<br/>Offsets und Indizes]
    end

    A -.-> C
    C -.-> H
    H -.-> J
```

### Spezielle Werte

- **Null-Wert**: `0xFEFEFEFEFEFEFEFE` (8 Bytes)
- **Separator**: `0xBBBBBBBBBBBBBBBB` (8 Bytes)
- **String-Terminator**: `[0, 0, 0, 0]` (4 Bytes null) - Muss 2-Byte aligned relativ zum String-Start sein

## 3. DATC64 Verarbeitungs-Pipeline

```mermaid
flowchart LR
    subgraph "Datei-Laden"
        A[*.datc64 Datei] --> B[File.ReadAllBytes]
        B --> C[byte Array]
    end

    subgraph "ReaderFactory"
        C --> D[GetReader<br/>Name + Data]
        D --> E[TypesFactory.StructsMap<br/>Suche Struct-Typ]
        E -->|Gefunden| F[DatReader.Read<br/>mit Struct-Typ]
        E -->|Nicht gefunden| G[DatReader.Read<br/>ohne Struct-Typ<br/>nur Bytes]
    end

    subgraph "DatReader Parsing"
        F --> H[Header lesen<br/>4 Bytes = Zeilenanzahl]
        H --> I[Separator suchen<br/>0xBBBBBBBBBBBBBBBB]
        I --> J[Zeilen extrahieren<br/>rowLength berechnen]
        J --> K[Datensektion extrahieren<br/>nach Separator]
        K --> L[Unsafe Memory Marshal<br/>byte zu C# Struct]
        L --> M[Liste von Objekten<br/>Rows]
    end

    G --> N[Nur RowBytes<br/>keine Structs]
```

### Memory Marshalling Details

```csharp
// Unsafe Memory Marshal
var pointer = Unsafe.AsPointer(ref MemoryMarshal.GetReference(itemData));
var item = Marshal.PtrToStructure(new IntPtr(pointer), underlyingType);
```

## 4. Metadata (.it) Verarbeitung

```mermaid
flowchart TD
    subgraph "Parsing"
        A[*.it Dateien finden] --> B[MetadataParser.Parse]
        B --> C[Zeile für Zeile lesen]
        C --> D{Zeile?}
        D -->|extends| E[Vererbung speichern]
        D -->|Öffnende Klammer| F[Neuen Block öffnen<br/>Child Node erstellen]
        D -->|Schließende Klammer| G[Block schließen<br/>zu Parent zurück]
        D -->|Key = Value| H[Key-Value-Paar<br/>als Node speichern]
    end

    subgraph "Vererbung auflösen"
        E --> I[Dictionary: Path -> Metadata]
        I --> J[Für jede Metadata:<br/>extends vorhanden?]
        J -->|Ja| K[Parent finden]
        K --> L[MetadataParser.Merge<br/>Parent + Child]
        L --> M[Children zusammenführen]
        J -->|Nein| N[Keine Änderung]
    end

    subgraph "Ausgabe"
        M --> O[JSON serialisieren]
        N --> O
        O --> P[.json Dateien schreiben<br/>data/xt/]
    end
```

## 5. CSD (.csd) Verarbeitung

```mermaid
flowchart TD
    subgraph "Parsing"
        A[*.csd Dateien finden] --> B[CsdParser.Parse]
        B --> C[Unicode Encoding lesen]
        C --> D{Zeile?}
        D -->|no_description| E[Entry ohne Beschreibung]
        D -->|description| F[Entry mit Beschreibungen]
        D -->|include| G[Ignorieren]

        F --> H[IDs lesen<br/>Teil-Anzahl + IDs]
        H --> I[Beschreibungen lesen<br/>Anzahl + Details]
        I --> J[Operator + Text + Parameter]
        J --> K[canonical_line Flag<br/>Parameter Paare]
    end

    subgraph "Ausgabe"
        E --> L[CSD-Objekt]
        K --> L
        L --> M[JSON serialisieren]
        M --> N[.json Dateien schreiben<br/>data/csd/]
    end
```

## 6. Serialisierungs-Prozess (DatStructSerializer)

```mermaid
flowchart TD
    subgraph "Initialisierung"
        A[DatReader + AllResults] --> B[DatStructSerializer erstellen]
        B --> C[includeRowIndex Flag setzen]
    end

    subgraph "Parallel Verarbeitung"
        C --> D[Parallel.ForEach<br/>Alle DatReader]
        D --> E[SerializeStructs<br/>für jeden Reader]
    end

    subgraph "Struct Serialisierung"
        E --> F[JSON Array starten]
        F --> G[Für jede Zeile:<br/>JSON Object starten]
        G -->|includeRowIndex| H[RowIndex schreiben]
        G --> I[Für jedes Feld:<br/>Feldtyp prüfen]
    end

    subgraph "Feldtyp-Behandlung"
        I --> J{Feldtyp?}
        J -->|StringReference| K[String aus Datensektion<br/>Offset auflösen]
        J -->|ArrayReference| L[Array-Werte extrahieren<br/>ElementType verwenden]
        J -->|TableReference| M[Cross-Table-Referenz<br/>siehe nächstes Diagramm]
        J -->|TBool| N[Boolean konvertieren<br/>1 = true, 0 = false]
        J -->|TEnum| O[Enum Index schreiben]
        J -->|byte| P[Hex String<br/>Convert.ToHexString]
        J -->|Standard| Q[Wert direkt schreiben]
    end

    subgraph "Ausgabe"
        K --> R[JSON schreiben]
        L --> R
        M --> R
        N --> R
        O --> R
        P --> R
        Q --> R
        R --> S[.json Dateien schreiben<br/>data/]
    end
```

## 7. Referenz-Auflösung (Cross-Table References)

```mermaid
flowchart TD
    subgraph "TableReference Verarbeitung"
        A[TableReference Wert] --> B{RowIndex == Null<br/>oder < 0?}
        B -->|Ja| C[null schreiben]
        B -->|Nein| D[Attribute lesen<br/>ReferenceTableAttribute]
        D --> E[TableName extrahieren]
    end

    subgraph "Target-Table Lookup"
        E --> F{Table in<br/>AllResults?}
        F -->|Nein| G[Nur TableName<br/>und RowIndex schreiben]
        F -->|Ja| H[Target DatReader holen]
        H --> I{RowIndex<br/>existiert?}
        I -->|Nein| G
        I -->|Ja| J[Zeile aus Target-Table<br/>holen]
    end

    subgraph "ID/Text Auflösung"
        J --> K{TableName ==<br/>'Words'?}
        K -->|Ja| L[Text-Feld suchen]
        K -->|Nein| M[Id-Feld suchen]
        L --> N{Feld<br/>gefunden?}
        M --> N
        N -->|Ja| O[Feldwert rekursiv<br/>serialisieren]
        N -->|Nein| P[Nur RowIndex]
    end

    subgraph "JSON Ausgabe"
        O --> Q{includeRowIndex<br/>oder kein ID?}
        Q -->|Ja| R[Object:<br/>TableName + Id/Text + RowIndex]
        Q -->|Nein| S[Object:<br/>TableName + Id/Text]
        P --> T[Object:<br/>TableName + RowIndex]
    end
```

### Beispiel TableReference

```json
{
  "TableName": "Characters",
  "Id": "Warrior",
  "RowIndex": 0
}
```

Bei `includeRowIndex = false`:

```json
{
  "TableName": "Characters",
  "Id": "Warrior"
}
```

## 8. String-Referenz-Auflösung

```mermaid
flowchart LR
    subgraph "StringReference"
        A[StringReference Wert] --> B[Offset in Datensektion]
        B --> C[reader.Data<br/>byte Array]
        C --> D[Ab Offset lesen<br/>bis String-Terminator]
        D --> E[Unicode Decoding]
        E --> F[String Wert]
    end

    subgraph "String-Terminator"
        D -.-> G[0x00 0x00 0x00 0x00<br/>4 Bytes null]
    end
```

## 9. Array-Referenz-Auflösung

```mermaid
flowchart TD
    subgraph "ArrayReference"
        A[ArrayReference Wert] --> B[Offset + Length<br/>in Datensektion]
        B --> C[ElementType Attribut<br/>lesen]
        C --> D[Für jedes Element:<br/>Bytes extrahieren]
    end

    subgraph "Element-Verarbeitung"
        D --> E{ElementType?}
        E -->|StringReference| F[String auflösen]
        E -->|TableReference| G[Table-Referenz auflösen]
        E -->|Primitiv| H[Direkt konvertieren<br/>int, bool, etc.]
    end

    subgraph "Ausgabe"
        F --> I[JSON Array]
        G --> I
        H --> I
    end
```

## 10. Output-Strukturen

```mermaid
graph TB
    subgraph "Output-Pfade"
        A[Config.DataOutputPath<br/>Entwicklung] --> B[data/<br/>Mit RowIndex]
        A --> C[csd/<br/>CSD JSON]
        A --> D[xt/<br/>Metadata JSON]

        E[Config.RepoDataOutputPath<br/>Repository] --> F[data/<br/>Ohne RowIndex]
        E --> G[data/English/<br/>Sprachspezifisch]
        E --> H[data/German/<br/>Sprachspezifisch]
        E --> I[csd/<br/>CSD JSON]
        E --> J[xt/<br/>Metadata JSON]
    end

    subgraph "Dateinamen"
        B --> K[TableName.json<br/>z.B. Characters.json]
        F --> K
        G --> K
        H --> K
        C --> L[Path.json<br/>z.B. metadata/path.csd.json]
        D --> M[Path.json<br/>z.B. metadata/path.it.json]
    end
```

## 11. Generierte Strukturen (TypesFactory)

```mermaid
flowchart TD
    subgraph "Code-Generierung"
        A[CSharpClassGenerator.Generate] --> B[Struct-Definitionen erstellen]
        B --> C[814 Structs<br/>Generated/Structs/]
        B --> D[144 Enums<br/>Generated/Enums/]
        C --> E[TypesFactory.StructsMap<br/>Dictionary]
    end

    subgraph "Struct-Attribute"
        C --> F[StructLayout<br/>LayoutKind.Explicit<br/>Pack = 1]
        C --> G[FieldOffset<br/>Exakte Byte-Positionen]
        C --> H[ReferenceTable<br/>Fremdschlüssel]
        C --> I[ElementType<br/>Array-Typen]
        C --> J[EnumName<br/>Enum-Zuordnung]
    end

    subgraph "Verwendung"
        E --> K[ReaderFactory<br/>Name -> Type]
        K --> L[Marshal.PtrToStructure<br/>byte zu Struct]
    end
```

## 12. Parallele Verarbeitung

```mermaid
flowchart LR
    subgraph "Sequential Phase"
        A[Alle DATC64-Dateien laden] --> B[ReaderFactory.GetReader<br/>für jede Datei]
        B --> C[ConcurrentDictionary<br/>results]
    end

    subgraph "Parallel Phase"
        C --> D[Parallel.ForEach<br/>results]
        D --> E[Thread 1:<br/>DatStructSerializer]
        D --> F[Thread 2:<br/>DatStructSerializer]
        D --> G[Thread 3:<br/>DatStructSerializer]
        D --> H[Thread N:<br/>DatStructSerializer]
    end

    subgraph "Thread-Safe Operations"
        E --> I[ConcurrentDictionary<br/>Cross-Table-Lookups]
        F --> I
        G --> I
        H --> I
    end

    subgraph "Output"
        E --> J[JSON Dateien schreiben]
        F --> J
        G --> J
        H --> J
    end
```

## Zusammenfassung: Kompletter Datenfluss

```mermaid
flowchart TD
    Start([Programmstart]) --> Gen[Code-Gen:<br/>Structs erstellen]

    Gen --> IT[IT-Dateien:<br/>Metadata parsen<br/>Vererbung auflösen]
    IT --> ITJSON[JSON: xt/]

    Gen --> CSD[CSD-Dateien:<br/>Descriptions parsen]
    CSD --> CSDJSON[JSON: csd/]

    Gen --> DAT[DATC64-Dateien:<br/>Binär-Parsing]
    DAT --> Factory[ReaderFactory:<br/>Struct-Typ finden]
    Factory --> Reader[DatReader:<br/>Header + Rows + Data]
    Reader --> Marshal[Memory Marshal:<br/>byte zu Struct]
    Marshal --> Dict[ConcurrentDictionary:<br/>Alle Reader sammeln]

    Dict --> Para[Parallel.ForEach:<br/>Jede Tabelle]
    Para --> Ser1[Thread 1: Serializer]
    Para --> Ser2[Thread 2: Serializer]
    Para --> SerN[Thread N: Serializer]

    Ser1 --> Ref[Referenzen auflösen:<br/>String, Array, Table]
    Ser2 --> Ref
    SerN --> Ref

    Ref --> Cross[Cross-Table-Lookups:<br/>Id/Text holen]
    Cross --> JSON[JSON serialisieren]
    JSON --> Output[JSON: data/]

    ITJSON --> Repo{UpdateRepo?}
    CSDJSON --> Repo
    Output --> Repo

    Repo -->|Ja| Repeat[Wiederholen ohne<br/>RowIndex für Repo]
    Repeat --> Lang[Für jede Sprache<br/>wiederholen]
    Lang --> End([Ende])

    Repo -->|Nein| End
```

## Performance-Optimierungen

1. **Parallele Serialisierung**: `Parallel.ForEach` für alle Tabellen
2. **ConcurrentDictionary**: Thread-sichere Cross-Table-Lookups
3. **Memory Marshal**: Direkter Speicherzugriff ohne Kopieren
4. **Unsafe Code**: Pointer-basierte Konvertierung für maximale Performance

## Wichtige Code-Stellen

| Komponente | Datei | Zeilen |
|------------|-------|--------|
| Haupt-Workflow | `src/Program.cs` | 9-31 |
| DATC64 Parsing | `src/Parsers/DatReader.cs` | 16-53 |
| Serialisierung | `src/Serialization/DatStructSerializer.cs` | 27-90 |
| Referenz-Auflösung | `src/Serialization/DatStructSerializer.cs` | 92-174 |
| Metadata Parsing | `src/Parsers/MetadataParser.cs` | 30-95 |
| CSD Parsing | `src/Parsers/CsdParser.cs` | 34-121 |
| Type Factory | `src/Parsers/ReaderFactory.cs` | 7-18 |
