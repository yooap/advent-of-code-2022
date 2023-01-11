package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func part1() (pairCount int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		pair := strings.Split(line, ",")
		firstRange, secondRange := splitIntoIntArray(pair[0]), splitIntoIntArray(pair[1])

		if fullyOverlaps(firstRange, secondRange) {
			pairCount++
		}
	}

	inputFile.Close()

	return
}
func splitIntoIntArray(sectionRange string) [2]int {
	split := strings.Split(sectionRange, "-")
	start, _ := strconv.Atoi(split[0])
	end, _ := strconv.Atoi(split[1])
	return [2]int{start, end}
}

func fullyOverlaps(firstRange, secondRange [2]int) bool {
	if firstRange[0] >= secondRange[0] && firstRange[1] <= secondRange[1] {
		return true
	}
	if secondRange[0] >= firstRange[0] && secondRange[1] <= firstRange[1] {
		return true
	}
	return false
}

func part2() (pairCount int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		pair := strings.Split(line, ",")
		firstRange, secondRange := splitIntoIntArray(pair[0]), splitIntoIntArray(pair[1])

		if hasAnyOverlap(firstRange, secondRange) {
			pairCount++
		}
	}

	inputFile.Close()

	return
}

func hasAnyOverlap(firstRange, secondRange [2]int) bool {
	return firstRange[1] >= secondRange[0] && firstRange[0] <= secondRange[1]
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
