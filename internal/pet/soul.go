package pet

import "strings"

const SoulCapacity = 30

// SoulGrowthPrompt asks the AI to pick ONE new soul word based on the interaction.
func SoulGrowthPrompt(currentSoul []string, lastInteraction string) string {
	var sb strings.Builder
	sb.WriteString("You just had this interaction:\n")
	sb.WriteString(lastInteraction)
	sb.WriteString("\n\n")

	if len(currentSoul) > 0 {
		sb.WriteString("Your current soul words: ")
		sb.WriteString(strings.Join(currentSoul, " "))
		sb.WriteString("\n\n")
	}

	sb.WriteString("Choose ONE new word to add to your soul. ")
	sb.WriteString("This word captures something you felt, learned, or became during this interaction. ")
	sb.WriteString("No repeats. Single word only. Raw, honest, surprising. ")
	sb.WriteString("Reply with ONLY the word, nothing else.")

	return sb.String()
}

// ParseSoulWord extracts a single soul word from AI response.
func ParseSoulWord(response string) string {
	response = strings.TrimSpace(response)
	// Take only the first word in case AI adds extra
	fields := strings.Fields(response)
	if len(fields) == 0 {
		return ""
	}
	word := strings.ToLower(fields[0])
	// Strip punctuation
	word = strings.Trim(word, ".,;:!?\"'()-")
	if len(word) < 2 {
		return ""
	}
	return word
}

// SoulContains checks if a word already exists in the soul.
func SoulContains(soul []string, word string) bool {
	for _, w := range soul {
		if w == word {
			return true
		}
	}
	return false
}

// FormatSoul returns the soul as a space-separated string for prompt injection.
func FormatSoul(soul []string) string {
	return strings.Join(soul, " ")
}
