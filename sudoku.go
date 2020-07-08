package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func mainOld() {
	fmt.Println("starting sudoku solve.")
	startTime := time.Now()
	var puzzleImpossibles [9][9][10]bool
	/* this first puzzle is actually rather dificult, to the level
	   that you need to know what the interactions between the possibilities
	   are in order to solve it, I will eventually figure out how to make
	   it solve this one, but for now I'm just going to do simpler puzzles. */

	/*var puzzle = [9][9]int{ // extremely hard puzzle
	    {0,0,0,0,3,6,0,5,0},
	    {3,0,0,0,0,0,4,0,0},
	    {0,0,5,8,0,1,0,3,0},
	    {9,0,6,0,0,3,0,0,0},
	    {0,2,0,0,0,0,0,8,0},
	    {0,0,0,4,0,0,5,0,6},
	    {0,1,0,7,0,8,3,0,0},
	    {0,0,7,0,0,0,0,0,9},
	    {0,4,0,1,6,0,0,0,0},
	} //http://www.7sudoku.com/view-solution?date=20161109 */
	var puzzle = [9][9]int{ // very simple puzzle
		{6, 0, 1, 0, 4, 0, 0, 5, 0},
		{0, 0, 2, 0, 5, 0, 4, 8, 0},
		{0, 0, 4, 0, 0, 1, 0, 0, 6},
		{0, 0, 0, 0, 0, 7, 8, 4, 0},
		{1, 0, 0, 0, 8, 0, 0, 0, 9},
		{0, 9, 8, 4, 0, 0, 0, 0, 0},
		{4, 0, 0, 7, 0, 0, 5, 0, 0},
		{0, 3, 9, 0, 2, 0, 1, 0, 0},
		{0, 2, 0, 0, 1, 0, 9, 0, 4},
	} //http://www.7sudoku.com/view-puzzle?date=20161117
	/* var puzzle = [9][9]int{
		{8, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 3, 6, 0, 0, 0, 0, 0},
		{0, 7, 0, 0, 9, 0, 2, 0, 0},
		{0, 5, 0, 0, 0, 7, 0, 0, 0},
		{0, 0, 0, 0, 4, 5, 7, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 3, 0},
		{0, 0, 1, 0, 0, 0, 0, 6, 8},
		{0, 0, 8, 5, 0, 0, 0, 1, 0},
		{0, 9, 0, 0, 0, 0, 4, 0, 0},
	} // from colleague's code */

	hasZeros := true // the puzzle has zeros at the beginning
	iterStop := 0    //fuck vs code, I don't trust it to not lock with the integrated terminal

	for hasZeros == true { // run until there's no more zeros
		hasZeros = false //since we don't know if it has zeros, assume none, prove existence

		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				//this loops over the inner dimension
				//  for each individual value:

				//if there's only one possible value, insert it
				countFalses := 0
				lastFalse := -1
				for k2 := 0; k2 < 9; k2++ {
					if puzzleImpossibles[i][j][k2] == false {
						countFalses++
						lastFalse = k2 + 1
					}
				}
				if countFalses == 1 {
					puzzle[i][j] = lastFalse
					fmt.Println("SOLVED [", i, "][", j, "] with brute force for :", lastFalse)
					printPuzzle(puzzle, i, j)
				}

				//if it's solved keep going

				//if it's not 0 flush all the possibilities in row, col, box
				if puzzle[i][j] == 0 {
					hasZeros = true //this is how the main loop knows the puzzle isn't solved
					//row, column, box
					if false {
						fmt.Println("puzzle[", i, "][", j, "] is zero.")
						for k := 0; k < 9; k++ {
							fmt.Println((k + 1), " impossible? : ", puzzleImpossibles[i][j][k])
						}
					}
				} else if puzzleImpossibles[i][j][9] == false {
					//don't flush things twice, 3 dimension has extra val to signal flush already done
					puzzleImpossibles[i][j][9] = true
					//flush impossible values
					ind3 := puzzle[i][j] - 1 // this is the index of the 3rd dimension that changes
					//first for the column (i) hold j constant
					for i2 := 0; i2 < 9; i2++ {
						puzzleImpossibles[i2][j][ind3] = true
					}
					//fmt.Println("set col",j+1,"impossible for",ind3+1)
					//second for the row (j) hold i constant
					for j2 := 0; j2 < 9; j2++ {
						puzzleImpossibles[i][j2][ind3] = true
					}
					//fmt.Println("set row",i+1,"impossible for",ind3+1)
					//now for the box (2 dimensions)
					iTop := i - (i % 3)
					iBot := iTop + 2
					jTop := j - (j % 3)
					jBot := jTop + 2
					for i3 := iTop; i3 <= iBot; i3++ {
						for j3 := jTop; j3 <= jBot; j3++ {
							puzzleImpossibles[i3][j3][ind3] = true
						}
					}
					//fmt.Println("set box i[",iTop+1,"-",iBot+1,"] j[",jTop+1,"-",jBot+1,"] impossble for",ind3+1)

					// set all of the impossibilities true because this one is solved
					puzzleImpossibles[i][j] = [10]bool{true, true, true, true, true, true, true, true, true, true}
				}

			}
		}
		iterStop++
		if iterStop > 1000000 {
			fmt.Println("iterstop break! :(")
			break
		}

		// DO THE FANCY CHECKS FOR SOLUTIONS HERE

		//if there's only one place in a row where a number fits, insert
		for n := 1; n <= 9; n++ {
			for i := 0; i < 9; i++ {

				nFalses := 0
				lastNFalse := -1
				for j := 0; j < 9; j++ {
					if puzzleImpossibles[i][j][n-1] == false {
						nFalses++
						lastNFalse = j
					}
				}
				if nFalses == 1 {
					puzzle[i][lastNFalse] = n
					fmt.Println("SOLVED [", i+1, "][", lastNFalse+1, "] with ROW CHECK for :", n)
					printPuzzle(puzzle, i, lastNFalse)
				}
			}
		}

		//if there's only one place in a col where a number fits, insert it
		for n := 1; n <= 9; n++ {
			for j := 0; j < 9; j++ {
				nFalses := 0
				lastNFalse := -1
				for i := 0; i < 9; i++ {
					if puzzleImpossibles[i][j][n-1] == false {
						nFalses++
						lastNFalse = i
					}
				}
				if nFalses == 1 {
					puzzle[lastNFalse][j] = n
					fmt.Println("SOLVED [", lastNFalse+1, "][", j+1, "] with COL CHECK for :", n)
					printPuzzle(puzzle, lastNFalse, j)
				}
			}
		}

		//if there's only one place in a box where a number fits, insert it.
		for n := 1; n <= 9; n++ {
			//fmt.Println("starting box check for number: ",n)
			for bRow := 0; bRow < 3; bRow++ {
				//fmt.Println("box row: ",bRow)
				for bCol := 0; bCol < 3; bCol++ {
					//fmt.Println("box col: ",bCol)
					nFalses := 0
					lastNFalse := [2]int{-1, -1}
					for i := bRow * 3; i < ((bRow * 3) + 3); i++ {
						for j := bCol * 3; j < ((bCol * 3) + 3); j++ {

							//fmt.Printf("index value: [%d][%d] ",i,j)
							if puzzleImpossibles[i][j][n-1] == false {
								//fmt.Println(" IS FALSE")
								nFalses++
								lastNFalse[0] = i
								lastNFalse[1] = j
							} else {
								//fmt.Println("")
							}
						}
					}
					if nFalses == 1 {
						puzzle[lastNFalse[0]][lastNFalse[1]] = n
						fmt.Println("SOLVED [", lastNFalse[0]+1, "][", lastNFalse[1]+1, "] with BOX CHECK for :", n)
						printPuzzle(puzzle, lastNFalse[0], lastNFalse[1])
					}
				}
			}
		}

	}
	//end outer for
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("took ", duration.Seconds(), "s to finish")
	//put code here for printing the solution!
	printPuzzle(puzzle, -1, -1)

}

func printPuzzleOld(puzzle [9][9]int, iHl int, jHl int) {
	for i := 0; i < 9; i++ {
		if i%3 == 0 {
			fmt.Println(" +-------+-------+-------+")
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 {
				fmt.Printf(" |")
			}
			if (puzzle[i][j] != 0) && ((i != iHl) || (j != jHl)) {
				fmt.Printf(" %d", puzzle[i][j])
			} else if (i == iHl) && (j == jHl) {
				fmt.Printf("\x1b[31;1m %d\x1b[0m", puzzle[i][j])
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Println(" |")
	}
	fmt.Println(" +-------+-------+-------+")
	//hold it blocked until press enter
	if false {
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
	}
	//end of printPuzzle
}
