package gmap

import (
	"fmt"
	"math"
)

type XY struct {
	x, y int
}

func NewXY(x, y int) XY {
	return XY{x, y}
}

func (c XY) InRange(b XY, rng int) bool {
	return b.Sub(c).Len() <= rng
}

func (c XY) Len() int {
	return int(math.Sqrt(math.Pow(float64(c.x), 2) + math.Pow(float64(c.y), 2)))
}

func (c XY) Sub(b XY) XY {
	return XY{c.x - b.x, c.y - b.y}
}

func (c XY) Add(b XY) XY {
	return XY{c.x + b.x, c.y + b.y}
}

func limit(a, b int) int {
	n := 1
	if a < 0 {
		n = -1
		a *= -1
	}
	if a > b {
		return b * n
	}
	return a * n
}

func (c XY) Limit(i int) XY {
	return XY{limit(c.x, i), limit(c.y, i)}
}

func (c XY) Equal(b XY) bool {
	return c.x == b.x && c.y == b.y
}

func (c XY) String() string {
	return fmt.Sprintf("[%d;%d]", c.x, c.y)
}
