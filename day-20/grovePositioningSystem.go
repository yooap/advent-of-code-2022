package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	groveCoordStartValue = 0
	groveCoordIndexStep1 = 1000
	groveCoordIndexStep2 = 2000
	groveCoordIndexStep3 = 3000

	decryptionKey = 811589153
)

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	numbers, groveCoordStartValueIndexInReferenceSet := getNumbers()
	result := mixNumbers(numbers, groveCoordStartValueIndexInReferenceSet, 1)
	return result
}

func part2() int {
	numbers, groveCoordStartValueIndexInReferenceSet := getNumbers()
	for i, number := range numbers {
		numbers[i] = number * decryptionKey
	}
	result := mixNumbers(numbers, groveCoordStartValueIndexInReferenceSet, 10)
	return result
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getNumbers() (numbers []int, groveCoordStartValueIndexInReferenceSet int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		number, _ := strconv.Atoi(line)
		numbers = append(numbers, number)
		if number == groveCoordStartValue {
			groveCoordStartValueIndexInReferenceSet = len(numbers) - 1
		}
	}
	inputFile.Close()
	return
}

func mixNumbers(numbers []int, groveCoordStartValueIndexInReferenceSet, loops int) int {
	indexReference := make([]int, len(numbers))
	for i, _ := range indexReference {
		indexReference[i] = i
	}

	for loopIdx := 0; loopIdx < loops; loopIdx++ {
		for refIdx, nIdx := range indexReference {
			number := numbers[nIdx]
			spacesToMove := number % (len(numbers) - 1)
			if spacesToMove < 0 && nIdx+spacesToMove <= 0 {
				spacesToMove = spacesToMove + len(numbers) - 1 // wrap
			} else if spacesToMove > 0 && nIdx+spacesToMove >= len(numbers) {
				spacesToMove = spacesToMove - len(numbers) + 1 // wrap
			} else if spacesToMove == 0 {
				continue
			}

			if spacesToMove < 0 {
				beforeNewIndex := numbers[:nIdx+spacesToMove]
				betweenNewAndCurrentIndexes := numbers[nIdx+spacesToMove : nIdx]
				afterCurrentIndex := numbers[nIdx+1:]
				numbers = append(beforeNewIndex, append([]int{number}, append(betweenNewAndCurrentIndexes, afterCurrentIndex...)...)...)

				for ii, v := range indexReference {
					if v >= nIdx+spacesToMove && v < nIdx {
						indexReference[ii]++
					}
				}
				indexReference[refIdx] += spacesToMove
			} else {
				beforeCurrentIndex := numbers[:nIdx]
				betweenCurrentAndNewIndexes := numbers[nIdx+1 : nIdx+spacesToMove+1]
				afterNewIndex := numbers[nIdx+spacesToMove+1:]
				numbers = append(beforeCurrentIndex, append(betweenCurrentAndNewIndexes, append([]int{number}, afterNewIndex...)...)...)

				for ii, v := range indexReference {
					if v > nIdx && v <= nIdx+spacesToMove {
						indexReference[ii]--
					}
				}
				indexReference[refIdx] += spacesToMove
			}
		}
	}

	indexInNumbersList := indexReference[groveCoordStartValueIndexInReferenceSet]
	groveCoordIndex1 := (indexInNumbersList + groveCoordIndexStep1) % len(numbers)
	groveCoordIndex2 := (indexInNumbersList + groveCoordIndexStep2) % len(numbers)
	groveCoordIndex3 := (indexInNumbersList + groveCoordIndexStep3) % len(numbers)
	return numbers[groveCoordIndex1] + numbers[groveCoordIndex2] + numbers[groveCoordIndex3]
}
