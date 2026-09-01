package game

import "time"

type World struct {
	Title     string `yaml:"title"`
	StartRoom string `yaml:"start_room"`
	Rooms     []Room `yaml:"rooms"`
	Items     []Item `yaml:"items"`
}

type Room struct {
	ID          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Image       string            `yaml:"image"`
	X           int               `yaml:"x"`
	Y           int               `yaml:"y"`
	Exits       map[string]string `yaml:"exits"` 
	Items       []string          `yaml:"items"` // List of Item IDs
}

type Slot string

const (
	SlotNone    Slot = ""
	SlotWeapon  Slot = "weapon"
	SlotOffhand Slot = "offhand"
)

type Item struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Slot        Slot   `yaml:"slot"`
}

type LogEntry struct {
	Turn    int
	At      time.Time
	Command string
	Output  string
	Kind    string // "system", "error", etc.
}

type Spell struct {
	ID          string
	Name        string
	Description string
	Cost        int // Mana cost
}
