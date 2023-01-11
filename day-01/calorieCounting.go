package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func part1() int {
	var currentCalorieTotal, highestCalorieTotal int

	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		if line == "" {
			if currentCalorieTotal > highestCalorieTotal {
				highestCalorieTotal = currentCalorieTotal
			}
			currentCalorieTotal = 0
		} else {
			calories, _ := strconv.Atoi(line)
			currentCalorieTotal += calories
		}
	}
	if currentCalorieTotal > highestCalorieTotal {
		highestCalorieTotal = currentCalorieTotal
	}

	inputFile.Close()

	return highestCalorieTotal
}

func part2() int {
	var currentCalorieTotal int
	topCalorieTotalArray := [3]int{}

	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		if line == "" {
			updateCalorieArray(&topCalorieTotalArray, currentCalorieTotal)
			currentCalorieTotal = 0
		} else {
			calories, _ := strconv.Atoi(line)
			currentCalorieTotal += calories
		}
	}
	updateCalorieArray(&topCalorieTotalArray, currentCalorieTotal)

	inputFile.Close()

	var sum int
	for _, calories := range topCalorieTotalArray {
		sum += calories
	}
	return sum
}

func updateCalorieArray(topCalorieTotalArray *[3]int, currentCalorieTotal int) {
	for i, topCalorieTotal := range *topCalorieTotalArray {
		if currentCalorieTotal <= topCalorieTotal {
			if i != 0 {
				(*topCalorieTotalArray)[i-1] = currentCalorieTotal // insert into previously shifted values place
			}
			break
		}

		if i != 0 {
			(*topCalorieTotalArray)[i-1] = topCalorieTotal // shift value left

		}

		if i == len(*topCalorieTotalArray)-1 {
			(*topCalorieTotalArray)[i] = currentCalorieTotal // add as last value
		}
	}
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
