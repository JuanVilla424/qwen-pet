package pet

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
)

// Profile holds the pet's identity chosen during onboarding.
type Profile struct {
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Traits []string `json:"traits"`
	Soul   []string `json:"soul"`
}

// LoadProfile reads a pet profile from a JSON file.
// Returns nil, nil if the file does not exist (onboarding needed).
func LoadProfile(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read profile: %w", err)
	}

	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse profile: %w", err)
	}

	if p.Type == "" || p.Name == "" || len(p.Traits) == 0 {
		return nil, fmt.Errorf("profile is incomplete: type, name, and traits are required")
	}

	return &p, nil
}

// RandomProfile generates a random pet profile (animal + name + 2-4 virtues + 1-2 defects).
func RandomProfile() *Profile {
	animals := ListAnimals()
	animal := animals[rand.IntN(len(animals))]

	virtues := ListVirtues()
	defects := ListDefects()

	rand.Shuffle(len(virtues), func(i, j int) { virtues[i], virtues[j] = virtues[j], virtues[i] })
	rand.Shuffle(len(defects), func(i, j int) { defects[i], defects[j] = defects[j], defects[i] })

	numVirtues := 2 + rand.IntN(3) // 2-4
	numDefects := 1 + rand.IntN(2) // 1-2

	var traitIDs []string
	for i := 0; i < numVirtues && i < len(virtues); i++ {
		traitIDs = append(traitIDs, virtues[i].ID)
	}
	for i := 0; i < numDefects && i < len(defects); i++ {
		traitIDs = append(traitIDs, defects[i].ID)
	}

	return &Profile{
		Type:   animal.Type,
		Name:   animal.DefaultName,
		Traits: traitIDs,
	}
}

// Summary returns a human-readable summary of the profile.
func (p *Profile) Summary() string {
	animal, _ := GetAnimal(p.Type)
	traits, _ := ResolveTraits(p.Traits)

	var virtueNames, defectNames []string
	for _, t := range traits {
		if t.IsDefect {
			defectNames = append(defectNames, t.Name)
		} else {
			virtueNames = append(virtueNames, t.Name)
		}
	}

	emoji := ""
	if animal != nil {
		emoji = animal.Emoji
	}

	return fmt.Sprintf("%s %s the %s\nVirtues: %s\nDefects: %s",
		emoji, p.Name, p.Type,
		strings.Join(virtueNames, ", "),
		strings.Join(defectNames, ", "),
	)
}

// SaveProfile writes the pet profile to a JSON file.
func (p *Profile) SaveProfile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create profile directory: %w", err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write profile: %w", err)
	}

	return nil
}
