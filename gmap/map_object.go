package gmap

import (
	"github.com/gdamore/tcell"
	"go.uber.org/atomic"
)

var nextMapObject atomic.Int64

type Object interface {
	Pic() rune
	Style() tcell.Style
	Hit(*MapObject, MapObject, int)
}

type MapObject struct {
	id     int64
	xy     XY
	object Object
}

func (mo MapObject) Object() Object {
	return mo.object
}

func (mo MapObject) InRange(target MapObject, rng int) bool {
	return target.xy.InRange(mo.xy, rng)
}

func (mo *MapObject) setXY(xy XY) {
	mo.xy = xy
}

func (mo *MapObject) SetObject(o Object) {
	mo.object = o
}

func (mo *MapObject) Hit(whom MapObject, i int) {
	mo.object.Hit(mo, whom, i)
}

func newMapObject(xy XY, object Object) MapObject {
	return MapObject{nextMapObject.Inc(), xy, object}
}
