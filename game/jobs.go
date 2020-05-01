package game

import (
	"github.com/dsociative/arena/gmap"
	"github.com/dsociative/arena/unit"
)

type KillJob struct {
	target gmap.MapObject
}

func NewKillJob(t gmap.MapObject) *KillJob {
	return &KillJob{t}
}

func (j *KillJob) Do(op unit.Operator) bool {
	return op.Hit(j.target)
}

type MoveJob struct {
	path []gmap.XY
}

func NewMoveJob(path []gmap.XY) *MoveJob {
	return &MoveJob{path}
}

func (j *MoveJob) Do(op unit.Operator) bool {
	if len(j.path) > 0 {
		var xy gmap.XY
		xy, j.path = j.path[0], j.path[1:]
		return op.Move(xy)
	}
	return false
}

type KillAnythingJob struct {
	target  gmap.MapObject
	killJob unit.Job
}

func NewKillAnythingJob(t gmap.MapObject) *KillAnythingJob {
	return &KillAnythingJob{target: t}
}

func (j *KillAnythingJob) Do(op unit.Operator) bool {
	if j.killJob != nil {
		return j.killJob.Do(op)
	}
	if target, ok := op.AnythingBTW(j.target); ok {
		j.killJob = &KillJob{target: target}
		return true
	}
	return false
}
