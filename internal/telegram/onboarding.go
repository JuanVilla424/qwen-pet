package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/JuanVilla424/qwen-pet/internal/pet"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type onboardingStep int

const (
	stepAnimal onboardingStep = iota
	stepName
	stepVirtues
	stepDefects
	stepConfirm
)

const (
	obMinVirtues = 2
	obMaxVirtues = 4
	obMinDefects = 1
	obMaxDefects = 2
)

type onboardingSession struct {
	step            onboardingStep
	chatID          int64
	animalType      string
	name            string
	selectedVirtues map[string]bool
	selectedDefects map[string]bool
	messageID       int
}

func newOnboardingSession(chatID int64) *onboardingSession {
	return &onboardingSession{
		step:            stepAnimal,
		chatID:          chatID,
		selectedVirtues: make(map[string]bool),
		selectedDefects: make(map[string]bool),
	}
}

// sendAnimalStep shows a grid of animals as inline keyboard buttons.
func (b *Bot) sendAnimalStep(ctx context.Context, session *onboardingSession) {
	animals := pet.ListAnimals()
	sort.Slice(animals, func(i, j int) bool {
		return animals[i].Type < animals[j].Type
	})

	var rows [][]models.InlineKeyboardButton
	var row []models.InlineKeyboardButton
	for _, a := range animals {
		label := fmt.Sprintf("%s %s", a.Emoji, a.Type)
		cbData := "ob_animal_" + a.Type
		row = append(row, models.InlineKeyboardButton{Text: label, CallbackData: cbData})
		if len(row) == 3 {
			rows = append(rows, row)
			row = nil
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}

	text := "Welcome to Qwen PET!\nChoose your pet:"
	msg, err := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: session.chatID,
		Text:   text,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: rows,
		},
	})
	if err != nil {
		slog.Error("failed to send animal step", "error", err)
		return
	}
	session.messageID = msg.ID
}

// sendNameStep asks for a name.
func (b *Bot) sendNameStep(ctx context.Context, session *onboardingSession) {
	animal, _ := pet.GetAnimal(session.animalType)
	defaultName := "Pet"
	if animal != nil {
		defaultName = animal.DefaultName
	}

	text := fmt.Sprintf("Name your %s (reply with a name, or tap the button for default):", session.animalType)
	cbData := "ob_name_default"
	label := fmt.Sprintf("Use default: %s", defaultName)

	msg, err := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: session.chatID,
		Text:   text,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: label, CallbackData: cbData}},
			},
		},
	})
	if err != nil {
		slog.Error("failed to send name step", "error", err)
		return
	}
	session.messageID = msg.ID
}

// sendVirtuesStep shows toggleable virtue buttons.
func (b *Bot) sendVirtuesStep(ctx context.Context, session *onboardingSession) {
	b.sendTraitStep(ctx, session, false)
}

// sendDefectsStep shows toggleable defect buttons.
func (b *Bot) sendDefectsStep(ctx context.Context, session *onboardingSession) {
	b.sendTraitStep(ctx, session, true)
}

