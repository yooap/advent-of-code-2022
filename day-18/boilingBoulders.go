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

type Coord struct {
	x, y, z int
}

func main() {
	start := time.Now()
	fmt.Println("Part 1:", part1(), " (", time.Since(start), ")")
	start = time.Now()
	fmt.Println("Part 2:", part2(), " (", time.Since(start), ")")
}

func part1() int {
	cubes := getCubes()
	surfaceArea := calcSurfaceArea(cubes)
	return surfaceArea
}

func part2() int {
	cubes := getCubes()
	outerSurfaceArea := fillRoomWithWaterAndCalcOuterSurfaceArea(cubes)
	return outerSurfaceArea
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getCubes() (cubes []Coord) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		coords := strings.Split(line, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])
		cubes = append(cubes, Coord{x, y, z})
	}
	inputFile.Close()
	return
}

func calcSurfaceArea(cubes []Coord) (surfaceArea int) {
	surfaceArea = 6 * len(cubes)
	sort.Slice(cubes, func(i, j int) bool { return cubes[i].x < cubes[j].x }) // sort by x
	for i, cubeLeft := range cubes {
		cubesThatCouldBeConnectedToTheRight := cubes[i+1:]
		for _, cubeRight := range cubesThatCouldBeConnectedToTheRight {
			if cubeLeft.x+1 < cubeRight.x {
				break
			}
			if cubeLeft.x+1 == cubeRight.x &&
				cubeLeft.y == cubeRight.y &&
				cubeLeft.z == cubeRight.z {
				surfaceArea -= 2
				break
			}
		}
	}

	sort.Slice(cubes, func(i, j int) bool { return cubes[i].y < cubes[j].y }) // sort by y
	for i, cubeLeft := range cubes {
		cubesThatCouldBeConnectedToTheRight := cubes[i+1:]
		for _, cubeRight := range cubesThatCouldBeConnectedToTheRight {
			if cubeLeft.y+1 < cubeRight.y {
				break
			}
			if cubeLeft.x == cubeRight.x &&
				cubeLeft.y+1 == cubeRight.y &&
				cubeLeft.z == cubeRight.z {
				surfaceArea -= 2
				break
			}
		}
	}

	sort.Slice(cubes, func(i, j int) bool { return cubes[i].z < cubes[j].z }) // sort by z
	for i, cubeLeft := range cubes {
		cubesThatCouldBeConnectedToTheRight := cubes[i+1:]
		for _, cubeRight := range cubesThatCouldBeConnectedToTheRight {
			if cubeLeft.z+1 < cubeRight.z {
				break
			}
			if cubeLeft.x == cubeRight.x &&
				cubeLeft.y == cubeRight.y &&
				cubeLeft.z+1 == cubeRight.z {
				surfaceArea -= 2
				break
			}
		}
	}

	return
}

func fillRoomWithWaterAndCalcOuterSurfaceArea(cubes []Coord) int {
	xMax, yMax, zMax := getBoundsOfRoom(cubes)
	max := Coord{xMax, yMax, zMax}
	xMin, yMin, zMin := -1, -1, -1
	min := Coord{xMin, yMin, zMin}
	visitedCoordsSet := map[Coord]bool{}
	lavaCubeSet := map[Coord]bool{}
	for _, cube := range cubes {
		lavaCubeSet[cube] = true
	}
	return traverseAndSumSurfaceArea(Coord{xMin, yMin, zMin}, visitedCoordsSet, lavaCubeSet, min, max)
}

func getBoundsOfRoom(cubes []Coord) (xMax, yMax, zMax int) {
	for _, cube := range cubes {
		if cube.x > xMax {
			xMax = cube.x
		}
		if cube.y > yMax {
			yMax = cube.y
		}
		if cube.z > zMax {
			zMax = cube.z
		}
	}
	xMax++
	yMax++
	zMax++
	return
}

func traverseAndSumSurfaceArea(point Coord, visitedCoordsSet map[Coord]bool, lavaCubeSet map[Coord]bool, min, max Coord) (surfaceArea int) {
	visitedCoordsSet[point] = true
	var nextPoint Coord

	// x + 1
	nextPoint = Coord{point.x + 1, point.y, point.z}
	surfaceArea += assesAndTraverseIfPossible(nextPoint, visitedCoordsSet, lavaCubeSet, min, max)

	// x - 1
	nextPoint = Coord{point.x - 1, point.y, point.z}
	surfaceArea += assesAndTraverseIfPossible(nextPoint, visitedCoordsSet, lavaCubeSet, min, max)

	// y + 1
	nextPoint = Coord{point.x, point.y + 1, point.z}
	surfaceArea += assesAndTraverseIfPossible(nextPoint, visitedCoordsSet, lavaCubeSet, min, max)

	// y - 1
	nextPoint = Coord{point.x, point.y - 1, point.z}
	surfaceArea += assesAndTraverseIfPossible(nextPoint, visitedCoordsSet, lavaCubeSet, min, max)

	// z + 1
	nextPoint = Coord{point.x, point.y, point.z + 1}
	surfaceArea += assesAndTraverseIfPossible(nextPoint, visitedCoordsSet, lavaCubeSet, min, max)

	// z - 1
	nextPoint = Coord{point.x, point.y, point.z - 1}
	surfaceArea += assesAndTraverseIfPossible(nextPoint, visitedCoordsSet, lavaCubeSet, min, max)

	return surfaceArea
}

func assesAndTraverseIfPossible(point Coord, visitedCoordsSet map[Coord]bool, lavaCubeSet map[Coord]bool, min, max Coord) int {
	if point.x < min.x || point.x > max.x ||
		point.y < min.y || point.y > max.y ||
		point.z < min.z || point.z > max.z {
		return 0
	}
	if _, ok := lavaCubeSet[point]; ok {
		return 1
	}
	if _, ok := visitedCoordsSet[point]; !ok {
		return traverseAndSumSurfaceArea(point, visitedCoordsSet, lavaCubeSet, min, max)
	}
	return 0
}
