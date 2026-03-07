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

const version = "0.1.2"

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

	// 3. Get animal from catalog
	animal, err := pet.GetAnimal(cfg.Pet.Type)
	if err != nil {
		slog.Error("invalid pet type", "error", err)
		os.Exit(1)
	}

	// 4. Resolve traits
	traits, err := pet.ResolveTraits(cfg.Pet.Traits)
	if err != nil {
		slog.Error("invalid trait", "error", err)
		os.Exit(1)
	}

	// 5. Build personality
	personality := pet.BuildPersonality(animal, traits, cfg.Pet.Name)
	slog.Info("pet initialized",
		"name", cfg.Pet.Name,
		"type", animal.Type,
		"emoji", animal.Emoji,
		"traits", personality.TraitNames(),
	)

	// 6. Load pet state
	petState, err := pet.LoadState("data/pet.json", cfg.Pet.MoodDecayHours)
	if err != nil {
		slog.Error("failed to load pet state", "error", err)
		os.Exit(1)
	}
	slog.Info("pet state loaded", "mood", petState.GetMood())

	// 7. Init embedding function
	embFunc, err := kb.NewEmbeddingFunc(cfg.Embedding)
	if err != nil {
		slog.Error("failed to create embedding function", "error", err)
		os.Exit(1)
	}

	// 8. Init KB store
	store, err := kb.NewStore(cfg.KB, embFunc)
	if err != nil {
		slog.Error("failed to init KB", "error", err)
		os.Exit(1)
	}
	slog.Info("knowledge base ready", "documents", store.Count())

	// 9. Init OpenRouter client
	aiClient := ai.NewClient(cfg.AI, secrets.OpenRouterAPIKey)

	// 10. Init Telegram bot
	tgBot, err := telegram.NewBot(
		secrets.TelegramBotToken,
		secrets.TelegramUserID,
		animal,
		petState,
		cfg.Telegram.TimeoutMinutes,
	)
	if err != nil {
		slog.Error("failed to init Telegram bot", "error", err)
		os.Exit(1)
	}

	// 11. Init Decision Engine
	engine := ai.NewEngine(aiClient, store, personality, petState, tgBot, cfg.Decision)

	// 12. Init MCP Server
	mcpServer := mcp.NewServer(engine, store, animal, petState)

	// 13. Context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// 14. Start Telegram bot in goroutine
	go tgBot.Start(ctx)

	// 15. Run MCP server (blocks on stdio)
	slog.Info("qwen-pet starting", "version", version, "pet", cfg.Pet.Name)
	if err := mcpServer.Run(ctx, version); err != nil {
		slog.Error("MCP server error", "error", err)
	}

	// 16. Graceful shutdown
	slog.Info("shutting down, saving pet state")
	if err := petState.SaveState(); err != nil {
		slog.Error("failed to save pet state", "error", err)
	}
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
