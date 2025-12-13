package character

import (
	"regexp"
	"strings"

	"github.com/r3dpixel/card-parser/property"
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/r3dpixel/toolkit/timestamp"
	"github.com/spf13/cast"
)

// Field names
const (
	NameField                    string = "name"
	DescriptionField             string = "description"
	PersonalityField             string = "personality"
	ScenarioField                string = "scenario"
	FirstMessageField            string = "first_mes"
	MessageExamplesField         string = "mes_example"
	CreatorNotesField            string = "creator_notes"
	PostHistoryInstructionsField string = "post_history_instructions"
	AlternateGreetingsField      string = "alternate_greetings"
	TagsField                    string = "tags"
	CreatorField                 string = "creator"
	DepthPromptKey               string = "depth_prompt"
	DepthPromptPromptKey         string = "prompt"
	DepthPromptDepthKey          string = "depth"
	DefaultDepth                 int    = 4
)

var (
	// Regexes to fix errors of the type {{{user}, {{char}, {char}}, {char} -> {{user}, {{char}}
	charRegex = regexp.MustCompile(`\{+char}+`)
	userRegex = regexp.MustCompile(`\{+user}+`)
)

// contentAlias alias for Content to avoid circular references
type contentAlias Content

// Content content of a V3 chara card
type Content struct {
	Title                   property.String      `json:"title"`
	Name                    property.String      `json:"name"`
	Description             property.String      `json:"description"`
	Personality             property.String      `json:"personality"`
	Scenario                property.String      `json:"scenario"`
	FirstMessage            property.String      `json:"first_mes"`
	MessageExamples         property.String      `json:"mes_example"`
	CreatorNotes            property.String      `json:"creator_notes"`
	SystemPrompt            property.String      `json:"system_prompt"`
	PostHistoryInstructions property.String      `json:"post_history_instructions"`
	AlternateGreetings      property.StringArray `json:"alternate_greetings"`
	CharacterBook           *Book                `json:"character_book,omitzero"`
	Tags                    property.StringArray `json:"tags"`
	Creator                 property.String      `json:"creator"`
	CharacterVersion        property.String      `json:"character_version"`
	DepthPrompt             DepthPrompt          `json:"-"`
	Extensions              map[string]any       `json:"extensions,omitzero"`

	Assets                   []Asset                    `json:"assets,omitzero"`
	Nickname                 property.String            `json:"nickname"`
	CreatorNotesMultilingual map[string]property.String `json:"creator_notes_multilingual,omitzero"`
	Source                   property.StringArray       `json:"source,omitzero"`
	GroupGreetings           property.StringArray       `json:"group_only_greetings,omitzero"`
	CreationDate             timestamp.Seconds          `json:"creation_date"`
	ModificationDate         timestamp.Seconds          `json:"modification_date"`

	SourceID    property.String `json:"source_id"`
	CharacterID property.String `json:"character_id"`
	PlatformID  property.String `json:"platform_id"`
	DirectLink  property.String `json:"direct_link"`
}

// DepthPrompt depth prompt structure of a V3 chara card
type DepthPrompt struct {
	Prompt string
	Depth  int
}

// MarshalJSON marshals Content into JSON format to respect Silly Tavern format using Sonic
func (c *Content) MarshalJSON() ([]byte, error) {
	// Insert depth prompt extension
	depthMap := c.insertDepthPrompt()
	// Purge depth prompt extension after marshaling (idempotent)
	defer c.purgeDepthPromptExtension(depthMap)
	// Delegate to Sonic encoder
	return sonicx.Config.Marshal((*contentAlias)(c))
}

// UnmarshalJSON unmarshals JSON into the Content, with fallbacks and best effort strategies using Sonic
func (c *Content) UnmarshalJSON(data []byte) error {
	// Unmarshal from JSON using Sonic
	if err := sonicx.Config.UnmarshalFromString(stringsx.FromBytes(data), (*contentAlias)(c)); err != nil {
		return err
	}
	c.extractDepthPrompt()

	// Decoding is complete
	return nil
}

// NormalizeSymbols replace all abnormal quotes, apostrophes or commas characters from ALL fields with the normal ASCII version (`"`, `,` `'`)
func (c *Content) NormalizeSymbols() {
	// Fix Quotes applied on every field
	c.Description.NormalizeSymbols()
	c.Personality.NormalizeSymbols()
	c.Scenario.NormalizeSymbols()
	c.FirstMessage.NormalizeSymbols()
	c.MessageExamples.NormalizeSymbols()
	c.CreatorNotes.NormalizeSymbols()
	c.SystemPrompt.NormalizeSymbols()
	c.PostHistoryInstructions.NormalizeSymbols()

	// Fix Quotes applied on each and every greeting
	greetings := c.AlternateGreetings
	for index := range greetings {
		greetings[index] = stringsx.NormalizeSymbols(greetings[index])
	}

	// Fix Quotes applied on every entry (name, comment, content)
	// Other fields ARE NOT affected (keywords, secondary keywords, etc.)
	if characterBook := c.CharacterBook; characterBook != nil {
		characterBook.NormalizeSymbols()
	}

	// Fix Quotes applied on the depth prompt content
	c.DepthPrompt.Prompt = stringsx.NormalizeSymbols(c.DepthPrompt.Prompt)
}

