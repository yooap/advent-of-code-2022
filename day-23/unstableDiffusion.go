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

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	elves := getElves()
	result, _ := doRounds(elves, 10)
	return result
}

func part2() int {
	elves := getElves()
	_, rounds := doRounds(elves, -1)
	return rounds
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getElves() map[Coord]bool {
	elves := map[Coord]bool{}
	inputFile, inputScanner := getInputScanner()
	y := 0
	for inputScanner.Scan() {
		line := inputScanner.Text()
		for x, val := range line {
			if val == '#' {
				elves[Coord{x, y}] = true
			}
		}
		y++
	}
	inputFile.Close()
	return elves
}

func doRounds(elves map[Coord]bool, rounds int) (int, int) {
	direction := 0
	hasAnyoneMoved := true
	round := 0
	for (round < rounds || rounds < 0) && hasAnyoneMoved {
		// printGrind(elves)
		// decide
		destinationToOrigin := map[Coord]Coord{}
		destinationsToNotGoTo := map[Coord]bool{}
		hasAnyoneMoved = false
		for elf, _ := range elves {
			destination, shouldMove := determineDestination(elf, elves, direction)
			if shouldMove {
				hasAnyoneMoved = true
				if _, ok := destinationToOrigin[destination]; ok {
					// someone already wants to come here
					destinationsToNotGoTo[destination] = true
				} else {
					destinationToOrigin[destination] = elf
				}
			}
		}

		//move
		for dest, orig := range destinationToOrigin {
			if _, ok := destinationsToNotGoTo[dest]; ok {
				continue
			}
			delete(elves, orig)
			elves[dest] = true
		}

		direction = incrementDirection(direction)
		round++
	}

	// calc result
	n, s, w, e := getBounds(elves)
	spaceCount := (s - n + 1) * (e - w + 1)
	return spaceCount - len(elves), round
}

func incrementDirection(direction int) int {
	switch direction {
	case 0:
		return 1 // S
	case 1:
		return 2 // W
	case 2:
		return 3 // E
	case 3:
		return 0 // N
	}
	panic("bad direction")
}

func determineDestination(elf Coord, elves map[Coord]bool, direction int) (Coord, bool) {
	if isNorthFree(elf, elves) && isSouthFree(elf, elves) && isWestFree(elf, elves) && isEastFree(elf, elves) {
		return Coord{}, false
	}
	for directionChange := 0; directionChange < 4; directionChange++ {
		if direction == 0 { // N
			if isNorthFree(elf, elves) {
				return Coord{elf.x, elf.y - 1}, true
			}
		} else if direction == 1 { // S
			if isSouthFree(elf, elves) {
				return Coord{elf.x, elf.y + 1}, true
			}
		} else if direction == 2 { // W
			if isWestFree(elf, elves) {
				return Coord{elf.x - 1, elf.y}, true
			}
		} else { // E
			if isEastFree(elf, elves) {
				return Coord{elf.x + 1, elf.y}, true
			}
		}
		direction = incrementDirection(direction)
	}
	return Coord{}, false
}

func isNorthFree(elf Coord, elves map[Coord]bool) bool {
	if _, ok := elves[Coord{elf.x - 1, elf.y - 1}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x, elf.y - 1}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x + 1, elf.y - 1}]; ok {
		return false
	}
	return true
}

func isSouthFree(elf Coord, elves map[Coord]bool) bool {
	if _, ok := elves[Coord{elf.x - 1, elf.y + 1}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x, elf.y + 1}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x + 1, elf.y + 1}]; ok {
		return false
	}
	return true
}

func isWestFree(elf Coord, elves map[Coord]bool) bool {
	if _, ok := elves[Coord{elf.x - 1, elf.y - 1}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x - 1, elf.y}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x - 1, elf.y + 1}]; ok {
		return false
	}
	return true
}

func isEastFree(elf Coord, elves map[Coord]bool) bool {
	if _, ok := elves[Coord{elf.x + 1, elf.y - 1}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x + 1, elf.y}]; ok {
		return false
	}
	if _, ok := elves[Coord{elf.x + 1, elf.y + 1}]; ok {
		return false
	}
	return true
}

func getBounds(elves map[Coord]bool) (n, s, w, e int) {
	firstValue := true
	for elf, _ := range elves {
		if firstValue {
			n, s, w, e = elf.y, elf.y, elf.x, elf.x
			firstValue = false
		}
		if elf.y < n {
			n = elf.y
		} else if elf.y > s {
			s = elf.y
		} else if elf.x < w {
			w = elf.x
		} else if elf.x > e {
			e = elf.x
		}
	}
	return
}

func printGrind(elves map[Coord]bool) {
	n, s, w, e := getBounds(elves)
	for y := n; y <= s; y++ {
		for x := w; x <= e; x++ {
			if _, ok := elves[Coord{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}
