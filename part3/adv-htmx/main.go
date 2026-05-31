package main

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
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

type LogEntry struct {
	Turn    int
	At      time.Time
	Command string
	Output  string
	Kind    string // "system", "user", ""
}

type RoomPageData struct {
	Room Room
	Log  []LogEntry
}

type GameState struct {
	mu     sync.Mutex
	RoomID string
	Turn   int
	Log    []LogEntry
}

func (g *GameState) SetRoom(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.RoomID = id
}

func (g *GameState) Snapshot() (roomID string, logCopy []LogEntry) {
	g.mu.Lock()
	defer g.mu.Unlock()
	roomID = g.RoomID
	logCopy = append([]LogEntry(nil), g.Log...)
	return
}

func (g *GameState) Append(command, output, kind string) LogEntry {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Turn++
	entry := LogEntry{
		Turn:    g.Turn,
		At:      time.Now(),
		Command: command,
		Output:  output,
		Kind:    kind,
	}
	g.Log = append(g.Log, entry)
	return entry
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[string]*GameState
}

func NewSessionStore() *SessionStore {
	return &SessionStore{sessions: make(map[string]*GameState)}
}

const sessionCookieName = "adv_session"

func (s *SessionStore) Get(w http.ResponseWriter, r *http.Request) *GameState {
	if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
		s.mu.Lock()
		gs := s.sessions[c.Value]
		s.mu.Unlock()
		if gs != nil {
			return gs
		}
	}

	// Create a new session
	sid := newSessionID()
	gs := &GameState{RoomID: "hallway"}
	gs.Append("", "You awaken in a place that smells faintly of dust and HTTP. Type `help`.", "system")

	s.mu.Lock()
	s.sessions[sid] = gs
	s.mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return gs
}

func newSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func interpretCommand(rooms map[string]Room, roomID string, raw string) (output string, kind string) {
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

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
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
	tmpl := template.Must(template.ParseFiles(
		"templates/room.html",
		"templates/log_entry.html",
	))

	sessions := NewSessionStore()
	mux := http.NewServeMux()

	// Redirect root to the starting area
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/room/hallway", http.StatusFound)
	})

	mux.HandleFunc("GET /room/{id}", func(w http.ResponseWriter, r *http.Request) {
		roomID := r.PathValue("id")

		room, ok := rooms[roomID]
		if !ok {
			http.NotFound(w, r)
			return
		}

		state := sessions.Get(w, r)
		state.SetRoom(roomID)

		_, logEntries := state.Snapshot()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "room.html", RoomPageData{
			Room: room,
			Log:  logEntries,
		}); err != nil {
			log.Printf("template execute error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /command", func(w http.ResponseWriter, r *http.Request) {
		state := sessions.Get(w, r)
		roomID, _ := state.Snapshot()

		cmd := r.FormValue("cmd")
		output, kind := interpretCommand(rooms, roomID, cmd)
		entry := state.Append(cmd, output, kind)

		if isHTMX(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Add("Vary", "HX-Request") // useful if you ever add caching
			if err := tmpl.ExecuteTemplate(w, "logEntry", entry); err != nil {
				log.Printf("logEntry template execute error: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// non-JS fallback: classic Post/Redirect/Get
		http.Redirect(w, r, "/room/"+roomID, http.StatusSeeOther)
	})

	// Start Server
	addr := ":4040"
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}