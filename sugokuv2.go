package main

import (
	"fmt"
)

// very simple puzzle
//http://www.7sudoku.com/view-puzzle?date=20161117
var puzzle = [9][9]int{
	{6, 0, 1, 0, 4, 0, 0, 5, 0},
	{0, 0, 2, 0, 5, 0, 4, 8, 0},
	{0, 0, 4, 0, 0, 1, 0, 0, 6},
	{0, 0, 0, 0, 0, 7, 8, 4, 0},
	{1, 0, 0, 0, 8, 0, 0, 0, 9},
	{0, 9, 8, 4, 0, 0, 0, 0, 0},
	{4, 0, 0, 7, 0, 0, 5, 0, 0},
	{0, 3, 9, 0, 2, 0, 1, 0, 0},
	{0, 2, 0, 0, 1, 0, 9, 0, 4}}

func main() {
	printPuzzle(puzzle, -1, -1)
}

func printPuzzle(puzzle [9][9]int, yHighlight, xHighlight int) {
	for y := 0; y < 9; y++ {
		if y%3 == 0 {
			fmt.Println(" +-------+-------+-------+")
		}
		for x := 0; x < 9; x++ {
			if x%3 == 0 {
				fmt.Printf(" |")
			}
			if (puzzle[y][x] != 0) && ((y != yHighlight) || (x != xHighlight)) {
				fmt.Printf(" %d", puzzle[y][x])
			} else if (y == yHighlight) && (x == xHighlight) {
				// print with red font
				fmt.Printf("\x1b[31;1m %d\x1b[0m", puzzle[y][x])
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Println(" |")
	}
	fmt.Println(" +-------+-------+-------+")
}
