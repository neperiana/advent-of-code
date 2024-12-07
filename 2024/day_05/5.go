package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Helper function to convert a slice of strings to a slice of ints
func convertToIntSlice(input []string) ([]int, error) {
	result := make([]int, len(input))
	for i, str := range input {
		value, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return nil, fmt.Errorf("error converting %q to int: %w", str, err)
		}
		result[i] = value
	}
	return result, nil
}

func createRuleMap(rules [][]int) map[int][]int {
	ruleMap := make(map[int][]int)
	for _, rule := range rules {
        ruleMap[rule[0]] = append(ruleMap[rule[0]], rule[1])
    }
	return ruleMap
}

func ReadFile(fileName string) (map[int][]int, [][]int, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Read all the rules and updates
	var rules [][]int
	var updates [][]int

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	blankFound := false
	for scanner.Scan() {
		// Read the line
		oneLine := scanner.Text()
		if blankFound { // add to updates
			update := strings.Split(oneLine, ",")
			intUpdate, err := convertToIntSlice(update)
			if err != nil {
				return nil, nil, err
			}
			updates = append(updates, intUpdate)
		} else if oneLine == "" {
			blankFound = true
		} else { // add to rules
			rule := strings.Split(oneLine, "|")
			intRule, err := convertToIntSlice(rule)
			if err != nil {
				return nil, nil, err
			}
			rules = append(rules, intRule)
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	// Transform rules into map
	ruleMap := createRuleMap(rules)

	return ruleMap, updates, nil
}

func Contains(slice []int, target int) (bool, int) {
	for i, v := range slice {
		if v == target {
			return true, i
		}
	}
	return false, -1
}

func IsUpdateSorted(update []int, rules map[int][]int) bool {
	// Check each rule
    for smaller, biggerList := range rules {
		for _, bigger := range biggerList {
			containsSmaller, smallerIndex := Contains(update, smaller)
			containsBigger, biggerIndex := Contains(update, bigger)
			if containsSmaller && containsBigger {
				if smallerIndex > biggerIndex {
					return false
				}
			}
		}
	}
	return true
}

func DivideUpdates(updates [][]int, rules map[int][]int) ([][]int, [][]int) {
	var sortedUpdates, unsortedUpdates [][]int 
	for _, update := range updates {
		if IsUpdateSorted(update, rules) {
			sortedUpdates = append(sortedUpdates, update)
		} else {
			unsortedUpdates = append(unsortedUpdates, update)
		}
	}
	return sortedUpdates, unsortedUpdates
} 

func AddUpMiddleElement(updates [][]int) int {
	var sum int = 0
	for _, update := range updates {
		l := len(update)
		middleElement := update[int(l/2)]
		sum += middleElement
	}
	return sum
}

func isSmaller(a int, b int, rules map[int][]int) bool {
	for smaller, biggerList := range rules {
		for _, bigger := range biggerList {
			if smaller == a && bigger == b {
				return true
			} else if smaller == b && bigger == a {
				return false
			}
		}
	}
	return a <= b
}

func SortUpdate(update []int, rules map[int][]int) []int {
	if len(update) <= 1 {
		return update  // Base case, slice is already sorted
	}
	
	// Use first element as pivot
	pivot := update[0]
	var less, greater []int

	// Compare the elements of the slice against the pivot
	for _, el := range update[1:] {
		if isSmaller(el, pivot, rules) {
			less = append(less, el)
		} else {
			greater = append(greater, el)
		}
	}

	// Recursively sort and combine
	lessSorted := SortUpdate(less, rules)
	greaterSorted := SortUpdate(greater, rules)
	return append(append(lessSorted, pivot), greaterSorted...)
}

func SortUpdates(updates [][]int, rules map[int][]int) [][]int {
	var sortedUpdates [][]int 
	for _, update := range updates {
		sortedUpdate := SortUpdate(update, rules)
		sortedUpdates = append(sortedUpdates, sortedUpdate)
	}
	return sortedUpdates
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]
	rules, updates, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		os.Exit(1)
	}
	
	// Check which updates are sorted
	sortedUpdates, unsortedUpdates := DivideUpdates(updates, rules)

	// Sort the unsorted ones
	sortedUnsortedUpdates := SortUpdates(unsortedUpdates, rules)

	// Sum middle page of sorted reports
	middleElementSortedSum := AddUpMiddleElement(sortedUpdates)
	middleElementUnsortedSum := AddUpMiddleElement(sortedUnsortedUpdates)
	fmt.Println("What's the sum of the middle element of sorted updates?", middleElementSortedSum)
	fmt.Println("What's the sum of the middle element of unsorted updates?", middleElementUnsortedSum)
}