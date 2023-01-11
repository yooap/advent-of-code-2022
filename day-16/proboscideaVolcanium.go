package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	name  string
	rate  int
	moves []*Node
}

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	nameResolutionMap := assembleCave()
	highscore := traverse(nameResolutionMap)
	return highscore
}

func part2() int {
	nameResolutionMap := assembleCave()
	highscore := traverseWithElephant(nameResolutionMap)
	return highscore
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func assembleCave() map[string]*Node {
	inputFile, inputScanner := getInputScanner()
	var currentValve *Node
	nameResolutionMap := map[string]*Node{}
	for inputScanner.Scan() {
		line := inputScanner.Text()
		values := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ',' || r == ';' || r == '='
		})
		name := values[1]
		currentValve, nameResolutionMap = getOrCreateNode(name, nameResolutionMap)
		rate, _ := strconv.Atoi(values[5])
		(*currentValve).rate = rate
		for i := 10; i < len(values); i++ {
			moveName := values[i]
			var moveValve *Node
			moveValve, nameResolutionMap = getOrCreateNode(moveName, nameResolutionMap)
			(*currentValve).moves = append((*currentValve).moves, moveValve)
		}
	}
	inputFile.Close()
	return nameResolutionMap
}

func getOrCreateNode(name string, nameResolutionMap map[string]*Node) (*Node, map[string]*Node) {
	if val, ok := nameResolutionMap[name]; ok {
		return val, nameResolutionMap
	} else {
		newNode := &(Node{name: name})
		nameResolutionMap[name] = newNode
		return newNode, nameResolutionMap
	}
}

func traverse(nameResolutionMap map[string]*Node) int {
	valvaesThatMatter := getValvesThatMatter(nameResolutionMap)

	shortestPaths := map[string]int{}
	pathsToNotTakeIfBetterFreeValveExists := map[string]map[string]bool{}
	for i, valve1 := range valvaesThatMatter {
		for _, valve2 := range valvaesThatMatter[i+1:] {
			shortestPath, nodesInBetweenWithBetterRate := getShortestPath(valve1, valve2, nameResolutionMap)
			shortestPaths[valve1+"-"+valve2] = shortestPath
			pathsToNotTakeIfBetterFreeValveExists[valve1+"-"+valve2] = nodesInBetweenWithBetterRate
			if valve1 == "AA" {
				continue // AA will never be a destination
			}

			// flip
			shortestPath, nodesInBetweenWithBetterRate = getShortestPath(valve2, valve1, nameResolutionMap)
			shortestPaths[valve2+"-"+valve1] = shortestPath
			pathsToNotTakeIfBetterFreeValveExists[valve2+"-"+valve1] = nodesInBetweenWithBetterRate
		}
	}

	// traverse all possible combinations of shortest paths
	possibleMoves := map[string]bool{}
	for _, valve := range valvaesThatMatter {
		possibleMoves[valve] = true
	}
	return getHighestScore("AA", possibleMoves, shortestPaths, pathsToNotTakeIfBetterFreeValveExists, 0, 0, 30, nameResolutionMap)
}

func getValvesThatMatter(nameResolutionMap map[string]*Node) []string {
	valvesThatMatter := []string{"AA"}
	for name, ref := range nameResolutionMap {
		if (*ref).rate > 0 {
			valvesThatMatter = append(valvesThatMatter, name)
		}
	}
	return valvesThatMatter
}

// shortest path to go to valve, and valves in between that have a better rate
func getShortestPath(start, end string, nameResolutionMap map[string]*Node) (int, map[string]bool) {
	shortestPath, nodesInBetweenWithBetterRate := traverseInSearchOfShortestPath(nameResolutionMap[start], nameResolutionMap[end], map[string]bool{})
	delete(nodesInBetweenWithBetterRate, start)
	return shortestPath, nodesInBetweenWithBetterRate
}

