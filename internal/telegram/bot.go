package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

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
	timeout       time.Duration

	mu      sync.Mutex
	pending *pendingRequest
}

// NewBot creates and configures a Telegram bot.
func NewBot(token string, userID int64, animal *pet.Animal, petState *pet.State, timeoutMinutes int) (*Bot, error) {
	b := &Bot{
		allowedUserID: userID,
		animal:        animal,
		petState:      petState,
		timeout:       time.Duration(timeoutMinutes) * time.Minute,
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

	return b, nil
}

// Start begins polling for updates (runs in goroutine).
func (b *Bot) Start(ctx context.Context) {
	slog.Info("starting Telegram bot polling")
	b.tgBot.Start(ctx)
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

	// Build message with pet personality
	emoji := b.animal.Emoji
	moodEmoji := b.animal.Moods[pet.MoodThinking]
	msg := fmt.Sprintf("%s %s *needs your help!*\n\n", emoji, moodEmoji)
	msg += fmt.Sprintf("*Question:* %s", question)
	if kbContext != "" {
		msg += fmt.Sprintf("\n\n*Context:*\n%s", kbContext)
	}
	msg += "\n\n_Reply with your answer..._"

	_, err := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    b.allowedUserID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
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
		if update.Message.From.ID != b.allowedUserID {
			return
		}
		next(ctx, tgBot, update)
	}
}
