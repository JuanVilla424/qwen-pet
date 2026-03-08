package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration.
type Config struct {
	Server    ServerCfg    `yaml:"server"`
	Pet       PetCfg       `yaml:"pet"`
	KB        KBCfg        `yaml:"kb"`
	Embedding EmbeddingCfg `yaml:"embeddings"`
	AI        AICfg        `yaml:"ai"`
	Telegram  TelegramCfg  `yaml:"telegram"`
	Decision  DecisionCfg  `yaml:"decision"`
	LogLevel  string       `yaml:"log_level"`
}

type ServerCfg struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type PetCfg struct {
	Type           string   `yaml:"type"`
	Name           string   `yaml:"name"`
	MoodDecayHours int      `yaml:"mood_decay_hours"`
	Traits         []string `yaml:"traits"`
}

type KBCfg struct {
	Path       string `yaml:"path"`
	Compress   bool   `yaml:"compress"`
	Collection string `yaml:"collection"`
}

type EmbeddingCfg struct {
	Provider    string `yaml:"provider"`
	OllamaModel string `yaml:"ollama_model"`
	OllamaURL   string `yaml:"ollama_url"`
}

type AICfg struct {
	BaseURL        string  `yaml:"base_url"`
	Model          string  `yaml:"model"`
	Temperature    float64 `yaml:"temperature"`
	MaxTokens      int     `yaml:"max_tokens"`
	TimeoutSeconds int     `yaml:"timeout_seconds"`
}

type TelegramCfg struct {
	TimeoutMinutes int `yaml:"timeout_minutes"`
}

type DecisionCfg struct {
	KBThreshold float32 `yaml:"kb_threshold"`
	AIThreshold float32 `yaml:"ai_threshold"`
}

// Secrets holds sensitive configuration from environment variables.
type Secrets struct {
	OpenRouterAPIKey string
	TelegramBotToken string
	TelegramUserID   int64
}

// Load reads configuration from a YAML file and applies environment variable overrides.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	applyEnvOverrides(&cfg)

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	return &cfg, nil
}

// LoadSecrets reads sensitive values from environment variables.
func LoadSecrets() (*Secrets, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY") // optional for local backends like Ollama

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	userIDStr := os.Getenv("TELEGRAM_USER_ID")
	if userIDStr == "" {
		return nil, fmt.Errorf("TELEGRAM_USER_ID is required")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("TELEGRAM_USER_ID must be a valid integer: %w", err)
	}

	return &Secrets{
		OpenRouterAPIKey: apiKey,
		TelegramBotToken: botToken,
		TelegramUserID:   userID,
	}, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("PET_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("PET_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("PET_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
}

func (c *Config) validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535, got %d", c.Server.Port)
	}
	if c.KB.Path == "" {
		return fmt.Errorf("kb.path is required")
	}
	if c.KB.Collection == "" {
		return fmt.Errorf("kb.collection is required")
	}
	if c.AI.BaseURL == "" {
		return fmt.Errorf("ai.base_url is required")
	}
	if c.AI.Model == "" {
		return fmt.Errorf("ai.model is required")
	}
	if c.Decision.KBThreshold <= 0 || c.Decision.KBThreshold > 1 {
		return fmt.Errorf("decision.kb_threshold must be between 0 and 1")
	}
	if c.Decision.AIThreshold <= 0 || c.Decision.AIThreshold > 1 {
		return fmt.Errorf("decision.ai_threshold must be between 0 and 1")
	}
	return nil
}
