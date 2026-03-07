package pet

import (
	"fmt"
	"strings"
)

// Personality is a pre-compiled personality description built from animal + traits.
// Created once at startup and reused for every prompt.
type Personality struct {
	Animal      *Animal
	Name        string
	Traits      []Trait
	Description string // pre-compiled personality description for prompts
	Virtues     string // pre-compiled virtues section
	Defects     string // pre-compiled defects section
}

// BuildPersonality compiles an animal + traits + name into a reusable Personality.
func BuildPersonality(animal *Animal, traits []Trait, name string) *Personality {
	p := &Personality{
		Animal: animal,
		Name:   name,
		Traits: traits,
	}

	var virtueHints []string
	var defectHints []string
	for _, t := range traits {
		if t.IsDefect {
			defectHints = append(defectHints, t.PromptHint)
		} else {
			virtueHints = append(virtueHints, t.PromptHint)
		}
	}

	p.Virtues = strings.Join(virtueHints, " ")
	p.Defects = strings.Join(defectHints, " ")

	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("You are %s, a %s %s.", name, animal.Emoji, animal.Type))
	if len(virtueHints) > 0 {
		desc.WriteString(fmt.Sprintf(" Your strengths: %s", p.Virtues))
	}
	if len(defectHints) > 0 {
		desc.WriteString(fmt.Sprintf(" Your weaknesses: %s", p.Defects))
	}
	p.Description = desc.String()

	return p
}

// WrapPrompt builds a complete prompt for the AI, injecting personality + mood + context.
func WrapPrompt(p *Personality, mood Mood, question, kbContext string) string {
	var prompt strings.Builder

	// Personality
	prompt.WriteString(p.Description)
	prompt.WriteString("\n\n")

	// Current mood
	moodEmoji := p.Animal.Moods[mood]
	sound := p.Animal.Sounds[mood]
	prompt.WriteString(fmt.Sprintf("Your current mood: %s %s %s\n\n", mood.String(), moodEmoji, sound))

	// Instructions
	prompt.WriteString("Instructions:\n")
	prompt.WriteString("- Respond in character with your personality traits (both strengths and weaknesses).\n")
	prompt.WriteString("- Use unicode emojis naturally in your responses.\n")
	prompt.WriteString("- Be concise but stay in character.\n")
	prompt.WriteString("- Your defects should manifest subtly, not as caricature.\n")
	prompt.WriteString("- Answer the user's question directly, your personality colors HOW you answer, not WHAT you answer.\n\n")

	// KB context if available
	if kbContext != "" {
		prompt.WriteString("Relevant knowledge from your memory:\n")
		prompt.WriteString(kbContext)
		prompt.WriteString("\n\n")
	}

	// Question
	prompt.WriteString(fmt.Sprintf("User's question: %s", question))

	return prompt.String()
}

// FormatResponse wraps a response with the pet's visual identity.
func FormatResponse(animal *Animal, mood Mood, text string) string {
	moodEmoji := animal.Moods[mood]
	sound := animal.Sounds[mood]

	return fmt.Sprintf("%s %s\n\n%s\n\n%s", animal.Emoji, moodEmoji, text, sound)
}

// TraitNames returns a comma-separated list of trait names.
func (p *Personality) TraitNames() string {
	names := make([]string, len(p.Traits))
	for i, t := range p.Traits {
		names[i] = t.Name
	}
	return strings.Join(names, ", ")
}

// HasDefects returns true if the personality has any defect traits.
func (p *Personality) HasDefects() bool {
	for _, t := range p.Traits {
		if t.IsDefect {
			return true
		}
	}
	return false
}
