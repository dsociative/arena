package unit

import (
	"github.com/dsociative/arena/gmap"
	"github.com/gdamore/tcell"
)

type Operator interface {
	KillClosestEnemy() Job
	Hit(target gmap.MapObject) bool
	Move(xy gmap.XY) bool
}

type Job interface {
	Do(Operator, *gmap.MapObject) bool
}

type Unit interface {
	gmap.Object
	WithFaction
	IsEnemy(WithFaction) bool
	Range() int
	IsAlive() bool
	Do(Operator) Job
	Done()
}

type JobQueue struct {
	job Job
}

func (j *JobQueue) AddJob(job Job) {
	j.job = job
}

func (j *JobQueue) Done() {
	j.job = nil
}

func (j *JobQueue) Job() Job {
	return j.job
}

type Goblin struct {
	Faction
	*Creature
	*JobQueue
}

func NewGoblin(faction Faction) *Goblin {
	return &Goblin{faction, &Creature{hp: 20}, &JobQueue{}}
}

func (g Goblin) Pic() rune {
	return 'G'
}

func (Goblin) Range() int {
	return 1
}

func (g *Goblin) Do(op Operator) Job {
	j := g.Job()
	if j == nil {
		g.AddJob(op.KillClosestEnemy())
		return g.Job()
	}
	return j
}

type Creature struct {
	hp int
}

func (c *Creature) Hit(me *gmap.MapObject, whom gmap.MapObject, amount int) {
	alive := c.IsAlive()
	if alive {
		c.hp -= amount
		if !c.IsAlive() {
			me.SetObject(NewCorpse(*me, whom))
		}
	}
}

func (c Creature) IsAlive() bool {
	return c.hp > 1
}

type Corpse struct {
	Source   Unit
	KilledBy Unit
}

func NewCorpse(source, killedBy gmap.MapObject) Corpse {
	c := Corpse{}
	c.KilledBy, _ = killedBy.Object().(Unit)
	c.Source, _ = source.Object().(Unit)
	return c
}

func (c Corpse) Hit(*gmap.MapObject, gmap.MapObject, int) {}

func (c Corpse) Do(me gmap.MapObject) gmap.MapObject {
	return me
}

func (Corpse) Pic() rune {
	return '☨'
}

func (Corpse) Style() tcell.Style {
	return tcell.StyleDefault.Background(tcell.ColorDimGray).Foreground(tcell.ColorWhiteSmoke)
}
