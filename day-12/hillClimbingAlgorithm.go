package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

func part1() int {
	start := time.Now()
	grid, startCoord, endCoord := parseGrid()
	result := calculateShortestPathSize(grid, startCoord, endCoord)
	elapsed := time.Since(start)
	fmt.Println("part1", elapsed)
	return result
}

func parseGrid() (grid [][]rune, startCoord, endCoord [2]int) {
	inputFile, inputScanner := getInputScanner()
	y := 0
	for inputScanner.Scan() {
		grid = append(grid, []rune{})
		for x, value := range inputScanner.Text() {
			switch value {
			case 'S':
				startCoord = [2]int{y, x}
				grid[y] = append(grid[y], 'a')
			case 'E':
				endCoord = [2]int{y, x}
				grid[y] = append(grid[y], 'z')
			default:
				grid[y] = append(grid[y], value)
			}
		}
		y++
	}
	inputFile.Close()
	return
}

func calculateShortestPathSize(grid [][]rune, startCoord, endCoord [2]int) int {
	shortestPathGrid := createShortestPathGrid(grid)
	calcualteForPosition(0, startCoord[0], startCoord[1], grid, &shortestPathGrid, startCoord, endCoord)
	return shortestPathGrid[endCoord[0]][endCoord[1]]
}

func createShortestPathGrid(grid [][]rune) (shortestPathGrid [][]int) {
	for i, row := range grid {
		shortestPathGrid = append(shortestPathGrid, []int{})
		for range row {
			shortestPathGrid[i] = append(shortestPathGrid[i], 0)
		}
	}
	return
}

func calcualteForPosition(stepCount, y, x int, grid [][]rune, shortestPathGrid *[][]int, startCoord, endCoord [2]int) {
	currentShortest := (*shortestPathGrid)[y][x]
	if currentShortest == 0 || currentShortest > stepCount {
		(*shortestPathGrid)[y][x] = stepCount
	} else {
		return
	}

	if y == endCoord[0] && x == endCoord[1] {
		return
	}

	currentHeight := grid[y][x]
	nextStepCount := stepCount + 1

	//up
	if y != 0 {
		upHeight := grid[y-1][x]
		if currentHeight-upHeight >= -1 {
			calcualteForPosition(nextStepCount, y-1, x, grid, shortestPathGrid, startCoord, endCoord)
		}
	}

	//right
	if x != len(grid[0])-1 {
		rightHeight := grid[y][x+1]
		if currentHeight-rightHeight >= -1 {
			calcualteForPosition(nextStepCount, y, x+1, grid, shortestPathGrid, startCoord, endCoord)
		}
	}

	//down
	if y != len(grid)-1 {
		downHeight := grid[y+1][x]
		if currentHeight-downHeight >= -1 {
			calcualteForPosition(nextStepCount, y+1, x, grid, shortestPathGrid, startCoord, endCoord)
		}
	}

	//left
	if x != 0 {
		leftHeight := grid[y][x-1]
		if currentHeight-leftHeight >= -1 {
			calcualteForPosition(nextStepCount, y, x-1, grid, shortestPathGrid, startCoord, endCoord)
		}
	}
}

func part2() int {
	start := time.Now()
	grid, _, endCoord := parseGrid()
	result := calculateShortestPathSizeWithAnyValidStartPoint(grid, endCoord)
	elapsed := time.Since(start)
	fmt.Println("part2", elapsed)
	return result
}

func calculateShortestPathSizeWithAnyValidStartPoint(grid [][]rune, endCoord [2]int) int {
	startingPoints := getAllPossibleStartingPoints(grid)
	paths := []int{}
	for _, startingPoint := range startingPoints {
		path := calculateShortestPathSize(grid, startingPoint, endCoord)
		if path != 0 {
			paths = append(paths, path)
		}
	}
	sort.Ints(paths)
	return paths[0]
}

func getAllPossibleStartingPoints(grid [][]rune) (startingPoints [][2]int) {
	for y, row := range grid {
		for x, value := range row {
			if value == 'a' {
				startingPoints = append(startingPoints, [2]int{y, x})
			}
		}
	}
	return
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	//inputFile, _ := os.Open("input-sample.txt")
	inputFile, _ := os.Open("input.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func main() {
	fmt.Println("Part 1:", part1())
	fmt.Println("Part 2:", part2())
}