// FixUserCharTemplates fixes the user character templates for all fields: {{{user}, {{char}, {char}}, {char} -> {{user}, {{char}}
func (c *Content) FixUserCharTemplates() {
	c.Description = c.fixUserCharTemplateProp(c.Description)
	c.Personality = c.fixUserCharTemplateProp(c.Personality)
	c.Scenario = c.fixUserCharTemplateProp(c.Scenario)
	c.FirstMessage = c.fixUserCharTemplateProp(c.FirstMessage)
	c.MessageExamples = c.fixUserCharTemplateProp(c.MessageExamples)
	c.SystemPrompt = c.fixUserCharTemplateProp(c.SystemPrompt)
	c.PostHistoryInstructions = c.fixUserCharTemplateProp(c.PostHistoryInstructions)
	for index := range c.AlternateGreetings {
		c.AlternateGreetings[index] = c.fixUserCharTemplate(c.AlternateGreetings[index])
	}

	c.DepthPrompt.Prompt = c.fixUserCharTemplate(c.DepthPrompt.Prompt)

}

// fixUserCharTemplateProp fixes the user character templates for a property field: {{{user}, {{char}, {char}}, {char} -> {{user}, {{char}}
func (c *Content) fixUserCharTemplateProp(input property.String) property.String {
	return property.String(c.fixUserCharTemplate(string(input)))
}

// fixUserCharTemplate fixes the user character templates for a string field: {{{user}, {{char}, {char}}, {char} -> {{user}, {{char}}
func (c *Content) fixUserCharTemplate(input string) string {
	if stringsx.IsBlank(input) {
		return ""
	}
	result := charRegex.ReplaceAllString(input, "{{char}}")
	return userRegex.ReplaceAllString(result, "{{user}}")
}

// insertDepthPrompt inserts the depth prompt extension into the Extensions map
func (c *Content) insertDepthPrompt() map[string]any {
	// Skip if no prompt
	if stringsx.IsBlank(c.DepthPrompt.Prompt) {
		return nil
	}

	// Create the Extensions map if needed
	if c.Extensions == nil {
		c.Extensions = make(map[string]any)
	}

	// Set the depth map in the Extensions map
	depthMap, ok := c.Extensions[DepthPromptKey].(map[string]any)
	if !ok {
		depthMap = make(map[string]any)
		c.Extensions[DepthPromptKey] = depthMap
	}

	// Populate the depth map with the prompt and depth values
	depthMap[DepthPromptPromptKey] = c.DepthPrompt.Prompt
	depthMap[DepthPromptDepthKey] = c.DepthPrompt.Depth

	// Return the depth map
	return depthMap
}

// extractDepthPrompt extracts the depth prompt extension from the Extensions map and populates the DepthPrompt field
// Reverse of the insertDepthPrompt method
func (c *Content) extractDepthPrompt() {
	// Skip if no Extensions map
	if c.Extensions == nil {
		return
	}

	// Extract the depth prompt extension from the Extensions map
	promptValue, ok := c.Extensions[DepthPromptKey]
	if !ok || promptValue == nil {
		return
	}

	// Check the type of the prompt value and populate the DepthPrompt field
	switch typedPromptValue := promptValue.(type) {
	// If the extension is a map
	case map[string]any:
		// Purge from the Extensions map after extracting the depth prompt
		defer c.purgeDepthPromptExtension(typedPromptValue)
		// Clean the prompt
		prompt := strings.TrimSpace(cast.ToString(typedPromptValue[DepthPromptPromptKey]))
		// Skip if the prompt is blank
		if stringsx.IsBlank(prompt) {
			return
		}
		// Populate the DepthPrompt content field
		c.DepthPrompt.Prompt = prompt

		// Populate the DepthPrompt depth field
		// Initialize to default depth
		c.DepthPrompt.Depth = DefaultDepth
		// Check if depth value is present
		if depthValue := typedPromptValue[DepthPromptDepthKey]; depthValue != nil {
			// Convert depth value to int (if error, the default remains set)
			if depth, err := cast.ToIntE(depthValue); err == nil {
				// Set the depth field
				c.DepthPrompt.Depth = depth
			}
		}
	// If the extension is an array
	case []any:
		// Convert the array to JSON string
		c.DepthPrompt.Prompt = jsonx.String(promptValue)
		// Set the depth to default
		c.DepthPrompt.Depth = DefaultDepth
		// Remove the extension
		delete(c.Extensions, DepthPromptKey)
	// If the extension is a string or any other type
	default:
		// Convert the value to string
		c.DepthPrompt.Prompt = cast.ToString(promptValue)
		// Set the depth to default
		c.DepthPrompt.Depth = DefaultDepth
		// Remove the extension
		delete(c.Extensions, DepthPromptKey)
	}
}

// purgeDepthPromptExtension removes the depth prompt extension from the Extensions map if it is empty
func (c *Content) purgeDepthPromptExtension(depthMap map[string]any) {
	// Remove the prompt and depth keys from the depth map
	delete(depthMap, DepthPromptPromptKey)
	delete(depthMap, DepthPromptDepthKey)
	// Remove the depth map from the Extensions map if it is empty
	if len(depthMap) == 0 {
		delete(c.Extensions, DepthPromptKey)
	}

	// Remove the Extensions map if it is empty
	if len(c.Extensions) == 0 {
		c.Extensions = nil
	}
}

// Integrity checks if the sheet is malformed (missing necessary fields)
func (c *Content) Integrity() bool {
	// Check if title, name, description, creator, nickname and source_id are not blank
	// CreationDate and ModificationDate must be strictly positive
	// ModificationDate must be greater or equal than CreationDate (ModificationDate >= CreationDate)
	return stringsx.IsNotBlank(string(c.Title)) &&
		stringsx.IsNotBlank(string(c.Name)) &&
		stringsx.IsNotBlank(string(c.Description)) &&
		stringsx.IsNotBlank(string(c.Creator)) &&
		stringsx.IsNotBlank(string(c.Nickname)) &&
		c.CreationDate > 0 &&
		c.ModificationDate > 0 &&
		c.ModificationDate >= c.CreationDate &&
		stringsx.IsNotBlank(string(c.SourceID))
}
