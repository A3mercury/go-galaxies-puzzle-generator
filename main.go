package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/A3mercury/galaxies-generator/structs"
)

const (
	CELL       = "cell"
	CROSS      = "cross"
	HORIZONTAL = "horizontal"
	VERTICAL   = "vertical"
)

var width, height, startingCenterCount int
var availableSpaces structs.AvailableSpaces
var galaxies structs.Galaxies
var lines structs.Lines

func main() {
	initStartingValues()
	initAvailableSpaces()
	selectStartingCenters()
	expandCenters()
	fillRemainingAvailableGalaxies()
	drawLines()

	fmt.Printf("Available Spaces (should be empty): %v\n", availableSpaces)
	fmt.Printf("Galaxies: %v\n", galaxies)
	fmt.Printf("Puzzles lines: %v\n", lines)
}

func initStartingValues() {
	width = 5
	height = 5
	startingCenterCount = 3
}

func initAvailableSpaces() {
	availableSpaces.FillSpaces(width, height)
}

func selectStartingCenters() {
	for range startingCenterCount {
		chooseFromAvailableSpaces()
	}
}

func fillRemainingAvailableGalaxies() {
	for len(availableSpaces.Cells) > 0 {
		chooseFromAvailableSpaces()
		expandCenters()
	}
}

func chooseFromAvailableSpaces() {
	randIndex := rand.Intn(len(availableSpaces.Cells))
	pick, err := availableSpaces.FindSpaceWithIndex(randIndex)
	if err != nil {
		log.Fatalf("There was a problem getting a random starting center: %v", err)
		return
	}
	startingGalaxy, err := startGalaxy(pick)
	if err != nil {
		log.Fatalf("There was a problem on starting a galaxy: %v", err)
		return
	}
	galaxies.AddGalaxy(startingGalaxy)
	unsetFromAvailableSpaces(pick)
}

func unsetFromAvailableSpaces(cell structs.Cell) {
	if _, found := availableSpaces.FindSpace(cell); !found {
		log.Fatalf("There was a problem finding cell in availableSpaces: %v", cell)
		return
	}
	offset, err := getOffsetsFromCurrent(cell)
	if err != nil {
		log.Fatalf("There was a problem getting offset from current cell: %v %v", cell, err)
		return
	}

	for x := cell.X - offset.X; x <= cell.X+offset.X; x++ {
		for y := cell.Y - offset.Y; y <= cell.Y+offset.Y; y++ {
			if !availableSpaces.RemoveSpace(structs.Cell{X: x, Y: y}) {
				continue
			}
		}
	}
}

func getOffsetsFromCurrent(current structs.Cell) (structs.Cell, error) {
	currentType, err := getCenterType(current)
	if err != nil {
		log.Fatalf("There was a problem getting this cell's center type: %v %v", current, err)
		var cell structs.Cell
		return cell, err
	}
	switch currentType {
	case CELL:
		return structs.Cell{X: 1, Y: 1}, nil
	case HORIZONTAL:
		return structs.Cell{X: 2, Y: 1}, nil
	case VERTICAL:
		return structs.Cell{X: 1, Y: 2}, nil
	case CROSS:
		return structs.Cell{X: 2, Y: 2}, nil
	default:
		var cell structs.Cell
		return cell, fmt.Errorf("something went wrong getting offets from current: %v", current)
	}
}

func getCenterType(cell structs.Cell) (string, error) {
	switch {
	case cell.X%2 == 0 && cell.Y%2 == 0:
		return CELL, nil
	case cell.X%2 != 0 && cell.Y%2 == 0:
		return HORIZONTAL, nil
	case cell.X%2 == 0 && cell.Y%2 != 0:
		return VERTICAL, nil
	case cell.X%2 != 0 && cell.Y%2 != 0:
		return CROSS, nil
	default:
		return "", fmt.Errorf("something went wrong getting center: %v", cell)
	}
}

