// contracts/Maze.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/presets/ERC721PresetMinterPauserAutoId.sol";

contract Maze is ERC721PresetMinterPauserAutoId {
    uint public constant MAZE_SIZE = 16;
    uint constant FIELD_SIZE = MAZE_SIZE * MAZE_SIZE;

    uint constant PRINTED_MAZE_COLS = (MAZE_SIZE * 4 + 2);
    uint constant PRINTED_MAZE_ROWS = (2 * MAZE_SIZE + 1);

    uint constant NORTH_WALL = 1;
    uint constant SOUTH_WALL = 2;
    uint constant EAST_WALL = 4;
    uint constant WEST_WALL = 8;
    uint constant VISITED = 16;

    bytes1 constant CHAR_PLUS = 0x2B;
    bytes1 constant CHAR_DASH = 0x2D;
    bytes1 constant CHAR_NEWLINE = 0x0A;
    bytes1 constant CHAR_PIPE = 0x7C;
    bytes1 constant CHAR_SPACE = 0x20;




    constructor() ERC721PresetMinterPauserAutoId("Maze", "MZE", "https://maze.j4.is/token/") {
    }

    function draw(uint id) public pure returns (string memory) {
        uint[FIELD_SIZE] memory bitfield = getMazeData(id);
        uint cell;
        uint cellIndex;

        bytes memory mazeOutput = new bytes( PRINTED_MAZE_COLS * PRINTED_MAZE_ROWS);
        // For each row
        for (uint i = 0; i < PRINTED_MAZE_ROWS; i++) {
            uint j = 0;

            if (i % 2 == 0) {
                while (j < PRINTED_MAZE_COLS) {
                    uint idx = (i * PRINTED_MAZE_COLS) + j;
                    if (j == 0) {
                        mazeOutput[idx] = CHAR_PLUS;
                        j++;
                        continue;
                    }
                    if (j + 1 == PRINTED_MAZE_COLS) {
                        mazeOutput[idx] = CHAR_NEWLINE;
                        j++;
                        continue;
                    }
                    cellIndex = rowColToBitIndex((i / 2) - 1 , (j - 1) / 4);

                    if (cellIndex >= FIELD_SIZE && i < 1) {
                        cell = SOUTH_WALL;
                    } else {
                        cell = bitfield[cellIndex];
                    }
                

                    if (cell & SOUTH_WALL == SOUTH_WALL) {
                        mazeOutput[idx++] = CHAR_DASH;
                        mazeOutput[idx++] = CHAR_DASH;
                        mazeOutput[idx++] = CHAR_DASH;
                    } else {
                        mazeOutput[idx++] = CHAR_SPACE;
                        mazeOutput[idx++] = CHAR_SPACE;
                        mazeOutput[idx++] = CHAR_SPACE;
                    }

                    mazeOutput[idx++] = CHAR_PLUS;
                    j = j + 4;
                }
            } else {
                while (j < PRINTED_MAZE_COLS) {
                    uint idx = (i * PRINTED_MAZE_COLS) + j;
                    if (j == 0) {
                        mazeOutput[idx] = CHAR_PIPE;
                        j++;
                        continue;
                    }
                    if (j + 1 == PRINTED_MAZE_COLS) {
                        mazeOutput[idx] = CHAR_NEWLINE;
                        j++;
                        continue;
                    }
                    if (j % 4 == 0) {

                        mazeOutput[idx] = CHAR_PIPE;
                    } else {
                        mazeOutput[idx] = CHAR_SPACE;
                    }
                    j++;
                }
            }

        }

        string memory finalMaze = string(mazeOutput);

        return finalMaze;
    }





    function getMazeData(uint id) public pure returns (uint[FIELD_SIZE] memory) {
        uint[FIELD_SIZE] memory bitfield;

        uint totalVisited = 1;
        uint[4] memory possibleNeighborIndexes;
        uint possibleNeighborCount;
        uint selectedNeighborIndex;
        uint cell;
        uint[2] memory cellRowCol;
        uint selectedNeighbor;
        bytes32 randState = randData(keccak256(abi.encode(id)));
        uint cellIndex = uint(randState) % FIELD_SIZE;
    
        // Initialize a bit field
        for (uint i = 0; i < FIELD_SIZE; i = i + 1) {
            bitfield[i] = 15;
        }

        // https://en.wikipedia.org/wiki/Maze_generation_algorithm#Aldous-Broder_algorithm
        // 1. Pick a random cell as the current cell and mark it as visited.
        cell = bitfield[cellIndex];
        cell = cell | VISITED;
        bitfield[cellIndex] = cell;

        uint a = 0;

        // 2. While there are unvisited cells:
        while (totalVisited < FIELD_SIZE) {

            // 2.1 Pick a random neighbour.
            possibleNeighborIndexes = [cellIndex, cellIndex, cellIndex, cellIndex];
            possibleNeighborCount = 0;
            cellRowCol = bitIndexToRowCol(cellIndex);

            // NORTH
            if (cellRowCol[0] - 1 >= 0) {
                possibleNeighborIndexes[0] = rowColToBitIndex(cellRowCol[0] - 1, cellRowCol[1]);
            }
            // EAST
            if (cellRowCol[1] + 1 < MAZE_SIZE) {
                possibleNeighborIndexes[1] = rowColToBitIndex(cellRowCol[0], cellRowCol[1] + 1);
            }
            // SOUTH
            if (cellRowCol[0] + 1 < MAZE_SIZE) {
                possibleNeighborIndexes[2] = rowColToBitIndex(cellRowCol[0] + 1, cellRowCol[1]);
            }
            // WEST
            if (cellRowCol[1] - 1 >= 0) {
                possibleNeighborIndexes[3] = rowColToBitIndex(cellRowCol[0], cellRowCol[1] - 1);
            }

            // if the PRNG is really messed up this could run forever.. e.g. if we end up getting the same cell index.
            randState = randData(randState);
            selectedNeighborIndex = uint(randState) % 4;
            while(cellIndex == possibleNeighborIndexes[selectedNeighborIndex]) {
                selectedNeighborIndex++;
                selectedNeighborIndex = selectedNeighborIndex % 4;
            }
            selectedNeighbor = bitfield[possibleNeighborIndexes[selectedNeighborIndex]];

            // 2.2 If the chosen neighbour has not been visited: 
            if ((selectedNeighbor & VISITED) != VISITED) {
                // 2.2.1 Remove the wall between the current cell and the chosen neighbour.
                if (selectedNeighborIndex == 0) {
                    selectedNeighbor = selectedNeighbor & ~SOUTH_WALL;
                    cell = cell & ~NORTH_WALL;
                }
                if (selectedNeighborIndex == 1) {
                    selectedNeighbor = selectedNeighbor & ~WEST_WALL;
                    cell = cell & ~EAST_WALL;
                }
                if (selectedNeighborIndex == 2) {
                    selectedNeighbor = selectedNeighbor & ~NORTH_WALL;
                    cell = cell & ~SOUTH_WALL;
                }
                if (selectedNeighborIndex == 3) {
                    selectedNeighbor = selectedNeighbor & ~EAST_WALL;
                    cell = cell & ~WEST_WALL;
                }
                bitfield[cellIndex] = cell;

                // 2.2.2 Mark the chosen neighbour as visited.
                selectedNeighbor = selectedNeighbor | VISITED;
                totalVisited += 1;
                bitfield[possibleNeighborIndexes[selectedNeighborIndex]] = selectedNeighbor;
            }
            // 2.3 Make the chosen neighbour the current cell.
            cell = selectedNeighbor;
            cellIndex = possibleNeighborIndexes[selectedNeighborIndex];
            a = a + 1;
            if (a > 8) {
            break;
            }
        }
        return bitfield;

    }

    function bitIndexToRowCol(uint idx) private pure returns (uint[2] memory) {
        uint row = idx / MAZE_SIZE;
        uint col = idx % MAZE_SIZE;
        return [row, col];
    }

    function rowColToBitIndex(uint row, uint col) private pure returns(uint) {
        return row * MAZE_SIZE + col;
    }

    function randData(bytes32 randState) public pure returns (bytes32) {
        bytes32 idHash = keccak256(abi.encode(randState));
        return idHash;
    }
}