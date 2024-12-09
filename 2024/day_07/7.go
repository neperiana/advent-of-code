package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Helper function to convert a slice of strings to a slice of ints
func convertToIntSlice(input []string) ([]int, error) {
	result := []int{}
	for _, str := range input {
		if str != "" {
			value, err := strconv.Atoi(strings.TrimSpace(str))
			if err != nil {
				return nil, fmt.Errorf("error converting %q to int: %w", str, err)
			}
			result = append(result, value)
		}
	}
	return result, nil
}

func ReadFile(fileName string) ([]int, [][]int, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Read all the rules and updates
	var results []int
	var operands [][]int

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read the line
		oneLine := scanner.Text()
		oneRow := strings.Split(oneLine, ":")

		result, err := convertToIntSlice([]string{oneRow[0]})
		if err != nil {
			return nil, nil, err
		}

		operand, err := convertToIntSlice(strings.Split(oneRow[1], " "))
		if err != nil {
			return nil, nil, err
		}

		results = append(results, result...)
		operands = append(operands, operand)
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return results, operands, nil
}

func RoundUpToPowerOf10(n int) int {
	if n <= 0 {
		return 0 // Handle zero or negative numbers gracefully
	}
	power := math.Ceil(math.Log10(float64(n + 1))) // Shift by 1 to handle edge cases
	return int(math.Pow(10, power))
}

func containsInt(slice []int, element int) bool {
	for _, s := range slice {
		if element == s {
			return true
		}
	}
	return false
}

func isFeasible(result int, operands []int, includeConcat bool) bool {
	possibleResults := []int{operands[0]}

	for i := 1; i < len(operands); i++ {
		possibleResultsNew := []int{}
		for _, val := range possibleResults {
			possibleResultsNew = append(possibleResultsNew, val*operands[i])
			possibleResultsNew = append(possibleResultsNew, val+operands[i])

			if includeConcat {
				factor := RoundUpToPowerOf10(operands[i])
				possibleResultsNew = append(possibleResultsNew, val*factor+operands[i])
			}

		}
		possibleResults = possibleResultsNew
	}

	return containsInt(possibleResults, result)
}

func ComputeFeasibility(results []int, operands [][]int, includeConcat bool) []int {
	feasibleIndices := []int{}
	for i := 0; i < len(results); i++ {
		if isFeasible(results[i], operands[i], includeConcat) {
			feasibleIndices = append(feasibleIndices, i)
		}
	}
	return feasibleIndices
}

func AddUp(results []int, indeces []int) int64 {
	var totalSum int64 = 0
	for _, i := range indeces {
		totalSum += int64(results[i])
	}
	return totalSum
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]

	// Read map
	results, operands, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		os.Exit(1)
	}

	// Return index for those that are feasible
	feasibleIndeces := ComputeFeasibility(results, operands, false)

	// Add up results
	totalSum := AddUp(results, feasibleIndeces)

	fmt.Println("What's the total?", totalSum)

	// Including concats
	feasibleIndeces = ComputeFeasibility(results, operands, true)
	totalSum = AddUp(results, feasibleIndeces)
	fmt.Println("What's the total (including concat)?", totalSum)

}
