package png

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/r3dpixel/card-parser/character"
	"github.com/r3dpixel/card-parser/property"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestCard(t *testing.T, revision character.Revision, name string) *character.Sheet {
	t.Helper()
	sheet := &character.Sheet{
		Revision: character.RevisionV2,
		Spec:     character.SpecV2,
		Version:  character.V2,
		Content:  character.Content{Name: property.String(name)},
	}
	if revision == character.RevisionV3 {
		sheet.Revision = character.RevisionV3
		sheet.Spec = character.SpecV3
		sheet.Version = character.V3
	}
	return sheet
}

func TestCard_IntermediateTransformations(t *testing.T) {
	pngBytes := createTestPNG(t, 4, 4)

	t.Run("RawCard to RawJsonCard", func(t *testing.T) {
		// Setup RawCard with base64 data
		rawCard, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		sheet := createTestCard(t, character.RevisionV2, "Test")
		jsonBytes, err := sheet.ToBytes()
		require.NoError(t, err)
		rawCard.RawCharaData = []byte(base64.StdEncoding.EncodeToString(jsonBytes))
		rawCard.Revision = character.RevisionV2

		// Test ToRawJson
		rawJsonCard, err := rawCard.ToRawJson()
		require.NoError(t, err)
		assert.Equal(t, jsonBytes, rawJsonCard.RawJsonData)
		assert.Equal(t, character.RevisionV2, rawJsonCard.Revision)
	})

	t.Run("RawCard to RawJsonCard with no data", func(t *testing.T) {
		rawCard, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawCard.RawCharaData = nil

		rawJsonCard, err := rawCard.ToRawJson()
		require.NoError(t, err)
		assert.Empty(t, rawJsonCard.RawJsonData)
	})

	t.Run("RawCard to RawJsonCard with invalid base64", func(t *testing.T) {
		rawCard, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawCard.RawCharaData = []byte("not valid base64!!!")

		_, err = rawCard.ToRawJson()
		assert.Error(t, err)
	})

	t.Run("RawJsonCard to CharacterCard", func(t *testing.T) {
		// Setup RawJsonCard
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		sheet := createTestCard(t, character.RevisionV2, "Test")
		jsonBytes, err := sheet.ToBytes()
		require.NoError(t, err)
		rawJsonCard := &RawJsonCard{
			pngData:     raw.pngData,
			RawJsonData: jsonBytes,
			Revision:    character.RevisionV2,
		}

		// Test ToCharacter
		charCard, err := rawJsonCard.ToCharacter()
		require.NoError(t, err)
		assert.Equal(t, "Test", string(charCard.Sheet.Content.Name))
		assert.Equal(t, character.RevisionV2, charCard.Sheet.Revision)
	})

	t.Run("RawJsonCard to CharacterCard with no data", func(t *testing.T) {
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawJsonCard := &RawJsonCard{
			pngData:     raw.pngData,
			RawJsonData: nil,
			Revision:    character.RevisionV2,
		}

		charCard, err := rawJsonCard.ToCharacter()
		require.NoError(t, err)
		assert.NotNil(t, charCard.Sheet)
	})

	t.Run("RawJsonCard to CharacterCard with invalid json", func(t *testing.T) {
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawJsonCard := &RawJsonCard{
			pngData:     raw.pngData,
			RawJsonData: []byte("{not valid json}"),
			Revision:    character.RevisionV2,
		}

		_, err = rawJsonCard.ToCharacter()
		assert.Error(t, err)
	})

	t.Run("CharacterCard to RawJsonCard", func(t *testing.T) {
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		sheet := createTestCard(t, character.RevisionV2, "Test")
		charCard := &CharacterCard{
			pngData: raw.pngData,
			Sheet:   sheet,
		}

		// Test ToRawJson
		rawJsonCard, err := charCard.ToRawJson()
		require.NoError(t, err)
		assert.NotEmpty(t, rawJsonCard.RawJsonData)
		assert.Equal(t, character.RevisionV2, rawJsonCard.Revision)
	})

	t.Run("CharacterCard to RawJsonCard with nil sheet", func(t *testing.T) {
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		charCard := &CharacterCard{
			pngData: raw.pngData,
			Sheet:   nil,
		}

		rawJsonCard, err := charCard.ToRawJson()
		require.NoError(t, err)
		assert.Empty(t, rawJsonCard.RawJsonData)
	})

	t.Run("RawJsonCard to RawCard", func(t *testing.T) {
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		sheet := createTestCard(t, character.RevisionV2, "Test")
		jsonBytes, err := sheet.ToBytes()
		require.NoError(t, err)
		rawJsonCard := &RawJsonCard{
			pngData:     raw.pngData,
			RawJsonData: jsonBytes,
			Revision:    character.RevisionV2,
		}

		// Test ToRaw
		rawCard := rawJsonCard.ToRaw()
		assert.NotEmpty(t, rawCard.RawCharaData)
		assert.Equal(t, character.RevisionV2, rawCard.Revision)
		// Verify it's valid base64
		decoded, err := base64.StdEncoding.DecodeString(string(rawCard.RawCharaData))
		require.NoError(t, err)
		assert.Equal(t, jsonBytes, decoded)
	})

	t.Run("RawJsonCard to RawCard with no data", func(t *testing.T) {
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawJsonCard := &RawJsonCard{
			pngData:     raw.pngData,
			RawJsonData: nil,
			Revision:    character.RevisionV2,
		}

		rawCard := rawJsonCard.ToRaw()
		assert.Empty(t, rawCard.RawCharaData)
	})

	t.Run("full pipeline round trip", func(t *testing.T) {
		// Start with RawCard
		rawCard, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		sheet := createTestCard(t, character.RevisionV3, "Pipeline Test")
		jsonBytes, err := sheet.ToBytes()
		require.NoError(t, err)
		rawCard.RawCharaData = []byte(base64.StdEncoding.EncodeToString(jsonBytes))
		rawCard.Revision = character.RevisionV3

		// RawCard → RawJsonCard
		rawJsonCard, err := rawCard.ToRawJson()
		require.NoError(t, err)

		// RawJsonCard → CharacterCard
		charCard, err := rawJsonCard.ToCharacter()
		require.NoError(t, err)
		assert.Equal(t, "Pipeline Test", string(charCard.Sheet.Content.Name))

		// CharacterCard → RawJsonCard
		rawJsonCard2, err := charCard.ToRawJson()
		require.NoError(t, err)

		// RawJsonCard → RawCard
		rawCard2 := rawJsonCard2.ToRaw()
		assert.Equal(t, rawCard.RawCharaData, rawCard2.RawCharaData)
		assert.Equal(t, character.RevisionV3, rawCard2.Revision)
	})
}

