package game

import (
	"github.com/dsociative/arena/gmap"
	"github.com/dsociative/arena/unit"
)

type operator struct {
	game *game
	mo   gmap.MapObject
	unit unit.Unit
}

func (o operator) KillClosestEnemy() (j unit.Job) {
	o.game.gameMap.ObjectFilter(
		o.mo,
		func(mo gmap.MapObject) bool {
			if unit, ok := mo.Object().(unit.Unit); ok {
				return unit.IsEnemy(o.unit)
			}
			return false
		},
		func(enemyMO gmap.MapObject) gmap.MapObject {
			if enemyMO.InRange(o.mo, o.unit.Range()) {
				j = &KillJob{enemyMO}
			} else {
				path := o.game.gameMap.BuildPath(o.mo, enemyMO)
				j = &MoveJob{path}
			}
			return enemyMO
		},
	)
	return
}

func (o operator) Hit(target gmap.MapObject) (ok bool) {
	o.game.gameMap.UpdateObject(target, func(mo gmap.MapObject) gmap.MapObject {
		if _, ok := mo.Object().(unit.Unit); ok {
			mo.Hit(o.mo, 1)
		}
		return mo
	})
	return
}

func (o operator) Move(target gmap.XY) (ok bool) {
	return o.game.gameMap.Move(o.mo, target)
}

type KillJob struct {
	target gmap.MapObject
}

func (j *KillJob) Do(op unit.Operator, executor *gmap.MapObject) bool {
	return op.Hit(j.target)
}

type MoveJob struct {
	path []gmap.XY
}

func (j *MoveJob) Do(op unit.Operator, executor *gmap.MapObject) bool {
	if len(j.path) > 1 {
		var xy gmap.XY
		xy, j.path = j.path[0], j.path[1:]
		return op.Move(xy)
	}
	return false
}
