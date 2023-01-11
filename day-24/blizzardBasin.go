package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Blizzard struct {
	x, y, dir int
}

type Coord struct {
	x, y int
}

type CoordAndCycle struct {
	coord Coord
	cycle int
}

var cyclicFactor int
var bounds Coord
var flipStartAndGoal bool = false

func main() {
	start := time.Now()
	partOneResult := part1()
	fmt.Println("Part 1:", partOneResult, " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(partOneResult), " (", time.Since(start), ")")
}

func part1() int {
	blizzards := getBlizzards()
	result := findShortestPath(blizzards)
	return result
}

func part2(partOneResult int) int {
	blizzards := getBlizzards()

	firstPath := partOneResult
	for i := 0; i < firstPath; i++ {
		blizzards = calculateBlizzardsAtNextMove(blizzards)
	}

	flipStartAndGoal = true
	secondPath := findShortestPath(blizzards)
	for i := 0; i < secondPath; i++ {
		blizzards = calculateBlizzardsAtNextMove(blizzards)
	}

	flipStartAndGoal = false
	thirdPath := findShortestPath(blizzards)

	return firstPath + secondPath + thirdPath
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getBlizzards() map[Blizzard]bool {
	blizzards := map[Blizzard]bool{}
	inputFile, inputScanner := getInputScanner()
	y, xMax := -1, 0
	for inputScanner.Scan() {
		line := inputScanner.Text()
		if y == -1 {
			xMax = len(line) - 3
		} else {
			for x, val := range line {
				dir := -1
				switch val {
				case '>':
					dir = 0
				case 'v':
					dir = 1
				case '<':
					dir = 2
				case '^':
					dir = 3
				}
				if dir != -1 {
					blizzards[Blizzard{x - 1, y, dir}] = true
				}
			}
		}
		y++
	}
	inputFile.Close()
	bounds = Coord{xMax, y - 2}
	cyclicFactor = getLeastCommonMultiple(bounds.x+1, bounds.y+1)
	return blizzards
}

func findShortestPath(blizzards map[Blizzard]bool) int {
	var me Coord
	if !flipStartAndGoal {
		me = Coord{0, -1}
	} else {
		me = Coord{bounds.x, bounds.y + 1}
	}
	visitedSpaceAtCycleToMinute := map[CoordAndCycle]int{} // use cyclic nature of blizzard momevemts, end traversal when space already visited under same conditions but sooner
	return traverse(me, blizzards, -1, 1000, visitedSpaceAtCycleToMinute)
}

func traverse(me Coord, blizzards map[Blizzard]bool, minutes, shortestPathSoFar int, visitedSpaceAtCycleToMinute map[CoordAndCycle]int) int {
	minutes++
	// printGrind(me, blizzards)

	// will reach goal in next move
	if flipStartAndGoal {
		if me.x == 0 && me.y == 0 {
			minutes++
			return minutes
		}
	} else {
		if me.x == bounds.x && me.y == bounds.y {
			minutes++
			return minutes
		}
	}

	// shorter path exists
	if shortestPathSoFar != -1 && minutes > shortestPathSoFar {
		return -1
	}

	cacheKey := CoordAndCycle{me, minutes % cyclicFactor}
	cacheValue, ok := visitedSpaceAtCycleToMinute[cacheKey]
	if ok {
		if cacheValue <= minutes { // already been here, faster or same time
			return -1
		}
	}
	visitedSpaceAtCycleToMinute[cacheKey] = minutes

	blizzards = calculateBlizzardsAtNextMove(blizzards)

	if canMoveRight(me, blizzards) {
		shortestPathForTraversal := traverse(Coord{me.x + 1, me.y}, blizzards, minutes, shortestPathSoFar, visitedSpaceAtCycleToMinute)
		if shouldUpdateShortestPath(shortestPathForTraversal, shortestPathSoFar) {
			shortestPathSoFar = shortestPathForTraversal
		}
	}
	if canMoveDown(me, blizzards) {
		shortestPathForTraversal := traverse(Coord{me.x, me.y + 1}, blizzards, minutes, shortestPathSoFar, visitedSpaceAtCycleToMinute)
		if shouldUpdateShortestPath(shortestPathForTraversal, shortestPathSoFar) {
			shortestPathSoFar = shortestPathForTraversal
		}
	}
	if canMoveLeft(me, blizzards) {
		shortestPathForTraversal := traverse(Coord{me.x - 1, me.y}, blizzards, minutes, shortestPathSoFar, visitedSpaceAtCycleToMinute)
		if shouldUpdateShortestPath(shortestPathForTraversal, shortestPathSoFar) {
			shortestPathSoFar = shortestPathForTraversal
		}
	}
	if canMoveUp(me, blizzards) {
		shortestPathForTraversal := traverse(Coord{me.x, me.y - 1}, blizzards, minutes, shortestPathSoFar, visitedSpaceAtCycleToMinute)
		if shouldUpdateShortestPath(shortestPathForTraversal, shortestPathSoFar) {
			shortestPathSoFar = shortestPathForTraversal
		}
	}
	if canWait(me, blizzards) {
		shortestPathForTraversal := traverse(Coord{me.x, me.y}, blizzards, minutes, shortestPathSoFar, visitedSpaceAtCycleToMinute)
		if shouldUpdateShortestPath(shortestPathForTraversal, shortestPathSoFar) {
			shortestPathSoFar = shortestPathForTraversal
		}
	}

	return shortestPathSoFar
}

func getLeastCommonMultiple(x, y int) int {
	for i := 1; i <= x*y; i++ {
		if i%x == 0 && i%y == 0 {
			return i
		}
	}
	panic("cant find LSM")
}

func shouldUpdateShortestPath(shortestPathForTraversal, shortestPathSoFar int) bool {
	if shortestPathForTraversal == -1 {
		return false
	}
	return shortestPathSoFar == -1 || shortestPathForTraversal < shortestPathSoFar
}

func calculateBlizzardsAtNextMove(blizzards map[Blizzard]bool) map[Blizzard]bool {
	blizzardsAtNextMove := map[Blizzard]bool{}
	for blizzard, _ := range blizzards {
		var blizzardAtNextMove Blizzard

		dir := blizzard.dir
		if dir == 0 {
			newX := blizzard.x + 1
			if newX > bounds.x {
				newX = 0
			}
			blizzardAtNextMove = Blizzard{newX, blizzard.y, dir}
		} else if dir == 1 {
			newY := blizzard.y + 1
			if newY > bounds.y {
				newY = 0
			}
			blizzardAtNextMove = Blizzard{blizzard.x, newY, dir}
		} else if dir == 2 {
			newX := blizzard.x - 1
			if newX < 0 {
				newX = bounds.x
			}
			blizzardAtNextMove = Blizzard{newX, blizzard.y, dir}
		} else if dir == 3 {
			newY := blizzard.y - 1
			if newY < 0 {
				newY = bounds.y
			}
			blizzardAtNextMove = Blizzard{blizzard.x, newY, dir}
		} else {
			panic("unkown direction")
		}

		blizzardsAtNextMove[blizzardAtNextMove] = true
	}

	return blizzardsAtNextMove
}

func canMoveRight(me Coord, blizzards map[Blizzard]bool) bool {
	if me.y == -1 || me.y == bounds.y+1 { // start pos
		return false
	}
	if me.x+1 > bounds.x {
		return false
	}

	me.x++
	return !coordOverlapsWithBlizzard(me, blizzards)
}

func canMoveDown(me Coord, blizzards map[Blizzard]bool) bool {
	if me.y+1 > bounds.y {
		return false
	}

	me.y++
	return !coordOverlapsWithBlizzard(me, blizzards)
}

func canMoveLeft(me Coord, blizzards map[Blizzard]bool) bool {
	if me.y == -1 || me.y == bounds.y+1 { // start pos
		return false
	}
	if me.x-1 < 0 {
		return false
	}

	me.x--
	return !coordOverlapsWithBlizzard(me, blizzards)
}

func canMoveUp(me Coord, blizzards map[Blizzard]bool) bool {
	if me.y-1 < 0 {
		return false
	}

	me.y--
	return !coordOverlapsWithBlizzard(me, blizzards)
}

func canWait(me Coord, blizzards map[Blizzard]bool) bool {
	if me.y == -1 || me.y == bounds.y+1 { // start pos
		return true
	}
	return !coordOverlapsWithBlizzard(me, blizzards)
}

func coordOverlapsWithBlizzard(coord Coord, blizzards map[Blizzard]bool) bool {
	for dir := 0; dir < 4; dir++ {
		if _, ok := blizzards[Blizzard{coord.x, coord.y, dir}]; ok {
			return true
		}
	}
	return false
}

func printGrind(me Coord, blizzards map[Blizzard]bool) {
	for y := 0; y <= bounds.y; y++ {
		for x := 0; x <= bounds.x; x++ {
			if me.x == x && me.y == y {
				fmt.Print("E")
			} else {
				if coordOverlapsWithBlizzard(Coord{x, y}, blizzards) {
					fmt.Print("@")
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Print("\n")
	}
}
