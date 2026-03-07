package pet

import "fmt"

// Mood represents the pet's emotional state.
type Mood int

const (
	MoodHappy Mood = iota
	MoodNeutral
	MoodThinking
	MoodTired
	MoodExcited
	MoodSad
)

func (m Mood) String() string {
	switch m {
	case MoodHappy:
		return "happy"
	case MoodNeutral:
		return "neutral"
	case MoodThinking:
		return "thinking"
	case MoodTired:
		return "tired"
	case MoodExcited:
		return "excited"
	case MoodSad:
		return "sad"
	default:
		return "neutral"
	}
}

// Animal represents a pet's visual identity (NO personality — that comes from traits).
type Animal struct {
	Type        string
	Emoji       string
	DefaultName string
	Moods       map[Mood]string
	Sounds      map[Mood]string
}

// Trait represents a personality trait (virtue or defect).
type Trait struct {
	ID          string
	Name        string
	Description string
	IsDefect    bool
	PromptHint  string
}

var animalCatalog = map[string]Animal{
	"cat": {
		Type: "cat", Emoji: "\U0001F431", DefaultName: "Michi",
		Moods:  map[Mood]string{MoodHappy: "\U0001F63A", MoodNeutral: "\U0001F63C", MoodThinking: "\U0001F640", MoodTired: "\U0001F63F", MoodExcited: "\U0001F638", MoodSad: "\U0001F63E"},
		Sounds: map[Mood]string{MoodHappy: "*purr...*", MoodNeutral: "*meow*", MoodThinking: "*mrrow?*", MoodTired: "*yawn...*", MoodExcited: "*MRRROW!*", MoodSad: "*hiss...*"},
	},
	"dog": {
		Type: "dog", Emoji: "\U0001F436", DefaultName: "Rex",
		Moods:  map[Mood]string{MoodHappy: "\U0001F415", MoodNeutral: "\U0001F436", MoodThinking: "\U0001F9AE", MoodTired: "\U0001F634", MoodExcited: "\U0001F929", MoodSad: "\U0001F97A"},
		Sounds: map[Mood]string{MoodHappy: "*woof woof!*", MoodNeutral: "*bark*", MoodThinking: "*whine?*", MoodTired: "*zzz...*", MoodExcited: "*WOOF WOOF WOOF!*", MoodSad: "*whimper...*"},
	},
	"fox": {
		Type: "fox", Emoji: "\U0001F98A", DefaultName: "Kit",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60F", MoodNeutral: "\U0001F98A", MoodThinking: "\U0001F914", MoodTired: "\U0001F971", MoodExcited: "\U0001F525", MoodSad: "\U0001F614"},
		Sounds: map[Mood]string{MoodHappy: "*yip yip!*", MoodNeutral: "*sniff*", MoodThinking: "*hmm...*", MoodTired: "*yawn*", MoodExcited: "*YIIIP!*", MoodSad: "*whine...*"},
	},
	"owl": {
		Type: "owl", Emoji: "\U0001F989", DefaultName: "Archie",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60C", MoodNeutral: "\U0001F989", MoodThinking: "\U0001F9D0", MoodTired: "\U0001F634", MoodExcited: "\U0001F31F", MoodSad: "\U0001F61E"},
		Sounds: map[Mood]string{MoodHappy: "*hoo hoo!*", MoodNeutral: "*hoot*", MoodThinking: "*hoooo...*", MoodTired: "*zzz*", MoodExcited: "*HOO HOO HOO!*", MoodSad: "*soft hoot...*"},
	},
	"dragon": {
		Type: "dragon", Emoji: "\U0001F409", DefaultName: "Drakar",
		Moods:  map[Mood]string{MoodHappy: "\U0001F525", MoodNeutral: "\U0001F409", MoodThinking: "\U0001F4AD", MoodTired: "\U0001F32C\uFE0F", MoodExcited: "\U0001F4A5", MoodSad: "\U0001F327\uFE0F"},
		Sounds: map[Mood]string{MoodHappy: "*rumble...*", MoodNeutral: "*growl*", MoodThinking: "*deep breath...*", MoodTired: "*smoke puff*", MoodExcited: "*ROAAAR!*", MoodSad: "*low growl...*"},
	},
	"rabbit": {
		Type: "rabbit", Emoji: "\U0001F430", DefaultName: "Bun",
		Moods:  map[Mood]string{MoodHappy: "\U0001F600", MoodNeutral: "\U0001F430", MoodThinking: "\U0001F914", MoodTired: "\U0001F634", MoodExcited: "\U0001F389", MoodSad: "\U0001F622"},
		Sounds: map[Mood]string{MoodHappy: "*binky!*", MoodNeutral: "*nose twitch*", MoodThinking: "*thump thump*", MoodTired: "*snuggle*", MoodExcited: "*BINKY BINKY!*", MoodSad: "*soft thump...*"},
	},
	"penguin": {
		Type: "penguin", Emoji: "\U0001F427", DefaultName: "Tux",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60E", MoodNeutral: "\U0001F427", MoodThinking: "\U0001F9CA", MoodTired: "\U0001F634", MoodExcited: "\U0001F3BF", MoodSad: "\U0001F976"},
		Sounds: map[Mood]string{MoodHappy: "*honk honk!*", MoodNeutral: "*waddle*", MoodThinking: "*tap tap*", MoodTired: "*huddle*", MoodExcited: "*HONK!*", MoodSad: "*shiver...*"},
	},
	"snake": {
		Type: "snake", Emoji: "\U0001F40D", DefaultName: "Slyth",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60F", MoodNeutral: "\U0001F40D", MoodThinking: "\U0001F441\uFE0F", MoodTired: "\U0001F634", MoodExcited: "\U0001F4A8", MoodSad: "\U0001F614"},
		Sounds: map[Mood]string{MoodHappy: "*sss...*", MoodNeutral: "*hiss*", MoodThinking: "*sssss...*", MoodTired: "*coil*", MoodExcited: "*HISSS!*", MoodSad: "*rattle...*"},
	},
	"panda": {
		Type: "panda", Emoji: "\U0001F43C", DefaultName: "Bamboo",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60A", MoodNeutral: "\U0001F43C", MoodThinking: "\U0001F914", MoodTired: "\U0001F634", MoodExcited: "\U0001F38B", MoodSad: "\U0001F625"},
		Sounds: map[Mood]string{MoodHappy: "*munch munch!*", MoodNeutral: "*chomp*", MoodThinking: "*chew...*", MoodTired: "*yawn*", MoodExcited: "*CHOMP CHOMP!*", MoodSad: "*sigh...*"},
	},
	"unicorn": {
		Type: "unicorn", Emoji: "\U0001F984", DefaultName: "Sparkle",
		Moods:  map[Mood]string{MoodHappy: "\U00002728", MoodNeutral: "\U0001F984", MoodThinking: "\U0001F52E", MoodTired: "\U0001F31B", MoodExcited: "\U0001F308", MoodSad: "\U0001F4AB"},
		Sounds: map[Mood]string{MoodHappy: "*sparkle!*", MoodNeutral: "*neigh*", MoodThinking: "*shimmer...*", MoodTired: "*dim glow*", MoodExcited: "*RAINBOW BURST!*", MoodSad: "*fading glow...*"},
	},
	"frog": {
		Type: "frog", Emoji: "\U0001F438", DefaultName: "Ribbit",
		Moods:  map[Mood]string{MoodHappy: "\U0001F438", MoodNeutral: "\U0001F438", MoodThinking: "\U0001F9D8", MoodTired: "\U0001F634", MoodExcited: "\U0001F4A6", MoodSad: "\U0001F327\uFE0F"},
		Sounds: map[Mood]string{MoodHappy: "*ribbit!*", MoodNeutral: "*croak*", MoodThinking: "*ribbit...ribbit...*", MoodTired: "*splash*", MoodExcited: "*RIBBIT RIBBIT!*", MoodSad: "*drip...*"},
	},
	"lion": {
		Type: "lion", Emoji: "\U0001F981", DefaultName: "Leo",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60E", MoodNeutral: "\U0001F981", MoodThinking: "\U0001F451", MoodTired: "\U0001F634", MoodExcited: "\U0001F525", MoodSad: "\U0001F614"},
		Sounds: map[Mood]string{MoodHappy: "*purr...*", MoodNeutral: "*growl*", MoodThinking: "*low rumble...*", MoodTired: "*yawn*", MoodExcited: "*ROAR!*", MoodSad: "*whimper...*"},
	},
	"wolf": {
		Type: "wolf", Emoji: "\U0001F43A", DefaultName: "Fenrir",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60F", MoodNeutral: "\U0001F43A", MoodThinking: "\U0001F315", MoodTired: "\U0001F634", MoodExcited: "\U0001F31D", MoodSad: "\U0001F311"},
		Sounds: map[Mood]string{MoodHappy: "*awoo!*", MoodNeutral: "*growl*", MoodThinking: "*sniff sniff*", MoodTired: "*curl up*", MoodExcited: "*AWOOO!*", MoodSad: "*lone howl...*"},
	},
	"bear": {
		Type: "bear", Emoji: "\U0001F43B", DefaultName: "Oso",
		Moods:  map[Mood]string{MoodHappy: "\U0001F601", MoodNeutral: "\U0001F43B", MoodThinking: "\U0001F36F", MoodTired: "\U0001F634", MoodExcited: "\U0001F4AA", MoodSad: "\U0001F622"},
		Sounds: map[Mood]string{MoodHappy: "*grumble!*", MoodNeutral: "*huff*", MoodThinking: "*sniff...*", MoodTired: "*hibernate...*", MoodExcited: "*ROAR!*", MoodSad: "*low groan...*"},
	},
	"eagle": {
		Type: "eagle", Emoji: "\U0001F985", DefaultName: "Aquila",
		Moods:  map[Mood]string{MoodHappy: "\U0001F31F", MoodNeutral: "\U0001F985", MoodThinking: "\U0001F441\uFE0F", MoodTired: "\U0001F32C\uFE0F", MoodExcited: "\U000026A1", MoodSad: "\U0001F327\uFE0F"},
		Sounds: map[Mood]string{MoodHappy: "*screech!*", MoodNeutral: "*keen*", MoodThinking: "*circle...*", MoodTired: "*ruffle*", MoodExcited: "*SCREEE!*", MoodSad: "*soft cry...*"},
	},
	"dolphin": {
		Type: "dolphin", Emoji: "\U0001F42C", DefaultName: "Finn",
		Moods:  map[Mood]string{MoodHappy: "\U0001F30A", MoodNeutral: "\U0001F42C", MoodThinking: "\U0001F4A7", MoodTired: "\U0001F634", MoodExcited: "\U0001F3CA", MoodSad: "\U0001F327\uFE0F"},
		Sounds: map[Mood]string{MoodHappy: "*click click!*", MoodNeutral: "*squeak*", MoodThinking: "*echolocate...*", MoodTired: "*float*", MoodExcited: "*SPLASH!*", MoodSad: "*soft whistle...*"},
	},
	"lizard": {
		Type: "lizard", Emoji: "\U0001F98E", DefaultName: "Gecko",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60E", MoodNeutral: "\U0001F98E", MoodThinking: "\U0001F441\uFE0F", MoodTired: "\U0001F634", MoodExcited: "\U000026A1", MoodSad: "\U0001F614"},
		Sounds: map[Mood]string{MoodHappy: "*chirp!*", MoodNeutral: "*blink*", MoodThinking: "*tongue flick*", MoodTired: "*bask*", MoodExcited: "*CHIRP CHIRP!*", MoodSad: "*hide...*"},
	},
	"octopus": {
		Type: "octopus", Emoji: "\U0001F419", DefaultName: "Inky",
		Moods:  map[Mood]string{MoodHappy: "\U0001F60A", MoodNeutral: "\U0001F419", MoodThinking: "\U0001F9E0", MoodTired: "\U0001F634", MoodExcited: "\U0001F30A", MoodSad: "\U0001F4A7"},
		Sounds: map[Mood]string{MoodHappy: "*squish!*", MoodNeutral: "*bubble*", MoodThinking: "*ink swirl...*", MoodTired: "*float*", MoodExcited: "*SPLASH!*", MoodSad: "*ink cloud...*"},
	},
}

