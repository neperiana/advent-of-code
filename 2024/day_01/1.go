// Pair up the numbers and measure how far apart they are. Pair up the smallest number in the left list with the smallest number in the right list, 
// then the second-smallest left number with the second-smallest right number, and so on.
// Within each pair, figure out how far apart the two numbers are; you'll need to add up 
// all of those distances. For example, if you pair up a 3 from the left list with a 7 from the right list, 
// the distance apart is 4; if you pair up a 9 with a 3, the distance apart is 6.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ReadColumns(fileName string) ([]int, []int) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil, nil
	}
	defer file.Close()

	// Arrays to store columns
	var column1 []int
	var column2 []int

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		// Get the current line
		line := scanner.Text()

		// Split the line into columns
		fields := strings.Fields(line) // Splits on whitespace

		// Convert and append each column's values to respective arrays
		if len(fields) == 2 { // Ensure there are exactly 2 columns
			col1, err1 := strconv.Atoi(fields[0])
			col2, err2 := strconv.Atoi(fields[1])
			if err1 == nil && err2 == nil { // Ensure conversion was successful
				column1 = append(column1, col1)
				column2 = append(column2, col2)
			} else {
				fmt.Printf("Error converting values: %v, %v\n", err1, err2)
			}
		} else {
			fmt.Printf("Unexpected format in line: %s\n", line)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	return column1, column2
}

func ComputeDifferencesAndAddUp(list1, list2 []int) ([]int, int) {
	// Loop through list and compute difference
	var l int = len(list1)
	var totalDiff int = 0
	var differences []int

	for i := 0; i < l; i++ {
        diff := int(math.Abs(float64(list1[i] - list2[i])))
		differences = append(differences, diff)
		totalDiff += diff
    }

	return differences, totalDiff
}

func countAllOccurrences(element int, list2 []int) int {
	var totalOccurrences int = 0

	for i := 0; i < len(list2); i++ {
		if list2[i] == element {
			totalOccurrences += 1
		}
    }
	return totalOccurrences
}

func ComputeSimilarityScore(list1, list2 []int) int {
	var l int = len(list1)
	var totalsimScore int = 0

	for i := 0; i < l; i++ {
		element := list1[i]
        elementCount := countAllOccurrences(element, list2)
		simScore := element * elementCount
		totalsimScore += simScore
    }

	return totalsimScore
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]
	list1, list2 := ReadColumns(fileName)
	
	// Sort lists
	sort.Ints(list1)
	sort.Ints(list2)

	// Calculate differences & add up
	_, totalDiff := ComputeDifferencesAndAddUp(list1, list2)

	// Calculate similarity score
	simScore := ComputeSimilarityScore(list1, list2)

	// Print result
	fmt.Println("Sum of differences:", totalDiff)
	fmt.Println("Similarity score:", simScore)
}