package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/robjsliwa/adv-htmx/internal/game"
)

type Server struct {
	Rooms    map[string]game.Room
	Sessions *game.SessionStore
	Tmpl     *template.Template
}

type RoomPageData struct {
	Room      game.Room
	Log       []game.LogEntry
	Inventory []game.Item
	ItemsHere []game.Item
	OOB       bool
}

type ListPartialData struct {
	OOB       bool
	Inventory []game.Item
	ItemsHere []game.Item
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.handleIndex)
	mux.HandleFunc("GET /room/{id}", s.handleRoom)
	mux.HandleFunc("POST /command", s.handleCommand)

	mux.HandleFunc("POST /take/{id}", s.handleTake)

	mux.HandleFunc("DELETE /inventory/{id}", s.handleInventoryDelete)
	mux.HandleFunc("POST /drop/{id}", s.handleDropFallback)

	return mux
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/room/hallway", http.StatusFound)
}

func (s *Server) handleRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("id")
	room, ok := s.Rooms[roomID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	state := s.Sessions.Get(w, r)
	state.SetRoom(roomID)
	snap := state.Snapshot()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = s.Tmpl.ExecuteTemplate(w, "room.html", RoomPageData{
		Room:      room,
		Log:       snap.Log,
		Inventory: snap.Inventory,
		ItemsHere: snap.ItemsHere,
	})
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	state := s.Sessions.Get(w, r)
	snap := state.Snapshot()

	cmd := r.FormValue("cmd")
	output, kind := game.InterpretCommand(s.Rooms, snap.RoomID, cmd)
	entry := state.Append(cmd, output, kind)

	if IsHTMX(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Add("Vary", "HX-Request")
		_ = s.Tmpl.ExecuteTemplate(w, "logEntry", entry)
		return
	}

	http.Redirect(w, r, "/room/"+snap.RoomID, http.StatusSeeOther)
}

func (s *Server) handleTake(w http.ResponseWriter, r *http.Request) {
	state := s.Sessions.Get(w, r)
	before := state.Snapshot()

	itemID := r.PathValue("id")
	it, ok := state.Take(itemID)

	var entry game.LogEntry
	if ok {
		entry = state.Append(
			"take "+itemID,
			fmt.Sprintf("You take the %s and tuck it into your backpack.", it.Name),
			"system",
		)
	} else {
		entry = state.Append(
			"take "+itemID,
			"You reach for it, but your hand closes on air.",
			"error",
		)
	}

	if !IsHTMX(r) {
		http.Redirect(w, r, "/room/"+before.RoomID, http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Vary", "HX-Request")

	// Main swap: append to #log (because the form targeted #log).
	if err := s.Tmpl.ExecuteTemplate(w, "logEntry", entry); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// OOB swaps: refresh both lists.
	after := state.Snapshot()

	_ = s.Tmpl.ExecuteTemplate(w, "inventoryList", ListPartialData{
		OOB:       true,
		Inventory: after.Inventory,
	})

	_ = s.Tmpl.ExecuteTemplate(w, "roomItemsList", ListPartialData{
		OOB:       true,
		ItemsHere: after.ItemsHere,
	})
}

func (s *Server) handleInventoryDelete(w http.ResponseWriter, r *http.Request) {
	state := s.Sessions.Get(w, r)

	itemID := r.PathValue("id")
	it, ok := state.Drop(itemID)

	if !ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	entry := state.Append(
		"drop "+itemID,
		fmt.Sprintf("You drop the %s. It lands with the dignity of a potato.", it.Name),
		"system",
	)

	if !IsHTMX(r) {
		snap := state.Snapshot()
		http.Redirect(w, r, "/room/"+snap.RoomID, http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Vary", "HX-Request")

	// OOB: append to log
	_ = s.Tmpl.ExecuteTemplate(w, "logEntryOOB", entry)

	// OOB: refresh "Items Here"
	after := state.Snapshot()
	_ = s.Tmpl.ExecuteTemplate(w, "roomItemsList", ListPartialData{
		OOB:       true,
		ItemsHere: after.ItemsHere,
	})
}

func (s *Server) handleDropFallback(w http.ResponseWriter, r *http.Request) {
	state := s.Sessions.Get(w, r)
	snap := state.Snapshot()

	itemID := r.PathValue("id")
	it, ok := state.Drop(itemID)
	if ok {
		state.Append(
			"drop "+itemID,
			fmt.Sprintf("You drop the %s. The dungeon accepts your offering.", it.Name),
			"system",
		)
	}

	http.Redirect(w, r, "/room/"+snap.RoomID, http.StatusSeeOther)
}