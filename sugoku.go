package main

import (
	"fmt"
)

// very simple puzzle
// https://www.7sudoku.com/view-puzzle?date=20161117
/*
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
*/
// moderate puzzle
// https://www.7sudoku.com/view-puzzle?date=20200703

var puzzle = [9][9]int{
	{0, 0, 0, 1, 0, 0, 2, 0, 4},
	{0, 5, 6, 4, 0, 0, 0, 3, 0},
	{3, 0, 0, 0, 0, 6, 0, 9, 0},
	{0, 0, 0, 5, 8, 0, 0, 0, 0},
	{8, 0, 0, 2, 0, 1, 0, 0, 9},
	{0, 0, 0, 0, 4, 3, 0, 0, 0},
	{0, 8, 0, 3, 0, 0, 0, 0, 6},
	{0, 2, 0, 0, 0, 7, 4, 1, 0},
	{7, 0, 1, 0, 0, 5, 0, 0, 0}}

var puzzlePossibles [9][9][9]bool

func main() {
	printPuzzle(&puzzle, -1, -1)
	for {
		fullEval(&puzzle, &puzzlePossibles)
		deductMove := deductMoves(&puzzle, &puzzlePossibles)
		fullEval(&puzzle, &puzzlePossibles)
		bruteMove := bruteMoves(&puzzle, &puzzlePossibles)
		if !bruteMove && !deductMove {
			break
		}
	}
	printPuzzle(&puzzle, -1, -1)
}

func fullEval(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) {
	basicEval(puzzle, puzzlePossibles)
	rowEval(puzzle, puzzlePossibles)
	colEval(puzzle, puzzlePossibles)
	boxEval(puzzle, puzzlePossibles)
}

// sets each 0 space with possibles = true and filled spaces = false
func basicEval(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			placeVal := puzzle[y][x]
			for p := 0; p < 9; p++ {
				if placeVal == 0 {
					puzzlePossibles[y][x][p] = true
				} else {
					puzzlePossibles[y][x][p] = false
				}
			}
		}
	}
}

// sets each row with impossibles for already set values
func rowEval(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			placeVal := puzzle[y][x]
			if placeVal == 0 {
				continue
			} else {
				for xi := 0; xi < 9; xi++ {
					puzzlePossibles[y][xi][placeVal-1] = false
				}
			}
		}
	}
}

// sets each column with impossibles for already set values
func colEval(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) {
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			placeVal := puzzle[y][x]
			if placeVal == 0 {
				continue
			} else {
				for yi := 0; yi < 9; yi++ {
					puzzlePossibles[yi][x][placeVal-1] = false
				}
			}
		}
	}
}

func boxEval(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) {
	for boxY := 0; boxY < 3; boxY++ {
		boxYPos := boxY * 3
		for boxX := 0; boxX < 3; boxX++ {
			boxXPos := boxX * 3
			for y := boxYPos; y < boxYPos+3; y++ {
				for x := boxXPos; x < boxXPos+3; x++ {
					placeVal := puzzle[y][x]
					if placeVal == 0 {
						continue
					} else {
						for yi := boxYPos; yi < boxYPos+3; yi++ {
							for xi := boxXPos; xi < boxXPos+3; xi++ {
								puzzlePossibles[yi][xi][placeVal-1] = false
							}
						}
					}
				}
			}
		}
	}
}

// finds positions with only 1 move left possible, makes move
func bruteMoves(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) bool {
	madeMove := false
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			possibleCount := 0
			possibleVal := 0
			for p := 0; p < 9; p++ {
				if puzzlePossibles[y][x][p] == true {
					possibleCount++
					possibleVal = p + 1
				} else {
					continue
				}
			}
			if possibleCount == 1 {
				puzzle[y][x] = possibleVal
				puzzlePossibles[y][x][possibleVal-1] = false
				madeMove = true
			}
		}
	}
	return madeMove
}

func deductMoves(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) bool {
	boolRow := deductRow(puzzle, puzzlePossibles)
	fullEval(puzzle, puzzlePossibles)
	boolColumn := deductColumn(puzzle, puzzlePossibles)
	fullEval(puzzle, puzzlePossibles)
	boolBox := deductBox(puzzle, puzzlePossibles)
	fullEval(puzzle, puzzlePossibles)
	if boolRow || boolColumn || boolBox {
		return true
	}
	return false
}

func deductRow(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) bool {
	// if a number only fits one space in a row, place it
	var xOpenIndex int
	madeMove := false
	for n := 1; n < 10; n++ {
		pIndex := n - 1
		for y := 0; y < 9; y++ {
			xOpenings := 0
			for x := 0; x < 9; x++ {
				if puzzlePossibles[y][x][pIndex] == true {
					xOpenings++
					xOpenIndex = x
				}
			}
			if xOpenings == 1 {
				puzzle[y][xOpenIndex] = n
				puzzlePossibles[y][xOpenIndex] = [9]bool{false, false, false, false,
					false, false, false, false, false}
				madeMove = true
			}
		}
	}
	return madeMove
}

func deductColumn(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) bool {
	// if a number only fits one space in a col, place it
	madeMove := false
	var yOpenIndex int
	for n := 1; n < 10; n++ {
		pIndex := n - 1
		for x := 0; x < 9; x++ {
			yOpenings := 0
			for y := 0; y < 9; y++ {
				if puzzlePossibles[y][x][pIndex] == true {
					yOpenings++
					yOpenIndex = y
				}
			}
			if yOpenings == 1 {
				puzzle[yOpenIndex][x] = n
				puzzlePossibles[yOpenIndex][x] = [9]bool{false, false, false, false,
					false, false, false, false, false}
				madeMove = true
			}
		}
	}
	return madeMove
}

func deductBox(puzzle *[9][9]int, puzzlePossibles *[9][9][9]bool) bool {
	// if a number only fits one space in a box, place it
	madeMove := false
	var nOpenIndex [2]int
	for n := 1; n < 10; n++ {
		pIndex := n - 1
		for bY := 0; bY < 3; bY++ {
			for bX := 0; bX < 3; bX++ {
				yTop := bY * 3
				yBot := yTop + 3
				xTop := bX * 3
				xBot := xTop + 3
				nOpenings := 0
				for y := yTop; y < yBot; y++ {
					for x := xTop; x < xBot; x++ {
						if puzzlePossibles[y][x][pIndex] == true {
							nOpenings++
							nOpenIndex = [2]int{y, x}
						}
					}
				}
				if nOpenings == 1 {
					puzzle[nOpenIndex[0]][nOpenIndex[1]] = n
					puzzlePossibles[nOpenIndex[0]][nOpenIndex[1]] = [9]bool{false, false, false, false,
						false, false, false, false, false}
					madeMove = true
				}
			}
		}
	}
	return madeMove
}

func printPuzzle(puzzle *[9][9]int, yHighlight, xHighlight int) {
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
