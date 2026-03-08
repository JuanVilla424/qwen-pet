package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/JuanVilla424/qwen-pet/internal/ai"
	"github.com/JuanVilla424/qwen-pet/internal/pet"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// pendingRequest represents an escalation waiting for user response.
type pendingRequest struct {
	question string
	ch       chan string
}

// Bot manages the Telegram bot lifecycle and escalation.
type Bot struct {
	tgBot         *bot.Bot
	allowedUserID int64
	animal        *pet.Animal
	petState      *pet.State
	personality   *pet.Personality
	aiClient      *ai.Client
	engine        *ai.Engine
	timeout       time.Duration

	profilePath string

	mu            sync.Mutex
	pending       *pendingRequest
	pendingReroll *pet.Profile
}

// NewBot creates and configures a Telegram bot.
func NewBot(token string, userID int64, animal *pet.Animal, petState *pet.State, personality *pet.Personality, aiClient *ai.Client, timeoutMinutes int, profilePath string) (*Bot, error) {
	b := &Bot{
		allowedUserID: userID,
		animal:        animal,
		petState:      petState,
		personality:   personality,
		aiClient:      aiClient,
		timeout:       time.Duration(timeoutMinutes) * time.Minute,
		profilePath:   profilePath,
	}

	opts := []bot.Option{
		bot.WithMiddlewares(b.userFilter),
		bot.WithDefaultHandler(b.handleDefault),
	}

	tgBot, err := bot.New(token, opts...)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	b.tgBot = tgBot

	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.handleStart)
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/status", bot.MatchTypeExact, b.handleStatus)
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/pet", bot.MatchTypeExact, b.handlePetInfo)
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, b.handleHelp)
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/reroll", bot.MatchTypeExact, b.handleReroll)

	return b, nil
}

// SetEngine binds the decision engine after construction (resolves circular dependency).
func (b *Bot) SetEngine(engine *ai.Engine) {
	b.engine = engine
}

// Start begins polling for updates (runs in goroutine).
func (b *Bot) Start(ctx context.Context) {
	slog.Info("starting Telegram bot polling")
	b.tgBot.Start(ctx)
}

// askAI generates a response using the AI with the pet's personality.
func (b *Bot) askAI(ctx context.Context, prompt string) string {
	mood := b.petState.GetMood()
	wrapped := pet.WrapPrompt(b.personality, mood, prompt, "")
	msgs := []ai.Message{{Role: "user", Content: wrapped}}

	response, _, err := b.aiClient.ChatCompletion(ctx, msgs)
	if err != nil {
		slog.Error("AI response failed", "error", err)
		return pet.FormatResponse(b.animal, mood, "I'm having trouble thinking right now. Try again in a moment.")
	}

	return pet.FormatResponse(b.animal, mood, response)
}

// growSoul asks the AI for a new soul word after an interaction.
// Runs async to not block the response.
func (b *Bot) growSoul(ctx context.Context, interaction string) {
	b.mu.Lock()
	currentSoul := b.personality.Soul
	profilePath := b.profilePath
	b.mu.Unlock()

	if len(currentSoul) >= pet.SoulCapacity {
		return
	}

	prompt := pet.SoulGrowthPrompt(currentSoul, interaction)
	msgs := []ai.Message{{Role: "user", Content: prompt}}

	response, _, err := b.aiClient.ChatCompletion(ctx, msgs)
	if err != nil {
		slog.Debug("soul growth failed", "error", err)
		return
	}

	word := pet.ParseSoulWord(response)
	if word == "" || pet.SoulContains(currentSoul, word) {
		return
	}

	b.mu.Lock()
	b.personality.Soul = append(b.personality.Soul, word)
	newSoul := make([]string, len(b.personality.Soul))
	copy(newSoul, b.personality.Soul)
	b.mu.Unlock()

	// Persist to profile
	profile := &pet.Profile{
		Type:   b.animal.Type,
		Name:   b.personality.Name,
		Traits: make([]string, len(b.personality.Traits)),
		Soul:   newSoul,
	}
	for i, t := range b.personality.Traits {
		profile.Traits[i] = t.ID
	}
	if err := profile.SaveProfile(profilePath); err != nil {
		slog.Error("failed to save soul growth", "error", err)
		return
	}

	slog.Info("soul grew", "word", word, "total", len(newSoul))
}

// Escalate sends a question to the user and waits for a response.
func (b *Bot) Escalate(ctx context.Context, question, kbContext string) (string, error) {
	b.mu.Lock()
	if b.pending != nil {
		b.mu.Unlock()
		return "", fmt.Errorf("another escalation is already pending")
	}

	ch := make(chan string, 1)
	b.pending = &pendingRequest{
		question: question,
		ch:       ch,
	}
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		b.pending = nil
		b.mu.Unlock()
	}()

	emoji := b.animal.Emoji
	moodEmoji := b.animal.Moods[pet.MoodThinking]
	msg := fmt.Sprintf("%s %s I need your help!\n\nQuestion: %s", emoji, moodEmoji, question)
	if kbContext != "" {
		msg += fmt.Sprintf("\n\nContext:\n%s", kbContext)
	}
	msg += "\n\nReply with your answer..."

	_, err := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: b.allowedUserID,
		Text:   msg,
	})
	if err != nil {
		return "", fmt.Errorf("send escalation message: %w", err)
	}

	select {
	case response := <-ch:
		return response, nil
	case <-time.After(b.timeout):
		return "", fmt.Errorf("escalation timed out after %v", b.timeout)
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// userFilter middleware rejects messages from unauthorized users.
func (b *Bot) userFilter(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.From == nil {
			return
		}
		slog.Info("telegram update received",
			"from_id", update.Message.From.ID,
			"allowed_id", b.allowedUserID,
			"text", update.Message.Text,
		)
		if update.Message.From.ID != b.allowedUserID {
			slog.Warn("unauthorized user blocked", "from_id", update.Message.From.ID)
			return
		}
		next(ctx, tgBot, update)
	}
}
