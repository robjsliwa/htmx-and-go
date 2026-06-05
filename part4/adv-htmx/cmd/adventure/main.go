package main

import (
	"log"
	"net/http"

	"github.com/robjsliwa/adv-htmx/internal/game"
	"github.com/robjsliwa/adv-htmx/internal/web"
)

func main() {
	rooms := map[string]game.Room{
		"hallway": {
			ID:          "hallway",
			Name:        "The Hallway",
			Description: "You are standing in a dimly lit hallway. Dust motes dance in the air.",
			Exits: []game.Exit{
				{Label: "Go North", To: "atrium"},
			},
		},
		"atrium": {
			ID:          "atrium",
			Name:        "The Atrium",
			Description: "You step into a grand atrium. The ceiling is glass, revealing a grey sky.",
			Exits: []game.Exit{
				{Label: "Go South", To: "hallway"},
			},
		},
	}

	seedItems := map[string][]game.Item{
		"hallway": {
			{ID: "lamp", Name: "Brass Lamp", Description: "A squat brass lamp. It's heavier than it looks."},
			{ID: "key", Name: "Rusty Key", Description: "A key that looks like it has opened exactly one door and regretted it."},
		},
		"atrium": {
			{ID: "coin", Name: "Ancient Coin", Description: "A coin stamped with a face you don't recognize, but it recognizes you."},
		},
	}

	templates := web.MustLoadTemplates()
	sessions := game.NewSessionStore(seedItems)

	srv := &web.Server{
		Rooms:    rooms,
		Sessions: sessions,
		Tmpl:     templates.T,
	}

	addr := ":4040"
	log.Printf("Server starting on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}