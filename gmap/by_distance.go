package gmap

type distance struct {
	d  int
	mo MapObject
}

type byDistance struct {
	xy  XY
	mos []distance
}

func (s byDistance) Len() int {
	return len(s.mos)
}
func (s byDistance) Swap(i, j int) {
	s.mos[i], s.mos[j] = s.mos[j], s.mos[i]
}

func (s byDistance) Less(i, j int) bool {
	return s.mos[i].d < s.mos[j].d
}

func (s *byDistance) append(mo MapObject) {
	s.mos = append(s.mos, distance{mo.xy.Sub(s.xy).Len(), mo})
}