func startGalaxy(center structs.Cell) (structs.Galaxy, error) {
	galaxyType, err := getCenterType(center)
	if err != nil {
		log.Fatalf("There was a problem starting the Galaxy: %v %v", center, err)
		var g structs.Galaxy
		return g, err
	}
	galaxy := structs.Galaxy{
		Center: structs.Center{
			Type: galaxyType,
			Cell: structs.Cell{
				X: center.X,
				Y: center.Y,
			},
		},
		Complete: false,
	}

	var nearCells []structs.Cell
	switch galaxyType {
	case CELL:
		nearCells = append(nearCells, structs.Cell{X: 0, Y: 0})
	case HORIZONTAL:
		nearCells = append(nearCells, structs.Cell{X: -1, Y: 0}, structs.Cell{X: 1, Y: 0})
	case VERTICAL:
		nearCells = append(nearCells, structs.Cell{X: 0, Y: -1}, structs.Cell{X: 0, Y: 1})
	case CROSS:
		nearCells = append(nearCells,
			structs.Cell{X: -1, Y: -1},
			structs.Cell{X: -1, Y: 1},
			structs.Cell{X: 1, Y: -1},
			structs.Cell{X: 1, Y: 1},
		)
	default:
		break
	}

	for _, cellOffset := range nearCells {
		targetCell := structs.Cell{
			X: galaxy.Center.Cell.X + cellOffset.X,
			Y: galaxy.Center.Cell.Y + cellOffset.Y,
		}
		galaxy.AddCell(targetCell)
	}
	return galaxy, nil
}

func expandCenters() {
	for _, galaxy := range galaxies.Galaxies {
		if galaxy.Complete {
			continue
		}
		center := galaxy.Center.Cell
		searchBoard(center, center, galaxy, 0)
		galaxies.MarkGalaxyComplete(galaxy)
	}
}

func searchBoard(current, center structs.Cell, galaxy structs.Galaxy, count int) {
	if count++; count >= 5 {
		return
	}

	var centerType string = ""
	var err error
	if center == current {
		centerType, err = getCenterType(center)
		if err != nil {
			log.Fatalf("There was a problem getting the center type: %v %v", center, err)
			return
		}
	}

	adjacent := getAdjacentCellsToUse(centerType)
	for k := range len(adjacent) {
		newCell := structs.Cell{
			X: current.X + adjacent[k].X,
			Y: current.Y + adjacent[k].Y,
		}

		if isValidSpace(newCell, center) {
			mirrored := getMirroredCoordinates(newCell, center)
			// TODO: This is not updating the galaxy inside the galaxies struct
			galaxies.AddCell(galaxy, newCell)
			galaxies.AddCell(galaxy, mirrored)

			unsetFromAvailableSpaces(newCell)
			unsetFromAvailableSpaces(mirrored)

			searchBoard(newCell, center, galaxy, count)
		}
	}
}

func getAdjacentCellsToUse(centerType string) []structs.Cell {
	switch centerType {
	case CELL:
		return []structs.Cell{{X: -2, Y: 0}, {X: 0, Y: 2}}
	case HORIZONTAL:
		return []structs.Cell{{X: -3, Y: 0}, {X: -1, Y: 2}, {X: 1, Y: 2}}
	case VERTICAL:
		return []structs.Cell{{X: -2, Y: -1}, {X: -2, Y: 1}, {X: 0, Y: 3}}
	case CROSS:
		return []structs.Cell{{X: -3, Y: -1}, {X: -3, Y: 1}, {X: -1, Y: 3}, {X: 1, Y: 3}}
	default:
		return []structs.Cell{{X: -2, Y: 0}, {X: 0, Y: 2}, {X: 2, Y: 0}, {X: 0, Y: -2}}
	}
}

func isValidSpace(current, center structs.Cell) bool {
	mirrored := getMirroredCoordinates(current, center)
	_, found := availableSpaces.FindSpace(current)
	_, foundMirrored := availableSpaces.FindSpace(mirrored)
	return isInBounds(current) && isInBounds(mirrored) && found && foundMirrored
}

func isInBounds(cell structs.Cell) bool {
	return cell.X >= 0 && cell.X <= width*2-2 && cell.Y >= 0 && cell.Y <= height*2-2
}

func getMirroredCoordinates(current, center structs.Cell) structs.Cell {
	return structs.Cell{
		X: ((current.X - center.X) * -1) + center.X,
		Y: ((current.Y - center.Y) * -1) + center.Y,
	}
}

func drawLines() {
	for _, galaxy := range galaxies.Galaxies {
		for _, cell := range galaxy.Cells {
			for _, adjacent := range getAdjacentCellsToUse("") {
				newCell := structs.Cell{
					X: cell.X + adjacent.X,
					Y: cell.Y + adjacent.Y,
				}

				if isInBounds(newCell) {
					if _, found := galaxy.FindCell(newCell); !found {
						line := structs.Cell{
							X: cell.X + (adjacent.X / 2),
							Y: cell.Y + (adjacent.Y / 2),
						}
						if _, foundLine := lines.FindCell(line); !foundLine {
							lines.AddCell(line)
						}
					}
				}
			}
		}
	}
	lines.SortLines()
}
