package game

import "time"

type Exit struct {
	Label string
	To    string
}

type Room struct {
	ID          string
	Name        string
	Description string
	Exits       []Exit
}

type Slot string

const (
	SlotNone    Slot = ""
	SlotWeapon  Slot = "weapon"
	SlotOffhand Slot = "offhand"
)

type Item struct {
	ID          string
	Name        string
	Description string
	Slot        Slot
}

type LogEntry struct {
	Turn    int
	At      time.Time
	Command string
	Output  string
	Kind    string // "system", "error", etc.
}