package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type dir struct {
	name      string
	subDirs   map[string]*dir
	files     map[string]int
	parent    *dir
	totalSize int
}

const totalDiskSpace int = 70000000
const diskSpaceNeededForUpdate int = 30000000

func part1() int {
	root := createFileSystem()
	result := sumDirTotalSizeWithThreshold(&root, 100000)
	return result
}

func createFileSystem() dir {
	root := dir{
		subDirs: map[string]*dir{},
		files:   map[string]int{},
	}

	var currentDirRef *dir

	inputFile, inputScanner := getInputScanner()
	for inputScanner.Scan() {
		line := inputScanner.Text()
		splitLine := strings.Split(line, " ")
		splitLineFirstElement := splitLine[0]
		if splitLineFirstElement == "$" {
			command := splitLine[1]
			if command == "cd" {
				newDirName := splitLine[2]
				if newDirName == "/" {
					currentDirRef = &root
				} else if newDirName == ".." {
					currentDirRef = currentDirRef.parent
				} else {
					currentDirRef = (*currentDirRef).subDirs[newDirName]
				}
			}
		} else if splitLineFirstElement == "dir" {
			dirName := splitLine[1]
			if _, ok := (*currentDirRef).subDirs[dirName]; !ok {
				newDir := dir{
					name:    dirName,
					subDirs: map[string]*dir{},
					files:   map[string]int{},
					parent:  currentDirRef}
				(*currentDirRef).subDirs[dirName] = &newDir
			}
		} else {
			fileName := splitLine[1]
			fileSize, _ := strconv.Atoi(splitLineFirstElement)
			if _, ok := (*currentDirRef).files[fileName]; !ok {
				(*currentDirRef).files[fileName] = fileSize
				updateDirTotalSizeRecursively(currentDirRef, fileSize)

			}
		}
	}
	inputFile.Close()
	return root
}

func updateDirTotalSizeRecursively(dirRef *dir, fileSize int) {
	(*dirRef).totalSize += fileSize
	parent := (*dirRef).parent
	if parent != nil {
		updateDirTotalSizeRecursively(parent, fileSize)
	}
}

func sumDirTotalSizeWithThreshold(dir *dir, threshold int) (sum int) {
	dirTotalSize := (*dir).totalSize
	if dirTotalSize <= threshold {
		sum += dirTotalSize
	}

	for _, subDirRef := range (*dir).subDirs {
		sum += sumDirTotalSizeWithThreshold(subDirRef, threshold)
	}

	return
}

func part2() int {
	root := createFileSystem()
	size := getSizeOfDirToDeleteToFreeUpSpace(&root)
	return size
}

func getSizeOfDirToDeleteToFreeUpSpace(root *dir) int {
	currentSpace := totalDiskSpace - (*root).totalSize
	if currentSpace > diskSpaceNeededForUpdate {
		panic("enough space already")
	}

	minSpaceToFreeUp := diskSpaceNeededForUpdate - currentSpace
	size := findBestCandidateForDeletion(root, minSpaceToFreeUp)
	if size == 0 {
		size = (*root).totalSize
	}
	return size
}

func findBestCandidateForDeletion(dir *dir, minSpaceToFreeUp int) (bestCandidateSize int) {
	for _, subDirRef := range (*dir).subDirs {
		subDirTotalSize := (*subDirRef).totalSize
		if subDirTotalSize >= minSpaceToFreeUp {
			if bestCandidateSize == 0 || subDirTotalSize < bestCandidateSize {
				bestCandidateSize = subDirTotalSize
			}

			candidateSize := findBestCandidateForDeletion(subDirRef, minSpaceToFreeUp)

			if candidateSize != 0 && candidateSize < bestCandidateSize {
				bestCandidateSize = candidateSize
			}
		}
	}

	return bestCandidateSize
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
