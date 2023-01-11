package main

import (
	"bufio"
	"fmt"
	"os"
)

func part1() (result int) {
	inputFile, inputScanner := getInputScanner()
	var packet1, packet2 string
	for i := 0; inputScanner.Scan(); i++ {
		if i%3 == 0 {
			packet1 = inputScanner.Text()
			continue
		} else if i%3 == 1 {
			packet2 = inputScanner.Text()
		} else {
			if compare(packet1, packet2) {
				result += (i + 1) / 3
			}
		}
	}
	inputFile.Close()
	return
}

func compare(packet1, packet2 string) bool {
	j := 0
	for i, value1 := range packet1 {
		// handle ran out of elements
		if value1 == ']' && packet2[j] == ']' {
			j++
			continue
		} else if value1 == ']' && packet2[j] != ']' {
			return true
		} else if value1 != ']' && packet2[j] == ']' {
			return false
		}

		if value1 == ',' {
			j++
			continue
		}

		if value1 == '[' {
			if packet2[j] == '[' {
				j++
				continue
			} else { // list vs int: convert
				packet2Converted := packet2[:j]
				packet2Converted += "["
				packet2Converted += packet2[j : j+1]
				if packet2[j+1] == '0' {
					j++
					packet2Converted += "0"
				}
				packet2Converted += "]"
				packet2Converted += packet2[j+1:]
				return compare(packet1, packet2Converted)
			}
		} else { // int
			for packet2[j] == '[' { // int vs list: convert
				packet1Converted := packet1[:i]
				packet1Converted += "["
				packet1Converted += packet1[i : i+1]
				if packet1[i+1] == '0' {
					i++
					packet1Converted += "0"
				}
				packet1Converted += "]"
				packet1Converted += packet1[i+1:]
				return compare(packet1Converted, packet2)
			}

			// as ints
			int1 := value1 - '0'
			int2 := []rune(packet2)[j] - '0'

			// handle double digit numbers
			if int1 == 1 {
				if packet1[i+1] == '0' {
					continue
				}
			}
			if int2 == 1 {
				if packet2[j+1] == '0' {
					j++
					int2 = 10
				}
			}
			if int1 == 0 {
				if packet1[i-1] == '1' {
					int1 = 10
				}
			}
			if int2 == 0 {
				if packet1[j-1] == '1' {
					int2 = 10
				}
			}

			// compare ints
			if int1 < int2 {
				return true
			} else if int2 < int1 {
				return false
			} else {
				j++
				continue
			}
		}
	}
	return true
}

func part2() int {
	inputFile, inputScanner := getInputScanner()
	sortedPackets := []string{"[[2]]", "[[6]]"}
	for inputScanner.Scan() {
		packet := inputScanner.Text()
		if len(packet) == 0 {
			continue
		}
		j := 0
		for _, packetFromSortedSlice := range sortedPackets {
			if compare(packet, packetFromSortedSlice) {
				break
			}
			j++
		}
		sortedPackets = append(sortedPackets[:j], append([]string{packet}, sortedPackets[j:]...)...)
	}

	inputFile.Close()
	return multiplyDividerPositions(sortedPackets)
}

func multiplyDividerPositions(sortedPackets []string) int {
	var index1, index2 int
	for i, packet := range sortedPackets {
		if packet == "[[2]]" {
			index1 = i + 1
			continue
		} else if packet == "[[6]]" {
			index2 = i + 1
			break
		}
	}
	return index1 * index2
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
