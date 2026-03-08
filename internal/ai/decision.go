package ai

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/JuanVilla424/qwen-pet/internal/config"
	"github.com/JuanVilla424/qwen-pet/internal/kb"
	"github.com/JuanVilla424/qwen-pet/internal/pet"
)

// Answer represents the decision engine's response.
type Answer struct {
	Text       string
	Source     string // "kb", "ai", "telegram"
	Confidence float32
	PetMood    pet.Mood
	KBResults  []kb.Result
}

// Escalator is the interface for Telegram escalation.
type Escalator interface {
	Escalate(ctx context.Context, question, kbContext string) (string, error)
}

// Engine orchestrates KB lookup, AI inference, and Telegram escalation.
type Engine struct {
	client      *Client
	store       *kb.Store
	personality *pet.Personality
	petState    *pet.State
	escalator   Escalator
	cfg         config.DecisionCfg
}

// NewEngine creates a decision engine.
func NewEngine(client *Client, store *kb.Store, personality *pet.Personality, petState *pet.State, escalator Escalator, cfg config.DecisionCfg) *Engine {
	return &Engine{
		client:      client,
		store:       store,
		personality: personality,
		petState:    petState,
		escalator:   escalator,
		cfg:         cfg,
	}
}

// Decide processes a question through the KB → AI → Telegram pipeline.
func (e *Engine) Decide(ctx context.Context, question, extraContext string) (Answer, error) {
	e.petState.SetMood(pet.MoodThinking)

	// Step 1: Query KB
	results, err := e.store.Query(ctx, question, 5)
	if err != nil {
		slog.Error("KB query failed", "error", err)
		results = nil
	}

	// Step 2: Check if KB has a confident answer
	if len(results) > 0 && results[0].Similarity >= e.cfg.KBThreshold {
		slog.Info("KB hit", "similarity", results[0].Similarity, "id", results[0].ID)
		text := results[0].Content
		e.petState.RecordInteraction("kb")
		_ = e.petState.SaveState()
		return Answer{
			Text:       text,
			Source:     "kb",
			Confidence: results[0].Similarity,
			PetMood:    pet.MoodHappy,
			KBResults:  results,
		}, nil
	}

	// Step 3: Build context from partial KB matches
	kbContext := buildKBContext(results)
	if extraContext != "" {
		kbContext = extraContext + "\n\n" + kbContext
	}

	// Step 4: Ask AI with personality
	mood := e.petState.GetMood()
	prompt := pet.WrapPrompt(e.personality, mood, question, kbContext)

	messages := []Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: question},
	}

	aiResponse, usage, err := e.client.ChatCompletion(ctx, messages)
	if err != nil {
		slog.Error("AI completion failed", "error", err)
	} else {
		slog.Info("AI response", "tokens", usage.TotalTokens)

		// If AI gives a reasonable response, use it
		if aiResponse != "" {
			text := aiResponse

			// Store the Q&A in KB for future reference
			e.storeInKB(ctx, question, aiResponse)

			e.petState.RecordInteraction("ai")
			_ = e.petState.SaveState()
			return Answer{
				Text:       text,
				Source:     "ai",
				Confidence: e.cfg.AIThreshold,
				PetMood:    pet.MoodHappy,
				KBResults:  results,
			}, nil
		}
	}

	// Step 5: Escalate to Telegram
	if e.escalator == nil {
		return Answer{
			Text:    "I don't know the answer and can't reach you via Telegram.",
			Source:  "none",
			PetMood: pet.MoodSad,
		}, nil
	}

	slog.Info("escalating to Telegram", "question", question)
	telegramResponse, err := e.escalator.Escalate(ctx, question, kbContext)
	if err != nil {
		return Answer{
			Text:    fmt.Sprintf("I tried to ask you via Telegram but: %v", err),
			Source:  "telegram",
			PetMood: pet.MoodSad,
		}, nil
	}

	// Store the user's response in KB
	e.storeInKB(ctx, question, telegramResponse)

	text := telegramResponse
	e.petState.RecordInteraction("telegram")
	_ = e.petState.SaveState()

	return Answer{
		Text:       text,
		Source:     "telegram",
		Confidence: 1.0,
		PetMood:    pet.MoodExcited,
		KBResults:  results,
	}, nil
}

func (e *Engine) storeInKB(ctx context.Context, question, answer string) {
	id := fmt.Sprintf("qa-%d", time.Now().UnixNano())
	content := fmt.Sprintf("Q: %s\nA: %s", question, answer)
	meta := map[string]string{
		"category": "qa",
		"type":     "auto-learned",
	}
	if err := e.store.Add(ctx, id, content, meta); err != nil {
		slog.Error("failed to store in KB", "error", err)
	}
}

func buildKBContext(results []kb.Result) string {
	if len(results) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, r := range results {
		if r.Similarity < 0.5 {
			continue
		}
		sb.WriteString(fmt.Sprintf("[%d] (%.0f%% match) %s\n", i+1, r.Similarity*100, r.Content))
	}
	return sb.String()
}
