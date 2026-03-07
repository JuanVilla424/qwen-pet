package pet

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	imageCacheDir  = "data/images"
	twemojiBaseURL = "https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/72x72"
	fetchTimeout   = 10 * time.Second
)

// emoji codepoint mappings for Twemoji CDN
var animalEmojiCodes = map[string]string{
	"cat":     "1f431",
	"dog":     "1f436",
	"fox":     "1f98a",
	"owl":     "1f989",
	"dragon":  "1f409",
	"rabbit":  "1f430",
	"penguin": "1f427",
	"snake":   "1f40d",
	"panda":   "1f43c",
	"unicorn": "1f984",
	"frog":    "1f438",
	"lion":    "1f981",
	"wolf":    "1f43a",
	"bear":    "1f43b",
	"eagle":   "1f985",
	"dolphin": "1f42c",
	"lizard":  "1f98e",
	"octopus": "1f419",
}

// FetchImage retrieves a pet image, using local cache first.
func FetchImage(animalType string) ([]byte, error) {
	cachePath := filepath.Join(imageCacheDir, animalType+".png")

	if data, err := os.ReadFile(cachePath); err == nil {
		return data, nil
	}

	code, ok := animalEmojiCodes[animalType]
	if !ok {
		return nil, fmt.Errorf("no emoji code for animal type: %q", animalType)
	}

	url := fmt.Sprintf("%s/%s.png", twemojiBaseURL, code)

	client := &http.Client{Timeout: fetchTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch image: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read image body: %w", err)
	}

	if err := os.MkdirAll(imageCacheDir, 0o755); err != nil {
		return data, nil // return data even if cache fails
	}
	_ = os.WriteFile(cachePath, data, 0o644)

	return data, nil
}

// ImagePath returns the cached image path if it exists.
func ImagePath(animalType string) (string, bool) {
	cachePath := filepath.Join(imageCacheDir, animalType+".png")
	if _, err := os.Stat(cachePath); err == nil {
		abs, _ := filepath.Abs(cachePath)
		return abs, true
	}
	return "", false
}
