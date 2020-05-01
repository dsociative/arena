package gmap

import (
	"fmt"
	"github.com/gdamore/tcell"
	"sort"
	"strconv"
)

type Arena struct {
	Height, Width int

	objects map[int64]MapObject
	busyLoc map[XY]int64
	screen  tcell.Screen
}

func NewArena(screen tcell.Screen, width, height int) Arena {
	return Arena{height, width, map[int64]MapObject{}, map[XY]int64{}, screen}
}

func (a Arena) Reset() {
	a.screen.Clear()
	a.objects = map[int64]MapObject{}
	a.busyLoc = map[XY]int64{}
	a.screen.Sync()
}

func (a Arena) IsBusy(xy XY) bool {
	i := a.busyLoc[xy]
	return a.inArea(xy) || i != 0
}

func (a Arena) inArea(xy XY) bool {
	return (xy.y > a.Height || xy.y < 0) || (xy.x > a.Width || xy.x < 0)
}

func (a Arena) NewMapObject(xy XY, o Object) {
	a.set(newMapObject(xy, o))
}

func (a Arena) set(mo MapObject) {
	a.objects[mo.id] = mo
	if !mo.object.Passable() {
		a.busyLoc[mo.xy] = mo.id
	}
	a.paint(mo.xy, mo.object.Pic(), mo.object.Style())
}

func (a Arena) update(old, new MapObject) {
	positionChanged := !old.xy.Equal(new.xy)
	if positionChanged {
		a.paint(old.xy, 0, tcell.StyleDefault)
	}
	if positionChanged || new.object.Passable() {
		delete(a.busyLoc, old.xy)
	}
	a.set(new)
}

func (a Arena) paint(xy XY, pic rune, st tcell.Style) {
	a.screen.SetContent(xy.x+2, xy.y+2, pic, nil, st)
}

func (a Arena) Sync() {
	a.screen.Sync()
}

func (a Arena) ExactObjectUpdate(target MapObject, f func(mo MapObject) MapObject) {
	if mo, ok := a.objects[target.id]; ok && mo.xy.Equal(target.xy) {
		a.update(mo, f(mo))
	}
}

func (a Arena) Object(id int64, f func(mo MapObject) MapObject) {
	if mo, ok := a.objects[id]; ok {
		f(mo)
	}
}

func (a Arena) ObjectFilter(mo MapObject, filterFun func(mo MapObject) bool, f func(mo MapObject) MapObject) {
	rt := byDistance{xy: mo.xy}
	for _, mo := range a.objects {
		if filterFun(mo) {
			rt.append(mo)
		}
	}
	sort.Sort(rt)
	if rt.Len() > 0 {
		f(rt.mos[0].mo)
	}
}

func (a Arena) AnythingBTW(position, target MapObject) (mo MapObject, ok bool) {
	next := position.xy.Add(target.xy.Sub(position.xy).Limit(1))
	if !next.Equal(position.xy) {
		var objectID int64
		if objectID, ok = a.busyLoc[next]; ok {
			mo, ok = a.objects[objectID]
			return
		}
	}
	return
}

func (a Arena) BuildPath(from, target MapObject) (path []XY) {
	path = a.straightPath(from.xy, target.xy, path, 5)
	if len(path) < 1 {
		path = a.pathAround(from, target)
	}
	return path
}

func (a Arena) straightPath(position XY, target XY, path []XY, limit int) []XY {
	direction := target.Sub(position).Limit(1)
	next := position.Add(direction)
	if !a.IsBusy(next) {
		path = append(path, next)
		position = next
		if len(path) >= limit {
			return path
		}
		return a.straightPath(next, target, path, limit)
	}
	return path
}

func (a Arena) reachable(position *path, target XY, pl *pathList) *pathList {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			next := position.xy.Add(XY{x, y})
			if !a.IsBusy(next) || next.Equal(target) {
				pl.append(&path{xy: next, prev: position})
			}
		}
	}
	return pl
}

func (a Arena) pathAround(from MapObject, target MapObject) []XY {
	reachable := a.reachable(newRootPath(from.xy), target.xy, newPathList())
	for reachable.Len() > 0 {
		node := reachable.pop()
		if node.xy.Equal(target.xy) {
			return node.WayBack()
		}
		reachable = a.reachable(node, target.xy, reachable)
	}
	return nil
}

func (a Arena) ObjectsIDs() (mos []int64) {
	l := len(mos)
	mos = make([]int64, 0, l)
	for _, mo := range a.objects {
		mos = append(mos, mo.id)
	}
	return
}

func (a Arena) Move(mo MapObject, target XY) (ok bool) {
	if !a.IsBusy(target) {
		a.ExactObjectUpdate(mo, func(old MapObject) MapObject {
			mo.xy = target
			ok = true
			return mo
		})
	}
	return
}

func (a Arena) Finish() {
	a.screen.Fini()
}

func (a Arena) DrawScore(x int, st tcell.Style, total int, alive int) {
	s := strconv.Itoa(total) + "/" + strconv.Itoa(alive)
	a.screen.SetContent(x, 0, ' ', []rune(s), st)
}

func (a Arena) DrawWinner(id int, st tcell.Style) {
	a.screen.SetContent(0, 1, ' ', []rune("Winner is "+strconv.Itoa(id)), st.Bold(true))
	a.screen.SetContent(0, 2, ' ', []rune("Frags"), tcell.StyleDefault.Foreground(tcell.ColorAntiqueWhite))

}

func (a Arena) DrawFrags(id int, st tcell.Style, frags int) {
	a.screen.SetContent(0, id+3, ' ', []rune(fmt.Sprintf("%d: %d", id, frags)), st)
}