func TestCard_EncodeDecode(t *testing.T) {
	pngBytes := createTestPNG(t, 4, 4)

	t.Run("successful round trip", func(t *testing.T) {
		initialRaw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)

		originalCardData := createTestCard(t, character.RevisionV2, "Test Sheet")
		initialCard := &CharacterCard{pngData: initialRaw.pngData, Sheet: originalCardData}

		encodedRaw, err := initialCard.Encode()
		require.NoError(t, err)
		assert.Equal(t, character.RevisionV2, encodedRaw.Revision)
		assert.NotEmpty(t, encodedRaw.RawCharaData)

		decodedCard, err := encodedRaw.Decode()
		require.NoError(t, err)
		assert.Equal(t, originalCardData.Content.Name, decodedCard.Sheet.Content.Name)
		assert.Equal(t, character.V2, decodedCard.Sheet.Version)
	})

	t.Run("decode with invalid base64 data", func(t *testing.T) {
		rawCard, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawCard.RawCharaData = []byte("this is not base64")
		_, err = rawCard.Decode()
		assert.Error(t, err)
	})

	t.Run("decode with invalid json data", func(t *testing.T) {
		rawCard, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawCard.RawCharaData = []byte(base64.StdEncoding.EncodeToString([]byte("{not json}")))
		_, err = rawCard.Decode()
		assert.Error(t, err)
	})

	t.Run("decode with no chara data", func(t *testing.T) {
		rawCard, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		rawCard.RawCharaData = nil
		decodedCard, err := rawCard.Decode()
		require.NoError(t, err)
		assert.NotNil(t, decodedCard.Sheet)
	})

	t.Run("encode with nil card data", func(t *testing.T) {
		raw, err := FromBytes(pngBytes).Get()
		require.NoError(t, err)
		characterCard := &CharacterCard{pngData: raw.pngData, Sheet: nil}
		encoded, err := characterCard.Encode()
		require.NoError(t, err)
		assert.Empty(t, encoded.RawCharaData)
	})
}

func TestRawCard_ToPngBytes_And_ToFile(t *testing.T) {
	pngBytes := createTestPNG(t, 4, 4)
	rawCard, err := FromBytes(pngBytes).Get()
	require.NoError(t, err)
	cardModel := createTestCard(t, character.RevisionV3, "V3 Sheet")
	cardJson, err := cardModel.ToBytes()
	require.NoError(t, err)
	rawCard.RawCharaData = make([]byte, base64.StdEncoding.EncodedLen(len(cardJson)))
	base64.StdEncoding.Encode(rawCard.RawCharaData, cardJson)
	rawCard.Revision = character.RevisionV3

	t.Run("ToImage creates valid png with chara chunkDetails", func(t *testing.T) {
		finalBytes, err := rawCard.ToBytes()
		require.NoError(t, err)
		assert.Greater(t, len(finalBytes), len(pngBytes))

		reparsedCard, err := FromBytes(finalBytes).Get()
		require.NoError(t, err)
		assert.Equal(t, rawCard.RawCharaData, reparsedCard.RawCharaData)
		assert.Equal(t, rawCard.Revision, reparsedCard.Revision)
	})

	t.Run("ToFile writes correct bytes to disk", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "output.png")
		err := rawCard.ToFile(filePath)
		require.NoError(t, err)

		fileBytes, err := os.ReadFile(filePath)
		require.NoError(t, err)
		reparsedCard, err := FromBytes(fileBytes).Get()
		require.NoError(t, err)
		assert.Equal(t, rawCard.RawCharaData, reparsedCard.RawCharaData)
	})

	t.Run("createCharaChunk uses fallback version", func(t *testing.T) {
		rawCard.Revision = 100
		finalBytes, err := rawCard.ToBytes()
		require.NoError(t, err)
		reparsedCard, err := FromBytes(finalBytes).Get()
		require.NoError(t, err)
		assert.Equal(t, character.RevisionV2, reparsedCard.Revision)
	})
}
