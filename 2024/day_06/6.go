package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadFile(fileName string) ([][]string, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read all the rules and updates
	var labMap [][]string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read the line
		oneLine := scanner.Text()
		onwRow := strings.Split(oneLine, "")
		labMap = append(labMap, onwRow)
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return labMap, nil
}

func RotateDirection90DegreesRight(dx, dy int) (int, int) {
	if dx == 0 && dy == -1 { // up
		return 1, 0 // right
	} else if dx == 1 && dy == 0 { // right
		return 0, 1 // down
	} else if dx == 0 && dy == 1 { // down
		return -1, 0 // left
	} else if dx == -1 && dy == 0 { // left
		return 0, -1 // up
	}
	return 9, 9
}

func FindNextStep(labMap [][]string, x, y, dx, dy int) ([]int, []int) {
	nx := x + dx
	ny := y + dy

	// Check if it's empty?
	if ny < 0 || ny >= len(labMap) || nx < 0 || nx >= len(labMap[0]) {
		return []int{nx, ny}, []int{dx, dy} // if outside of bounds return same cell
	} else if labMap[ny][nx] != "#" {
		return []int{nx, ny}, []int{dx, dy} // if cell is unblocked
	}

	// rotate 90 degrees and try again
	ndx, ndy := RotateDirection90DegreesRight(dx, dy)
	return FindNextStep(labMap, x, y, ndx, ndy)
}

func WhereAmI(pathSoFar, directionSoFar [][]int) (int, int, int, int) {
	// Where am I?
	var y, x, dy, dx int
	x, y = pathSoFar[len(pathSoFar)-1][0], pathSoFar[len(pathSoFar)-1][1]
	dx, dy = directionSoFar[len(directionSoFar)-1][0], directionSoFar[len(directionSoFar)-1][1]
	return x, y, dx, dy
}

func isLoop(pathSoFar, directionSoFar [][]int) bool {
	x, y, dx, dy := WhereAmI(pathSoFar, directionSoFar)

	for i, loc := range pathSoFar[:len(pathSoFar)-1] {
		loc_x, loc_y := loc[0], loc[1]
		if loc_x == x && loc_y == y {
			loc_dx, loc_dy := directionSoFar[i][0], directionSoFar[i][1]
			if loc_dx == dx && loc_dy == dy {
				return true
			}
		}
	}

	return false
}

func PredictPathRecursive(labMap [][]string, pathSoFar [][]int, directionSoFar [][]int) ([][]int, bool) {
	x, y, dx, dy := WhereAmI(pathSoFar, directionSoFar)

	// Have we found a loop?
	if isLoop(pathSoFar, directionSoFar) {
		return nil, true
	}

	// Make a decision on where to go
	nextPosition, nextDirection := FindNextStep(labMap, x, y, dx, dy)

	// Base case, is next step outside of bounds? We are done.
	nx, ny := nextPosition[0], nextPosition[1]
	if ny < 0 || ny >= len(labMap) || nx < 0 || nx >= len(labMap[0]) {
		return pathSoFar, false
	}

	// Recursively create the path otherwise
	newPathSoFar := append(pathSoFar, nextPosition)
	newDirectionSoFar := append(directionSoFar, nextDirection)
	return PredictPathRecursive(labMap, newPathSoFar, newDirectionSoFar)
}

func FindGuard(labMap [][]string) (int, int) {
	for y, row := range labMap {
		for x, cell := range row {
			if cell == "^" {
				return x, y
			}
		}
	}
	return 0, 0
}

func PredictPath(labMap [][]string) ([][]int, bool) {
	// Let's find the guard first
	x, y := FindGuard(labMap)

	return PredictPathRecursive(labMap, [][]int{{x, y}}, [][]int{{0, -1}})
}

func sliceToString(slice []int) string {
	var sb strings.Builder
	for i, num := range slice {
		if i > 0 {
			sb.WriteString(",") // Delimiter
		}
		sb.WriteString(strconv.Itoa(num))
	}
	return sb.String()
}

func dedupeNestedSlices(nested [][]int) [][]int {
	seen := make(map[string]bool)
	var result [][]int

	for _, slice := range nested {
		key := sliceToString(slice)
		if !seen[key] {
			seen[key] = true
			result = append(result, slice)
		}
	}

	return result
}

func FindBlocksThatLoop(labMap [][]string) int {
	blockCounts := 0
	for y, row := range labMap {
		for x, cell := range row {
			if cell == "." {
				labMap[y][x] = "#"
				_, isLoop := PredictPath(labMap)
				if isLoop {
					blockCounts += 1
				}
				labMap[y][x] = "."
			}
		}
	}
	return blockCounts
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]

	// Read map
	labMap, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		os.Exit(1)
	}

	// Map path
	var path [][]int
	path, _ = PredictPath(labMap)

	// Compute unique positions
	dedupedPath := dedupeNestedSlices(path)
	distinctPositions := len(dedupedPath)

	fmt.Println("How many distinct positions?", distinctPositions)

	// Loop through map and find blockers that will create loops
	numberOfBlocks := FindBlocksThatLoop(labMap)
	fmt.Println("How many blocks that loop?", numberOfBlocks)
}
