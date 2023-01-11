package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Coord struct {
	x, y int
}

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	grid, directions := getGridAndDirections()
	result := walkGridAndGetPassword(grid, directions)
	return result
}

func part2() int {
	grid, directions := getGridAndDirections()
	result := walkCubeAndGetPassword(grid, directions)
	return result
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getGridAndDirections() (grid [][]rune, directions []int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()

		if len(line) == 0 {
			inputScanner.Scan()
			directionsLine := inputScanner.Text()
			directions = getDirections(directionsLine)
			break
		}

		grid = append(grid, make([]rune, len(line)))
		for i, char := range line {
			grid[len(grid)-1][i] = char
		}
	}
	inputFile.Close()
	return
}

func getDirections(directionsLine string) (directions []int) {
	numberStr := ""
	for _, char := range directionsLine {
		if char == 'R' || char == 'L' {
			var rotation int
			switch char {
			case 'R':
				rotation++
			case 'L':
				rotation--
			}
			number, _ := strconv.Atoi(numberStr)
			directions = append(directions, number, rotation)
			numberStr = ""
		} else {
			numberStr += string(char)
		}
	}
	number, _ := strconv.Atoi(numberStr)
	directions = append(directions, number)
	return
}

func walkGridAndGetPassword(grid [][]rune, directions []int) int {
	pos := getStartingPos(grid)
	rotation := 0 // right
	for i, direction := range directions {
		if i%2 == 0 { // move
			for direction > 0 {
				if rotation == 0 { // right
					direction = goRight(pos, grid, direction)
				} else if rotation == 1 { // down
					direction = goDown(pos, grid, direction)
				} else if rotation == 2 { // left
					direction = goLeft(pos, grid, direction)
				} else { // up
					direction = goUp(pos, grid, direction)
				}
				grid[(*pos).y][(*pos).x] = 'O'
			}
		} else { // rotate
			rotation += direction
			switch rotation {
			case -1:
				rotation = 3
			case 4:
				rotation = 0
			}
		}
	}

	return 1000*((*pos).y+1) + 4*((*pos).x+1) + rotation
}

func getStartingPos(grid [][]rune) *Coord {
	for i, value := range grid[0] {
		if value == '.' {
			return &Coord{i, 0}
		}
	}
	panic("no start pos")
}

func goRight(pos *Coord, grid [][]rune, direction int) int {
	if (*pos).x+1 >= len(grid[(*pos).y]) { // wrap
		oldX := (*pos).x
		(*pos).x = 0
		for {
			if isSpaceWall(grid, pos) {
				(*pos).x = oldX
				break
			}
			if isSpaceFree(grid, pos) {
				break
			}
			(*pos).x++
		}
		if oldX == (*pos).x {
			return 0
		}
	} else if grid[(*pos).y][(*pos).x+1] == '#' { // wall
		return 0
	} else { // go
		(*pos).x++
	}

	direction--
	return direction
}

func goDown(pos *Coord, grid [][]rune, direction int) int {
	if (*pos).y+1 > getMaxRowIndexForX(grid, (*pos).x) || grid[(*pos).y+1][(*pos).x] == ' ' { // wrap
		oldY := (*pos).y
		(*pos).y = 0
		for {
			if isSpaceWall(grid, pos) {
				(*pos).y = oldY
				break
			}
			if isSpaceFree(grid, pos) {
				break
			}
			(*pos).y++
		}
		if oldY == (*pos).y {
			return 0
		}
	} else if grid[(*pos).y+1][(*pos).x] == '#' { // wall
		return 0
	} else { // go
		(*pos).y++
	}

	direction--
	return direction
}

func goLeft(pos *Coord, grid [][]rune, direction int) int {
	if (*pos).x-1 < 0 || grid[(*pos).y][(*pos).x-1] == ' ' { // wrap
		oldX := (*pos).x
		(*pos).x = len(grid[(*pos).y]) - 1
		for {
			if isSpaceWall(grid, pos) {
				(*pos).x = oldX
				break
			}
			if isSpaceFree(grid, pos) {
				break
			}
			(*pos).x--
		}
		if oldX == (*pos).x {
			return 0
		}
	} else if grid[(*pos).y][(*pos).x-1] == '#' { // wall
		return 0
	} else { // go
		(*pos).x--
	}

	direction--
	return direction
}

func goUp(pos *Coord, grid [][]rune, direction int) int {
	if (*pos).y-1 < 0 || grid[(*pos).y-1][(*pos).x] == ' ' { // wrap
		oldY := (*pos).y
		(*pos).y = getMaxRowIndexForX(grid, (*pos).x)
		for {
			if isSpaceWall(grid, pos) {
				(*pos).y = oldY
				break
			}
			if isSpaceFree(grid, pos) {
				break
			}
			(*pos).y--
		}
		if oldY == (*pos).y {
			return 0
		}
	} else if grid[(*pos).y-1][(*pos).x] == '#' { // wall
		return 0
	} else { // go
		(*pos).y--
	}

	direction--
	return direction
}

func getMaxRowIndexForX(grid [][]rune, x int) int {
	for y := len(grid) - 1; y >= 0; y-- {
		if x < len(grid[y]) {
			return y
		}
	}
	panic("can not wrap")
}

func isSpaceFree(grid [][]rune, pos *Coord) bool {
	return isSpace(grid, pos, '.') || isSpace(grid, pos, 'O')
}

