# Card Parser

A Go library for parsing and manipulating character card data embedded in PNG images.

## Features

- Extract character data from PNG images (V2 and V3 format support)
- Multiple scan modes (First, LastVersion, LastLongest)
- Support for character sheets with lorebooks and entries
- Property system with strong typing (String, Integer, Float, Bool, etc.)
- Image format conversion (JPEG, WebP, etc. to PNG)
- URL fetching support
- JSON serialization/deserialization with Sonic

## Installation

```bash
go get github.com/r3dpixel/card-parser
```

## Usage

### Extract Character Data from PNG

```go
import "github.com/r3dpixel/card-parser/png"

// From file
processor := png.FromFile("character.png")
card, err := processor.Get()

// From URL
processor := png.FromURL(client, "https://example.com/character.png")
card, err := processor.Get()

// From bytes
processor := png.FromBytes(imageData)
card, err := processor.Get()
```

### Scan Modes

```go
// Get the first card found
processor.First()

// Get the card with the highest version
processor.LastVersion()

// Get the longest card data
processor.LastLongest()
```

### Work with Character Sheets

```go
import "github.com/r3dpixel/card-parser/character"

// Parse character sheet from JSON
sheet, err := character.FromJSON(reader)

// Export to JSON
err = sheet.ToFile("output.json")

// Access character data
name := sheet.Name
description := sheet.Description
lorebook := sheet.CharacterBook
```

## Project Structure

- `png/` - PNG image parsing and character data extraction
- `character/` - Character sheet, lorebook, and entry structures
- `property/` - Typed property system for character attributes

## Requirements

- Go 1.25.4+
