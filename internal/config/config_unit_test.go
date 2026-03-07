package config

import (
	"os"
	"path/filepath"
	"testing"
)

const validYAML = `server:
  host: "0.0.0.0"
  port: 19999
pet:
  type: "fox"
  name: "Kit"
  mood_decay_hours: 24
  traits:
    - curious
    - sarcastic
kb:
  path: "./data/chromem"
  compress: true
  collection: "knowledge"
embeddings:
  provider: "ollama"
  ollama_model: "nomic-embed-text"
  ollama_url: "http://localhost:11434/api"
ai:
  base_url: "https://openrouter.ai/api/v1/chat/completions"
  model: "qwen/qwen3.5-flash-02-23"
  temperature: 0.7
  max_tokens: 1024
  timeout_seconds: 30
telegram:
  timeout_minutes: 5
decision:
  kb_threshold: 0.82
  ai_threshold: 0.70
log_level: "debug"
`

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeConfig(t, validYAML)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 19999 {
		t.Errorf("Server.Port = %d, want 19999", cfg.Server.Port)
	}
	if cfg.Pet.Type != "fox" {
		t.Errorf("Pet.Type = %q, want fox", cfg.Pet.Type)
	}
	if cfg.Pet.Name != "Kit" {
		t.Errorf("Pet.Name = %q, want Kit", cfg.Pet.Name)
	}
	if len(cfg.Pet.Traits) != 2 {
		t.Errorf("Pet.Traits count = %d, want 2", len(cfg.Pet.Traits))
	}
	if cfg.KB.Path != "./data/chromem" {
		t.Errorf("KB.Path = %q", cfg.KB.Path)
	}
	if cfg.Decision.KBThreshold != 0.82 {
		t.Errorf("Decision.KBThreshold = %f, want 0.82", cfg.Decision.KBThreshold)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want debug", cfg.LogLevel)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Load(nonexistent) should return error")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeConfig(t, "{{{{invalid yaml")

	_, err := Load(path)
	if err == nil {
		t.Error("Load(invalid yaml) should return error")
	}
}

func TestLoad_MissingPetType(t *testing.T) {
	yaml := `server:
  host: "0.0.0.0"
  port: 19999
pet:
  type: ""
  name: "Kit"
kb:
  path: "./data"
  collection: "kb"
ai:
  base_url: "https://api.example.com"
  model: "test"
decision:
  kb_threshold: 0.8
  ai_threshold: 0.7
`
	path := writeConfig(t, yaml)

	_, err := Load(path)
	if err == nil {
		t.Error("Load(missing pet type) should return error")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	yaml := `server:
  host: "0.0.0.0"
  port: 99999
pet:
  type: "fox"
  name: "Kit"
kb:
  path: "./data"
  collection: "kb"
ai:
  base_url: "https://api.example.com"
  model: "test"
decision:
  kb_threshold: 0.8
  ai_threshold: 0.7
`
	path := writeConfig(t, yaml)

	_, err := Load(path)
	if err == nil {
		t.Error("Load(invalid port) should return error")
	}
}

func TestLoad_InvalidThreshold(t *testing.T) {
	yaml := `server:
  host: "0.0.0.0"
  port: 19999
pet:
  type: "fox"
  name: "Kit"
kb:
  path: "./data"
  collection: "kb"
ai:
  base_url: "https://api.example.com"
  model: "test"
decision:
  kb_threshold: 1.5
  ai_threshold: 0.7
`
	path := writeConfig(t, yaml)

	_, err := Load(path)
	if err == nil {
		t.Error("Load(threshold > 1) should return error")
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	path := writeConfig(t, validYAML)

	t.Setenv("PET_HOST", "127.0.0.1")
	t.Setenv("PET_PORT", "18888")
	t.Setenv("PET_LOG_LEVEL", "error")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Server.Host = %q after env override, want 127.0.0.1", cfg.Server.Host)
	}
	if cfg.Server.Port != 18888 {
		t.Errorf("Server.Port = %d after env override, want 18888", cfg.Server.Port)
	}
	if cfg.LogLevel != "error" {
		t.Errorf("LogLevel = %q after env override, want error", cfg.LogLevel)
	}
}

func TestLoadSecrets_AllPresent(t *testing.T) {
	t.Setenv("OPENROUTER_API_KEY", "sk-test-123")
	t.Setenv("TELEGRAM_BOT_TOKEN", "bot-token-456")
	t.Setenv("TELEGRAM_USER_ID", "12345678")

	secrets, err := LoadSecrets()
	if err != nil {
		t.Fatalf("LoadSecrets() error: %v", err)
	}

	if secrets.OpenRouterAPIKey != "sk-test-123" {
		t.Errorf("OpenRouterAPIKey = %q", secrets.OpenRouterAPIKey)
	}
	if secrets.TelegramBotToken != "bot-token-456" {
		t.Errorf("TelegramBotToken = %q", secrets.TelegramBotToken)
	}
	if secrets.TelegramUserID != 12345678 {
		t.Errorf("TelegramUserID = %d, want 12345678", secrets.TelegramUserID)
	}
}

func TestLoadSecrets_MissingAPIKey(t *testing.T) {
	t.Setenv("OPENROUTER_API_KEY", "")
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("TELEGRAM_USER_ID", "123")

	_, err := LoadSecrets()
	if err == nil {
		t.Error("LoadSecrets(missing API key) should return error")
	}
}

func TestLoadSecrets_MissingBotToken(t *testing.T) {
	t.Setenv("OPENROUTER_API_KEY", "key")
	t.Setenv("TELEGRAM_BOT_TOKEN", "")
	t.Setenv("TELEGRAM_USER_ID", "123")

	_, err := LoadSecrets()
	if err == nil {
		t.Error("LoadSecrets(missing bot token) should return error")
	}
}

func TestLoadSecrets_InvalidUserID(t *testing.T) {
	t.Setenv("OPENROUTER_API_KEY", "key")
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("TELEGRAM_USER_ID", "not-a-number")

	_, err := LoadSecrets()
	if err == nil {
		t.Error("LoadSecrets(invalid user ID) should return error")
	}
}
