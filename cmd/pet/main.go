package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/JuanVilla424/qwen-pet/internal/ai"
	"github.com/JuanVilla424/qwen-pet/internal/config"
	"github.com/JuanVilla424/qwen-pet/internal/kb"
	"github.com/JuanVilla424/qwen-pet/internal/mcp"
	"github.com/JuanVilla424/qwen-pet/internal/pet"
	"github.com/JuanVilla424/qwen-pet/internal/telegram"
)

const version = "0.3.0"

func main() {
	// 1. Resolve paths
	dataDir := resolveDataDir()
	cfgPath := resolveConfigPath()
	profilePath := filepath.Join(dataDir, "pet-profile.json")
	statePath := filepath.Join(dataDir, "pet.json")

	// 2. Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		slog.Error("failed to create data directory", "path", dataDir, "error", err)
		os.Exit(1)
	}

	// 3. Load config
	cfg, err := config.Load(cfgPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Resolve KB path relative to data dir if not absolute
	if !filepath.IsAbs(cfg.KB.Path) {
		cfg.KB.Path = filepath.Join(dataDir, cfg.KB.Path)
	}

	secrets, err := config.LoadSecrets()
	if err != nil {
		slog.Error("failed to load secrets", "error", err)
		os.Exit(1)
	}

	// 4. Init logging
	initLogger(cfg.LogLevel)
	slog.Info("paths resolved", "data_dir", dataDir, "config", cfgPath)

	// 5. Load or create pet profile
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

	// 6. Get animal from catalog
	animal, err := pet.GetAnimal(profile.Type)
	if err != nil {
		slog.Error("invalid pet type in profile", "error", err)
		os.Exit(1)
	}

	// 7. Resolve traits
	traits, err := pet.ResolveTraits(profile.Traits)
	if err != nil {
		slog.Error("invalid trait in profile", "error", err)
		os.Exit(1)
	}

	// 8. Build personality
	personality := pet.BuildPersonality(animal, traits, profile.Name, profile.Soul)
	slog.Info("pet initialized",
		"name", profile.Name,
		"type", animal.Type,
		"emoji", animal.Emoji,
		"traits", personality.TraitNames(),
		"soul", len(profile.Soul),
	)

	// 9. Load pet state
	petState, err := pet.LoadState(statePath, cfg.Pet.MoodDecayHours)
	if err != nil {
		slog.Error("failed to load pet state", "error", err)
		os.Exit(1)
	}
	slog.Info("pet state loaded", "mood", petState.GetMood())

	// 10. Init embedding function
	embFunc, err := kb.NewEmbeddingFunc(cfg.Embedding)
	if err != nil {
		slog.Error("failed to create embedding function", "error", err)
		os.Exit(1)
	}

	// 11. Init KB store
	store, err := kb.NewStore(cfg.KB, embFunc)
	if err != nil {
		slog.Error("failed to init KB", "error", err)
		os.Exit(1)
	}
	slog.Info("knowledge base ready", "documents", store.Count())

	// 12. Init OpenRouter client
	aiClient := ai.NewClient(cfg.AI, secrets.OpenRouterAPIKey)

	// 13. Init Telegram bot
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

	// 14. Init Decision Engine + bind to bot
	engine := ai.NewEngine(aiClient, store, personality, petState, tgBot, cfg.Decision)
	tgBot.SetEngine(engine)

	// 15. Init MCP Server
	mcpServer := mcp.NewServer(engine, store, animal, petState)

	// 16. Context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// 17. Run based on mode (env override > auto-detect)
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
		go tgBot.Start(ctx)
		<-ctx.Done()
	default:
		if err := mcpServer.Run(ctx, version); err != nil {
			slog.Error("MCP server error", "error", err)
		}
	}

	// 18. Graceful shutdown
	slog.Info("shutting down, saving pet state")
	if err := petState.SaveState(); err != nil {
		slog.Error("failed to save pet state", "error", err)
	}
}

func resolveDataDir() string {
	if v := os.Getenv("PET_DATA_DIR"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err == nil {
		xdg := filepath.Join(home, ".local", "share", "qwen-pet")
		if _, err := os.Stat(xdg); err == nil {
			return xdg
		}
	}
	return "data"
}

func resolveConfigPath() string {
	if v := os.Getenv("PET_CONFIG"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err == nil {
		xdg := filepath.Join(home, ".config", "qwen-pet", "pet.yaml")
		if _, err := os.Stat(xdg); err == nil {
			return xdg
		}
	}
	return "configs/pet.yaml"
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
