package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/JuanVilla424/qwen-pet/internal/pet"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) handleStart(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	b.mu.Lock()
	needsOb := b.needsOnboarding
	b.mu.Unlock()

	if needsOb {
		b.startOnboarding(ctx, update.Message.Chat.ID)
		return
	}

	msg := b.askAI(ctx, "A new user just started a conversation with you. Greet them warmly, introduce yourself briefly, and mention the available commands: /status, /pet, /help")

	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
	if err != nil {
		slog.Error("failed to send /start response", "error", err)
	}
}

func (b *Bot) handleStatus(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	stats := b.petState.Stats()
	prompt := fmt.Sprintf("Here are your current stats:\n%s\n\nComment on your own status with personality. Be brief.", stats)
	msg := b.askAI(ctx, prompt)

	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

func (b *Bot) handlePetInfo(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	msg := b.askAI(ctx, "The user wants to know about you. Describe yourself: your animal type, your name, your personality traits (both strengths and weaknesses), and what you do as a knowledge pet.")

	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

func (b *Bot) handleHelp(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	msg := b.askAI(ctx, "List your available commands and what each does: /start (greeting), /status (your mood and stats), /pet (about you), /help (command list). Also mention that when you need help, you'll ask a question and the user should reply directly.")

	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

func (b *Bot) handleReroll(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	b.startOnboarding(ctx, update.Message.Chat.ID)
}

func (b *Bot) handleDefault(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// Check if onboarding is active and waiting for name input
	b.mu.Lock()
	session := b.onboarding
	needsOb := b.needsOnboarding
	b.mu.Unlock()

	if needsOb && session == nil {
		b.startOnboarding(ctx, update.Message.Chat.ID)
		return
	}

	if session != nil && session.step == stepName {
		name := strings.TrimSpace(update.Message.Text)
		if name != "" {
			b.mu.Lock()
			session.name = name
			session.step = stepVirtues
			b.mu.Unlock()
			b.sendVirtuesStep(ctx, session)
			return
		}
	}

	if session != nil {
		tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please use the buttons to complete setup.",
		})
		return
	}

	// Check for pending escalation response
	b.mu.Lock()
	pending := b.pending
	b.mu.Unlock()

	if pending != nil {
		slog.Info("received escalation response", "text_len", len(update.Message.Text))
		pending.ch <- update.Message.Text

		mood := pet.MoodExcited
		b.petState.SetMood(mood)
		ack := b.askAI(ctx, "The user just answered a question you asked them. Thank them briefly and enthusiastically.")
		tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   ack,
		})
		return
	}

	// No pending escalation — process through decision engine
	if b.engine != nil {
		answer, err := b.engine.Decide(ctx, update.Message.Text, "")
		if err != nil {
			slog.Error("engine decide failed", "error", err)
			msg := b.askAI(ctx, "You tried to answer but something went wrong internally. Apologize briefly and ask the user to try again.")
			tgBot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   msg,
			})
			return
		}

		b.petState.RecordInteraction(answer.Source)
		tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   answer.Text,
		})
		interaction := fmt.Sprintf("Q: %s\nA: %s", update.Message.Text, answer.Text)
		go b.growSoul(context.Background(), interaction)
		return
	}

	// Fallback if engine not set
	msg := b.askAI(ctx, fmt.Sprintf("The user said: %s\n\nRespond naturally in character.", update.Message.Text))
	tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}
