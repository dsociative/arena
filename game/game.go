package game

import (
	"context"
	"github.com/dsociative/arena/gmap"
	"github.com/dsociative/arena/unit"
	"github.com/gdamore/tcell"
	"math/rand"
	"time"
)

type gameState int

const (
	maxRandomSpawnRetry = 10

	gameNew gameState = iota
	gameRunning
	gameEnd
)

type game struct {
	state      gameState
	spawnCount int
	gameMap    gmap.Arena
	total      map[unit.Faction]int
	factions   []unit.Faction
}

func NewGame(gm gmap.Arena, spawnCount int) *game {
	return &game{
		gameNew,
		spawnCount,
		gm,
		map[unit.Faction]int{},
		[]unit.Faction{
			unit.NewFaction(tcell.StyleDefault.Foreground(tcell.NewRGBColor(240, 10, 10))),
			unit.NewFaction(tcell.StyleDefault.Foreground(tcell.NewRGBColor(10, 10, 240))),
			unit.NewFaction(tcell.StyleDefault.Foreground(tcell.NewRGBColor(240, 240, 10))),
			unit.NewFaction(tcell.StyleDefault.Foreground(tcell.NewRGBColor(10, 240, 10))),
		},
	}
}

func (g *game) SpawnRandom(u unit.Unit) {
	g.spawnRandom(u, 0)
}

func (g *game) spawnRandom(u unit.Unit, retry int) bool {
	if !g.Spawn(rand.Intn(g.gameMap.Width), rand.Intn(g.gameMap.Height), u) && retry < maxRandomSpawnRetry {
		g.spawnRandom(u, retry+1)
	}
	return true
}

func (g *game) Spawn(x, y int, u unit.Unit) bool {
	c := gmap.NewXY(x, y)
	if g.gameMap.IsBusy(c) {
		return false
	}
	g.total[u.FactionData()]++
	g.gameMap.NewMapObject(c, u)
	return true
}

func (g *game) Run(ctx context.Context, keys chan *tcell.EventKey) {
	ticker := time.NewTicker(time.Millisecond * 50)
	for {
		switch g.state {
		case gameNew:
			g.gameMap.Reset()
			for i := 0; i < g.spawnCount; i++ {
				for _, f := range g.factions {
					g.SpawnRandom(unit.NewGoblin(f))
				}
			}
			g.state = gameRunning
		case gameRunning:
			select {
			case <-ctx.Done():
				g.gameMap.Finish()
				return
			case <-ticker.C:
				g.tick()
				g.gameMap.Sync()
			}
		case gameEnd:
			<-keys
			g.state = gameNew
		}
	}
}

func (g *game) tick() {
	alive := map[unit.Faction]int{}

	for _, id := range g.gameMap.ObjectsIDs() {
		g.gameMap.Object(id, func(mo gmap.MapObject) gmap.MapObject {
			if unit, ok := mo.Object().(unit.Unit); ok {
				alive[unit.FactionData()] ++
				o := operator{g, mo, unit}
				if j := unit.Do(o); j != nil {
					if !j.Do(o, &mo) {
						unit.Done()
					}
				}
			}
			return mo
		})
	}

	for faction, total := range g.total {
		g.gameMap.DrawScore(faction.FactionID(), faction.Style(), total, alive[faction])
	}

	if len(alive) <= 1 {
		g.state = gameEnd

		var winner unit.Faction
		for f, _ := range alive {
			winner = f
		}

		killCount := map[unit.Faction]int{}
		for _, id := range g.gameMap.ObjectsIDs() {
			g.gameMap.Object(id, func(mo gmap.MapObject) gmap.MapObject {
				if unit, ok := mo.Object().(unit.Corpse); ok {
					killCount[unit.KilledBy.FactionData()] ++
				}
				return mo
			})
		}

		g.gameMap.DrawWinner(winner.FactionID(), winner.Style())
		for f, frags := range killCount {
			g.gameMap.DrawFrags(f.FactionID(), f.Style(), frags)
		}
	}
}
