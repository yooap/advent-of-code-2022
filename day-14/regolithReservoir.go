package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const sandSourceX, sandSourceY = 500, 0

func part1() (result int) {
	rockPaths := getRockPaths()
	grid := drawGrid(rockPaths)
	sandUnits := simulateSand(grid)
	return sandUnits
}

func getRockPaths() (rockPaths [][][2]int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		rockPathString := inputScanner.Text()
		rockPathPointsString := strings.Split(rockPathString, " -> ")
		rockPath := [][2]int{}
		for _, rockPathPointString := range rockPathPointsString {
			xy := strings.Split(rockPathPointString, ",")
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			rockPath = append(rockPath, [2]int{x, y})
		}
		rockPaths = append(rockPaths, rockPath)
	}
	inputFile.Close()
	return
}

func drawGrid(rockPaths [][][2]int) (grid [][]rune) {
	xMaxIdx, yMaxIdx := getGridBounds(rockPaths)

	// empty grid
	for x := 0; x < xMaxIdx+1; x++ {
		column := make([]rune, yMaxIdx+1)
		for y := range column {
			column[y] = '.'
		}
		grid = append(grid, column)
	}

	grid[sandSourceX][sandSourceY] = '+'

	for _, rockPath := range rockPaths {
		for i := 1; i < len(rockPath); i++ {
			point1, point2 := rockPath[i-1], rockPath[i]
			for x := point1[0]; shouldContinue(point1[0], point2[0], x); x += getIncrement(point2[0], point1[0]) {
				grid[x][point1[1]] = '#'
			}
			for y := point1[1]; shouldContinue(point1[1], point2[1], y); y += getIncrement(point2[1], point1[1]) {
				grid[point1[0]][y] = '#'
			}
		}
	}

	return
}

func shouldContinue(i1, i2, i int) bool {
	if i1 < i2 {
		return i <= i2
	}
	if i1 > i2 {
		return i >= i2
	}
	return false
}

func getIncrement(i1, i2 int) int {
	subtractionResult := float64(i1 - i2)
	return int(subtractionResult / math.Abs(subtractionResult))
}

func getGridBounds(rockPaths [][][2]int) (xMax, yMax int) {
	for _, rockPath := range rockPaths {
		for _, rockPathCoord := range rockPath {
			x, y := rockPathCoord[0], rockPathCoord[1]
			if x > xMax {
				xMax = x
			}
			if y > yMax {
				yMax = y
			}
		}
	}
	return xMax, yMax
}

func simulateSand(grid [][]rune) (restingSandUnits int) {
	for {
		sandUnit := [2]int{sandSourceX, sandSourceY}
		for {
			// flows into the abyss
			if sandUnit[1]+1 == len(grid[0]) {
				return
			}

			down := grid[sandUnit[0]][sandUnit[1]+1]
			if down == '.' {
				// fall down
				sandUnit = [2]int{sandUnit[0], sandUnit[1] + 1}
				continue
			}
			left := grid[sandUnit[0]-1][sandUnit[1]+1]
			if left == '.' {
				// fall left
				sandUnit = [2]int{sandUnit[0] - 1, sandUnit[1] + 1}
				continue
			}
			right := grid[sandUnit[0]+1][sandUnit[1]+1]
			if right == '.' {
				// fall right
				sandUnit = [2]int{sandUnit[0] + 1, sandUnit[1] + 1}
				continue
			}

			// stops at source
			if sandUnit[0] == sandSourceX && sandUnit[1] == sandSourceY {
				restingSandUnits++
				return
			}

			// rest
			grid[sandUnit[0]][sandUnit[1]] = 'o'
			break
		}
		restingSandUnits++
		// printGrid(grid)
	}
}

func printGrid(grid [][]rune) {
	for y := 0; y < len(grid[0]); y++ {
		for x := 490; x < len(grid); x++ {
			fmt.Print(string(grid[x][y]))
		}
		fmt.Print("\n")
	}
}

func part2() int {
	rockPaths := getRockPaths()
	grid := drawGridWithFloor(rockPaths)
	sandUnits := simulateSand(grid)
	return sandUnits
}

func drawGridWithFloor(rockPaths [][][2]int) (grid [][]rune) {
	xMaxIdx, yMaxIdx := getGridBounds(rockPaths)

	// empty grid
	for x := 0; x < xMaxIdx+yMaxIdx; x++ {
		column := make([]rune, yMaxIdx+3)
		for y := range column {
			column[y] = '.'
		}
		grid = append(grid, column)
	}

	grid[sandSourceX][sandSourceY] = '+'

	for _, rockPath := range rockPaths {
		for i := 1; i < len(rockPath); i++ {
			point1, point2 := rockPath[i-1], rockPath[i]
			for x := point1[0]; shouldContinue(point1[0], point2[0], x); x += getIncrement(point2[0], point1[0]) {
				grid[x][point1[1]] = '#'
			}
			for y := point1[1]; shouldContinue(point1[1], point2[1], y); y += getIncrement(point2[1], point1[1]) {
				grid[point1[0]][y] = '#'
			}
		}
	}

	// floor
	for x := 0; x < len(grid); x++ {
		grid[x][len(grid[x])-1] = '#'
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
