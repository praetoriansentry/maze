// contracts/Maze.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/presets/ERC721PresetMinterPauserAutoId.sol";

contract Maze is ERC721PresetMinterPauserAutoId {
    uint public constant MAZE_SIZE = 15;
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
        uint row;
        uint col;

        bytes memory mazeOutput = new bytes(PRINTED_MAZE_COLS * PRINTED_MAZE_ROWS);
        // For each row
        for (uint i = 0; i < PRINTED_MAZE_ROWS; i++) {
            uint j = 0;
            while (j < PRINTED_MAZE_COLS) {
                uint idx = i * PRINTED_MAZE_COLS + j;
                // END OF THE ROW
                if (j + 1 == PRINTED_MAZE_COLS) {
                    mazeOutput[idx] = CHAR_NEWLINE;
                    j = j + 1;
                    continue;
                }

                // START OF THE ROW
                if (j == 0 && i % 2 == 1) {
                    mazeOutput[idx] = CHAR_PIPE;
                    j = j + 1;
                    continue;
                } else if (j == 0 && i % 2 == 0) {
                    mazeOutput[idx] = CHAR_PLUS;
                    j = j + 1;
                    continue;
                }


                if (i == 0) {
                    cell = 31;
                } else {
                    row = (i - 1) / 2;
                    col = j / 4;
                    cellIndex = rowColToBitIndex(row, col);
                    if (cellIndex < FIELD_SIZE && cellIndex >= 0) {
                        cell = bitfield[cellIndex];
                    } else {
                        cell = 31;
                    }
                }

                // What cell are we in?
                if (i % 2 == 1) {// EAST WEST
                    mazeOutput[idx++] = CHAR_SPACE;
                    mazeOutput[idx++] = CHAR_SPACE;
                    mazeOutput[idx++] = CHAR_SPACE;
                    if ((cell & EAST_WALL) == EAST_WALL) {
                        mazeOutput[idx++] = CHAR_PIPE;
                    } else {
                        mazeOutput[idx++] = CHAR_SPACE;
                    }
                    j = j + 4;
                } else {// NORTH SOUTH
                    if ((cell & SOUTH_WALL) == SOUTH_WALL) {
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
            }
        }

        string memory finalMaze = string(mazeOutput);
        return finalMaze;
    }

    function getMazeData(uint id) public pure returns (uint[FIELD_SIZE] memory) {
        uint[FIELD_SIZE] memory bitfield;
        uint cell;
        bytes32 randState = randData(keccak256(abi.encode(id)));
        bool goSouth;
        // Initialize a bit field
        for (uint i = 0; i < FIELD_SIZE; i = i + 1) {
            cell = 31;
            randState = randData(randState);
            goSouth = (uint(randState) % 2 == 0);
            uint[2] memory rowCol = bitIndexToRowCol(i);
            if (goSouth) {
                if (rowCol[0] < (MAZE_SIZE - 1)) {
                    cell = cell & ~SOUTH_WALL;
                } else if (rowCol[1] < (MAZE_SIZE -1)) {
                    cell = cell & ~EAST_WALL;
                }
            } else {
                if (rowCol[1] < (MAZE_SIZE - 1)) {
                    cell = cell & ~EAST_WALL;
                } else {
                    cell = cell & ~SOUTH_WALL;

                }
            }
            bitfield[i] = cell;
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