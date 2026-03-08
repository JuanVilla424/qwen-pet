package pet

import (
	"strings"
	"testing"
)

func newTestAnimal() *Animal {
	return &Animal{
		Type: "fox", Emoji: "\U0001F98A", DefaultName: "Kit",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60F", MoodNeutral: "\U0001F98A", MoodThinking: "\U0001F914", MoodTired: "\U0001F971", MoodExcited: "\U0001F525", MoodSad: "\U0001F614"},
		Sounds: map[Mood]string{MoodHappy: "*yip yip!*", MoodNeutral: "*sniff*", MoodThinking: "*hmm...*", MoodTired: "*yawn*", MoodExcited: "*YIIIP!*", MoodSad: "*whine...*"},
	}
}

func newTestTraits() []Trait {
	return []Trait{
		{ID: "curious", Name: "Curious", IsDefect: false, PromptHint: "You investigate thoroughly."},
		{ID: "sarcastic", Name: "Sarcastic", IsDefect: true, PromptHint: "You use sarcasm."},
	}
}

func TestBuildPersonality(t *testing.T) {
	animal := newTestAnimal()
	traits := newTestTraits()

	p := BuildPersonality(animal, traits, "Kit", nil)

	if p.Name != "Kit" {
		t.Errorf("Name = %q, want %q", p.Name, "Kit")
	}
	if p.Animal != animal {
		t.Error("Animal pointer mismatch")
	}
	if len(p.Traits) != 2 {
		t.Errorf("Traits count = %d, want 2", len(p.Traits))
	}
	if !strings.Contains(p.Description, "Kit") {
		t.Error("Description should contain pet name")
	}
	if !strings.Contains(p.Description, "fox") {
		t.Error("Description should contain animal type")
	}
	if !strings.Contains(p.Virtues, "investigate") {
		t.Error("Virtues should contain virtue prompt hint")
	}
	if !strings.Contains(p.Defects, "sarcasm") {
		t.Error("Defects should contain defect prompt hint")
	}
}

func TestBuildPersonality_OnlyVirtues(t *testing.T) {
	animal := newTestAnimal()
	traits := []Trait{{ID: "honest", Name: "Honest", IsDefect: false, PromptHint: "You speak the truth."}}

	p := BuildPersonality(animal, traits, "Kit", nil)

	if !strings.Contains(p.Description, "strengths") {
		t.Error("Description should contain strengths section")
	}
	if strings.Contains(p.Description, "weaknesses") {
		t.Error("Description should NOT contain weaknesses when no defects")
	}
}

func TestBuildPersonality_OnlyDefects(t *testing.T) {
	animal := newTestAnimal()
	traits := []Trait{{ID: "lazy", Name: "Lazy", IsDefect: true, PromptHint: "You take shortcuts."}}

	p := BuildPersonality(animal, traits, "Kit", nil)

	if strings.Contains(p.Description, "strengths") {
		t.Error("Description should NOT contain strengths when no virtues")
	}
	if !strings.Contains(p.Description, "weaknesses") {
		t.Error("Description should contain weaknesses section")
	}
}

func TestBuildPersonality_NoTraits(t *testing.T) {
	animal := newTestAnimal()

	p := BuildPersonality(animal, []Trait{}, "Kit", nil)

	if !strings.Contains(p.Description, "Kit") {
		t.Error("Description should still contain pet name")
	}
	if strings.Contains(p.Description, "strengths") || strings.Contains(p.Description, "weaknesses") {
		t.Error("No traits should produce no strengths/weaknesses sections")
	}
}

func TestWrapPrompt(t *testing.T) {
	animal := newTestAnimal()
	traits := newTestTraits()
	p := BuildPersonality(animal, traits, "Kit", nil)

	prompt := WrapPrompt(p, MoodHappy, "what is Go?", "Go is a language")

	if !strings.Contains(prompt, "Kit") {
		t.Error("Prompt should contain personality description")
	}
	if !strings.Contains(prompt, "happy") {
		t.Error("Prompt should contain mood string")
	}
	if !strings.Contains(prompt, "Go is a language") {
		t.Error("Prompt should contain KB context")
	}
	if !strings.Contains(prompt, "what is Go?") {
		t.Error("Prompt should contain the question")
	}
	if !strings.Contains(prompt, "Instructions:") {
		t.Error("Prompt should contain instructions section")
	}
}

func TestWrapPrompt_NoKBContext(t *testing.T) {
	animal := newTestAnimal()
	p := BuildPersonality(animal, newTestTraits(), "Kit", nil)

	prompt := WrapPrompt(p, MoodThinking, "test question", "")

	if strings.Contains(prompt, "Relevant knowledge") {
		t.Error("Prompt should not contain KB section when no context")
	}
}

func TestFormatResponse(t *testing.T) {
	animal := newTestAnimal()
	result := FormatResponse(animal, MoodHappy, "This is the answer")

	if !strings.Contains(result, animal.Emoji) {
		t.Error("FormatResponse should contain animal emoji")
	}
	if !strings.Contains(result, animal.Moods[MoodHappy]) {
		t.Error("FormatResponse should contain mood emoji")
	}
	if !strings.Contains(result, "This is the answer") {
		t.Error("FormatResponse should contain the text")
	}
	if !strings.Contains(result, animal.Sounds[MoodHappy]) {
		t.Error("FormatResponse should contain mood sound")
	}
}

func TestPersonality_TraitNames(t *testing.T) {
	animal := newTestAnimal()
	traits := newTestTraits()
	p := BuildPersonality(animal, traits, "Kit", nil)

	names := p.TraitNames()

	if !strings.Contains(names, "Curious") {
		t.Errorf("TraitNames() = %q, should contain Curious", names)
	}
	if !strings.Contains(names, "Sarcastic") {
		t.Errorf("TraitNames() = %q, should contain Sarcastic", names)
	}
}

func TestPersonality_HasDefects(t *testing.T) {
	animal := newTestAnimal()

	withDefects := BuildPersonality(animal, newTestTraits(), "Kit", nil)
	if !withDefects.HasDefects() {
		t.Error("HasDefects() should return true when defects present")
	}

	noDefects := BuildPersonality(animal, []Trait{
		{ID: "honest", IsDefect: false},
	}, "Kit", nil)
	if noDefects.HasDefects() {
		t.Error("HasDefects() should return false when no defects")
	}
}