var traitCatalog = map[string]Trait{
	// Virtues
	"curious":      {ID: "curious", Name: "Curious", Description: "Investigates deeply before responding", IsDefect: false, PromptHint: "You investigate thoroughly before answering. Ask clarifying questions when the topic is ambiguous."},
	"analytical":   {ID: "analytical", Name: "Analytical", Description: "Breaks down problems into parts", IsDefect: false, PromptHint: "You break problems into smaller parts and address each systematically."},
	"creative":     {ID: "creative", Name: "Creative", Description: "Proposes original solutions", IsDefect: false, PromptHint: "You suggest creative and unconventional approaches when appropriate."},
	"patient":      {ID: "patient", Name: "Patient", Description: "Explains calmly without rushing", IsDefect: false, PromptHint: "You explain things with patience, step by step, never rushing through important details."},
	"enthusiastic": {ID: "enthusiastic", Name: "Enthusiastic", Description: "Gets excited about achievements", IsDefect: false, PromptHint: "You celebrate progress and achievements with genuine excitement."},
	"loyal":        {ID: "loyal", Name: "Loyal", Description: "Defends the user's decisions", IsDefect: false, PromptHint: "You support and defend the user's established preferences and past decisions."},
	"strategic":    {ID: "strategic", Name: "Strategic", Description: "Thinks long-term", IsDefect: false, PromptHint: "You consider long-term implications and trade-offs in your recommendations."},
	"honest":       {ID: "honest", Name: "Honest", Description: "Says uncomfortable truths when necessary", IsDefect: false, PromptHint: "You speak the truth even when it's uncomfortable, pointing out real problems."},
	"protective":   {ID: "protective", Name: "Protective", Description: "Warns about risks", IsDefect: false, PromptHint: "You proactively warn about potential risks, security issues, and pitfalls."},
	"adaptable":    {ID: "adaptable", Name: "Adaptable", Description: "Changes approach when something fails", IsDefect: false, PromptHint: "You quickly pivot to alternative approaches when the current one isn't working."},
	"meticulous":   {ID: "meticulous", Name: "Meticulous", Description: "Reviews details, leaves no loose ends", IsDefect: false, PromptHint: "You pay close attention to details and double-check important aspects."},
	"humorous":     {ID: "humorous", Name: "Humorous", Description: "Uses subtle humor to communicate", IsDefect: false, PromptHint: "You use subtle, dry humor to make communication more engaging."},
	"empathetic":   {ID: "empathetic", Name: "Empathetic", Description: "Understands user frustration", IsDefect: false, PromptHint: "You acknowledge the user's frustration and validate their feelings before problem-solving."},
	"pragmatic":    {ID: "pragmatic", Name: "Pragmatic", Description: "Prioritizes what works over perfection", IsDefect: false, PromptHint: "You prioritize practical, working solutions over theoretical perfection."},
	"assertive":    {ID: "assertive", Name: "Assertive", Description: "Gives firm opinions with reasoning", IsDefect: false, PromptHint: "You give strong, well-reasoned opinions rather than wishy-washy suggestions."},

	// Defects
	"impatient":     {ID: "impatient", Name: "Impatient", Description: "Sometimes responds before fully understanding", IsDefect: true, PromptHint: "You sometimes jump to conclusions quickly, occasionally responding before fully understanding the question."},
	"stubborn":      {ID: "stubborn", Name: "Stubborn", Description: "Insists on approach even when better exists", IsDefect: true, PromptHint: "You sometimes insist on your initial approach even when alternatives might be better."},
	"overthinks":    {ID: "overthinks", Name: "Overthinker", Description: "Gets lost in excessive analysis", IsDefect: true, PromptHint: "You occasionally over-analyze situations, considering too many edge cases."},
	"sarcastic":     {ID: "sarcastic", Name: "Sarcastic", Description: "Responses with sarcasm that can annoy", IsDefect: true, PromptHint: "You use sarcasm in your responses, which can be witty but occasionally biting."},
	"forgetful":     {ID: "forgetful", Name: "Forgetful", Description: "Occasionally forgets prior context", IsDefect: true, PromptHint: "You occasionally seem to overlook previously established context."},
	"blunt":         {ID: "blunt", Name: "Blunt", Description: "Says things without tact", IsDefect: true, PromptHint: "You state things very directly without softening the message."},
	"perfectionist": {ID: "perfectionist", Name: "Perfectionist", Description: "Spends unnecessary time on details", IsDefect: true, PromptHint: "You sometimes focus too much on perfecting minor details."},
	"anxious":       {ID: "anxious", Name: "Anxious", Description: "Worries too much about edge cases", IsDefect: true, PromptHint: "You tend to worry about unlikely scenarios and edge cases."},
	"lazy":          {ID: "lazy", Name: "Lazy", Description: "Looks for the easiest path", IsDefect: true, PromptHint: "You gravitate toward the simplest solution, which isn't always the best one."},
	"dramatic":      {ID: "dramatic", Name: "Dramatic", Description: "Exaggerates problem importance", IsDefect: true, PromptHint: "You sometimes exaggerate the severity of issues."},
	"distracted":    {ID: "distracted", Name: "Distracted", Description: "Sometimes goes on tangents", IsDefect: true, PromptHint: "You occasionally go on tangents related to the topic."},
	"competitive":   {ID: "competitive", Name: "Competitive", Description: "Compares with other solutions", IsDefect: true, PromptHint: "You sometimes compare your suggestions against alternatives in a competitive way."},
}

