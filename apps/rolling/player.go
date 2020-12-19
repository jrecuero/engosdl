package main

import "github.com/jrecuero/engosdl"

// Player represents entity for the player.
type Player struct {
	*engosdl.Entity
	team string
}

// NewPlayer creates a new player instance.
func NewPlayer(name string) *Player {
	return &Player{
		Entity: engosdl.NewEntity(name),
		team:   "player-team",
	}
}

// UpdateActions updates player possible action buttons.
func (p *Player) UpdateActions(actions []string) {
	for _, child := range p.GetChildren() {
		matched := false
		for _, action := range actions {
			if child.GetName() == action {
				matched = true
				break
			}
		}
		child.SetEnabled(matched)
	}
}
