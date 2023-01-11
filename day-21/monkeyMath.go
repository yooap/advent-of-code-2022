package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Monkey struct {
	value              int
	jobPart1, jobPart2 string
	operation          string
}

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	monkies, monkiesWithUnresolvedJobs := getMonkeys()
	result := resolveRoot(monkies, monkiesWithUnresolvedJobs)
	return result
}

func part2() int {
	monkies, monkiesWithUnresolvedJobs := getMonkeys()
	result := resolveHumn(monkies, monkiesWithUnresolvedJobs)
	return result
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getMonkeys() (monkies map[string]*Monkey, monkiesWithUnresolvedJobs map[string]bool) {
	monkies, monkiesWithUnresolvedJobs = make(map[string]*Monkey), make(map[string]bool)
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		values := strings.Split(line, " ")
		name := strings.TrimRight(values[0], ":")
		if len(values) == 2 {
			value, _ := strconv.Atoi(values[1])
			monkies[name] = &Monkey{value: value}
		} else {
			monkies[name] = &Monkey{jobPart1: values[1], jobPart2: values[3], operation: values[2]}
			monkiesWithUnresolvedJobs[name] = true
		}
	}
	inputFile.Close()
	return
}

func resolveRoot(monkies map[string]*Monkey, monkiesWithUnresolvedJobs map[string]bool) int {
	for {
		for monkeyName, _ := range monkiesWithUnresolvedJobs {
			monkey := monkies[monkeyName]
			_, jobPart1_Unresolvable := monkiesWithUnresolvedJobs[(*monkey).jobPart1]
			if jobPart1_Unresolvable {
				continue
			}
			_, jobPart2_Unresolvable := monkiesWithUnresolvedJobs[(*monkey).jobPart2]
			if jobPart2_Unresolvable {
				continue
			}

			value1, value2 := (*monkies[(*monkey).jobPart1]).value, (*monkies[(*monkey).jobPart2]).value
			op := (*monkey).operation
			value := doOp(value1, value2, op)

			if monkeyName == "root" {
				return value
			}

			(*monkey).value = value
			delete(monkiesWithUnresolvedJobs, monkeyName)
		}
	}
}

func doOp(value1, value2 int, op string) int {
	switch op {
	case "+":
		return value1 + value2
	case "-":
		return value1 - value2
	case "*":
		return value1 * value2
	case "/":
		return value1 / value2
	default:
		panic("bad op")
	}
}

func resolveHumn(monkies map[string]*Monkey, monkiesWithUnresolvedJobs map[string]bool) int {
	rootMonkey := monkies["root"]
	rootEquationMonkey1, rootEquationMonkey2 := (*rootMonkey).jobPart1, (*rootMonkey).jobPart2

	delete(monkiesWithUnresolvedJobs, "root")
	monkiesWithUnresolvedJobs["humn"] = true

	opsInLoop := 1
	for opsInLoop > 0 { // resolve what is possible
		opsInLoop = 0
		for monkeyName, _ := range monkiesWithUnresolvedJobs {
			if monkeyName == "humn" {
				continue
			}
			monkey := monkies[monkeyName]
			_, jobPart1_Unresolvable := monkiesWithUnresolvedJobs[(*monkey).jobPart1]
			if jobPart1_Unresolvable {
				continue
			}
			_, jobPart2_Unresolvable := monkiesWithUnresolvedJobs[(*monkey).jobPart2]
			if jobPart2_Unresolvable {
				continue
			}

			value1, value2 := (*monkies[(*monkey).jobPart1]).value, (*monkies[(*monkey).jobPart2]).value
			op := (*monkey).operation
			value := doOp(value1, value2, op)
			(*monkey).value = value
			delete(monkiesWithUnresolvedJobs, monkeyName)
			opsInLoop++
		}
	}
	// fmt.Println((*monkies[rootEquationMonkey1]).value)
	// fmt.Println((*monkies[rootEquationMonkey2]).value)

	var rootEquationValue int
	var rootEquationUnresolvedSide string
	if _, ok := monkiesWithUnresolvedJobs[rootEquationMonkey1]; ok {
		rootEquationValue = (*monkies[rootEquationMonkey2]).value
		rootEquationUnresolvedSide = rootEquationMonkey1
	} else {
		rootEquationValue = (*monkies[rootEquationMonkey1]).value
		rootEquationUnresolvedSide = rootEquationMonkey2
	}

	return resolveRecursivly(monkies[rootEquationUnresolvedSide], rootEquationValue, monkies, monkiesWithUnresolvedJobs)
}

func resolveRecursivly(monkey *Monkey, valueToCompareTo int, monkies map[string]*Monkey, monkiesWithUnresolvedJobs map[string]bool) int {
	jobPart1 := (*monkey).jobPart1
	jobPart2 := (*monkey).jobPart2
	op := (*monkey).operation

	if _, ok := monkiesWithUnresolvedJobs[jobPart1]; ok { // x op {number} == valueToCompareTo
		unresolvableMonkey := monkies[jobPart1]
		monkeyWithValue := monkies[jobPart2]
		number := (*monkeyWithValue).value
		valueForUnresolvableMonkey := solveEquationWithNumberAsSecondValue(number, op, valueToCompareTo)
		if jobPart1 == "humn" {
			return valueForUnresolvableMonkey
		}
		(*unresolvableMonkey).value = valueForUnresolvableMonkey
		delete(monkiesWithUnresolvedJobs, jobPart1)
		return resolveRecursivly(unresolvableMonkey, valueForUnresolvableMonkey, monkies, monkiesWithUnresolvedJobs)
	} else { // {number} op x == valueToCompareTo
		unresolvableMonkey := monkies[jobPart2]
		monkeyWithValue := monkies[jobPart1]
		number := (*monkeyWithValue).value
		valueForUnresolvableMonkey := solveEquationWithNumberAsFirstValue(number, op, valueToCompareTo)
		if jobPart2 == "humn" {
			return valueForUnresolvableMonkey
		}
		(*unresolvableMonkey).value = valueForUnresolvableMonkey
		delete(monkiesWithUnresolvedJobs, jobPart2)
		return resolveRecursivly(unresolvableMonkey, valueForUnresolvableMonkey, monkies, monkiesWithUnresolvedJobs)
	}
}

func solveEquationWithNumberAsFirstValue(number int, op string, valueToCompareTo int) int {
	switch op {
	case "+":
		return valueToCompareTo - number
	case "-":
		return number - valueToCompareTo
	case "*":
		return valueToCompareTo / number
	case "/":
		return number / valueToCompareTo
	default:
		panic("bad op")
	}
}

func solveEquationWithNumberAsSecondValue(number int, op string, valueToCompareTo int) int {
	switch op {
	case "+":
		return valueToCompareTo - number
	case "-":
		return valueToCompareTo + number
	case "*":
		return valueToCompareTo / number
	case "/":
		return valueToCompareTo * number
	default:
		panic("bad op")
	}
}
