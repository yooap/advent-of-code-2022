package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type monkey struct {
	items       []int
	op          [3]string
	test        int
	testResult  map[bool]int
	inspections int
}

func newMonkey() monkey {
	return monkey{testResult: map[bool]int{}}
}

func (m monkey) updateWorryLvl(item int) int {
	var firstOpParam, secondOpParam int

	if m.op[0] == "old" {
		firstOpParam = item
	} else {
		firstOpParam, _ = strconv.Atoi(m.op[0])
	}

	if m.op[2] == "old" {
		secondOpParam = item
	} else {
		secondOpParam, _ = strconv.Atoi(m.op[2])
	}

	var opResult int
	switch m.op[1] {
	case "+":
		opResult = firstOpParam + secondOpParam
	case "*":
		opResult = firstOpParam * secondOpParam
	default:
		panic("unexpected op")
	}

	return opResult / 3
}

func (m monkey) updateWorryLvlUsingMagicNumber(item, magicNumber int) int {
	var firstOpParam, secondOpParam int

	if m.op[0] == "old" {
		firstOpParam = item
	} else {
		firstOpParam, _ = strconv.Atoi(m.op[0])
	}

	if m.op[2] == "old" {
		secondOpParam = item
	} else {
		secondOpParam, _ = strconv.Atoi(m.op[2])
	}

	var opResult int
	switch m.op[1] {
	case "+":
		opResult = firstOpParam + secondOpParam
	case "*":
		opResult = firstOpParam * secondOpParam
	default:
		panic("unexpected op")
	}

	return opResult % magicNumber
}

func (m monkey) doTest(item int) int {
	return m.testResult[(item%m.test == 0)]
}

func part1() int {
	monkeys := parseMonkeys()
	result := calculateMonkeyBusiness(monkeys)
	return result
}

func parseMonkeys() (monkeys []*monkey) {
	inputFile, inputScanner := getInputScanner()
	var currentMonkey *monkey
	for inputScanner.Scan() {
		line := strings.TrimSpace(inputScanner.Text())
		if strings.HasPrefix(line, "Monkey") {
			newMonkey := newMonkey()
			monkeys = append(monkeys, &newMonkey)
			currentMonkey = &newMonkey
		} else if strings.HasPrefix(line, "Starting items: ") {
			itemsString := strings.TrimPrefix(line, "Starting items: ")
			items := strings.Split(itemsString, ", ")
			for _, item := range items {
				itemAsInt, _ := strconv.Atoi(item)
				currentMonkey.items = append(currentMonkey.items, itemAsInt)
			}
		} else if strings.HasPrefix(line, "Operation: ") {
			opString := strings.TrimPrefix(line, "Operation: new = ")
			opParts := strings.Split(opString, " ")
			currentMonkey.op = [3]string{opParts[0], opParts[1], opParts[2]}
		} else if strings.HasPrefix(line, "Test") {
			testString := strings.TrimPrefix(line, "Test: divisible by ")
			testValue, _ := strconv.Atoi(testString)
			currentMonkey.test = testValue
		} else if strings.HasPrefix(line, "If true") {
			testResultString := strings.TrimPrefix(line, "If true: throw to monkey ")
			testResultValue, _ := strconv.Atoi(testResultString)
			currentMonkey.testResult[true] = testResultValue
		} else if strings.HasPrefix(line, "If false") {
			testResultString := strings.TrimPrefix(line, "If false: throw to monkey ")
			testResultValue, _ := strconv.Atoi(testResultString)
			currentMonkey.testResult[false] = testResultValue
		}
	}
	inputFile.Close()
	return
}

func calculateMonkeyBusiness(monkeys []*monkey) int {
	for round := 1; round <= 20; round++ {
		for _, monkeyRef := range monkeys {
			for _, item := range (*monkeyRef).items {
				newItem := (*monkeyRef).updateWorryLvl(item)
				destinationMonkey := (*monkeyRef).doTest(newItem)
				(*monkeys[destinationMonkey]).items = append((*monkeys[destinationMonkey]).items, newItem)
			}
			(*monkeyRef).inspections += len((*monkeyRef).items)
			(*monkeyRef).items = []int{}
		}
	}

	return getInspectionResult(monkeys)
}

func getInspectionResult(monkeys []*monkey) int {
	inspections := []int{}
	for _, monkeyRef := range monkeys {
		inspections = append(inspections, (*monkeyRef).inspections)
	}
	sort.Ints(inspections)
	return inspections[len(inspections)-1] * inspections[len(inspections)-2]
}

func part2() int {
	monkeys := parseMonkeys()
	result := calculateMonkeyBusinessWithNoWorryDivision(monkeys)
	return result
}

func calculateMonkeyBusinessWithNoWorryDivision(monkeys []*monkey) int {
	magicNumber := 1
	for _, monkeyRef := range monkeys {
		magicNumber *= (*monkeyRef).test
	}

	for round := 1; round <= 10000; round++ {
		for _, monkeyRef := range monkeys {
			for _, item := range (*monkeyRef).items {
				newItem := (*monkeyRef).updateWorryLvlUsingMagicNumber(item, magicNumber)
				destinationMonkey := (*monkeyRef).doTest(newItem)
				(*monkeys[destinationMonkey]).items = append((*monkeys[destinationMonkey]).items, newItem)
			}
			(*monkeyRef).inspections += len((*monkeyRef).items)
			(*monkeyRef).items = []int{}
		}
	}

	return getInspectionResult(monkeys)
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
