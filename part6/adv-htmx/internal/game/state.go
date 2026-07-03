package game

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	ErrNotInInventory = errors.New("you don't have that item")
	ErrNotEquippable  = errors.New("that item cannot be equipped")
	ErrUnknownSlot    = errors.New("unknown equipment slot")
)

type EquipResult struct {
	Item       Item
	Slot       Slot
	Equipped   bool // true if equipped after the toggle, false if unequipped
	Replaced   Item
	ReplacedOK bool
}

type State struct {
	mu sync.Mutex

	RoomID string
	Turn   int
	Log    []LogEntry

	Inventory map[string]Item
	RoomItems map[string]map[string]Item

	WeaponID  string
	OffhandID string
}

type Snapshot struct {
	RoomID    string
	Log       []LogEntry
	Inventory []Item
	ItemsHere []Item

	WeaponID  string
	OffhandID string

	Weapon    Item
	WeaponOK  bool
	Offhand   Item
	OffhandOK bool
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

	var weapon Item
	var weaponOK bool
	if s.WeaponID != "" {
		if it, ok := s.Inventory[s.WeaponID]; ok {
			weapon = it
			weaponOK = true
		}
	}

	var offhand Item
	var offhandOK bool
	if s.OffhandID != "" {
		if it, ok := s.Inventory[s.OffhandID]; ok {
			offhand = it
			offhandOK = true
		}
	}

	return Snapshot{
		RoomID:    s.RoomID,
		Log:       logCopy,
		Inventory: inv,
		ItemsHere: here,

		WeaponID:  s.WeaponID,
		OffhandID: s.OffhandID,

		Weapon:    weapon,
		WeaponOK:  weaponOK,
		Offhand:   offhand,
		OffhandOK: offhandOK,
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

	// If the item was equipped, unequip it.
	if s.WeaponID == itemID {
		s.WeaponID = ""
	}
	if s.OffhandID == itemID {
		s.OffhandID = ""
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

func (s *State) ToggleEquip(itemID string) (EquipResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	it, ok := s.Inventory[itemID]
	if !ok {
		return EquipResult{}, ErrNotInInventory
	}
	if it.Slot == SlotNone {
		return EquipResult{}, ErrNotEquippable
	}

	var slotPtr *string
	switch it.Slot {
	case SlotWeapon:
		slotPtr = &s.WeaponID
	case SlotOffhand:
		slotPtr = &s.OffhandID
	default:
		return EquipResult{}, ErrUnknownSlot
	}

	// If already equipped -> unequip.
	if *slotPtr == itemID {
		*slotPtr = ""
		return EquipResult{
			Item:     it,
			Slot:     it.Slot,
			Equipped: false,
		}, nil
	}

	// Otherwise equip, possibly replacing.
	var replaced Item
	var replacedOK bool
	if *slotPtr != "" {
		if prev, ok := s.Inventory[*slotPtr]; ok {
			replaced = prev
			replacedOK = true
		}
	}

	*slotPtr = itemID
	return EquipResult{
		Item:       it,
		Slot:       it.Slot,
		Equipped:   true,
		Replaced:   replaced,
		ReplacedOK: replacedOK,
	}, nil
}