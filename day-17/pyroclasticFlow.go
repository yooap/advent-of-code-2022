package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Coord struct {
	x, y int
}

type Shape struct {
	shapeIndex       int
	pointsOfInterest []Coord
}

var referencePointsOfInterestHorizontalShape = []Coord{{0, 0}, {1, 0}, {2, 0}, {3, 0}}
var referencePointsOfInterestPlusShape = []Coord{{1, 0}, {0, 1}, {2, 1}, {1, 2}}
var referencePointsOfInterestReverseLShape = []Coord{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}}
var referencePointsOfInterestVerticalShape = []Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
var referencePointsOfInterestSquareShape = []Coord{{0, 0}, {1, 0}, {0, 1}, {1, 1}}

var relevantPushLeftPointIndexesHorizontalShape = []int{0}
var relevantPushLeftPointIndexesPlusShape = []int{1, 0}
var relevantPushLeftPointIndexesReverseLShape = []int{0}
var relevantPushLeftPointIndexesVerticalShape = []int{0, 1, 2, 3}
var relevantPushLeftPointIndexesSquareShape = []int{0, 2}

var relevantPushRightPointIndexesHorizontalShape = []int{3}
var relevantPushRightPointIndexesPlusShape = []int{2, 0}
var relevantPushRightPointIndexesReverseLShape = []int{2, 3, 4}
var relevantPushRightPointIndexesVerticalShape = []int{0, 1, 2, 3}
var relevantPushRightPointIndexesSquareShape = []int{1, 3}

var relevantDropPointIndexesHorizontalShape = []int{0, 1, 2, 3}
var relevantDropPointIndexesPlusShape = []int{0, 1, 2}
var relevantDropPointIndexesReverseLShape = []int{0, 1, 2}
var relevantDropPointIndexesVerticalShape = []int{0}
var relevantDropPointIndexesSquareShape = []int{0, 1}

var shapeInitData = [5][]Coord{
	referencePointsOfInterestHorizontalShape,
	referencePointsOfInterestPlusShape,
	referencePointsOfInterestReverseLShape,
	referencePointsOfInterestVerticalShape,
	referencePointsOfInterestSquareShape,
}

var shapePushLeftRelevantIndexes = [5][]int{
	relevantPushLeftPointIndexesHorizontalShape,
	relevantPushLeftPointIndexesPlusShape,
	relevantPushLeftPointIndexesReverseLShape,
	relevantPushLeftPointIndexesVerticalShape,
	relevantPushLeftPointIndexesSquareShape,
}

var shapePushRightRelevantIndexes = [5][]int{
	relevantPushRightPointIndexesHorizontalShape,
	relevantPushRightPointIndexesPlusShape,
	relevantPushRightPointIndexesReverseLShape,
	relevantPushRightPointIndexesVerticalShape,
	relevantPushRightPointIndexesSquareShape,
}

var shapeDropRelevantIndexes = [5][]int{
	relevantDropPointIndexesHorizontalShape,
	relevantDropPointIndexesPlusShape,
	relevantDropPointIndexesReverseLShape,
	relevantDropPointIndexesVerticalShape,
	relevantDropPointIndexesSquareShape,
}

const xMax = 7

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	pushes := getPushes()
	height := run(pushes, 2022)
	return height
}

func part2() int {
	pushes := getPushes()
	height := run(pushes, 1000000000000)
	return height
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getPushes() []rune {
	inputFile, inputScanner := getInputScanner()
	defer func() {
		inputFile.Close()
	}()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		return []rune(line)
	}
	return nil
}

func run(pushes []rune, maxShapeCount int) (height int) {
	grid := make([][]bool, xMax)
	for x, _ := range grid {
		grid[x] = []bool{}
	}

	pushCount := 0
	var heightAtStartOfLoop, shapeCountAtStartOfLoop, shapeIndexAtStartOfLoop, loopHeight, loopShapeCount, sumOfAllLoopIterations int // loop cache
	for shapeCount := 0; shapeCount < maxShapeCount; shapeCount++ {
		shapeIndex := shapeCount % len(shapeInitData)

		const prePush bool = false
		var shape *Shape
		if prePush {
			// pre push 3 moves
			push1 := pushAsInt(pushes[pushCount%len(pushes)])
			push2 := pushAsInt(pushes[(pushCount+1)%len(pushes)])
			push3 := pushAsInt(pushes[(pushCount+2)%len(pushes)])
			pushCount += 3
			shape = spawnWithPreCalcualtedStep(shapeIndex, height, push1, push2, push3)
		} else {
			shape = spawn(shapeIndex, height)
		}

		for {
			// detect loop
			if pushCount%len(pushes) == 0 && pushCount != 0 {
				if heightAtStartOfLoop == 0 {
					heightAtStartOfLoop = height
					shapeCountAtStartOfLoop = shapeCount
					shapeIndexAtStartOfLoop = shapeIndex
				} else if shapeIndex == shapeIndexAtStartOfLoop {
					// fast foward the loops
					loopHeight = height - heightAtStartOfLoop
					loopShapeCount = shapeCount - shapeCountAtStartOfLoop
					timesToLoop := (maxShapeCount - shapeCount) / loopShapeCount
					shapeCount += loopShapeCount * timesToLoop
					sumOfAllLoopIterations += loopHeight * timesToLoop
				}
			}

			// push
			pushDirection := pushes[pushCount%len(pushes)]
			xPush := pushAsInt(pushDirection)
			if canPushX(shape, grid, xPush) {
				push(shape, xPush, 0)
			}
			// tryPush(shape, grid, xPush, 0)

			pushCount++

			// drop
			if canPushDown(shape, grid) {
				push(shape, 0, -1)
			} else {
				grid, height = drawShapeOnGrid(shape, grid, height)
				break
			}
			// if !tryPush(shape, grid, 0, -1) {
			// 	grid = drawShapeOnGrid(shape, grid)
			// 	height = len(grid[0])
			// 	break
			// }
		}
	}

	// printGrid(grid)
	return height + sumOfAllLoopIterations
}

