package gmap

import (
	"fmt"
	"github.com/gdamore/tcell"
	"go.uber.org/atomic"
	"sort"
	"strconv"
)

var Debug = atomic.Bool{}

type Arena struct {
	Height, Width int

	objects map[int64]MapObject
	busyLoc map[XY]bool
	screen  tcell.Screen
}

func NewArena(screen tcell.Screen, width, height int) Arena {
	return Arena{height, width, map[int64]MapObject{}, map[XY]bool{}, screen}
}

func (a Arena) Reset() {
	a.screen.Clear()
	a.objects = map[int64]MapObject{}
	a.busyLoc = map[XY]bool{}
	a.screen.Sync()
}

func (a Arena) IsBusy(xy XY) bool {
	if xy.y > a.Height || xy.y < 0 {
		return true
	}
	if xy.x > a.Width || xy.x < 0 {
		return true
	}
	return a.busyLoc[xy] == true
}

func (a Arena) NewMapObject(xy XY, o Object) {
	a.set(newMapObject(xy, o))
}

func (a Arena) set(mo MapObject) {
	a.objects[mo.id] = mo
	a.busyLoc[mo.xy] = true
	a.paint(mo.xy, mo.object.Pic(), mo.object.Style())
}

func (a Arena) paint(xy XY, pic rune, st tcell.Style) {
	a.screen.SetContent(xy.x+2, xy.y+2, pic, nil, st)
}

func (a Arena) Sync() {
	a.screen.Sync()
}

func (a Arena) UpdateObject(target MapObject, f func(mo MapObject) MapObject) {
	if mo, ok := a.objects[target.id]; ok && mo.xy.Equal(target.xy) {
		//f(mo)
		a.set(f(mo))
	}
}

func (a Arena) Update(id int64, f func(mo MapObject) MapObject) {
	if mo, ok := a.objects[id]; ok {
		f(mo)
		//a.set(f(mo))
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

func (a Arena) getReachableByLine(position *path, target XY, pl *pathList) (*pathList, *path) {
	direction := target.Sub(position.xy).Limit(1)
	next := position.xy.Add(direction)
	if !a.IsBusy(next) {
		p := &path{xy: next, prev: position}
		pl.append(p)
		position = p
		return a.getReachableByLine(p, target, pl)
	}
	return pl, position
}

func (a Arena) getReachable(position *path, target XY, pl *pathList) *pathList {
	pl, position = a.getReachableByLine(position, target, pl)
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			xy := XY{x, y}
			next := position.xy.Add(xy)
			if !a.IsBusy(next) || next.Equal(target) {
				pl.append(&path{xy: next, prev: position})
			}
		}
	}
	return pl
}

func (a Arena) BuildPath(from, target MapObject) []XY {
	sp := a.pathStraight(from, target, 5)
	if sp.Len() > 0 {
		return sp.WayBack()
	}
	return a.pathAround(from, target)
}

func (a Arena) pathAround(from MapObject, target MapObject) []XY {
	reachable := a.getReachable(&path{xy: from.xy}, target.xy, &pathList{explored: map[XY]bool{}})
	for reachable.Len() > 0 {
		node := reachable.pop()
		if node.xy.Equal(target.xy) {
			return node.WayBack()
		}
		reachable = a.getReachable(node, target.xy, reachable)
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
		a.UpdateObject(mo, func(old MapObject) MapObject {
			delete(a.busyLoc, old.xy)
			a.paint(old.xy, 0, tcell.StyleDefault)
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
	a.screen.Clear()
	a.screen.SetContent(0, 0, ' ', []rune("Winner is " + strconv.Itoa(id)), st.Bold(true))
	a.screen.SetContent(0, 1, ' ', []rune("Frags"), tcell.StyleDefault.Foreground(tcell.ColorAntiqueWhite))

}

func (a Arena) DrawFrags(id int, st tcell.Style, frags int) {
	a.screen.SetContent(0, id+2, ' ', []rune(fmt.Sprintf("%d: %d", id, frags)), st)
}
