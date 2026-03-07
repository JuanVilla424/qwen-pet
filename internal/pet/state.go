package pet

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// State holds the pet's mutable state that persists between restarts.
type State struct {
	mu sync.RWMutex

	Mood             Mood      `json:"mood"`
	InteractionCount int64     `json:"interaction_count"`
	LastInteraction  time.Time `json:"last_interaction"`
	TotalQuestions   int64     `json:"total_questions"`
	CorrectAnswers   int64     `json:"correct_answers"`
	KBAnswers        int64     `json:"kb_answers"`
	AIAnswers        int64     `json:"ai_answers"`
	TelegramAnswers  int64     `json:"telegram_answers"`
	CreatedAt        time.Time `json:"created_at"`

	path           string
	moodDecayHours int
}

// LoadState reads pet state from a JSON file or creates a new one.
func LoadState(path string, moodDecayHours int) (*State, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create state directory: %w", err)
	}

	s := &State{
		path:           path,
		moodDecayHours: moodDecayHours,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			s.Mood = MoodNeutral
			s.CreatedAt = time.Now()
			return s, nil
		}
		return nil, fmt.Errorf("read state file: %w", err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}

	s.path = path
	s.moodDecayHours = moodDecayHours
	s.applyMoodDecay()

	return s, nil
}

// SaveState writes the current state to disk.
func (s *State) SaveState() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write state: %w", err)
	}

	return nil
}

// RecordInteraction updates counters and mood after an interaction.
func (s *State) RecordInteraction(source string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.InteractionCount++
	s.TotalQuestions++
	s.LastInteraction = time.Now()

	switch source {
	case "kb":
		s.KBAnswers++
		s.CorrectAnswers++
	case "ai":
		s.AIAnswers++
	case "telegram":
		s.TelegramAnswers++
		s.CorrectAnswers++
	}

	s.updateMoodAfterInteraction()
}

// GetMood returns the current mood (thread-safe).
func (s *State) GetMood() Mood {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Mood
}

// SetMood explicitly sets the mood (thread-safe).
func (s *State) SetMood(mood Mood) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Mood = mood
}

func (s *State) updateMoodAfterInteraction() {
	switch {
	case s.InteractionCount%10 == 0:
		s.Mood = MoodExcited
	case s.Mood == MoodSad || s.Mood == MoodTired:
		s.Mood = MoodNeutral
	default:
		s.Mood = MoodHappy
	}
}

func (s *State) applyMoodDecay() {
	if s.LastInteraction.IsZero() {
		return
	}

	hours := time.Since(s.LastInteraction).Hours()
	decayThreshold := float64(s.moodDecayHours)

	switch {
	case hours > decayThreshold*2:
		s.Mood = MoodSad
	case hours > decayThreshold:
		s.Mood = MoodTired
	}
}

// Stats returns a formatted string with pet statistics.
func (s *State) Stats() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var lastSeen string
	if s.LastInteraction.IsZero() {
		lastSeen = "never"
	} else {
		lastSeen = time.Since(s.LastInteraction).Round(time.Minute).String() + " ago"
	}

	resolutionRate := float64(0)
	if s.TotalQuestions > 0 {
		resolutionRate = float64(s.KBAnswers) / float64(s.TotalQuestions) * 100
	}

	return fmt.Sprintf(
		"Mood: %s | Interactions: %d | Questions: %d\n"+
			"KB: %d | AI: %d | Telegram: %d\n"+
			"Resolution rate: %.0f%% | Last seen: %s\n"+
			"Born: %s",
		s.Mood, s.InteractionCount, s.TotalQuestions,
		s.KBAnswers, s.AIAnswers, s.TelegramAnswers,
		resolutionRate, lastSeen,
		s.CreatedAt.Format("2006-01-02"),
	)
}
