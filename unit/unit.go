package unit

import (
	"github.com/dsociative/arena/gmap"
)

type Unit interface {
	gmap.Object
	WithFaction
	IsEnemy(WithFaction) bool
	Range() int
	IsAlive() bool
	Do(Operator) Job
	Done()
}

type Creature struct {
	hp int
}

func (c *Creature) Hit(me *gmap.MapObject, whom gmap.MapObject, amount int) bool {
	alive := c.IsAlive()
	if alive {
		c.hp -= amount
		alive = c.IsAlive()
		if !alive {
			me.SetObject(NewCorpse(*me, whom))
		}
	}
	return alive
}

func (c Creature) IsAlive() bool {
	return c.hp > 1
}
