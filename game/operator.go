package game

import (
	"github.com/dsociative/arena/gmap"
	"github.com/dsociative/arena/unit"
)

var _ unit.Operator = operator{}

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
				if len(path) > 0 {
					j = &MoveJob{path}
				} else {
					j = &KillAnythingBTW{target: enemyMO}
				}
			}
			return enemyMO
		},
	)
	return
}

func (o operator) AnythingBTW(target gmap.MapObject) (gmap.MapObject, bool) {
	return o.game.gameMap.AnythingAround(o.mo, target)
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
	if len(j.path) > 0 {
		var xy gmap.XY
		xy, j.path = j.path[0], j.path[1:]
		return op.Move(xy)
	}
	return false
}

type KillAnythingBTW struct {
	target gmap.MapObject
	killJob unit.Job
}

func (j *KillAnythingBTW) Do(op unit.Operator, executor *gmap.MapObject) bool {
	if j.killJob != nil {
		return j.killJob.Do(op, executor)
	}
	if target, ok := op.AnythingBTW(j.target); ok {
		j.killJob = &KillJob{target: target}
		return true
	}
	return false
}
