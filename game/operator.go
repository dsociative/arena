package game

import (
	"github.com/dsociative/arena/gmap"
	"github.com/dsociative/arena/unit"
)

var _ unit.Operator = operator{}

type operator struct {
	game *game
	mo   gmap.MapObject
}

func (o operator) KillClosestEnemy(u unit.WithFaction, rng int) (j unit.Job) {
	o.game.gameMap.ObjectFilter(
		o.mo,
		func(mo gmap.MapObject) bool {
			if unit, ok := mo.Object().(unit.Unit); ok {
				return unit.IsEnemy(u)
			}
			return false
		},
		func(enemy gmap.MapObject) gmap.MapObject {
			if enemy.InRange(o.mo, rng) {
				j = NewKillJob(enemy)
			} else {
				path := o.game.gameMap.BuildPath(o.mo, enemy)
				if len(path) > 0 {
					j = NewMoveJob(path)
				} else {
					j = NewKillAnythingJob(enemy)
				}
			}
			return enemy
		},
	)
	return
}

func (o operator) AnythingBTW(target gmap.MapObject) (gmap.MapObject, bool) {
	return o.game.gameMap.AnythingBTW(o.mo, target)
}

func (o operator) Hit(target gmap.MapObject) (ok bool) {
	o.game.gameMap.ExactObjectUpdate(target, func(mo gmap.MapObject) gmap.MapObject {
		mo.Hit(o.mo, 1)
		return mo
	})
	return
}

func (o operator) Move(target gmap.XY) (ok bool) {
	return o.game.gameMap.Move(o.mo, target)
}
