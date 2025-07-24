package main

import "fmt"

type Coordinate struct{ X, Y int }

type AvailableSpaces struct{ Coordinates []Coordinate }

func (availableSpaces *AvailableSpaces) FillSpaces(width, height int) {
	for i := range width*2 - 1 {
		for j := range height*2 - 1 {
			availableSpaces.AddSpace(Coordinate{X: i, Y: j})
		}
	}
}

func (availableSpaces *AvailableSpaces) AddSpace(coordinate Coordinate) {
	availableSpaces.Coordinates = append(availableSpaces.Coordinates, coordinate)
}

var width, height, startingCenterCount int
var availableSpaces AvailableSpaces

func main() {
	width = 5
	height = 5
	startingCenterCount = 3
	availableSpaces.FillSpaces(width, height)

	fmt.Println(availableSpaces)
}