func traverseInSearchOfShortestPath(current *Node, end *Node, visitedValves map[string]bool) (int, map[string]bool) {
	currentName := (*current).name
	if currentName == (*end).name {
		return len(visitedValves), map[string]bool{}
	}

	if _, ok := visitedValves[currentName]; ok {
		return 1000, map[string]bool{} // loop
	}

	visitedValves[currentName] = true

	shortestPath := 1000
	nodesInBetweenWithBetterRateForChosenPath := map[string]bool{}
	for _, nextMove := range (*current).moves {
		newPath, nodesInBetweenWithBetterRateForThisPath := traverseInSearchOfShortestPath(nextMove, end, copyMap(visitedValves))
		if newPath == 1 {
			// can not get shorter
			return 1, map[string]bool{}
		}
		if newPath < shortestPath {
			nodesInBetweenWithBetterRateForChosenPath = nodesInBetweenWithBetterRateForThisPath
			shortestPath = newPath
		} else if newPath == shortestPath {
			for nodeToAdd, _ := range nodesInBetweenWithBetterRateForThisPath {
				nodesInBetweenWithBetterRateForChosenPath[nodeToAdd] = true
			}
		}
	}
	if (*current).rate > (*end).rate {
		nodesInBetweenWithBetterRateForChosenPath[(*current).name] = true
	}
	return shortestPath, nodesInBetweenWithBetterRateForChosenPath
}

func copyMap(oldMap map[string]bool) map[string]bool {
	newMap := map[string]bool{}
	for k, v := range oldMap {
		newMap[k] = v
	}
	return newMap
}

func getHighestScore(currentNode string, possibleMoves map[string]bool, shortestPaths map[string]int, pathsToNotTakeIfBetterFreeValveExists map[string]map[string]bool, rate, total, movesLeft int, nameResolutionMap map[string]*Node) int {
	delete(possibleMoves, currentNode)
	scores := []int{}

	if len(possibleMoves) == 0 {
		// wait it out
		scores = append(scores, total+(rate*movesLeft))
	}

	for dest, _ := range possibleMoves {
		if shouldNotTryPath(currentNode, dest, possibleMoves, pathsToNotTakeIfBetterFreeValveExists) {
			continue
		}
		pathLength := shortestPaths[currentNode+"-"+dest]
		pathLengthAndOpeningValve := pathLength + 1
		if movesLeft-pathLengthAndOpeningValve <= 0 { // end
			scores = append(scores, total+(rate*movesLeft))
		} else {
			newRate := rate + (*nameResolutionMap[dest]).rate
			newTotal := total + (pathLengthAndOpeningValve * rate)
			newMovesLeft := movesLeft - pathLengthAndOpeningValve
			scores = append(scores, getHighestScore(dest, copyMap(possibleMoves), shortestPaths, pathsToNotTakeIfBetterFreeValveExists, newRate, newTotal, newMovesLeft, nameResolutionMap))
		}
	}
	sort.Ints(scores)
	return scores[len(scores)-1]
}

func shouldNotTryPath(node1, node2 string, possibleMoves map[string]bool, pathsToNotTakeIfBetterFreeValveExists map[string]map[string]bool) bool {
	otherBetterValves := pathsToNotTakeIfBetterFreeValveExists[node1+"-"+node2]
	for betterValve, _ := range otherBetterValves {
		if _, ok := possibleMoves[betterValve]; ok {
			return true // taking path does not make sense, a closer unopened valve with a better rate exists
		}
	}
	return false
}

func traverseWithElephant(nameResolutionMap map[string]*Node) int {
	valvesThatMatter := getValvesThatMatter(nameResolutionMap)

	shortestPaths := map[string]int{}
	pathsToNotTakeIfBetterFreeValveExists := map[string]map[string]bool{}
	for i, valve1 := range valvesThatMatter {
		for _, valve2 := range valvesThatMatter[i+1:] {
			shortestPath, nodesInBetweenWithBetterRate := getShortestPath(valve1, valve2, nameResolutionMap)
			shortestPaths[valve1+"-"+valve2] = shortestPath
			pathsToNotTakeIfBetterFreeValveExists[valve1+"-"+valve2] = nodesInBetweenWithBetterRate
			if valve1 == "AA" {
				continue // AA will never be a destination
			}

			// flip
			shortestPath, nodesInBetweenWithBetterRate = getShortestPath(valve2, valve1, nameResolutionMap)
			shortestPaths[valve2+"-"+valve1] = shortestPath
			pathsToNotTakeIfBetterFreeValveExists[valve2+"-"+valve1] = nodesInBetweenWithBetterRate
		}
	}

	// traverse all possible combinations of shortest paths
	possibleMoves := map[string]bool{}
	for _, valve := range valvesThatMatter {
		if valve == "AA" {
			continue
		}
		possibleMoves[valve] = true
	}
	return getHighestScoreWithElephant("AA", "AA", possibleMoves, shortestPaths, pathsToNotTakeIfBetterFreeValveExists, 0, 0, 26, nameResolutionMap, 0, 0)
}

