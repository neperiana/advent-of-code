package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func ExecuteMult(command string) (int, error) {
	// Extract integers
	openIdx := strings.Index(command, "(")
	closeIdx := strings.Index(command, ")")
	content := command[openIdx+1 : closeIdx]
	parts := strings.Split(content, ",")

	// Convert to integer and multiply
	int1, err1 := strconv.Atoi(parts[0])
	int2, err2 := strconv.Atoi(parts[1])
	if err1 == nil && err2 == nil { // Ensure conversion was successful
		return int1 * int2, nil
	} else {
		return fmt.Printf("Error converting values: %v, %v\n", err1, err2)
	}
}

func RemoveDonts(commands []string) []string {
	keep := true
	clean_commands := make([]string, len(commands))
	copy(clean_commands, commands)
	removed := 0

	for i := 0; i < len(commands); i++ {
		command := commands[i]

		if command == "don't()" {
			keep = false
		} else if command == "do()" {
			keep = true
		}

		if !keep {
			// clean_commands = append(clean_commands, clean_commands[:i]...)
			j := i - removed
			clean_commands = append(clean_commands[:j], clean_commands[j+1:]...)
			removed++
		}
	}
	return clean_commands
}

func ExecuteAndAddUp(commands []string) (int, error) {
	var totalResult int = 0

	for i := 0; i < len(commands); i++ {
		if commands[i][:3] == "mul" {
			result, err := ExecuteMult(commands[i])
			if err == nil {
				totalResult += result
			} else {
				return 0, err
			}
		}
	}

	return totalResult, nil
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]
	gibberish, err := ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading reports: %v\n", err)
		os.Exit(1)
	}

	// Extract valid commands
	pattern := `(mul\((-?\d+),(-?\d+)\)|do\(\)|don't\(\))`
	commands, err := ExtractRegex(gibberish, pattern)
	if err != nil {
		fmt.Printf("Error extracting regex: %v\n", err)
		os.Exit(1)
	}

	// for part b), remove dont's
	do_commands := RemoveDonts(commands)

	// Execute and add up, part a
	totalResultA, err := ExecuteAndAddUp(commands)
	if err != nil {
		fmt.Printf("Error executing commands: %v\n", err)
		os.Exit(1)
	}
	// part b
	totalResultB, err := ExecuteAndAddUp(do_commands)
	if err != nil {
		fmt.Printf("Error executing commands: %v\n", err)
		os.Exit(1)
	}

	// // Print result
	fmt.Println("What's the result, part a?:", totalResultA)
	fmt.Println("What's the result, part b?:", totalResultB)
}
