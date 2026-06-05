package game

import "strings"

func InterpretCommand(rooms map[string]Room, roomID string, raw string) (output string, kind string) {
	cmd := strings.TrimSpace(raw)
	if cmd == "" {
		return "You open your mouth... and say nothing. (Try typing a command.)", "error"
	}

	switch strings.ToLower(cmd) {
	case "help", "?":
		return "Commands: look, wait, help", "system"
	case "look", "l":
		room, ok := rooms[roomID]
		if !ok {
			return "You look around, but reality fails to load. (Unknown room id.)", "error"
		}
		return room.Description, "system"
	case "wait":
		return "You wait. Somewhere, a pipe ticks. The world does not feel rushed.", "system"
	default:
		return "You try that. Nothing happens. The dungeon remains unimpressed.", "error"
	}
}