func (b *Bot) sendTraitStep(ctx context.Context, session *onboardingSession, isDefect bool) {
	var traits []pet.Trait
	var selected map[string]bool
	var label, prefix, doneCallback string
	var minCount, maxCount int

	if isDefect {
		traits = pet.ListDefects()
		selected = session.selectedDefects
		label = "defects"
		prefix = "ob_defect_"
		doneCallback = "ob_defects_done"
		minCount = obMinDefects
		maxCount = obMaxDefects
	} else {
		traits = pet.ListVirtues()
		selected = session.selectedVirtues
		label = "virtues"
		prefix = "ob_virtue_"
		doneCallback = "ob_virtues_done"
		minCount = obMinVirtues
		maxCount = obMaxVirtues
	}

	sort.Slice(traits, func(i, j int) bool {
		return traits[i].ID < traits[j].ID
	})

	var rows [][]models.InlineKeyboardButton
	var row []models.InlineKeyboardButton
	for _, t := range traits {
		check := "  "
		if selected[t.ID] {
			check = "✓ "
		}
		btnLabel := check + t.Name
		cbData := prefix + t.ID
		row = append(row, models.InlineKeyboardButton{Text: btnLabel, CallbackData: cbData})
		if len(row) == 2 {
			rows = append(rows, row)
			row = nil
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}

	count := len(selected)
	confirmLabel := fmt.Sprintf("Confirm (%d selected)", count)
	rows = append(rows, []models.InlineKeyboardButton{
		{Text: confirmLabel, CallbackData: doneCallback},
	})

	text := fmt.Sprintf("Choose %d-%d %s (tap to toggle):", minCount, maxCount, label)
	keyboard := &models.InlineKeyboardMarkup{InlineKeyboard: rows}

	if session.messageID != 0 {
		_, err := b.tgBot.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      session.chatID,
			MessageID:   session.messageID,
			Text:        text,
			ReplyMarkup: keyboard,
		})
		if err != nil {
			slog.Debug("edit failed, sending new message", "error", err)
			msg, _ := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      session.chatID,
				Text:        text,
				ReplyMarkup: keyboard,
			})
			if msg != nil {
				session.messageID = msg.ID
			}
		}
	} else {
		msg, _ := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      session.chatID,
			Text:        text,
			ReplyMarkup: keyboard,
		})
		if msg != nil {
			session.messageID = msg.ID
		}
	}
}

// sendConfirmStep shows the final summary.
func (b *Bot) sendConfirmStep(ctx context.Context, session *onboardingSession) {
	animal, _ := pet.GetAnimal(session.animalType)
	emoji := ""
	if animal != nil {
		emoji = animal.Emoji
	}

	var virtueNames, defectNames []string
	for id := range session.selectedVirtues {
		t, _ := pet.GetTrait(id)
		if t != nil {
			virtueNames = append(virtueNames, t.Name)
		}
	}
	for id := range session.selectedDefects {
		t, _ := pet.GetTrait(id)
		if t != nil {
			defectNames = append(defectNames, t.Name)
		}
	}
	sort.Strings(virtueNames)
	sort.Strings(defectNames)

	text := fmt.Sprintf("%s %s the %s\n\nVirtues: %s\nDefects: %s\n\nCreate this pet?",
		emoji, session.name, session.animalType,
		strings.Join(virtueNames, ", "),
		strings.Join(defectNames, ", "),
	)

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Confirm", CallbackData: "ob_confirm_yes"},
				{Text: "Start over", CallbackData: "ob_confirm_restart"},
			},
		},
	}

	if session.messageID != 0 {
		b.tgBot.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      session.chatID,
			MessageID:   session.messageID,
			Text:        text,
			ReplyMarkup: keyboard,
		})
	} else {
		msg, _ := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      session.chatID,
			Text:        text,
			ReplyMarkup: keyboard,
		})
		if msg != nil {
			session.messageID = msg.ID
		}
	}
}