func isSpaceWall(grid [][]rune, pos *Coord) bool {
	return isSpace(grid, pos, '#')
}

func isSpace(grid [][]rune, pos *Coord, value rune) bool {
	return grid[(*pos).y][(*pos).x] == value
}

func walkCubeAndGetPassword(grid [][]rune, directions []int) int {
	pos := getStartingPos(grid)
	rotation := 0 // right
	for i, direction := range directions {
		if i%2 == 0 { // move
			for direction > 0 {
				if rotation == 0 { // right
					direction, rotation = goRightCube(pos, grid, direction, rotation)
				} else if rotation == 1 { // down
					direction, rotation = goDownCube(pos, grid, direction, rotation)
				} else if rotation == 2 { // left
					direction, rotation = goLeftCube(pos, grid, direction, rotation)
				} else { // up
					direction, rotation = goUpCube(pos, grid, direction, rotation)
				}
				grid[(*pos).y][(*pos).x] = 'O'
			}
		} else { // rotate
			rotation += direction
			switch rotation {
			case -1:
				rotation = 3
			case 4:
				rotation = 0
			}
		}
	}

	return 1000*((*pos).y+1) + 4*((*pos).x+1) + rotation
}

func goRightCube(pos *Coord, grid [][]rune, direction, rotation int) (int, int) {
	if (*pos).x+1 >= len(grid[(*pos).y]) { // wrap
		oldR := rotation
		oldX := (*pos).x
		oldY := (*pos).y

		if (*pos).y < 50 {
			rotation = 2
			(*pos).x = 99
			(*pos).y = 100 + (50 - oldY)
		} else if (*pos).y < 100 {
			rotation = 3
			(*pos).x = 100 + (oldY - 50)
			(*pos).y = 49
		} else if (*pos).y < 150 {
			rotation = 2
			(*pos).x = 149
			(*pos).y = 50 - (oldY - 100)
		} else {
			rotation = 3
			(*pos).x = 50 + (oldY - 150)
			(*pos).y = 149
		}

		if isSpaceWall(grid, pos) {
			(*pos).x = oldX
			(*pos).y = oldY
			return 0, oldR
		}
	} else if grid[(*pos).y][(*pos).x+1] == '#' { // wall
		return 0, rotation
	} else { // go
		(*pos).x++
	}

	direction--
	return direction, rotation
}

func goDownCube(pos *Coord, grid [][]rune, direction, rotation int) (int, int) {
	if (*pos).y+1 > getMaxRowIndexForX(grid, (*pos).x) || grid[(*pos).y+1][(*pos).x] == ' ' { // wrap
		oldR := rotation
		oldX := (*pos).x
		oldY := (*pos).y

		if (*pos).x < 50 {
			rotation = 1
			(*pos).x = 100 + oldX
			(*pos).y = 0
		} else if (*pos).x < 100 {
			rotation = 2
			(*pos).x = 49
			(*pos).y = 100 + oldX
		} else {
			rotation = 2
			(*pos).x = 99
			(*pos).y = oldX - 50
		}

		if isSpaceWall(grid, pos) {
			(*pos).x = oldX
			(*pos).y = oldY
			return 0, oldR
		}
	} else if grid[(*pos).y+1][(*pos).x] == '#' { // wall
		return 0, rotation
	} else { // go
		(*pos).y++
	}

	direction--
	return direction, rotation
}

func goLeftCube(pos *Coord, grid [][]rune, direction, rotation int) (int, int) {
	if (*pos).x-1 < 0 || grid[(*pos).y][(*pos).x-1] == ' ' { // wrap
		oldR := rotation
		oldX := (*pos).x
		oldY := (*pos).y

		if (*pos).y < 50 {
			rotation = 0
			(*pos).x = 0
			(*pos).y = (49 - oldY) + 100
		} else if (*pos).y < 100 {
			rotation = 1
			(*pos).x = oldY - 50
			(*pos).y = 100
		} else if (*pos).y < 150 {
			rotation = 0
			(*pos).x = 50
			(*pos).y = (50 - (oldY - 99))
		} else {
			rotation = 1
			(*pos).x = oldY - 100
			(*pos).y = 0
		}

		if isSpaceWall(grid, pos) {
			(*pos).x = oldX
			(*pos).y = oldY
			return 0, oldR
		}
	} else if grid[(*pos).y][(*pos).x-1] == '#' { // wall
		return 0, rotation
	} else { // go
		(*pos).x--
	}

	direction--
	return direction, rotation
}

func goUpCube(pos *Coord, grid [][]rune, direction, rotation int) (int, int) {
	if (*pos).y-1 < 0 || grid[(*pos).y-1][(*pos).x] == ' ' { // wrap
		oldR := rotation
		oldX := (*pos).x
		oldY := (*pos).y

		if (*pos).x < 50 {
			rotation = 0
			(*pos).x = 50
			(*pos).y = oldX + 50
		} else if (*pos).x < 100 {
			rotation = 0
			(*pos).x = 0
			(*pos).y = oldX + 100
		} else {
			rotation = 3
			(*pos).x = oldX - 100
			(*pos).y = 199
		}

		if isSpaceWall(grid, pos) {
			(*pos).x = oldX
			(*pos).y = oldY
			return 0, oldR
		}
	} else if grid[(*pos).y-1][(*pos).x] == '#' { // wall
		return 0, rotation
	} else { // go
		(*pos).y--
	}

	direction--
	return direction, rotation
}
