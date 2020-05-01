package gmap

type path struct {
	xy   XY
	prev *path
}

func newPath(xy XY, prev *path) *path {
	return &path{xy, prev}
}

func newRootPath(xy XY) *path {
	return newPath(xy, nil)
}

func (p *path) String() string {
	return p.xy.String()
}

func (p *path) WayBack() (way []XY) {
	node := p
	for node.prev != nil {
		way = append(way, node.xy)
		node = node.prev
	}
	// beautiful slice reverse in golang world
	for left, right := 0, len(way)-1; left < right; left, right = left+1, right-1 {
		way[left], way[right] = way[right], way[left]
	}
	return way
}

type pathList struct {
	paths    []*path
	explored map[XY]bool
}

func newPathList() *pathList {
	return &pathList{explored: map[XY]bool{}}
}

func (pl *pathList) append(p *path) {
	if !pl.explored[p.xy] {
		pl.paths = append(pl.paths, p)
		pl.explored[p.xy] = true
	}
}

func (pl *pathList) pop() (p *path) {
	p, pl.paths = pl.paths[0], pl.paths[1:]
	return p
}

func (pl *pathList) Len() int {
	if pl == nil {
		return 0
	}
	return len(pl.paths)
}
