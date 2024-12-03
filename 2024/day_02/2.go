package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
    MinGap = 1
    MaxGap = 3
)

func ReadReports(fileName string) ([][]int, error)  {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
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
				return nil, err
			}
		}

		// Add to reports
		reports = append(reports, this_report)
	}
	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}


func isSafe(report []int) bool {
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

func isAlmostSafe(report []int) bool {
	// Remove one level at a time and check safety
	for i:=0; i<len(report); i++ {

		newReport := make([]int, 0)
		newReport = append(newReport, report[:i]...)
		newReport = append(newReport, report[i+1:]...)

		if isSafe(newReport) {
			return true
		}
	}
	return false
}


func CountSafety(reports [][]int) (int, int) {
	var numberSafeReports, numberAlmostSafeReports = 0, 0

	for i:=0; i<len(reports); i++ {
		if isSafe(reports[i]) {
			numberSafeReports++
        	numberAlmostSafeReports++
		} else if isAlmostSafe(reports[i]) {
			numberAlmostSafeReports++
		}
		
	}

	return numberSafeReports, numberAlmostSafeReports
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]
	reports, err:= ReadReports(fileName)
	if err != nil {
		fmt.Printf("Error reading reports: %v\n", err)
		os.Exit(1)
	}

	// Calculate safety & add up
	numberSafeReports, numberAlmostSafeReports := CountSafety(reports)

	// Print result
	fmt.Println("Number of safe reports?:", numberSafeReports)
	fmt.Println("Number of almost safe reports?:", numberAlmostSafeReports)
}