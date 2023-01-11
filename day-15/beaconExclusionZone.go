package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type SensorAndBeacon struct {
	sensor, beacon Coord
}

type Coord struct {
	x, y int
}

const (
	// rowToCount = 10
	rowToCount = 2000000

	// tuningFrequencyMaxX = 20
	tuningFrequencyMaxX = 4000000
	tuningFrequencyMinX = 0

	// tuningFrequencyMaxY = 20
	tuningFrequencyMaxY = 4000000
	tuningFrequencyMinY = 0

	tuningFrequencyXMultiplier = 4000000
)

func main() {
	fmt.Println("Part 1:", part1())
	fmt.Println("Part 2:", part2())
}

func part1() (result int) {
	start := time.Now()
	sensorsAndBeacons := getSensorsAndBeacons()
	count := countAtLine(sensorsAndBeacons, rowToCount)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	return count
}

func part2() int {
	start := time.Now()
	sensorsAndBeacons := getSensorsAndBeacons()
	freq := getTuningFrequency(sensorsAndBeacons, rowToCount)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	return freq
}

func getInputScanner() (*os.File, *bufio.Scanner) {
	inputFile, _ := os.Open("input.txt")
	// inputFile, _ := os.Open("input-sample.txt")
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	return inputFile, inputScanner
}

func getSensorsAndBeacons() (sensorAndBeaconSlice []SensorAndBeacon) {
	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		values := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ',' || r == ':' || r == '='
		})
		sensorX, _ := strconv.Atoi(values[3])
		sensorY, _ := strconv.Atoi(values[5])
		beaconX, _ := strconv.Atoi(values[11])
		beaconY, _ := strconv.Atoi(values[13])
		sensorAndBeacon := SensorAndBeacon{
			Coord{sensorX, sensorY},
			Coord{beaconX, beaconY}}
		sensorAndBeaconSlice = append(sensorAndBeaconSlice, sensorAndBeacon)
	}
	inputFile.Close()
	return
}

func countAtLine(sensorAndBeaconSlice []SensorAndBeacon, row int) (count int) {
	xMax, _, xMin, _ := getGridBounds(sensorAndBeaconSlice)
	xOffset := (xMax - xMin) / 2

	rowSlice := make([]rune, (xMax-xMin)*3)

	for _, sensorAndBeacon := range sensorAndBeaconSlice {
		sensor, beacon := sensorAndBeacon.sensor, sensorAndBeacon.beacon
		if sensor.y == row {
			rowSlice[sensor.x+xOffset] = 'S'
		}
		if beacon.y == row {
			rowSlice[beacon.x+xOffset] = 'B'
		}
		radius := int(math.Abs(float64(beacon.x-sensor.x))) + int(math.Abs(float64(beacon.y-sensor.y)))
		if row >= sensor.y-radius && row <= sensor.y+radius {
			for x := sensor.x + xOffset - radius; x <= sensor.x+xOffset+radius; x++ {
				yTravel := radius - int(math.Abs(float64(sensor.x+xOffset-x)))
				if row >= sensor.y-yTravel && row <= sensor.y+yTravel {
					if rowSlice[x] == '\x00' {
						rowSlice[x] = '#'
					}
				}
			}
		}
	}

	for _, value := range rowSlice {
		if value == '#' || value == 'S' {
			count++
		}
	}

	return
}

func getGridBounds(sensorAndBeaconSlice []SensorAndBeacon) (xMax, yMax, xMin, yMin int) {
	for _, sensorAndBeacon := range sensorAndBeaconSlice {
		coords := []Coord{sensorAndBeacon.sensor, sensorAndBeacon.beacon}
		for _, coord := range coords {
			x, y := coord.x, coord.y
			if x > xMax {
				xMax = x
			}
			if y > yMax {
				yMax = y
			}
			if x < xMin {
				xMin = x
			}
			if y < yMin {
				yMin = y
			}
		}
	}
	return xMax, yMax, xMin, yMin
}

func getTuningFrequency(sensorAndBeaconSlice []SensorAndBeacon, rowToCount int) int {
	for x := tuningFrequencyMinX; x <= tuningFrequencyMaxX; x++ {
		for y := tuningFrequencyMinY; y <= tuningFrequencyMaxY; y++ {
			found := true
			for _, sensorAndBeacon := range sensorAndBeaconSlice {
				sensor, beacon := sensorAndBeacon.sensor, sensorAndBeacon.beacon
				if sensor.x == x && sensor.y == y {
					found = false
					break
				}
				if beacon.x == x && sensor.y == y {
					found = false
					break
				}
				radius := math.Abs(float64(beacon.x-sensor.x)) + math.Abs(float64(beacon.y-sensor.y))
				toPoint := radiusToPoint(x, y, sensor.x, sensor.y)
				if toPoint <= radius {
					found = false
					y += int(radius - toPoint)
					for y > tuningFrequencyMaxY {
						y = 0
						x++
						toPoint = radiusToPoint(x, y, sensor.x, sensor.y)
						if toPoint <= radius {
							y += int(radius - toPoint)
						}
					}
					break
				}
			}
			if found {
				return calculateTuningFreqValue(x, y)
			}
		}
	}
	panic("value missing")
}

func radiusToPoint(pointX, pointY, sensorX, sensorY int) float64 {
	return math.Abs(float64(pointY-sensorY)) + math.Abs(float64(pointX-sensorX))
}

func calculateTuningFreqValue(x, y int) int {
	return x*tuningFrequencyXMultiplier + y
}
