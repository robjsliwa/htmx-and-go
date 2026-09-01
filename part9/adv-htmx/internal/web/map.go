package web

import (
	"fmt"
	"strings"

	"github.com/robjsliwa/adv-htmx/internal/game"
)

// RenderMap generates a raw SVG string based on the player's knowledge of the world.
func RenderMap(currentRoomID string, visited map[string]bool, rooms map[string]game.Room, goblinRoomID string) string {
	var sb strings.Builder

	// Start the SVG canvas
	sb.WriteString(`<svg width="200" height="200" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">`)
	
	// Draw a dark background
	sb.WriteString(`<rect width="100%" height="100%" fill="#1a1a1a" />`)

	// Draw the rooms
	for id, room := range rooms {
		// If the player hasn't seen it, it doesn't exist on the map!
		if !visited[id] {
			continue
		}

		// Apply our Coordinate Transformation
		x := 100 + (room.X * 40)
		y := 100 - (room.Y * 40)

		// Determine the color: Bright green for "You are here", grey for "Visited"
		color := "#555"
		if id == currentRoomID {
			color = "#33ff00"
		}

		// Draw a room box centered on our calculated point
		// (We offset by -15 so the 30x30 box is centered)
		sb.WriteString(fmt.Sprintf(
			`<rect x="%d" y="%d" width="30" height="30" rx="4" fill="%s" stroke="#000" stroke-width="2" />`,
			x-15, y-15, color,
		))

        // Draw the Goblin if he is here!
        if id == goblinRoomID {
            // A scary red circle with a CSS animation
            sb.WriteString(fmt.Sprintf(
                `<circle cx="%d" cy="%d" r="8" fill="#ff0000" class="pulse-danger" />`,
                x, y,
            ))
        }
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}