// handleOnboardingCallback dispatches inline keyboard callbacks during onboarding.
func (b *Bot) handleOnboardingCallback(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	cb := update.CallbackQuery
	data := cb.Data

	b.mu.Lock()
	session := b.onboarding
	b.mu.Unlock()

	if session == nil {
		tgBot.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
			Text:            "No onboarding in progress.",
		})
		return
	}

	// Answer callback to stop loading indicator
	tgBot.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: cb.ID})

	switch {
	case strings.HasPrefix(data, "ob_animal_"):
		animalType := strings.TrimPrefix(data, "ob_animal_")
		if !pet.AnimalExists(animalType) {
			return
		}
		b.mu.Lock()
		session.animalType = animalType
		session.step = stepName
		b.mu.Unlock()
		b.sendNameStep(ctx, session)

	case data == "ob_name_default":
		animal, _ := pet.GetAnimal(session.animalType)
		defaultName := "Pet"
		if animal != nil {
			defaultName = animal.DefaultName
		}
		b.mu.Lock()
		session.name = defaultName
		session.step = stepVirtues
		b.mu.Unlock()
		b.sendVirtuesStep(ctx, session)

	case strings.HasPrefix(data, "ob_virtue_"):
		traitID := strings.TrimPrefix(data, "ob_virtue_")
		b.mu.Lock()
		if session.selectedVirtues[traitID] {
			delete(session.selectedVirtues, traitID)
		} else if len(session.selectedVirtues) < obMaxVirtues {
			session.selectedVirtues[traitID] = true
		}
		b.mu.Unlock()
		b.sendVirtuesStep(ctx, session)

	case data == "ob_virtues_done":
		b.mu.Lock()
		count := len(session.selectedVirtues)
		b.mu.Unlock()
		if count < obMinVirtues || count > obMaxVirtues {
			tgBot.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: cb.ID,
				Text:            fmt.Sprintf("Select %d-%d virtues", obMinVirtues, obMaxVirtues),
				ShowAlert:       true,
			})
			return
		}
		b.mu.Lock()
		session.step = stepDefects
		b.mu.Unlock()
		b.sendDefectsStep(ctx, session)

	case strings.HasPrefix(data, "ob_defect_"):
		traitID := strings.TrimPrefix(data, "ob_defect_")
		b.mu.Lock()
		if session.selectedDefects[traitID] {
			delete(session.selectedDefects, traitID)
		} else if len(session.selectedDefects) < obMaxDefects {
			session.selectedDefects[traitID] = true
		}
		b.mu.Unlock()
		b.sendDefectsStep(ctx, session)

	case data == "ob_defects_done":
		b.mu.Lock()
		count := len(session.selectedDefects)
		b.mu.Unlock()
		if count < obMinDefects || count > obMaxDefects {
			tgBot.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: cb.ID,
				Text:            fmt.Sprintf("Select %d-%d defects", obMinDefects, obMaxDefects),
				ShowAlert:       true,
			})
			return
		}
		b.mu.Lock()
		session.step = stepConfirm
		b.mu.Unlock()
		b.sendConfirmStep(ctx, session)

	case data == "ob_confirm_yes":
		b.finalizeOnboarding(ctx, session)

	case data == "ob_confirm_restart":
		b.mu.Lock()
		b.onboarding = newOnboardingSession(session.chatID)
		session = b.onboarding
		b.mu.Unlock()
		b.sendAnimalStep(ctx, session)
	}
}

// finalizeOnboarding creates the profile and initializes the bot.
func (b *Bot) finalizeOnboarding(ctx context.Context, session *onboardingSession) {
	var traitIDs []string
	for id := range session.selectedVirtues {
		traitIDs = append(traitIDs, id)
	}
	for id := range session.selectedDefects {
		traitIDs = append(traitIDs, id)
	}
	sort.Strings(traitIDs)

	profile := &pet.Profile{
		Type:   session.animalType,
		Name:   session.name,
		Traits: traitIDs,
	}

	if err := profile.SaveProfile(b.profilePath); err != nil {
		slog.Error("failed to save profile during onboarding", "error", err)
		b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: session.chatID,
			Text:   "Error saving profile. Try /reroll.",
		})
		return
	}

	// Initialize the bot with the new profile
	b.completeSetup(profile)

	b.mu.Lock()
	b.onboarding = nil
	b.needsOnboarding = false
	b.mu.Unlock()

	slog.Info("onboarding completed", "name", profile.Name, "type", profile.Type)

	// Send welcome message
	msg := b.askAI(ctx, "You were just created! Greet your owner warmly, introduce yourself, and mention /status, /pet, /help, /reroll commands.")
	b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: session.chatID,
		Text:   msg,
	})
}
