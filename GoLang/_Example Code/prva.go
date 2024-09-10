package main

import (
	"fmt"
	"sort"
	"strings"
)

func transpose(a [][]string) [][]string {
	result := make([][]string, len(a[0]))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			result[j] = append(result[j], a[i][j])
		}
	}
	return result
}

func findWords(rows, columns int, matrix [][]string, words []string) []string {
	var current string
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if matrix[i][j] == "#" {
				if len(current) > 1 {
					words = append(words, current)
				}
				current = ""
			} else if j == columns-1 {
				if len(current) > 0 {
					words = append(words, current+matrix[i][j])
				}
				current = ""
			} else {
				current = current + matrix[i][j]
			}
		}
	}
	return words
}

func main() {
	var rows, columns int
	fmt.Scanln(&rows, &columns) // read first line (matrix size)

	matrix := make([][]string, rows)

	for i := 0; i < rows; i++ { // we scan all lines
		var s string
		matrix[i] = make([]string, columns)
		fmt.Scanln(&s)
		for j := 0; j < columns; j++ {
			var temp []string
			temp = strings.Split(s, "")
			matrix[i][j] = temp[j]
		}
	}

	// now we find all words
	var words []string
	words = findWords(rows, columns, matrix, words)
	words = findWords(columns, rows, (transpose(matrix)), words)

	// finally we sort
	sort.Strings(words)
	fmt.Println(words[0])
}
