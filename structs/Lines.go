package structs

import "sort"

type Lines struct {
	Cells []Cell
}

func (lines *Lines) AddCell(cell Cell) {
	lines.Cells = append(lines.Cells, cell)
}

func (lines *Lines) FindCell(cell Cell) (int, bool) {
	for i, c := range lines.Cells {
		if c.X == cell.X && c.Y == cell.Y {
			return i, true
		}
	}
	return -1, false
}

type LinesSlice []Cell

func (l LinesSlice) Len() int {
	return len(l)
}

func (l LinesSlice) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l LinesSlice) Less(i, j int) bool {
	if l[i].X == l[j].X {
		return l[i].Y < l[j].Y
	}
	return l[i].X < l[j].X
}

func (lines *Lines) SortLines() {
	sort.Sort(LinesSlice(lines.Cells))
}
