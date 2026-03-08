package pet

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const (
	minVirtues = 2
	maxVirtues = 4
	minDefects = 1
	maxDefects = 2
)

// RunOnboarding runs the interactive CLI onboarding flow.
// It reads from r (typically os.Stdin) and writes to w (typically os.Stdout).
// Returns the created Profile and saves it to profilePath.
func RunOnboarding(r io.Reader, w io.Writer, profilePath string) (*Profile, error) {
	scanner := bufio.NewScanner(r)

	fmt.Fprint(w, "\n  Welcome to Qwen PET!\n")
	fmt.Fprint(w, "  Your personal knowledge tamagotchi.\n\n")

	// Step 1: Choose animal
	animalType, err := chooseAnimal(scanner, w)
	if err != nil {
		return nil, fmt.Errorf("choose animal: %w", err)
	}

	animal, _ := GetAnimal(animalType)

	// Step 2: Choose name
	name, err := chooseName(scanner, w, animal)
	if err != nil {
		return nil, fmt.Errorf("choose name: %w", err)
	}

	// Step 3: Choose virtues
	virtueIDs, err := chooseTraits(scanner, w, ListVirtues(), "virtues", minVirtues, maxVirtues)
	if err != nil {
		return nil, fmt.Errorf("choose virtues: %w", err)
	}

	// Step 4: Choose defects
	defectIDs, err := chooseTraits(scanner, w, ListDefects(), "defects", minDefects, maxDefects)
	if err != nil {
		return nil, fmt.Errorf("choose defects: %w", err)
	}

	allTraits := append(virtueIDs, defectIDs...)

	profile := &Profile{
		Type:   animalType,
		Name:   name,
		Traits: allTraits,
	}

	if err := profile.SaveProfile(profilePath); err != nil {
		return nil, fmt.Errorf("save profile: %w", err)
	}

	// Summary
	resolved, _ := ResolveTraits(allTraits)
	names := make([]string, len(resolved))
	for i, t := range resolved {
		names[i] = t.Name
	}

	fmt.Fprintf(w, "\n  %s %s the %s is ready!\n", animal.Emoji, name, animal.Type)
	fmt.Fprintf(w, "  Traits: %s\n", strings.Join(names, ", "))
	fmt.Fprintf(w, "  Profile saved to %s\n\n", profilePath)

	return profile, nil
}

func chooseAnimal(scanner *bufio.Scanner, w io.Writer) (string, error) {
	animals := ListAnimals()
	sort.Slice(animals, func(i, j int) bool {
		return animals[i].Type < animals[j].Type
	})

	fmt.Fprint(w, "  Choose your pet:\n\n")
	for i, a := range animals {
		fmt.Fprintf(w, "  %2d. %s %-10s  %s\n", i+1, a.Emoji, a.Type, a.DefaultName)
	}

	for {
		fmt.Fprint(w, "\n  > ")
		if !scanner.Scan() {
			return "", fmt.Errorf("unexpected end of input")
		}

		input := strings.TrimSpace(scanner.Text())
		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(animals) {
			fmt.Fprintf(w, "  Pick a number between 1 and %d\n", len(animals))
			continue
		}

		return animals[idx-1].Type, nil
	}
}

func chooseName(scanner *bufio.Scanner, w io.Writer, animal *Animal) (string, error) {
	fmt.Fprintf(w, "\n  Name your %s (default: %s):\n", animal.Type, animal.DefaultName)
	fmt.Fprint(w, "  > ")

	if !scanner.Scan() {
		return "", fmt.Errorf("unexpected end of input")
	}

	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		name = animal.DefaultName
	}

	return name, nil
}

func chooseTraits(scanner *bufio.Scanner, w io.Writer, available []Trait, label string, min, max int) ([]string, error) {
	sort.Slice(available, func(i, j int) bool {
		return available[i].ID < available[j].ID
	})

	fmt.Fprintf(w, "\n  Choose %s (%d-%d, comma-separated numbers):\n\n", label, min, max)
	for i, t := range available {
		fmt.Fprintf(w, "  %2d. %-15s %s\n", i+1, t.ID, t.Description)
	}

	for {
		fmt.Fprint(w, "\n  > ")
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected end of input")
		}

		input := strings.TrimSpace(scanner.Text())
		parts := strings.Split(input, ",")

		var selected []string
		valid := true
		seen := make(map[int]bool)

		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}

			idx, err := strconv.Atoi(p)
			if err != nil || idx < 1 || idx > len(available) {
				fmt.Fprintf(w, "  Invalid number: %s (pick 1-%d)\n", p, len(available))
				valid = false
				break
			}

			if seen[idx] {
				continue
			}
			seen[idx] = true
			selected = append(selected, available[idx-1].ID)
		}

		if !valid {
			continue
		}

		if len(selected) < min || len(selected) > max {
			fmt.Fprintf(w, "  Pick %d to %d %s\n", min, max, label)
			continue
		}

		return selected, nil
	}
}
