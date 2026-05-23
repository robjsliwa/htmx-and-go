package main

import (
	"html/template"
	"log"
	"net/http"
)

// Exit represents a navigable path.
type Exit struct {
	Label string 
	To    string 
}

// Room represents a node in our graph.
type Room struct {
	ID          string
	Name        string
	Description string
	Exits       []Exit
}

// RoomPageData acts as the ViewModel passed to the template.
type RoomPageData struct {
	Room Room
}

func main() {
	// In a real app, this would be a Repository connecting to a DB.
	rooms := map[string]Room{
		"hallway": {
			ID:          "hallway",
			Name:        "The Hallway",
			Description: "You are standing in a dimly lit hallway. Dust motes dance in the air.",
			Exits: []Exit{
				{Label: "Go North", To: "atrium"},
			},
		},
		"atrium": {
			ID:          "atrium",
			Name:        "The Atrium",
			Description: "You step into a grand atrium. The ceiling is glass, revealing a grey sky.",
			Exits: []Exit{
				{Label: "Go South", To: "hallway"},
			},
		},
	}

	// We parse once at startup for performance and to catch syntax errors early.
	tmpl := template.Must(template.ParseFiles("templates/room.html"))

	mux := http.NewServeMux()

	// Redirect root to the starting area
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/room/hallway", http.StatusFound)
	})

	// The Room Handler
	// Serves the state of a specific room based on URL path.
	mux.HandleFunc("GET /room/{id}", func(w http.ResponseWriter, r *http.Request) {
		roomID := r.PathValue("id") // Go 1.22 path value extraction

		room, exists := rooms[roomID]
		if !exists {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		
		// Inject the Go struct into the HTML template
		if err := tmpl.Execute(w, RoomPageData{Room: room}); err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Start Server
	addr := ":4040"
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}