func pushAsInt(direction rune) int {
	if direction == '>' {
		return 1
	} else {
		return -1
	}
}

func spawn(shapeIndex int, height int) *Shape {
	pointRef := shapeInitData[shapeIndex]
	points := []Coord{}
	for _, point := range pointRef {
		points = append(points, Coord{point.x, point.y})
	}
	shape := Shape{shapeIndex, points}
	push(&shape, 2, height+3)
	return &shape
}

func spawnWithPreCalcualtedStep(shapeIndex int, height int, push1, push2, push3 int) *Shape {
	pointRef := shapeInitData[shapeIndex]
	points := make([]Coord, len(pointRef))

	x := 0
	push1and2 := push1 + push2
	if push1and2 == 0 {
		x = 2 + push3
	} else if push1and2 == -2 {
		if push3 == 1 {
			x = 1
		}
	} else if push1and2 == 2 {
		if shapeIndex == 0 {
			if push3 == 1 {
				x = 3
			} else {
				x = 2
			}
		} else if shapeIndex == 1 || shapeIndex == 2 {
			if push3 == 1 {
				x = 4
			} else {
				x = 3
			}
		} else {
			x = 4 + push3
		}
	}

	for i, point := range pointRef {
		points[i] = Coord{point.x + x, point.y + height}
	}
	shape := Shape{shapeIndex, points}
	return &shape
}

func push(shape *Shape, x, y int) {
	for i, coord := range (*shape).pointsOfInterest {
		(*shape).pointsOfInterest[i] = Coord{coord.x + x, coord.y + y}
	}
}

func canPushX(shape *Shape, grid [][]bool, x int) bool {
	if x == -1 {
		return canPushLeft(shape, grid)
	} else {
		return canPushRight(shape, grid)
	}
}

func canPushLeft(shape *Shape, grid [][]bool) bool {
	for _, i := range shapePushLeftRelevantIndexes[(*shape).shapeIndex] {
		coord := (*shape).pointsOfInterest[i]
		coord = Coord{coord.x - 1, coord.y}
		if !validateNewCoord(coord, grid) {
			return false
		}
	}
	return true
}

func canPushRight(shape *Shape, grid [][]bool) bool {
	for _, i := range shapePushRightRelevantIndexes[(*shape).shapeIndex] {
		coord := (*shape).pointsOfInterest[i]
		coord = Coord{coord.x + 1, coord.y}
		if !validateNewCoord(coord, grid) {
			return false
		}
	}
	return true
}

func canPushDown(shape *Shape, grid [][]bool) bool {
	for _, i := range shapeDropRelevantIndexes[(*shape).shapeIndex] {
		coord := (*shape).pointsOfInterest[i]
		coord = Coord{coord.x, coord.y - 1}
		if !validateNewCoord(coord, grid) {
			return false
		}
	}
	return true
}

func tryPush(shape *Shape, grid [][]bool, x, y int) bool {
	updatedPointsOfInterest := make([]Coord, len((*shape).pointsOfInterest))
	for i, coord := range (*shape).pointsOfInterest {
		coord = Coord{coord.x + x, coord.y + y}
		if !validateNewCoord(coord, grid) {
			return false
		} else {
			updatedPointsOfInterest[i] = coord
		}
	}
	(*shape).pointsOfInterest = updatedPointsOfInterest
	return true
}

func validateNewCoord(coord Coord, grid [][]bool) bool {
	if coord.x < 0 || coord.x >= len(grid) { // out of x bounds
		return false
	}
	if coord.y < 0 { // below y lower bound
		return false
	}
	if coord.y >= len(grid[0]) { // above grid height
		return true
	}
	return !grid[coord.x][coord.y]
}

func drawShapeOnGrid(shape *Shape, grid [][]bool, height int) ([][]bool, int) {
	// highest point of shape
	var yMax int
	for _, coord := range (*shape).pointsOfInterest {
		if coord.y > yMax {
			yMax = coord.y
		}
	}

	// extend grid
	for height <= yMax {
		for x, column := range grid {
			grid[x] = append(column, false)
		}
		height++
	}
	// if yMax+1 > len(grid[0]) {
	// 	for x, column := range grid {
	// 		newColumn := make([]bool, yMax+1)
	// 		copy(newColumn, column)
	// 		grid[x] = newColumn
	// 	}
	// }

	// draw
	for _, coord := range (*shape).pointsOfInterest {
		grid[coord.x][coord.y] = true
	}

	return grid, height
}

func printGrid(grid [][]bool) {
	for y := len(grid[0]) - 1; y >= 0; y-- {
		for x := 0; x < len(grid); x++ {
			if grid[x][y] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}

		}
		fmt.Print("\n")
	}
}
