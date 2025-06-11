package model

import (
	"errors"
	"fmt"
)

const (
	BoardWidth  = 10
	BoardHeight = 20
)

var (
	ErrOutOfBounds      = errors.New("ボード範囲外です")
	ErrInvalidBoardSize = errors.New("無効なボードサイズです")
	ErrBlockOccupied    = errors.New("ブロックが既に配置されています")
)

type Board struct {
	Grid   [][]bool
	Width  int
	Height int
}

func NewBoard(width, height int) (*Board, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("%w: 幅=%d, 高さ=%d", ErrInvalidBoardSize, width, height)
	}

	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}

	return &Board{
		Grid:   grid,
		Width:  width,
		Height: height,
	}, nil
}

func (b *Board) IsValidPosition(point Point) bool {
	return point.X >= 0 && point.X < b.Width && point.Y >= 0 && point.Y < b.Height
}

func (b *Board) IsOccupied(point Point) (bool, error) {
	if !b.IsValidPosition(point) {
		return false, fmt.Errorf("%w: 座標(%d, %d)", ErrOutOfBounds, point.X, point.Y)
	}
	return b.Grid[point.Y][point.X], nil
}

func (b *Board) SetBlock(point Point, occupied bool) error {
	if !b.IsValidPosition(point) {
		return fmt.Errorf("%w: 座標(%d, %d)", ErrOutOfBounds, point.X, point.Y)
	}
	b.Grid[point.Y][point.X] = occupied
	return nil
}

func (b *Board) CanPlaceTetromino(tetromino *Tetromino) bool {
	if tetromino == nil {
		return false
	}

	blocks := tetromino.GetBlocks()
	for _, block := range blocks {
		if !b.IsValidPosition(block) {
			return false
		}
		if occupied, err := b.IsOccupied(block); err != nil || occupied {
			return false
		}
	}
	return true
}

func (b *Board) PlaceTetromino(tetromino *Tetromino) error {
	if tetromino == nil {
		return errors.New("テトロミノがnilです")
	}

	if !b.CanPlaceTetromino(tetromino) {
		return fmt.Errorf("%w: テトロミノを配置できません", ErrBlockOccupied)
	}

	blocks := tetromino.GetBlocks()
	for _, block := range blocks {
		if err := b.SetBlock(block, true); err != nil {
			return fmt.Errorf("ブロック配置エラー: %w", err)
		}
	}
	return nil
}

func (b *Board) GetCompletedLines() []int {
	var completedLines []int
	for y := 0; y < b.Height; y++ {
		if b.isLineFull(y) {
			completedLines = append(completedLines, y)
		}
	}
	return completedLines
}

func (b *Board) isLineFull(y int) bool {
	if y < 0 || y >= b.Height {
		return false
	}
	for x := 0; x < b.Width; x++ {
		if !b.Grid[y][x] {
			return false
		}
	}
	return true
}

func (b *Board) ClearLines(lines []int) error {
	if len(lines) == 0 {
		return nil
	}

	for _, line := range lines {
		if line < 0 || line >= b.Height {
			return fmt.Errorf("%w: ライン番号=%d", ErrOutOfBounds, line)
		}
	}

	for i := len(lines) - 1; i >= 0; i-- {
		lineIndex := lines[i]
		for y := lineIndex; y > 0; y-- {
			copy(b.Grid[y], b.Grid[y-1])
		}
		for x := 0; x < b.Width; x++ {
			b.Grid[0][x] = false
		}
	}
	return nil
}

func (b *Board) IsGameOver() bool {
	for x := 0; x < b.Width; x++ {
		if b.Grid[0][x] {
			return true
		}
	}
	return false
}
