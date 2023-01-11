package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func part1() int {
	instructions := getInstructions()
	result := executeInstructions(instructions)
	return result
}

func getInstructions() (instructions []string) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		instruction := inputScanner.Text()
		instructions = append(instructions, instruction)
	}
	inputFile.Close()
	return
}

func executeInstructions(instructions []string) (result int) {
	cycle, x := 0, 1
	for _, instruction := range instructions {
		if strings.HasPrefix(instruction, "noop") {
			cycle, result = runCycle(cycle, x, result)
		} else {
			splitInstruction := strings.Split(instruction, " ")
			xIncrementInNextCycle, _ := strconv.Atoi(splitInstruction[1])
			cycle, result = runCycle(cycle, x, result)
			cycle, result = runCycle(cycle, x, result)
			x += xIncrementInNextCycle
		}
	}
	return
}

func runCycle(cycle, x, result int) (int, int) {
	cycle++

	if (cycle+20)%40 == 0 {
		result += x * cycle
	}

	return cycle, result
}

func part2() string {
	instructions := getInstructions()
	result := drawImage(instructions)
	return result
}

func drawImage(instructions []string) (result string) {
	cycle, x := 0, 1
	for _, instruction := range instructions {
		if strings.HasPrefix(instruction, "noop") {
			cycle, result = runCycleWhileDrawing(cycle, x, result)
		} else {
			splitInstruction := strings.Split(instruction, " ")
			xIncrementInNextCycle, _ := strconv.Atoi(splitInstruction[1])
			cycle, result = runCycleWhileDrawing(cycle, x, result)
			cycle, result = runCycleWhileDrawing(cycle, x, result)
			x += xIncrementInNextCycle
		}
	}
	return
}

func runCycleWhileDrawing(cycle, x int, result string) (int, string) {
	cycle++

	crtPos := (cycle - 1) % 40

	if crtPos == 0 {
		result += "\n"

	}

	if crtPos >= x-1 && crtPos <= x+1 {
		result += "#"
	} else {
		result += "."
	}

	return cycle, result
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
