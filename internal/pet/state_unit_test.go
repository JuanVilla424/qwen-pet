package pet

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLoadState_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "pet.json")

	s, err := LoadState(path, 24)
	if err != nil {
		t.Fatalf("LoadState() error: %v", err)
	}

	if s.Mood != MoodNeutral {
		t.Errorf("New state mood = %v, want MoodNeutral", s.Mood)
	}
	if s.InteractionCount != 0 {
		t.Errorf("InteractionCount = %d, want 0", s.InteractionCount)
	}
	if s.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestSaveAndLoadState(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pet.json")

	s, err := LoadState(path, 24)
	if err != nil {
		t.Fatalf("LoadState() error: %v", err)
	}

	s.RecordInteraction("kb")
	s.RecordInteraction("ai")
	s.RecordInteraction("telegram")

	if err := s.SaveState(); err != nil {
		t.Fatalf("SaveState() error: %v", err)
	}

	s2, err := LoadState(path, 24)
	if err != nil {
		t.Fatalf("LoadState(existing) error: %v", err)
	}

	if s2.InteractionCount != 3 {
		t.Errorf("InteractionCount = %d, want 3", s2.InteractionCount)
	}
	if s2.KBAnswers != 1 {
		t.Errorf("KBAnswers = %d, want 1", s2.KBAnswers)
	}
	if s2.AIAnswers != 1 {
		t.Errorf("AIAnswers = %d, want 1", s2.AIAnswers)
	}
	if s2.TelegramAnswers != 1 {
		t.Errorf("TelegramAnswers = %d, want 1", s2.TelegramAnswers)
	}
}

func TestRecordInteraction_CountersAndMood(t *testing.T) {
	s := &State{
		Mood:           MoodNeutral,
		moodDecayHours: 24,
	}

	s.RecordInteraction("kb")
	if s.KBAnswers != 1 {
		t.Errorf("KBAnswers = %d after kb interaction", s.KBAnswers)
	}
	if s.CorrectAnswers != 1 {
		t.Errorf("CorrectAnswers = %d after kb interaction, want 1", s.CorrectAnswers)
	}
	if s.Mood != MoodHappy {
		t.Errorf("Mood = %v after interaction, want MoodHappy", s.Mood)
	}

	s.RecordInteraction("telegram")
	if s.TelegramAnswers != 1 {
		t.Errorf("TelegramAnswers = %d after telegram interaction", s.TelegramAnswers)
	}
	if s.CorrectAnswers != 2 {
		t.Errorf("CorrectAnswers = %d after telegram, want 2", s.CorrectAnswers)
	}
}

func TestRecordInteraction_ExcitedEvery10(t *testing.T) {
	s := &State{
		InteractionCount: 9,
		Mood:             MoodNeutral,
		moodDecayHours:   24,
	}

	s.RecordInteraction("ai")

	if s.InteractionCount != 10 {
		t.Errorf("InteractionCount = %d, want 10", s.InteractionCount)
	}
	if s.Mood != MoodExcited {
		t.Errorf("Mood at 10 interactions = %v, want MoodExcited", s.Mood)
	}
}

func TestRecordInteraction_RecoverFromSad(t *testing.T) {
	s := &State{
		Mood:           MoodSad,
		moodDecayHours: 24,
	}

	s.RecordInteraction("kb")

	if s.Mood != MoodNeutral {
		t.Errorf("Mood after interaction from sad = %v, want MoodNeutral", s.Mood)
	}
}

func TestApplyMoodDecay(t *testing.T) {
	tests := []struct {
		name    string
		hours   float64
		decay   int
		want    Mood
		initial Mood
	}{
		{"no decay", 1, 24, MoodHappy, MoodHappy},
		{"tired after threshold", 25, 24, MoodTired, MoodHappy},
		{"sad after 2x threshold", 49, 24, MoodSad, MoodHappy},
		{"zero last interaction", 0, 24, MoodHappy, MoodHappy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				Mood:           tt.initial,
				moodDecayHours: tt.decay,
			}
			if tt.hours > 0 {
				s.LastInteraction = time.Now().Add(-time.Duration(tt.hours) * time.Hour)
			}

			s.applyMoodDecay()

			if s.Mood != tt.want {
				t.Errorf("mood = %v, want %v", s.Mood, tt.want)
			}
		})
	}
}

func TestGetSetMood_ThreadSafe(t *testing.T) {
	s := &State{Mood: MoodNeutral, moodDecayHours: 24}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.SetMood(MoodHappy)
		}()
		go func() {
			defer wg.Done()
			_ = s.GetMood()
		}()
	}
	wg.Wait()
}

func TestRecordInteraction_ThreadSafe(t *testing.T) {
	s := &State{Mood: MoodNeutral, moodDecayHours: 24}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.RecordInteraction("kb")
		}()
	}
	wg.Wait()

	if s.InteractionCount != 50 {
		t.Errorf("InteractionCount = %d after 50 concurrent interactions, want 50", s.InteractionCount)
	}
}

func TestStats_Format(t *testing.T) {
	s := &State{
		Mood:             MoodHappy,
		InteractionCount: 10,
		TotalQuestions:   10,
		KBAnswers:        7,
		AIAnswers:        2,
		TelegramAnswers:  1,
		LastInteraction:  time.Now(),
		CreatedAt:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		moodDecayHours:   24,
	}

	stats := s.Stats()

	if stats == "" {
		t.Fatal("Stats() returned empty string")
	}
	if !contains(stats, "70%") {
		t.Errorf("Stats should show 70%% resolution rate, got: %s", stats)
	}
	if !contains(stats, "2026-01-01") {
		t.Errorf("Stats should show born date, got: %s", stats)
	}
}

func TestStats_NoInteractions(t *testing.T) {
	s := &State{
		Mood:           MoodNeutral,
		CreatedAt:      time.Now(),
		moodDecayHours: 24,
	}

	stats := s.Stats()
	if !contains(stats, "never") {
		t.Errorf("Stats with no interactions should say 'never', got: %s", stats)
	}
	if !contains(stats, "0%") {
		t.Errorf("Stats with no questions should show 0%%, got: %s", stats)
	}
}

func TestLoadState_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pet.json")

	if err := os.WriteFile(path, []byte("{invalid}"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadState(path, 24)
	if err == nil {
		t.Error("LoadState(invalid JSON) should return error")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && searchSubstring(s, substr))
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
