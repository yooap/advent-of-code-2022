package main

import (
	"bufio"
	"fmt"
	"os"
)

func part1() int {
	signal := readInput()
	return solveForMarkerLength(signal, 4)
}

func solveForMarkerLength(signal string, length int) int {
	lastChars := make([]rune, length)
	for i, char := range signal {
		if i < length {
			lastChars[i] = char
		} else {
			lastChars = shiftLeft(lastChars)
			lastChars[len(lastChars)-1] = char
		}

		if i > length-2 && allUnique(lastChars) {
			return i + 1
		}
	}

	panic("not found")
}

func shiftLeft(slice []rune) []rune {
	for i, value := range slice {
		if i == 0 {
			continue
		}
		slice[i-1] = value
	}
	return slice
}

func allUnique(slice []rune) bool {
	m := map[rune]bool{}
	for _, char := range slice {
		if m[char] {
			return false
		}
		m[char] = true
	}
	return true
}

func part2() int {
	signal := readInput()
	return solveForMarkerLength(signal, 14)
}

func readInput() (input string) {
	//inputFile, _ := os.Open("input-sample.txt")
	inputFile, _ := os.Open("input.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)

	for inputScanner.Scan() {
		input = inputScanner.Text()
	}
	inputFile.Close()

	return
}

func main() {
	fmt.Println("Part 1:", part1())
	fmt.Println("Part 2:", part2())
}
