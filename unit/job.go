package unit

import "github.com/dsociative/arena/gmap"

type Operator interface {
	KillClosestEnemy(WithFaction, int) Job
	Hit(target gmap.MapObject) bool
	Move(xy gmap.XY) bool
	AnythingBTW(target gmap.MapObject) (gmap.MapObject, bool)
}

type Job interface {
	Do(Operator) bool
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

