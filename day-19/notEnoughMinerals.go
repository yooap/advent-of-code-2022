package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Blueprint struct {
	oreRobotOre                                                int
	clayRobotOre                                               int
	obsidianRobotOre, obsidianRobotClay                        int
	geodeRobotOre, geodeRobotObsidian                          int
	maxNeededOreRate, maxNeededClayRate, maxNeededObsidianRate int
}

type Status struct {
	ore, clay, obsidian, geodes                        int
	oreRobots, clayRobots, obsidianRobots, geodeRobots int
	minutesLeft                                        int
}

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	blueprints := getBlueprints()
	result := sumQuality(blueprints)
	return result
}

func part2() int {
	blueprints := getBlueprints()[0:3]
	result := multiplyMaximums(blueprints)
	return result
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getBlueprints() (blueprints []Blueprint) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		values := strings.Split(line, " ")
		oreRobotOre, _ := strconv.Atoi(values[6])
		clayRobotOre, _ := strconv.Atoi(values[12])
		obsidianRobotOre, _ := strconv.Atoi(values[18])
		obsidianRobotClay, _ := strconv.Atoi(values[21])
		geodeRobotOre, _ := strconv.Atoi(values[27])
		geodeRobotObsidian, _ := strconv.Atoi(values[30])
		blueprint := Blueprint{oreRobotOre,
			clayRobotOre,
			obsidianRobotOre, obsidianRobotClay,
			geodeRobotOre, geodeRobotObsidian,
			maxOf(clayRobotOre, obsidianRobotOre, geodeRobotOre), obsidianRobotClay, geodeRobotObsidian}
		blueprints = append(blueprints, blueprint)
	}
	inputFile.Close()
	return
}

func maxOf(vars ...int) int {
	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}

func sumQuality(blueprints []Blueprint) (quality int) {
	for i, blueprint := range blueprints {
		quality += (i + 1) * getGeodeMax(blueprint, 24)
	}
	return
}

func getGeodeMax(blueprint Blueprint, minutes int) int {
	status := Status{
		0, 0, 0, 0,
		1, 0, 0, 0,
		minutes}

	geodeMax := run(status, blueprint, 0)
	return geodeMax
}

func run(status Status, blueprint Blueprint, maxGeodesForAnAlreadyFoundPath int) int {
	status.minutesLeft--
	if status.minutesLeft == 0 {
		status.ore += status.oreRobots
		status.clay += status.clayRobots
		status.obsidian += status.obsidianRobots
		status.geodes += status.geodeRobots
		return status.geodes
	}

	if status.oreRobots >= blueprint.geodeRobotOre && status.obsidianRobots >= blueprint.geodeRobotObsidian {
		// keep purchasing 1 geodeo robot every next round and skip the rest
		for status.minutesLeft >= 0 {
			status.geodes += status.geodeRobots
			status.geodeRobots++
			status.minutesLeft--
		}
		return status.geodes
	} else {
		purchaseOptionSet := getPurchaseOptions(status, blueprint)

		for purchaseOption, _ := range purchaseOptionSet {
			purchaseOption.ore += status.oreRobots
			purchaseOption.clay += status.clayRobots
			purchaseOption.obsidian += status.obsidianRobots
			purchaseOption.geodes += status.geodeRobots

			maxGeodesForRun := run(purchaseOption, blueprint, maxGeodesForAnAlreadyFoundPath)
			if maxGeodesForRun > maxGeodesForAnAlreadyFoundPath {
				maxGeodesForAnAlreadyFoundPath = maxGeodesForRun
			}
		}
	}

	return maxGeodesForAnAlreadyFoundPath
}

// including option to not purchase anything
func getPurchaseOptions(status Status, blueprint Blueprint) map[Status]bool {
	purchaseOptions := attemptBuy(status, blueprint)
	purchaseOptionSet := map[Status]bool{}
	for _, purchaseOption := range purchaseOptions {
		purchaseOptionSet[purchaseOption] = true
	}
	return purchaseOptionSet
}

func attemptBuy(status Status, blueprint Blueprint) []Status {
	purchases := []Status{}
	for i := 3; i >= 0; i-- {
		if i == 0 { // oreRobot
			if blueprint.oreRobotOre <= status.ore && blueprint.maxNeededOreRate > status.oreRobots {
				updatedStatus := status
				updatedStatus.ore -= blueprint.oreRobotOre
				updatedStatus.oreRobots++
				purchases = append(purchases, updatedStatus)

			}
		} else if i == 1 { // clayRobot
			if blueprint.clayRobotOre <= status.ore && blueprint.maxNeededClayRate > status.clayRobots {
				updatedStatus := status
				updatedStatus.ore -= blueprint.clayRobotOre
				updatedStatus.clayRobots++
				purchases = append(purchases, updatedStatus)
			}
		} else if i == 2 { // obsidianRobot
			if blueprint.obsidianRobotOre <= status.ore && blueprint.obsidianRobotClay <= status.clay &&
				blueprint.maxNeededObsidianRate > status.obsidianRobots {
				updatedStatus := status
				updatedStatus.ore -= blueprint.obsidianRobotOre
				updatedStatus.clay -= blueprint.obsidianRobotClay
				updatedStatus.obsidianRobots++
				purchases = append(purchases, updatedStatus)
			}
		} else { // geodeRobot
			if blueprint.geodeRobotOre <= status.ore && blueprint.geodeRobotObsidian <= status.obsidian {
				updatedStatus := status
				updatedStatus.ore -= blueprint.geodeRobotOre
				updatedStatus.obsidian -= blueprint.geodeRobotObsidian
				updatedStatus.geodeRobots++
				purchases = append(purchases, updatedStatus)
			}
		}
	}

	purchases = append(purchases, status)
	return purchases
}

func multiplyMaximums(blueprints []Blueprint) int {
	result := 1
	for _, blueprint := range blueprints {
		geodeMaxForBlueprint := getGeodeMax(blueprint, 32)
		println(geodeMaxForBlueprint)
		result *= geodeMaxForBlueprint
	}
	return result
}
