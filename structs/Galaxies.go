package structs

type Galaxies struct {
	Galaxies []Galaxy
}

type Galaxy struct {
	Center   Center
	Cells    []Cell
	Type     string
	Complete bool
}

func (galaxies *Galaxies) AddGalaxy(targetGalaxy Galaxy) {
	galaxies.Galaxies = append(galaxies.Galaxies, targetGalaxy)
}

func (galaxies *Galaxies) FindGalaxy(targetGalaxy Galaxy) (int, bool) {
	for i, g := range galaxies.Galaxies {
		if g.Center.Cell.X == targetGalaxy.Center.Cell.X && g.Center.Cell.Y == targetGalaxy.Center.Cell.Y {
			return i, true
		}
	}
	return -1, false
}

func (galaxies *Galaxies) MarkGalaxyComplete(targetGalaxy Galaxy) bool {
	if i, f := galaxies.FindGalaxy(targetGalaxy); f {
		galaxies.Galaxies[i].Complete = true
		return true
	}
	return false
}

func (galaxies *Galaxies) AddCell(galaxy Galaxy, cell Cell) {
	if galaxyIndex, found := galaxies.FindGalaxy(galaxy); found {
		galaxies.Galaxies[galaxyIndex].Cells = append(galaxies.Galaxies[galaxyIndex].Cells, cell)
	}
}

func (galaxy *Galaxy) AddCell(cell Cell) {
	galaxy.Cells = append(galaxy.Cells, cell)
}

func (galaxy *Galaxy) FindCell(cell Cell) (int, bool) {
	for i, c := range galaxy.Cells {
		if c.X == cell.X && c.Y == cell.Y {
			return i, true
		}
	}
	return -1, false
}
