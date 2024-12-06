package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadMatrix(fileName string) ([][]string, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read all the content
	var M [][]string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read the line
		one_report := scanner.Text()
		row := strings.Split(one_report, "")
		M = append(M, row)
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return M, nil
}

func findStartOfChain(M [][]string, target string) [][]int {
	// Find all the Xs in the matrix
	var Xs [][]int
	for i, row := range M {
		for j, cell := range row {
			if cell == target {
				Xs = append(Xs, []int{i, j})
			}
		}
	}
	return Xs
}

var directions = [][]int{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {1, 1}, {-1, -1}, {1, -1}, {-1, 1},
}

func scanChain(M [][]string, i int, j int, chain string, chainCoords string, direction []int, targetChain string) int {
	// Recursive function that scans sorroundings of i,j to see how many target chains are present
	// We keep track of the current chain in the chain variable
	if chain == targetChain {
		return 1
	} // Assuming we are in the middle of the chain, check if the current cell is part of the chain
	// Check if i, j is within bounds
	if i < 0 || i >= len(M) || j < 0 || j >= len(M[0]) {
		return 0
	} // Check if we are at the end of the chain
	if M[i][j][0] != targetChain[len(chain)] {
		return 0
		// We could be at the start of the chain
	} else if M[i][j][0] == targetChain[0] {
		// Check all the sorroundings
		var count int = 0
		var chainCoords string = fmt.Sprintf("[%d,%d]", i, j)

		for _, dir := range directions {
			count += scanChain(
				M,
				i+dir[0],
				j+dir[1],
				chain+M[i][j],
				chainCoords,
				dir,
				targetChain,
			)
		}
		return count
		// otherwise, we are in the middle of the chain so let's keep scanning in the same direction
	} else {
		var newChainCoords string = fmt.Sprintf("[%d,%d]", i, j)
		return scanChain(
			M,
			i+direction[0],
			j+direction[1],
			chain+M[i][j],
			chainCoords+newChainCoords,
			direction,
			targetChain,
		)
	}

}

func CountChainInstances(M [][]string, Xs [][]int, targetChain string) int {
	// Scan if X is part of XMAS
	var count int = 0
	for _, X := range Xs {
		// Check if X is part of XMAS
		i, j := X[0], X[1]
		count += scanChain(M, i, j, "", "", nil, targetChain)
	}
	return count
}

func scanCell(M [][]string, i int, j int, target string) bool {
	// Function that checks if the content of cell i,j is equal to target
	// Check if i, j is within bounds
	if i < 0 || i >= len(M) || j < 0 || j >= len(M[0]) {
		return false
	}
	// Check if cell is equal to target
	if M[i][j] == target {
		return true
	}
	return false
}

func scanDiagonalForMas(M [][]string, i int, j int, direction []int) int {
	// Function that checks if a MAS chain is present in the diagonal specified by direction
	// Returns 1 if so, 0 otherwise

	// Check if i,j + direction is M
	if scanCell(M, i+direction[0], j+direction[1], "M") {
		// Check if i,j - direction is S
		if scanCell(M, i-direction[0], j-direction[1], "S") {
			return 1
		}
	}
	return 0
}

var masDirections = [][]int{
	{1, 1}, {-1, -1}, {1, -1}, {-1, 1},
}

func CountMasCrossInstances(M [][]string, As [][]int) int {
	// Scan if A is part of a MAS cross
	var count int = 0
	for _, A := range As {
		// Check if A is part of XMAS
		masCount := 0

		// Check all the sorroundings
		for _, dir := range masDirections {
			i, j := A[0], A[1]
			masCount += scanDiagonalForMas(
				M,
				i,
				j,
				dir,
			)
		}
		if masCount == 2 {
			count++
		}
	}
	return count
}

func main() {
	// Read the file name from command-line arguments & read input file
	fileName := os.Args[1]
	M, err := ReadMatrix(fileName)
	if err != nil {
		fmt.Printf("Error reading data: %v\n", err)
		os.Exit(1)
	}

	// Find all the Xs in the matrix
	Xs := findStartOfChain(M, "X")

	// Scan if X is part of XMAS
	numberOfXmas := CountChainInstances(M, Xs, "XMAS")

	fmt.Printf("Number of XMAS chains: %d\n", numberOfXmas)

	// Find all A's in the matrix
	As := findStartOfChain(M, "A")

	// Scan if A is part of a X-MAS
	numberOfMasCross := CountMasCrossInstances(M, As)

	fmt.Printf("Number of MAS crosses: %d\n", numberOfMasCross)
}
