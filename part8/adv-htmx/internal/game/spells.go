package game

import "strings"

var Grimoire = []Spell{
	{ID: "light", Name: "Light", Description: "Illuminates dark places.", Cost: 1},
	{ID: "missile", Name: "Magic Missile", Description: "Always hits. Always annoying.", Cost: 3},
	{ID: "fireball", Name: "Fireball", Description: "Solves most problems, creates new structural ones.", Cost: 10},
	{ID: "ice", Name: "Ice Spike", Description: "Cool off your enemies.", Cost: 5},
	{ID: "heal", Name: "Lesser Heal", Description: "Stops the bleeding, mostly.", Cost: 4},
	{ID: "levitate", Name: "Levitate", Description: "Walk over traps with style.", Cost: 6},
	{ID: "silence", Name: "Silence", Description: "Good for librarians and evil wizards.", Cost: 8},
	{ID: "identify", Name: "Identify", Description: "What does this button do?", Cost: 2},
}

// SearchSpells is our "database query"
func SearchSpells(query string) []Spell {
	if query == "" {
		return nil
	}
	
	q := strings.ToLower(query)
	var results []Spell
	
	for _, s := range Grimoire {
		if strings.Contains(strings.ToLower(s.Name), q) {
			results = append(results, s)
		}
	}
	return results
}