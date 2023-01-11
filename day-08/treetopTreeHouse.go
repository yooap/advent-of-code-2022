package main

import (
	"bufio"
	"fmt"
	"os"
)

func part1() int {
	grid := assembleGrid()
	return countVisibleTrees(grid)
}

func assembleGrid() (grid [][]int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		row := []int{}
		for _, height := range line {
			row = append(row, int(height-'0'))
		}
		grid = append(grid, row)
	}
	inputFile.Close()
	return
}

func countVisibleTrees(grid [][]int) (count int) {
	maxRowIndex := len(grid) - 1
	maxColumnIndex := len(grid[0]) - 1
	for i, row := range grid {
		if i == 0 || i == maxRowIndex { // first & last row
			count += len(row)
			continue
		}

		for j, height := range row {
			if j == 0 || j == maxColumnIndex { // first & last column
				count++
				continue
			}

			if isVisible(grid, i, j, height, maxRowIndex, maxColumnIndex) {
				count++
			}
		}
	}
	return
}

func isVisible(grid [][]int, i, j, height, maxRowIndex, maxColumnIndex int) bool {
	top, right, bottom, left := true, true, true, true

	// is visible from top
	for step := 1; i-step >= 0; step++ {
		if grid[i-step][j] >= height {
			top = false
			break
		}
	}
	if top {
		return true
	}

	// is visible from right
	for step := 1; j+step <= maxColumnIndex; step++ {
		if grid[i][j+step] >= height {
			right = false
			break
		}
	}
	if right {
		return true
	}

	// is visible from bottom
	for step := 1; i+step <= maxRowIndex; step++ {
		if grid[i+step][j] >= height {
			bottom = false
			break
		}
	}
	if bottom {
		return true
	}

	// is visible from left
	for step := 1; j-step >= 0; step++ {
		if grid[i][j-step] >= height {
			left = false
			break
		}
	}
	if left {
		return true
	}

	return false
}

func part2() int {
	grid := assembleGrid()
	return getHighestSceanicScore(grid)
}

func getHighestSceanicScore(grid [][]int) (highestScore int) {
	maxRowIndex := len(grid) - 1
	maxColumnIndex := len(grid[0]) - 1
	for i, row := range grid {
		for j, height := range row {
			score := getSceanicScore(grid, i, j, height, maxRowIndex, maxColumnIndex)
			if score > highestScore {
				highestScore = score
			}
		}
	}
	return
}

func getSceanicScore(grid [][]int, i, j, height, maxRowIndex, maxColumnIndex int) int {
	if i == 0 || i == maxRowIndex || j == 0 || j == maxColumnIndex {
		return 0
	}

	top, right, bottom, left := 1, 1, 1, 1

	// top score
	for step := 1; i-step > 0; step++ {
		if grid[i-step][j] >= height {
			break
		}
		top++
	}

	// right score
	for step := 1; j+step < maxColumnIndex; step++ {
		if grid[i][j+step] >= height {
			break
		}
		right++
	}

	// bottom score
	for step := 1; i+step < maxRowIndex; step++ {
		if grid[i+step][j] >= height {
			break
		}
		bottom++
	}

	// left score
	for step := 1; j-step > 0; step++ {
		if grid[i][j-step] >= height {
			break
		}
		left++
	}

	return top * right * bottom * left
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
