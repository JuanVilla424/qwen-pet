package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/JuanVilla424/qwen-pet/internal/pet"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleStart(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	emoji := b.animal.Emoji
	mood := b.petState.GetMood()
	sound := b.animal.Sounds[mood]

	msg := fmt.Sprintf(
		"%s *Hey there!* %s\n\n"+
			"I'm your personal knowledge pet. "+
			"I help your AI tools remember your preferences and decisions.\n\n"+
			"%s\n\n"+
			"Commands:\n"+
			"/status \u2014 my current state\n"+
			"/pet \u2014 about me\n"+
			"/help \u2014 all commands",
		emoji, emoji, sound,
	)

	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (b *Bot) handleStatus(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	mood := b.petState.GetMood()
	moodEmoji := b.animal.Moods[mood]
	sound := b.animal.Sounds[mood]
	stats := b.petState.Stats()

	msg := fmt.Sprintf(
		"%s %s %s\n\n```\n%s\n```\n\n%s",
		b.animal.Emoji, moodEmoji, mood.String(),
		stats,
		sound,
	)

	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (b *Bot) handlePetInfo(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	emoji := b.animal.Emoji

	msg := fmt.Sprintf(
		"%s *About me*\n\n"+
			"Type: %s %s\n"+
			"Name: configured in pet.yaml\n\n"+
			"I live in your MCP server and help Claude Code / OpenCode "+
			"remember your preferences, conventions, and past decisions.",
		emoji, b.animal.Type, emoji,
	)

	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (b *Bot) handleHelp(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	msg := fmt.Sprintf(
		"%s *Commands*\n\n"+
			"/start \u2014 welcome message\n"+
			"/status \u2014 pet mood + stats\n"+
			"/pet \u2014 pet info\n"+
			"/help \u2014 this message\n\n"+
			"_When I need your help, I'll send you a question. "+
			"Just reply with your answer!_",
		b.animal.Emoji,
	)

	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
	})
}

func (b *Bot) handleDefault(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	b.mu.Lock()
	pending := b.pending
	b.mu.Unlock()

	if pending != nil {
		slog.Info("received escalation response", "text_len", len(update.Message.Text))
		pending.ch <- update.Message.Text

		mood := pet.MoodExcited
		b.petState.SetMood(mood)
		tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("%s %s Got it! Thanks! %s", b.animal.Emoji, b.animal.Moods[mood], b.animal.Sounds[mood]),
			ParseMode: models.ParseModeMarkdown,
		})
		return
	}

	// No pending escalation — just acknowledge
	mood := b.petState.GetMood()
	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("%s %s I'm not waiting for any answer right now. Use /help to see commands.", b.animal.Emoji, b.animal.Moods[mood]),
		ParseMode: models.ParseModeMarkdown,
	})
}
