package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type moveStruct struct {
	count int
	from  int
	to    int
}

func part1() string {
	stacks := [][]rune{}
	moves := []moveStruct{}

	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "[") {
			for i, char := range line {
				if (i-1)%4 == 0 && !unicode.IsSpace(char) {
					stackIndex := (i - 1) / 4
					if len(stacks) <= stackIndex {
						stacks = extendStackSlice(stackIndex, stacks)
					}
					stacks[stackIndex] = append([]rune{char}, stacks[stackIndex]...)
				}
			}
		} else if strings.HasPrefix(line, "move") {
			moveInstance := createMoveInstance(line)
			moves = append(moves, moveInstance)
		}
	}
	inputFile.Close()

	stacks = doMoves(stacks, moves)
	result := getResult(stacks)

	return result
}

func extendStackSlice(requiredStackIndex int, stacks [][]rune) [][]rune {
	stacksLengthCurrently := len(stacks)
	for i := 0; i <= (requiredStackIndex + 1 - stacksLengthCurrently); i++ {
		stacks = append(stacks, []rune{})
	}
	return stacks
}

func createMoveInstance(line string) moveStruct {
	moveLine := strings.Split(line, " ")
	count, _ := strconv.Atoi(moveLine[1])
	from, _ := strconv.Atoi(moveLine[3])
	to, _ := strconv.Atoi(moveLine[5])
	moveInstance := moveStruct{
		count: count,
		from:  from,
		to:    to}
	return moveInstance
}

func doMoves(stacks [][]rune, moves []moveStruct) [][]rune {
	for _, move := range moves {
		fromIndex, toIndex := move.from-1, move.to-1
		for i := 0; i < move.count; i++ {
			objectToMove := stacks[fromIndex][len(stacks[fromIndex])-1]
			stacks[fromIndex] = stacks[fromIndex][:len(stacks[fromIndex])-1] // remove
			stacks[toIndex] = append(stacks[toIndex], objectToMove)          // add
		}
	}
	return stacks
}

func getResult(stacks [][]rune) (result string) {
	for _, stack := range stacks {
		if len(stack) != 0 {
			result += string(stack[len(stack)-1])
		}
	}
	return
}

func part2() string {
	stacks := [][]rune{}
	moves := []moveStruct{}

	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "[") {
			for i, char := range line {
				if (i-1)%4 == 0 && !unicode.IsSpace(char) {
					stackIndex := (i - 1) / 4
					if len(stacks) <= stackIndex {
						stacks = extendStackSlice(stackIndex, stacks)
					}
					stacks[stackIndex] = append([]rune{char}, stacks[stackIndex]...)
				}
			}
		} else if strings.HasPrefix(line, "move") {
			moveInstance := createMoveInstance(line)
			moves = append(moves, moveInstance)
		}
	}
	inputFile.Close()

	stacks = doMoves9001(stacks, moves)
	result := getResult(stacks)

	return result
}

func doMoves9001(stacks [][]rune, moves []moveStruct) [][]rune {
	for _, move := range moves {
		fromIndex, toIndex := move.from-1, move.to-1
		objectsToMove := stacks[fromIndex][len(stacks[fromIndex])-move.count:]
		stacks[toIndex] = append(stacks[toIndex], objectsToMove...)               // add
		stacks[fromIndex] = stacks[fromIndex][:len(stacks[fromIndex])-move.count] // remove
	}
	return stacks
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
