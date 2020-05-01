package unit

import "github.com/gdamore/tcell"

var nextID int

type WithFaction interface {
	FactionID() int
	FactionData() Faction
}

type Faction struct {
	id    int
	style tcell.Style
}

func NewFaction(style tcell.Style) Faction {
	id := nextID
	nextID++
	return Faction{id: id, style: style}
}

func (f Faction) Style() tcell.Style {
	return f.style
}

func (f Faction) Equal(b WithFaction) bool {
	return f.FactionID() == b.FactionID()
}

func (f Faction) IsEnemy(b WithFaction) bool {
	return !f.Equal(b)
}

func (f Faction) FactionID() int {
	return f.id
}

func (f Faction) FactionData() Faction {
	return f
}
