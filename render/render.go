package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"
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

	cubeSize = "1.0"
)

type (
	Prop struct {
		TraitType   string      `json:"trait_type"`
		DisplayType string      `json:"display_type"`
		Value       interface{} `json:"value"`
	}

	ERC721 struct {
		Name         string `json:"name"`
		Num          string `json:"num"`
		ID           string `json:"id"`
		Description  string `json:"description"`
		Image        string `json:"image"`
		AnimationUrl string `json:"animation_url"`
		ExternalUrl  string `json:"external_url"`
		MazeString   string `json:"maze_string"`
		MazeSHA      string `json:"maze_sha"`
		Attributes   []Prop `json:"attributes"`
	}
)

func main() {
	templateFile, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New("mze").Parse(string(templateFile))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 256; i = i + 1 {
		draw(i, tmpl)
	}
}
func render3D(bitfield [FIELD_SIZE]int) string {
	maze3dSize := 41 // 20 * 2 + 1
	output := ""
	for i := 0; i < maze3dSize; i += 1 {
		for j := 0; j < maze3dSize; j += 1 {
			// Always print the base layer. For all blocks, I want to have a floor level
			output += fmt.Sprintf("translate([%d, %d, 0]) {color([0,0,0]) cube(%s);};\n", i, j, cubeSize)

			if i == 0 || i+1 == maze3dSize || j == 0 || j+1 == maze3dSize {
				output += renderBlock(i, j)
				continue
			}

			if i%2 == 0 || j%2 == 0 {
				row := (i - 1) / 2
				col := (j - 1) / 2
				var cell int
				cellIndex := rowColToBitIndex(row, col)
				if cellIndex < FIELD_SIZE && cellIndex >= 0 {
					cell = bitfield[cellIndex]
				} else {
					cell = 31
				}

				// south
				if i%2 == 0 && (cell&SOUTH_WALL) == SOUTH_WALL {
					output += renderBlock(i, j)
					continue
				}

				// east
				if j%2 == 0 && (cell&EAST_WALL) == EAST_WALL {
					output += renderBlock(i, j)
					continue
				}
			}
		}
	}
	return output
}
func renderBlock(i, j int) string {
	output := ""
	output += fmt.Sprintf("translate([%d, %d, 1]) {cube(%s);};\n", i, j, cubeSize)
	output += fmt.Sprintf("translate([%d, %d, 2]) {cube(%s);};\n", i, j, cubeSize)
	output += fmt.Sprintf("translate([%d, %d, 3]) {cube(%s);};\n", i, j, cubeSize)
	output += fmt.Sprintf("translate([%d, %d, 4]) {cube(%s);};\n", i, j, cubeSize)
	output += fmt.Sprintf("translate([%d, %d, 5]) {cube(%s);};\n", i, j, cubeSize)
	return output
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
func draw(id int, tmpl *template.Template) string {
	os.Mkdir(fmt.Sprintf("out/%d/", id), os.ModePerm)
	copy(fmt.Sprintf("masks/%d.jpg", id%30), fmt.Sprintf("out/%d/mask.jpg", id))
	var bitfield [FIELD_SIZE]int = getMazeData(id)
	mazeString := render(bitfield)
	mazeString3d := render3D(bitfield)
	mazeSha := fmt.Sprintf("%x\n", sha1.Sum([]byte(mazeString)))

	err := ioutil.WriteFile(fmt.Sprintf("out/%d/%d-final.txt.sha1", id, id), []byte(mazeSha), os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = ioutil.WriteFile(fmt.Sprintf("out/%d/maze.scad", id), []byte(mazeString3d), os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}

	south := 0
	east := 0

	for _, v := range bitfield {
		if (v & EAST_WALL) != EAST_WALL {
			east += 1
		}
		if (v & SOUTH_WALL) != SOUTH_WALL {
			south += 1
		}
	}

	mzData := new(ERC721)
	attrs := make([]Prop, 0)

	mzData.Num = fmt.Sprintf("%03d", id)
	mzData.ID = fmt.Sprintf("%d", id)
	mzData.Name = fmt.Sprintf("BlockMazing #%03d", id)
	mzData.Description = fmt.Sprintf("BlockMazing #%03d", id)
	mzData.Image = fmt.Sprintf("https://blockmazing.com/m/%d/maze.png", id)
	mzData.AnimationUrl = fmt.Sprintf("https://blockmazing.com/m/%d/maze.mp4", id)
	mzData.ExternalUrl = fmt.Sprintf("https://blockmazing.com/m/%d/", id)
	mzData.MazeString = mazeString
	mzData.MazeSHA = mazeSha
	attrs = append(attrs, Prop{TraitType: "South turns", Value: south, DisplayType: "number"})
	attrs = append(attrs, Prop{TraitType: "East turns", Value: east, DisplayType: "number"})
	mzData.Attributes = attrs
	mzJson, err := json.Marshal(mzData)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("out/%d/%d.json", id, id), mzJson, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mzData)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("out/%d/index.html", id), buf.Bytes(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return mazeString
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
		err := ioutil.WriteFile(fmt.Sprintf("out/%d/%d-%03d.txt", id, id, i), []byte(render(bitfield)), os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}
		bitfield[i] = cell
	}
	finalOut := render(bitfield)
	err := ioutil.WriteFile(fmt.Sprintf("out/%d/%d-final.txt", id, id), []byte(finalOut), os.ModePerm)
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