func getHighestScoreWithElephant(currentMeNode, currentElephantNode string,
	possibleMoves map[string]bool,
	shortestPaths map[string]int,
	pathsToNotTakeIfBetterFreeValveExists map[string]map[string]bool,
	rate, total, movesLeft int,
	nameResolutionMap map[string]*Node,
	pathLeftForMe, pathLeftForElephant int) int {

	delete(possibleMoves, currentMeNode)
	delete(possibleMoves, currentElephantNode)

	if pathLeftForMe == 0 && len(possibleMoves) > 0 { // need new dest
		scores := []int{}
		for dest, _ := range possibleMoves {
			if shouldNotTryPath(currentMeNode, dest, possibleMoves, pathsToNotTakeIfBetterFreeValveExists) {
				continue
			}
			pathLength := shortestPaths[currentMeNode+"-"+dest]
			pathLeftForMe = pathLength + 1
			scores = append(scores, getHighestScoreWithElephant(dest, currentElephantNode, copyMap(possibleMoves), shortestPaths, pathsToNotTakeIfBetterFreeValveExists, rate, total, movesLeft, nameResolutionMap, pathLeftForMe, pathLeftForElephant))
		}

		if len(scores) > 0 {
			sort.Ints(scores)
			highScore := scores[len(scores)-1]
			return highScore
		} else {
			return 0
		}
	}

	if pathLeftForElephant == 0 && len(possibleMoves) > 0 { // need new dest
		scores := []int{}
		for dest, _ := range possibleMoves {
			if shouldNotTryPath(currentElephantNode, dest, possibleMoves, pathsToNotTakeIfBetterFreeValveExists) {
				continue
			}
			pathLength := shortestPaths[currentElephantNode+"-"+dest]
			pathLeftForElephant = pathLength + 1
			scores = append(scores, getHighestScoreWithElephant(currentMeNode, dest, copyMap(possibleMoves), shortestPaths, pathsToNotTakeIfBetterFreeValveExists, rate, total, movesLeft, nameResolutionMap, pathLeftForMe, pathLeftForElephant))
		}

		if len(scores) > 0 {
			sort.Ints(scores)
			highScore := scores[len(scores)-1]
			return highScore
		} else {
			return 0
		}
	}

	if pathLeftForMe == 0 && pathLeftForElephant == 0 {
		// wait it out
		return total + (rate * movesLeft)
	} else if pathLeftForMe == pathLeftForElephant {
		// both reach dest at the same time
		newRate := rate + (*nameResolutionMap[currentMeNode]).rate + (*nameResolutionMap[currentElephantNode]).rate
		newTotal := total + (pathLeftForMe * rate)
		newMovesLeft := movesLeft - pathLeftForMe
		if newMovesLeft <= 0 {
			return total + (rate * movesLeft)
		}
		return getHighestScoreWithElephant(currentMeNode, currentElephantNode, copyMap(possibleMoves), shortestPaths, pathsToNotTakeIfBetterFreeValveExists, newRate, newTotal, newMovesLeft, nameResolutionMap, 0, 0)
	} else if pathLeftForElephant == 0 || (pathLeftForMe < pathLeftForElephant && pathLeftForMe != 0) {
		// I will finish first
		newRate := rate + (*nameResolutionMap[currentMeNode]).rate
		newTotal := total + (pathLeftForMe * rate)
		newMovesLeft := movesLeft - pathLeftForMe
		if newMovesLeft <= 0 {
			return total + (rate * movesLeft)
		}
		if pathLeftForElephant != 0 {
			pathLeftForElephant = pathLeftForElephant - pathLeftForMe
		}
		return getHighestScoreWithElephant(currentMeNode, currentElephantNode, copyMap(possibleMoves), shortestPaths, pathsToNotTakeIfBetterFreeValveExists, newRate, newTotal, newMovesLeft, nameResolutionMap, 0, pathLeftForElephant)
	} else {
		// Elephant will finish first
		newRate := rate + (*nameResolutionMap[currentElephantNode]).rate
		newTotal := total + (pathLeftForElephant * rate)
		newMovesLeft := movesLeft - pathLeftForElephant
		if newMovesLeft <= 0 {
			return total + (rate * movesLeft)
		}
		if pathLeftForMe != 0 {
			pathLeftForMe = pathLeftForMe - pathLeftForElephant
		}
		return getHighestScoreWithElephant(currentMeNode, currentElephantNode, copyMap(possibleMoves), shortestPaths, pathsToNotTakeIfBetterFreeValveExists, newRate, newTotal, newMovesLeft, nameResolutionMap, pathLeftForMe, 0)
	}
}
