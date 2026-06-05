package game

import (
	"sort"
	"sync"
	"time"
)

type State struct {
	mu sync.Mutex

	RoomID string
	Turn   int
	Log    []LogEntry

	Inventory map[string]Item
	RoomItems map[string]map[string]Item
}

type Snapshot struct {
	RoomID    string
	Log       []LogEntry
	Inventory []Item
	ItemsHere []Item
}

func NewState(seedItems map[string][]Item) *State {
	s := &State{
		RoomID:     "hallway",
		Inventory:  make(map[string]Item),
		RoomItems:  make(map[string]map[string]Item),
		Log:        nil,
		Turn:       0,
	}

	for roomID, items := range seedItems {
		roomInv := make(map[string]Item, len(items))
		for _, it := range items {
			roomInv[it.ID] = it
		}
		s.RoomItems[roomID] = roomInv
	}

	s.Append("", "You awaken in a place that smells faintly of dust and HTTP. Type `help` or pick something up.", "system")
	return s
}

func (s *State) SetRoom(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RoomID = id
}

func (s *State) Append(command, output, kind string) LogEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Turn++
	entry := LogEntry{
		Turn:    s.Turn,
		At:      time.Now(),
		Command: command,
		Output:  output,
		Kind:    kind,
	}
	s.Log = append(s.Log, entry)
	return entry
}

func (s *State) Snapshot() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	logCopy := append([]LogEntry(nil), s.Log...)
	inv := sortedItems(s.Inventory)

	var here []Item
	if roomInv, ok := s.RoomItems[s.RoomID]; ok {
		here = sortedItems(roomInv)
	}

	return Snapshot{
		RoomID:    s.RoomID,
		Log:       logCopy,
		Inventory: inv,
		ItemsHere: here,
	}
}

func (s *State) Take(itemID string) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	roomInv := s.RoomItems[s.RoomID]
	if roomInv == nil {
		return Item{}, false
	}

	it, ok := roomInv[itemID]
	if !ok {
		return Item{}, false
	}

	delete(roomInv, itemID)
	s.Inventory[itemID] = it
	return it, true
}

func (s *State) Drop(itemID string) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.Inventory[itemID]
	if !ok {
		return Item{}, false
	}
	delete(s.Inventory, itemID)

	roomInv := s.RoomItems[s.RoomID]
	if roomInv == nil {
		roomInv = make(map[string]Item)
		s.RoomItems[s.RoomID] = roomInv
	}
	roomInv[itemID] = it

	return it, true
}

func sortedItems(m map[string]Item) []Item {
	if len(m) == 0 {
		return nil
	}
	out := make([]Item, 0, len(m))
	for _, it := range m {
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name == out[j].Name {
			return out[i].ID < out[j].ID
		}
		return out[i].Name < out[j].Name
	})
	return out
}