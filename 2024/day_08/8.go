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
	var antennaMap [][]string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read the line
		oneLine := scanner.Text()
		onwRow := strings.Split(oneLine, "")
		antennaMap = append(antennaMap, onwRow)
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return antennaMap, nil
}

func findAntennas(antennaMap [][]string) map[string][][]int {
	antennaLocations := make(map[string][][]int) // Use make for initialization

	// Iterate through the map
	for y, row := range antennaMap {
		for x, freq := range row {
			if freq != "." {
				// Append the location to the corresponding frequency
				antennaLocations[freq] = append(antennaLocations[freq], []int{x, y})
			}
		}
	}

	return antennaLocations
}

func computeDistance(x, y []int) []int {
	return []int{x[0] - y[0], x[1] - y[1]}
}

func addVectors(operation string, x, y []int) []int {
	if operation == "+" {
		return []int{x[0] + y[0], x[1] + y[1]}
	} else if operation == "-" {
		return []int{x[0] - y[0], x[1] - y[1]}
	}
	return nil
}

func multVectors(factor int, x []int) []int {
	return []int{x[0] * factor, x[1] * factor}
}

func isWithinBoundaries(location []int, lenX, lenY int) bool {
	if location[0] < 0 || location[0] >= lenX {
		return false
	} else if location[1] < 0 || location[1] >= lenY {
		return false
	}
	return true
}

func placeAntinodesForTwoAntennas(antennaA, antennaB []int, lenX, lenY, distFactor int) [][]int {
	distance := computeDistance(antennaA, antennaB)
	antinodesLocation := make([][]int, 0, 2) // Preallocate space for up to 2 antinodes

	// Compute potential antinodes
	antinodeA := addVectors("+", antennaA, multVectors(distFactor, distance))
	antinodeB := addVectors("-", antennaB, multVectors(distFactor, distance))

	// Add valid antinodes to the result
	if isWithinBoundaries(antinodeA, lenX, lenY) {
		antinodesLocation = append(antinodesLocation, antinodeA)
	}
	if isWithinBoundaries(antinodeB, lenX, lenY) {
		antinodesLocation = append(antinodesLocation, antinodeB)
	}

	return antinodesLocation
}

func placeAntinodesForSingleFreq(locations [][]int, lenX, lenY int, infinite bool) [][]int {
	antinodesLocations := make([][]int, 0) // Use make for better memory allocation

	for _, pivot := range locations {
		for _, loc := range locations {
			// Skip if pivot and loc are aligned
			if pivot[0] == loc[0] || pivot[1] == loc[1] {
				continue
			}

			if infinite {
				// Generate antinodes for infinite case
				for distance := 0; ; distance++ {
					antinodes := placeAntinodesForTwoAntennas(pivot, loc, lenX, lenY, distance)
					if len(antinodes) == 0 {
						break
					}
					antinodesLocations = append(antinodesLocations, antinodes...)
				}
			} else {
				// Generate antinodes for finite case
				antinodes := placeAntinodesForTwoAntennas(pivot, loc, lenX, lenY, 1)
				antinodesLocations = append(antinodesLocations, antinodes...)
			}
		}
	}

	return antinodesLocations
}

func placeAntinodes(antennaLocations map[string][][]int, lenX, lenY int, infinite bool) map[string][][]int {
	antinodesLocations := make(map[string][][]int)

	for frequency, locations := range antennaLocations {
		antinodesLocations[frequency] = placeAntinodesForSingleFreq(locations, lenX, lenY, infinite)
	}

	return antinodesLocations
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

func dedupeNestedSlices(nested map[string][][]int) [][]int {
	seen := make(map[string]struct{}) // Use struct{} for a more efficient memory footprint
	result := make([][]int, 0)

	for _, slices := range nested {
		for _, loc := range slices {
			key := sliceToString(loc)
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				result = append(result, loc)
			}
		}
	}

	return result
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]

	// Read map
	antennaMap, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		os.Exit(1)
	}

	// Find frequencies and locations
	antennaLocations := findAntennas(antennaMap)

	// Map antinodes
	lenX, lenY := len(antennaMap[0]), len(antennaMap)
	antinodesLocations := placeAntinodes(antennaLocations, lenX, lenY, false)
	// fmt.Println(antinodesLocations)

	// Count unique locations
	uniqueAntinodes := dedupeNestedSlices(antinodesLocations)
	uniqueAntinodeCount := len(uniqueAntinodes)

	fmt.Println("How many antinodes?", uniqueAntinodeCount)

	// Map ALL antinodes
	antinodesInfLocations := placeAntinodes(antennaLocations, lenX, lenY, true)
	// fmt.Println(antinodesLocations)

	// Count unique locations
	uniqueInfAntinodes := dedupeNestedSlices(antinodesInfLocations)
	uniqueInfAntinodeCount := len(uniqueInfAntinodes)

	fmt.Println("How many antinodes (infinite)?", uniqueInfAntinodeCount)
}
