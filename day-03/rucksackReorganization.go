package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func part1() (sum int) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		halfWayPoint := len(line) / 2
		first, second := line[:halfWayPoint], line[halfWayPoint:]
		itemType := singleIntersection(first, second)
		sum += toValue(itemType)
	}

	inputFile.Close()

	return
}

func singleIntersection(first string, others ...string) rune {
	for _, char := range first {
		found := true
		for _, other := range others {
			if !strings.ContainsRune(other, char) {
				found = false
				break
			}
		}
		if found {
			return char
		}
	}
	return rune(0)
}

func toValue(itemType rune) int {
	if unicode.IsUpper(itemType) {
		return int(itemType) - 38
	} else {
		return int(itemType) - 96
	}
}

func part2() (sum int) {
	inputFile, inputScanner := getInputScanner()
	i := 0
	group := [3]string{}
	for inputScanner.Scan() {
		line := inputScanner.Text()
		group[i] = line

		if i == 2 {
			itemType := singleIntersection(group[0], group[1], group[2])
			sum += toValue(itemType)
			i = 0
		} else {
			i++
		}
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
