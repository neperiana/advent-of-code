package main

import (
	// "bufio"
	"fmt"
	"io"
	// "math"
	"os"
	"regexp"
	// "strconv"
	// "strings"
)

func ReadFile(fileName string) (string, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read all the content
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Convert to string and return
	return string(data), nil
}

func ExtractRegex(rawInput string, regexPattern string) ([]string, error) {
	// Compile the regex
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}

	// Find all matches
	matches := re.FindAllString(rawInput, -1)
	return matches, nil
}

func ExecuteAndAddUp(commands []string) int {
	
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]
	gibberish, err:= ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading reports: %v\n", err)
		os.Exit(1)
	}

	// Extract valid commands
	pattern := `mul\((-?\d+),(-?\d+)\)`
	multCommands, err := ExtractRegex(gibberish, pattern)
	if err != nil {
		fmt.Printf("Error extracting regex: %v\n", err)
		os.Exit(1)
	}

	// Execute and add up
	totalResult := ExecuteAndAddUp(multCommands)

	// Print result
	fmt.Println("What's the result?:", totalResult)
}