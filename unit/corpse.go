package unit

import (
	"github.com/dsociative/arena/gmap"
	"github.com/gdamore/tcell"
)

type Corpse struct {
	Source   Unit
	KilledBy Unit
	*Creature
}

func NewCorpse(source, killedBy gmap.MapObject) Corpse {
	c := Corpse{Creature: &Creature{hp: 5}}
	c.KilledBy, _ = killedBy.Object().(Unit)
	c.Source, _ = source.Object().(Unit)
	return c
}

func (c Corpse) Hit(me *gmap.MapObject, whom gmap.MapObject, amount int) bool {
	c.Creature.hp -= amount
	return c.Creature.IsAlive()
}

func (c Corpse) Do(me gmap.MapObject) gmap.MapObject {
	return me
}

func (c Corpse) Passable() bool {
	return !c.Creature.IsAlive()
}

func (c Corpse) Pic() rune {
	if c.Creature.IsAlive() {
		return '☨'
	} else {
		return tcell.RuneBoard
	}
}

func (Corpse) Style() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorDimGray).Foreground(tcell.ColorWhiteSmoke)
}
