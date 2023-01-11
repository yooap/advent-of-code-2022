package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type SnafuNumber struct {
	digits []int // reverse order
}

func (baseNumberRef *SnafuNumber) Add(numberToAdd SnafuNumber) SnafuNumber {
	digits1 := (*baseNumberRef).digits
	digits2 := numberToAdd.digits
	digitsResult := []int{}
	max := int(math.Max(float64(len(digits1)), float64(len(digits2))))
	var transfer int
	for i := 0; i < max+1; i++ {
		var digitResult int
		var digit1, digit2 int
		if i < len(digits1) {
			digit1 = digits1[i]
		}
		if i < len(digits2) {
			digit2 = digits2[i]
		}
		digitResult = digit1 + digit2 + transfer
		transfer = 0
		if digitResult < -2 {
			transfer = -1
			digitResult = 5 + digitResult
		} else if digitResult > 2 {
			transfer = 1
			digitResult = -5 + digitResult
		}

		digitsResult = append(digitsResult, digitResult)
		if i == max-1 && transfer == 0 {
			break
		}
	}

	return SnafuNumber{digitsResult}
}

func (baseNumberRef *SnafuNumber) Print() (stringRepresentation string) {
	digits := (*baseNumberRef).digits
	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]
		switch digit {
		case -2:
			stringRepresentation += "="
		case -1:
			stringRepresentation += "-"
		default:
			stringRepresentation += strconv.Itoa(digit)
		}
	}
	return stringRepresentation
}

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() string {
	blizzards := getFuelSnafuNumbers()
	result := sum(blizzards)
	return result.Print()
}

func part2() int {
	return 0
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getFuelSnafuNumbers() (numberes []SnafuNumber) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		number := SnafuNumber{}
		for _, char := range line {
			var digit int
			switch char {
			case '=':
				digit = -2
			case '-':
				digit = -1
			default:
				digit, _ = strconv.Atoi(string(char))
			}

			number.digits = append([]int{digit}, number.digits...)
		}
		numberes = append(numberes, number)
	}
	inputFile.Close()
	return
}

func sum(numbers []SnafuNumber) SnafuNumber {
	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = result.Add(numbers[i])
	}
	return result
}
