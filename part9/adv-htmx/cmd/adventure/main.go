package main

import (
	"log"
	"net/http"
	"os"

	"github.com/robjsliwa/adv-htmx/internal/game"
	"github.com/robjsliwa/adv-htmx/internal/web"
	"gopkg.in/yaml.v3"
)

func main() {
	// Load the raw data
	data, err := os.ReadFile("world.yaml")
	if err != nil {
		log.Fatal("Failed to read world.yaml: ", err)
	}

	// Parse YAML into our World struct
	var world game.World
	if err := yaml.Unmarshal(data, &world); err != nil {
		log.Fatal("Failed to parse YAML: ", err)
	}

	// We convert lists to maps for fast lookup during gameplay
	roomsMap := make(map[string]game.Room)
	for i := range world.Rooms {
		r := world.Rooms[i]
		roomsMap[r.ID] = r
	}

	itemsMap := make(map[string]game.Item)
	for i := range world.Items {
		item := world.Items[i]
		itemsMap[item.ID] = item
	}

	// Verify that every exit points to a real room
	for _, room := range roomsMap {
		for dir, targetID := range room.Exits {
			if _, exists := roomsMap[targetID]; !exists {
				log.Fatalf("DATA ERROR: Room '%s' has exit '%s' pointing to unknown room '%s'", room.ID, dir, targetID)
			}
		}
		// Verify items exist
		for _, itemID := range room.Items {
			if _, exists := itemsMap[itemID]; !exists {
				log.Fatalf("DATA ERROR: Room '%s' contains unknown item '%s'", room.ID, itemID)
			}
		}
	}

	// Place items into their starting rooms
	seedItems := make(map[string][]game.Item)
	for _, r := range world.Rooms {
		var roomItems []game.Item
		for _, itemID := range r.Items {
			if item, exists := itemsMap[itemID]; exists {
				roomItems = append(roomItems, item)
			}
		}
		seedItems[r.ID] = roomItems
	}

	// Start the server
	templates := web.MustLoadTemplates()
	sessions := game.NewSessionStore(seedItems)

	srv := &web.Server{
		Rooms:    roomsMap,
		Sessions: sessions,
		Tmpl:     templates.T,
	}

	log.Printf("Booting '%s'...", world.Title)
	log.Printf("Loaded %d rooms and %d items.", len(roomsMap), len(itemsMap))
	log.Println("Server starting on http://localhost:4040")
	log.Fatal(http.ListenAndServe(":4040", srv.Routes()))
}