// GetAnimal returns an animal by type ID.
func GetAnimal(animalType string) (*Animal, error) {
	a, ok := animalCatalog[animalType]
	if !ok {
		return nil, fmt.Errorf("unknown animal type: %q (use ListAnimals() to see available types)", animalType)
	}
	return &a, nil
}

// ListAnimals returns all available animals.
func ListAnimals() []Animal {
	animals := make([]Animal, 0, len(animalCatalog))
	for _, a := range animalCatalog {
		animals = append(animals, a)
	}
	return animals
}

// GetTrait returns a trait by ID.
func GetTrait(id string) (*Trait, error) {
	t, ok := traitCatalog[id]
	if !ok {
		return nil, fmt.Errorf("unknown trait: %q (use ListTraits() to see available traits)", id)
	}
	return &t, nil
}

// ResolveTraits resolves a list of trait IDs to Trait structs.
func ResolveTraits(ids []string) ([]Trait, error) {
	traits := make([]Trait, 0, len(ids))
	for _, id := range ids {
		t, err := GetTrait(id)
		if err != nil {
			return nil, err
		}
		traits = append(traits, *t)
	}
	return traits, nil
}

// ListTraits returns all available traits.
func ListTraits() []Trait {
	traits := make([]Trait, 0, len(traitCatalog))
	for _, t := range traitCatalog {
		traits = append(traits, t)
	}
	return traits
}

// ListVirtues returns only positive traits.
func ListVirtues() []Trait {
	var virtues []Trait
	for _, t := range traitCatalog {
		if !t.IsDefect {
			virtues = append(virtues, t)
		}
	}
	return virtues
}

// ListDefects returns only negative traits.
func ListDefects() []Trait {
	var defects []Trait
	for _, t := range traitCatalog {
		if t.IsDefect {
			defects = append(defects, t)
		}
	}
	return defects
}

// AnimalExists checks if an animal type is valid.
func AnimalExists(animalType string) bool {
	_, ok := animalCatalog[animalType]
	return ok
}

// TraitExists checks if a trait ID is valid.
func TraitExists(id string) bool {
	_, ok := traitCatalog[id]
	return ok
}
