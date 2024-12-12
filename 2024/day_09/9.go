package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadFile(fileName string) ([]int, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure the file is closed when the function exits

	var numbers []int

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first (or only) line of the file
	if scanner.Scan() {
		line := scanner.Text() // Get the input as a string

		// Split the line into individual strings
		parts := strings.Split(line, "")

		// Convert strings to integers and store them in a slice
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			numbers = append(numbers, num)
		}
	}
	return numbers, err
}

func SumIntSlice(nums []int) (total int) {
	for _, num := range nums {
		total += num
	}
	return
}

func CreateDisk(diskMap []int) []string {
	diskSize := SumIntSlice(diskMap)
	disk := make([]string, 0, diskSize)

	j := 0
	for i, times := range diskMap {
		var char string
		if i%2 == 0 {
			char = strconv.Itoa(j)
			j++
		} else {
			char = "."
		}
		for k := 0; k < times; k++ {
			disk = append(disk, char)
		}
	}
	return disk
}

func FindFileInDisk(disk []string, fragment bool, rightToLeftMarker int) (int, int, int) {
	var size, emptySpaces, fileStarts int
	var fileContent string
	size, emptySpaces = 0, 0
	fileFound := false
	fileEnded := false

	for i := rightToLeftMarker; i >= 0; i-- {
		if !fileFound && disk[i] != "." {
			fileFound = true
			fileContent = disk[i]
			size++
		} else if !fileFound && disk[i] == "." {
			emptySpaces++
		} else if !fragment && fileFound && !fileEnded && disk[i] == fileContent {
			size++
		} else if !fragment && fileFound && !fileEnded && disk[i] != fileContent {
			fileEnded = true
			fileStarts = i + 1
		} else if fragment && fileFound && !fileEnded {
			fileEnded = true
			fileStarts = i + 1
		} else if fileEnded {
			break
		}
	}
	return fileStarts, size, emptySpaces
}

func FindEmptySpaceInDisk(disk []string, sizeNeeded int) (int, bool) {
	var size, spaceStarts int
	size = 0
	spaceFound := false
	spaceBigEnough := false

	for i := 0; i < len(disk); i++ {
		if !spaceFound && disk[i] == "." {
			spaceFound = true
			spaceStarts = i
			size++
		} else if spaceFound && size < sizeNeeded && disk[i] == "." {
			size++
		} else if spaceFound && size < sizeNeeded && disk[i] != "." {
			spaceFound = false
			size = 0
		} else if spaceFound && size >= sizeNeeded {
			spaceBigEnough = true
			break
		}
	}
	return spaceStarts, spaceBigEnough
}

func Optimise(disk []string, fragment bool, rightToLeftMarker int) []string {
	// base case: if we are done with our swip, return optimised disk
	if rightToLeftMarker < 0 {
		return disk
	}

	lastFileLoc, fileSize, emptySpaces := FindFileInDisk(disk, fragment, rightToLeftMarker)
	firstEmptyLoc, isThereSpace := FindEmptySpaceInDisk(disk, fileSize)

	if isThereSpace {
		if firstEmptyLoc < lastFileLoc {
			for i := 0; i < fileSize; i++ {
				disk[firstEmptyLoc+i], disk[lastFileLoc+i] = disk[lastFileLoc+i], disk[firstEmptyLoc+i]
			}
		}
	}

	// we update the marker and keep going
	rightToLeftMarker = rightToLeftMarker - fileSize - emptySpaces
	return Optimise(disk, fragment, rightToLeftMarker)
}

func ComputeCheckSum(disk []string) int {
	checkSum := 0
	for i, v := range disk {
		if v != "." {
			intVal, _ := strconv.Atoi(v)
			checkSum += i * intVal
		}
	}
	return checkSum
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]

	// Read map
	diskMap, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		os.Exit(1)
	}

	// Create disk
	disk := CreateDisk(diskMap)

	// Optimise disk
	optDisk := make([]string, len(disk)) // Create a new slice with the same length
	copy(optDisk, disk)
	optDisk = Optimise(optDisk, true, len(disk)-1)

	optDiskNoFragment := make([]string, len(disk)) // Create a new slice with the same length
	copy(optDiskNoFragment, disk)
	optDiskNoFragment = Optimise(optDiskNoFragment, false, len(disk)-1)

	// Calculate checksum
	checkSum := ComputeCheckSum(optDisk)
	fmt.Println("Check sum?", checkSum)

	checkSumNoFragment := ComputeCheckSum(optDiskNoFragment)
	fmt.Println("Check sum (no fragmenting)?", checkSumNoFragment)
}
