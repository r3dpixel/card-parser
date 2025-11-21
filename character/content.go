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
	charRegex = regexp.MustCompile(`\{+char}+`)
	userRegex = regexp.MustCompile(`\{+user}+`)
)

// Content content of a V3 chara card
type contentAlias Content
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
	depthMap := c.insertDepthPrompt()
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

	c.DepthPrompt.Prompt = stringsx.NormalizeSymbols(c.DepthPrompt.Prompt)
}

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

func (c *Content) fixUserCharTemplateProp(input property.String) property.String {
	return property.String(c.fixUserCharTemplate(string(input)))
}

func (c *Content) fixUserCharTemplate(input string) string {
	if stringsx.IsBlank(input) {
		return stringsx.Empty
	}
	result := charRegex.ReplaceAllString(input, "{{char}}")
	return userRegex.ReplaceAllString(result, "{{user}}")
}

func (c *Content) insertDepthPrompt() map[string]any {
	if stringsx.IsBlank(c.DepthPrompt.Prompt) {
		return nil
	}

	if c.Extensions == nil {
		c.Extensions = make(map[string]any)
	}

	depthMap, ok := c.Extensions[DepthPromptKey].(map[string]any)
	if !ok {
		depthMap = make(map[string]any)
		c.Extensions[DepthPromptKey] = depthMap
	}

	depthMap[DepthPromptPromptKey] = c.DepthPrompt.Prompt
	depthMap[DepthPromptDepthKey] = c.DepthPrompt.Depth

	return depthMap
}

func (c *Content) extractDepthPrompt() {
	if c.Extensions == nil {
		return
	}

	promptValue, ok := c.Extensions[DepthPromptKey]
	if !ok || promptValue == nil {
		return
	}

	switch typedPromptValue := promptValue.(type) {
	case map[string]any:
		defer c.purgeDepthPromptExtension(typedPromptValue)
		prompt := strings.TrimSpace(cast.ToString(typedPromptValue[DepthPromptPromptKey]))
		if stringsx.IsBlank(prompt) {
			return
		}
		c.DepthPrompt.Prompt = prompt
		c.DepthPrompt.Depth = DefaultDepth
		if depthValue := typedPromptValue[DepthPromptDepthKey]; depthValue != nil {
			if depth, err := cast.ToIntE(depthValue); err == nil {
				c.DepthPrompt.Depth = depth
			}
		}
	case []any:
		c.DepthPrompt.Prompt = jsonx.String(promptValue)
		c.DepthPrompt.Depth = DefaultDepth
		delete(c.Extensions, DepthPromptKey)
	default:
		c.DepthPrompt.Prompt = cast.ToString(promptValue)
		c.DepthPrompt.Depth = DefaultDepth
		delete(c.Extensions, DepthPromptKey)
	}
}

func (c *Content) purgeDepthPromptExtension(depthMap map[string]any) {
	delete(depthMap, DepthPromptPromptKey)
	delete(depthMap, DepthPromptDepthKey)
	if len(depthMap) == 0 {
		delete(c.Extensions, DepthPromptKey)
	}

	if len(c.Extensions) == 0 {
		c.Extensions = nil
	}
}

// Integrity checks if the sheet is malformed (missing necessary fields)
func (c *Content) Integrity() bool {
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
