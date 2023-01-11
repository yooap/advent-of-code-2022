package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type move struct {
	direction rune
	moveCount int
}

type coord struct {
	x, y int
}

func part1() int {
	moves := getMoves()
	visitedSpaceCount := executeMoves(moves)
	return visitedSpaceCount
}

func getMoves() (moves []move) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		moveArguemnts := strings.Split(inputScanner.Text(), " ")
		direction := []rune(moveArguemnts[0])[0]
		moveCount, _ := strconv.Atoi(moveArguemnts[1])
		moves = append(moves, move{direction, moveCount})
	}
	inputFile.Close()
	return
}

func executeMoves(moves []move) int {
	visitedSpaces := map[coord]bool{}
	head, tail := coord{0, 0}, coord{0, 0}
	for _, move := range moves {
		for i := 0; i < move.moveCount; i++ {
			head = updateHeadPosition(move.direction, head)
			tail = updateTailPosition(tail, head)
			visitedSpaces[tail] = true
		}
	}

	return len(visitedSpaces)
}

func updateHeadPosition(direction rune, head coord) coord {
	switch direction {
	case 'U':
		head = coord{head.x, head.y + 1}
	case 'R':
		head = coord{head.x + 1, head.y}
	case 'D':
		head = coord{head.x, head.y - 1}
	case 'L':
		head = coord{head.x - 1, head.y}
	}
	return head
}

func updateTailPosition(tail, head coord) coord {
	tailX, tailY, headX, headY := tail.x, tail.y, head.x, head.y
	moveX, moveY := 0, 0

	if math.Abs(float64(tailX-headX)) == 2 {
		moveX = (headX - tailX) / 2
		moveY = headY - tailY
	}

	if math.Abs(float64(tailY-headY)) == 2 {
		moveY = (headY - tailY) / 2
		if moveX == 0 {
			moveX = headX - tailX
		}
	}

	return coord{tailX + moveX, tailY + moveY}
}

func part2() int {
	moves := getMoves()
	visitedSpaceCount := executeMovesWithTenKnots(moves)
	return visitedSpaceCount
}

func executeMovesWithTenKnots(moves []move) int {
	visitedSpaces := map[coord]bool{}
	knots := [10]coord{}
	for _, move := range moves {
		for i := 0; i < move.moveCount; i++ {
			knots[0] = updateHeadPosition(move.direction, knots[0])
			for j := 1; j < len(knots); j++ {
				knots[j] = updateTailPosition(knots[j], knots[j-1])
			}
			visitedSpaces[knots[len(knots)-1]] = true
		}
	}

	return len(visitedSpaces)
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	//inputFile, _ := os.Open("input-sample-2.txt")
	inputFile, _ := os.Open("input.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func main() {
	fmt.Println("Part 1:", part1())
	fmt.Println("Part 2:", part2())
}
