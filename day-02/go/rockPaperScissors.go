package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var shapeToValueMap = map[string]int{
	"X": 1,
	"Y": 2,
	"Z": 3,
}

var oponentVsMyShapeResultValueMap = map[string]map[string]int{
	"A": { // rock
		"X": 3, // rock
		"Y": 6, // paper
		"Z": 0, // scissors
	},
	"B": { // paper
		"X": 0, // rock
		"Y": 3, // paper
		"Z": 6, // scissors
	},
	"C": { // scissors
		"X": 6, // rock
		"Y": 0, // paper
		"Z": 3, // scissors
	},
}

func part1() (score int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		values := strings.Split(line, " ")
		oponentValue, myValue := values[0], values[1]
		score += shapeToValueMap[myValue]
		score += oponentVsMyShapeResultValueMap[oponentValue][myValue]
	}

	inputFile.Close()

	return
}

var resultToValueMap = map[string]int{
	"X": 0,
	"Y": 3,
	"Z": 6,
}

var oponentShapeAndResultToMyShapeValueMap = map[string]map[string]int{
	"A": { // rock
		"X": 3, // lose with scissors
		"Y": 1, // draw with rock
		"Z": 2, // win with paper
	},
	"B": { // paper
		"X": 1, // lose with rock
		"Y": 2, // draw with paper
		"Z": 3, // win with scissors
	},
	"C": { // scissors
		"X": 2, // lose with paper
		"Y": 3, // draw with scissors
		"Z": 1, // win with rock
	},
}

func part2() (score int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		values := strings.Split(line, " ")
		oponentValue, result := values[0], values[1]
		score += resultToValueMap[result]
		score += oponentShapeAndResultToMyShapeValueMap[oponentValue][result]
	}

	inputFile.Close()

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
