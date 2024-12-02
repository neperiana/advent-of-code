package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func ReadReports(fileName string) [][]int {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	
	// Slice to store reports/levels
	var reports [][]int

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Define report variable
		var this_report []int

		// Read line and extract levels
		one_report := scanner.Text()
		levels := strings.Fields(one_report)

		// Convert to numeric
		for i := 0; i < len(levels); i++ {
			level, err := strconv.Atoi(levels[i])
			if err == nil {
				this_report = append(this_report, level)
			} else {
				fmt.Printf("Error converting values: %v\n", err)
			}
		}

		// Add to reports
		reports = append(reports, this_report)
	}

	return reports
}


func isThisReportSafe(report []int) bool {
	// Safety conditions to check
	var isNotAsc, isNotDesc, gapIsNotOK = false, false, false

	// Loop through report and check
	for i:=0; i<len(report)-1; i++ {
		diff := report[i+1]-report[i]
		absDiff := int(math.Abs(float64(diff)))

		if diff > 0 {
			isNotDesc = true
		}
		if diff < 0 {
			isNotAsc = true
		}
		if absDiff < 1 || absDiff > 3 {
			gapIsNotOK = true
		}
	}

	// Return True is Safe (all conditions false), False otherwise
	return !((isNotAsc && isNotDesc) || gapIsNotOK)
}


func ComputeSafetyAndAddUp(reports [][]int) int {
	var numberSafeReports = 0

	for i:=0; i<len(reports); i++ {
		if isThisReportSafe(reports[i]) {
			numberSafeReports += 1
		}
	}

	return numberSafeReports
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]
	reports:= ReadReports(fileName)

	// Calculate safety & add up
	numberSafeReports := ComputeSafetyAndAddUp(reports)

	// Print result
	fmt.Println("Number of safe reports?:", numberSafeReports)
}