package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/JuanVilla424/qwen-pet/internal/ai"
	"github.com/JuanVilla424/qwen-pet/internal/config"
	"github.com/JuanVilla424/qwen-pet/internal/kb"
	"github.com/JuanVilla424/qwen-pet/internal/mcp"
	"github.com/JuanVilla424/qwen-pet/internal/pet"
	"github.com/JuanVilla424/qwen-pet/internal/telegram"
)

const (
	version     = "0.2.0"
	profilePath = "data/pet-profile.json"
	statePath   = "data/pet.json"
)

func main() {
	// 1. Load config
	cfgPath := "configs/pet.yaml"
	if v := os.Getenv("PET_CONFIG"); v != "" {
		cfgPath = v
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	secrets, err := config.LoadSecrets()
	if err != nil {
		slog.Error("failed to load secrets", "error", err)
		os.Exit(1)
	}

	// 2. Init logging
	initLogger(cfg.LogLevel)

	// 3. Load or create pet profile
	profile, err := pet.LoadProfile(profilePath)
	if err != nil {
		slog.Error("failed to load profile", "error", err)
		os.Exit(1)
	}

	if profile == nil {
		slog.Info("no pet profile found, starting onboarding")
		profile, err = pet.RunOnboarding(os.Stdin, os.Stdout, profilePath)
		if err != nil {
			slog.Error("onboarding failed", "error", err)
			os.Exit(1)
		}
	}

	// 4. Get animal from catalog
	animal, err := pet.GetAnimal(profile.Type)
	if err != nil {
		slog.Error("invalid pet type in profile", "error", err)
		os.Exit(1)
	}

	// 5. Resolve traits
	traits, err := pet.ResolveTraits(profile.Traits)
	if err != nil {
		slog.Error("invalid trait in profile", "error", err)
		os.Exit(1)
	}

	// 6. Build personality
	personality := pet.BuildPersonality(animal, traits, profile.Name, profile.Soul)
	slog.Info("pet initialized",
		"name", profile.Name,
		"type", animal.Type,
		"emoji", animal.Emoji,
		"traits", personality.TraitNames(),
		"soul", len(profile.Soul),
	)

	// 7. Load pet state
	petState, err := pet.LoadState(statePath, cfg.Pet.MoodDecayHours)
	if err != nil {
		slog.Error("failed to load pet state", "error", err)
		os.Exit(1)
	}
	slog.Info("pet state loaded", "mood", petState.GetMood())

	// 8. Init embedding function
	embFunc, err := kb.NewEmbeddingFunc(cfg.Embedding)
	if err != nil {
		slog.Error("failed to create embedding function", "error", err)
		os.Exit(1)
	}

	// 9. Init KB store
	store, err := kb.NewStore(cfg.KB, embFunc)
	if err != nil {
		slog.Error("failed to init KB", "error", err)
		os.Exit(1)
	}
	slog.Info("knowledge base ready", "documents", store.Count())

	// 10. Init OpenRouter client
	aiClient := ai.NewClient(cfg.AI, secrets.OpenRouterAPIKey)

	// 11. Init Telegram bot
	tgBot, err := telegram.NewBot(
		secrets.TelegramBotToken,
		secrets.TelegramUserID,
		animal,
		petState,
		personality,
		aiClient,
		cfg.Telegram.TimeoutMinutes,
		profilePath,
	)
	if err != nil {
		slog.Error("failed to init Telegram bot", "error", err)
		os.Exit(1)
	}

	// 12. Init Decision Engine + bind to bot
	engine := ai.NewEngine(aiClient, store, personality, petState, tgBot, cfg.Decision)
	tgBot.SetEngine(engine)

	// 13. Init MCP Server
	mcpServer := mcp.NewServer(engine, store, animal, petState)

	// 14. Context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// 15. Start Telegram bot in goroutine
	go tgBot.Start(ctx)

	// 16. Run based on mode (env override > auto-detect)
	mode := os.Getenv("PET_MODE")
	if mode == "" {
		if stdinIsPipe() {
			mode = "mcp"
		} else {
			mode = "standalone"
		}
	}

	slog.Info("qwen-pet starting", "version", version, "pet", profile.Name, "mode", mode)

	switch mode {
	case "standalone":
		<-ctx.Done()
	default:
		if err := mcpServer.Run(ctx, version); err != nil {
			slog.Error("MCP server error", "error", err)
		}
	}

	// 17. Graceful shutdown
	slog.Info("shutting down, saving pet state")
	if err := petState.SaveState(); err != nil {
		slog.Error("failed to save pet state", "error", err)
	}
}

func stdinIsPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

func initLogger(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))
}
