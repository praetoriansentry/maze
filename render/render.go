package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	MAZE_SIZE       = 20
	MAX_TOKEN_COUNT = 256
	FIELD_SIZE      = MAZE_SIZE * MAZE_SIZE

	PRINTED_MAZE_COLS = (MAZE_SIZE*4 + 2)
	PRINTED_MAZE_ROWS = (2*MAZE_SIZE + 1)

	NORTH_WALL = 1
	SOUTH_WALL = 2
	EAST_WALL  = 4
	WEST_WALL  = 8
	VISITED    = 16

	CHAR_PLUS    = 0x2B
	CHAR_DASH    = 0x2D
	CHAR_NEWLINE = 0x0A
	CHAR_PIPE    = 0x7C
	CHAR_SPACE   = 0x20
)

func main() {
	for i := 0; i < 256; i = i + 1 {
		draw(i)
	}
}

func render(bitfield [FIELD_SIZE]int) string {
	var cell int
	var cellIndex int
	var row int
	var col int

	mazeOutput := make([]byte, PRINTED_MAZE_COLS*PRINTED_MAZE_ROWS)
	// For each row
	for i := 0; i < PRINTED_MAZE_ROWS; i++ {
		j := 0
		for j < PRINTED_MAZE_COLS {
			idx := i*PRINTED_MAZE_COLS + j
			// END OF THE ROW
			if j+1 == PRINTED_MAZE_COLS {
				mazeOutput[idx] = CHAR_NEWLINE
				j = j + 1
				continue
			}

			// START OF THE ROW
			if j == 0 && i%2 == 1 {
				mazeOutput[idx] = CHAR_PIPE
				j = j + 1
				continue
			} else if j == 0 && i%2 == 0 {
				mazeOutput[idx] = CHAR_PLUS
				j = j + 1
				continue
			}

			if i == 0 {
				cell = 31
			} else {
				row = (i - 1) / 2
				col = j / 4
				cellIndex = rowColToBitIndex(row, col)
				if cellIndex < FIELD_SIZE && cellIndex >= 0 {
					cell = bitfield[cellIndex]
				} else {
					cell = 31
				}
			}

			// What cell are we in?
			if i%2 == 1 { // EAST WEST
				mazeOutput[idx] = CHAR_SPACE
				idx += 1
				mazeOutput[idx] = CHAR_SPACE
				idx += 1
				mazeOutput[idx] = CHAR_SPACE
				idx += 1
				if (cell & EAST_WALL) == EAST_WALL {
					mazeOutput[idx] = CHAR_PIPE
				} else {
					mazeOutput[idx] = CHAR_SPACE
				}
				j = j + 4
			} else { // NORTH SOUTH
				if (cell & SOUTH_WALL) == SOUTH_WALL {
					mazeOutput[idx] = CHAR_DASH
					idx += 1
					mazeOutput[idx] = CHAR_DASH
					idx += 1
					mazeOutput[idx] = CHAR_DASH
					idx += 1
				} else {
					mazeOutput[idx] = CHAR_SPACE
					idx += 1
					mazeOutput[idx] = CHAR_SPACE
					idx += 1
					mazeOutput[idx] = CHAR_SPACE
					idx += 1
				}
				mazeOutput[idx] = CHAR_PLUS
				j = j + 4
			}
		}
	}

	finalMaze := string(mazeOutput)
	return finalMaze
}

// Take a token id and draw the maze that's unique to that token.
func draw(id int) string {
	os.Mkdir(fmt.Sprintf("out/%03d/", id), os.ModePerm)
	copy(fmt.Sprintf("masks/%d.jpg", id), fmt.Sprintf("out/%03d/mask.jpg", id))
	var bitfield [FIELD_SIZE]int = getMazeData(id)
	return render(bitfield)
}

// Implementing the fastest and most trivial algorithm I could find.
// https://weblog.jamisbuck.org/2011/2/1/maze-generation-binary-tree-algorithm
func getMazeData(id int) [FIELD_SIZE]int {
	var bitfield [FIELD_SIZE]int
	var cell int
	randInt := semiRandomData(id)
	var goSouth bool
	// Initialize a bit field
	for i := 0; i < FIELD_SIZE; i = i + 1 {
		cell = 31
		randInt = semiRandomData(randInt)
		// coin filt to decide if I should go south or east
		goSouth = (randInt % 1000) >= 500
		rowCol := bitIndexToRowCol(int(i))
		if goSouth {
			if rowCol[0] < (MAZE_SIZE - 1) {
				cell = cell & (^SOUTH_WALL)
			} else if rowCol[1] < (MAZE_SIZE - 1) {
				cell = cell & (^EAST_WALL)
			}
		} else {
			if rowCol[1] < (MAZE_SIZE - 1) {
				cell = cell & (^EAST_WALL)
			} else if rowCol[0] < (MAZE_SIZE)-1 {
				cell = cell & (^SOUTH_WALL)
			}
		}
		err := ioutil.WriteFile(fmt.Sprintf("out/%03d/%03d-%03d.txt", id, id, i), []byte(render(bitfield)), os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}
		bitfield[i] = cell
	}
	finalOut := render(bitfield)
	err := ioutil.WriteFile(fmt.Sprintf("out/%03d/%03d-final.txt", id, id), []byte(finalOut), os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = ioutil.WriteFile(fmt.Sprintf("out/%03d/%03d-final.txt.sha1", id, id), []byte(fmt.Sprintf("%x\n", sha1.Sum([]byte(finalOut)))), os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}
	return bitfield
}

// take an index that stored the flat array and return an array corresponding to the row and column in question
func bitIndexToRowCol(idx int) [2]int {
	row := idx / MAZE_SIZE
	col := idx % MAZE_SIZE
	return [2]int{row, col}
}

// Given a row and column, at the index in the flat array.
func rowColToBitIndex(row, col int) int {
	return row*MAZE_SIZE + col
}

func semiRandomData(seed int) int {
	return ((seed * 1103515245) + 12345) % 2147483648
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
