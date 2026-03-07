package kb

import (
	"fmt"

	"github.com/JuanVilla424/qwen-pet/internal/config"
	chromem "github.com/philippgille/chromem-go"
)

// NewEmbeddingFunc creates an embedding function based on configuration.
func NewEmbeddingFunc(cfg config.EmbeddingCfg) (chromem.EmbeddingFunc, error) {
	switch cfg.Provider {
	case "ollama":
		return chromem.NewEmbeddingFuncOllama(cfg.OllamaModel, cfg.OllamaURL), nil
	default:
		return nil, fmt.Errorf("unsupported embedding provider: %q (supported: ollama)", cfg.Provider)
	}
}
