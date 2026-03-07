package pet

import "testing"

func TestGetAnimal_ValidTypes(t *testing.T) {
	types := []string{
		"cat", "dog", "fox", "owl", "dragon", "rabbit", "penguin", "snake",
		"panda", "unicorn", "frog", "lion", "wolf", "bear", "eagle",
		"dolphin", "lizard", "octopus",
	}

	for _, typ := range types {
		a, err := GetAnimal(typ)
		if err != nil {
			t.Errorf("GetAnimal(%q) returned error: %v", typ, err)
			continue
		}
		if a.Type != typ {
			t.Errorf("GetAnimal(%q).Type = %q, want %q", typ, a.Type, typ)
		}
		if a.Emoji == "" {
			t.Errorf("GetAnimal(%q).Emoji is empty", typ)
		}
		if a.DefaultName == "" {
			t.Errorf("GetAnimal(%q).DefaultName is empty", typ)
		}
	}
}

func TestGetAnimal_InvalidType(t *testing.T) {
	_, err := GetAnimal("invalid_animal")
	if err == nil {
		t.Error("GetAnimal(invalid) should return error")
	}
}

func TestGetAnimal_MoodsComplete(t *testing.T) {
	moods := []Mood{MoodHappy, MoodNeutral, MoodThinking, MoodTired, MoodExcited, MoodSad}

	for _, a := range ListAnimals() {
		for _, m := range moods {
			if _, ok := a.Moods[m]; !ok {
				t.Errorf("animal %q missing mood emoji for %s", a.Type, m)
			}
			if _, ok := a.Sounds[m]; !ok {
				t.Errorf("animal %q missing sound for %s", a.Type, m)
			}
		}
	}
}

func TestListAnimals_MinimumCount(t *testing.T) {
	animals := ListAnimals()
	if len(animals) < 15 {
		t.Errorf("ListAnimals() returned %d animals, want at least 15", len(animals))
	}
}

func TestAnimalExists(t *testing.T) {
	if !AnimalExists("cat") {
		t.Error("AnimalExists(cat) = false, want true")
	}
	if AnimalExists("dinosaur") {
		t.Error("AnimalExists(dinosaur) = true, want false")
	}
}

func TestGetTrait_ValidVirtues(t *testing.T) {
	virtues := []string{
		"curious", "analytical", "creative", "patient", "enthusiastic",
		"loyal", "strategic", "honest", "protective", "adaptable",
		"meticulous", "humorous", "empathetic", "pragmatic", "assertive",
	}

	for _, id := range virtues {
		tr, err := GetTrait(id)
		if err != nil {
			t.Errorf("GetTrait(%q) returned error: %v", id, err)
			continue
		}
		if tr.IsDefect {
			t.Errorf("GetTrait(%q).IsDefect = true, want false (it's a virtue)", id)
		}
		if tr.PromptHint == "" {
			t.Errorf("GetTrait(%q).PromptHint is empty", id)
		}
	}
}

func TestGetTrait_ValidDefects(t *testing.T) {
	defects := []string{
		"impatient", "stubborn", "overthinks", "sarcastic", "forgetful",
		"blunt", "perfectionist", "anxious", "lazy", "dramatic",
		"distracted", "competitive",
	}

	for _, id := range defects {
		tr, err := GetTrait(id)
		if err != nil {
			t.Errorf("GetTrait(%q) returned error: %v", id, err)
			continue
		}
		if !tr.IsDefect {
			t.Errorf("GetTrait(%q).IsDefect = false, want true (it's a defect)", id)
		}
	}
}

func TestGetTrait_InvalidTrait(t *testing.T) {
	_, err := GetTrait("nonexistent")
	if err == nil {
		t.Error("GetTrait(nonexistent) should return error")
	}
}

func TestResolveTraits(t *testing.T) {
	ids := []string{"curious", "sarcastic", "honest"}
	traits, err := ResolveTraits(ids)
	if err != nil {
		t.Fatalf("ResolveTraits() returned error: %v", err)
	}
	if len(traits) != 3 {
		t.Fatalf("ResolveTraits() returned %d traits, want 3", len(traits))
	}
}

func TestResolveTraits_InvalidID(t *testing.T) {
	ids := []string{"curious", "nonexistent"}
	_, err := ResolveTraits(ids)
	if err == nil {
		t.Error("ResolveTraits(with invalid) should return error")
	}
}

func TestResolveTraits_Empty(t *testing.T) {
	traits, err := ResolveTraits([]string{})
	if err != nil {
		t.Fatalf("ResolveTraits(empty) returned error: %v", err)
	}
	if len(traits) != 0 {
		t.Errorf("ResolveTraits(empty) returned %d traits, want 0", len(traits))
	}
}

func TestListVirtues_AllPositive(t *testing.T) {
	virtues := ListVirtues()
	if len(virtues) < 10 {
		t.Errorf("ListVirtues() returned %d, want at least 10", len(virtues))
	}
	for _, v := range virtues {
		if v.IsDefect {
			t.Errorf("ListVirtues() returned defect: %q", v.ID)
		}
	}
}

func TestListDefects_AllNegative(t *testing.T) {
	defects := ListDefects()
	if len(defects) < 10 {
		t.Errorf("ListDefects() returned %d, want at least 10", len(defects))
	}
	for _, d := range defects {
		if !d.IsDefect {
			t.Errorf("ListDefects() returned virtue: %q", d.ID)
		}
	}
}

func TestTraitExists(t *testing.T) {
	if !TraitExists("curious") {
		t.Error("TraitExists(curious) = false, want true")
	}
	if TraitExists("nonexistent") {
		t.Error("TraitExists(nonexistent) = true, want false")
	}
}

func TestMoodString(t *testing.T) {
	tests := []struct {
		mood Mood
		want string
	}{
		{MoodHappy, "happy"},
		{MoodNeutral, "neutral"},
		{MoodThinking, "thinking"},
		{MoodTired, "tired"},
		{MoodExcited, "excited"},
		{MoodSad, "sad"},
		{Mood(99), "neutral"},
	}

	for _, tt := range tests {
		got := tt.mood.String()
		if got != tt.want {
			t.Errorf("Mood(%d).String() = %q, want %q", tt.mood, got, tt.want)
		}
	}
}
