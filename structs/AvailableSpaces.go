package structs

import "fmt"

type AvailableSpaces struct {
	Cells []Cell
}

func (availableSpaces *AvailableSpaces) FillSpaces(width, height int) {
	for i := range width*2 - 1 {
		for j := range height*2 - 1 {
			availableSpaces.AddSpace(Cell{X: i, Y: j})
		}
	}
}

func (availableSpaces *AvailableSpaces) AddSpace(cell Cell) {
	availableSpaces.Cells = append(availableSpaces.Cells, cell)
}

func (availableSpaces *AvailableSpaces) FindSpaceWithIndex(index int) (Cell, error) {
	if index >= 0 && index < len(availableSpaces.Cells) {
		return availableSpaces.Cells[index], nil
	}
	var cell Cell
	return cell, fmt.Errorf("out of bounds")
}

func (availableSpaces *AvailableSpaces) FindSpace(targetCell Cell) (int, bool) {
	for i, cell := range availableSpaces.Cells {
		if cell.X == targetCell.X && cell.Y == targetCell.Y {
			return i, true
		}
	}
	return -1, false
}

func (availableSpaces *AvailableSpaces) RemoveSpace(targetCell Cell) bool {
	if i, found := availableSpaces.FindSpace(targetCell); found {
		availableSpaces.Cells = append(availableSpaces.Cells[:i], availableSpaces.Cells[i+1:]...)
		return true
	}
	return false
}
