package unit

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

func (Goblin) Passable() bool {
	return false
}

func (g *Goblin) Do(op Operator) Job {
	j := g.Job()
	if j == nil {
		g.AddJob(op.KillClosestEnemy(g, g.Range()))
		return g.Job()
	}
	return j
}
