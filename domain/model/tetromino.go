package model

import (
	"errors"
	"fmt"
)

type TetrominoType int

const (
	I TetrominoType = iota
	O
	T
	S
	Z
	J
	L
)

var (
	ErrInvalidTetrominoType = errors.New("無効なテトロミノタイプです")
	ErrInvalidPosition      = errors.New("無効な位置です")
	ErrRotationFailed       = errors.New("回転に失敗しました")
)

type Tetromino struct {
	Type     TetrominoType
	Shape    [][]bool
	Position Point
	size     int
}

var tetrominoShapes = map[TetrominoType][][][]bool{
	I: {
		{
			{false, false, false, false},
			{true, true, true, true},
			{false, false, false, false},
			{false, false, false, false},
		},
		{
			{false, false, true, false},
			{false, false, true, false},
			{false, false, true, false},
			{false, false, true, false},
		},
	},
	O: {
		{
			{false, false, false, false},
			{false, true, true, false},
			{false, true, true, false},
			{false, false, false, false},
		},
	},
	T: {
		{
			{false, false, false, false},
			{false, true, false, false},
			{true, true, true, false},
			{false, false, false, false},
		},
		{
			{false, false, false, false},
			{false, true, false, false},
			{false, true, true, false},
			{false, true, false, false},
		},
		{
			{false, false, false, false},
			{false, false, false, false},
			{true, true, true, false},
			{false, true, false, false},
		},
		{
			{false, false, false, false},
			{false, true, false, false},
			{true, true, false, false},
			{false, true, false, false},
		},
	},
	S: {
		{
			{false, false, false, false},
			{false, true, true, false},
			{true, true, false, false},
			{false, false, false, false},
		},
		{
			{false, false, false, false},
			{false, true, false, false},
			{false, true, true, false},
			{false, false, true, false},
		},
	},
	Z: {
		{
			{false, false, false, false},
			{true, true, false, false},
			{false, true, true, false},
			{false, false, false, false},
		},
		{
			{false, false, false, false},
			{false, false, true, false},
			{false, true, true, false},
			{false, true, false, false},
		},
	},
	J: {
		{
			{false, false, false, false},
			{true, false, false, false},
			{true, true, true, false},
			{false, false, false, false},
		},
		{
			{false, false, false, false},
			{false, true, true, false},
			{false, true, false, false},
			{false, true, false, false},
		},
		{
			{false, false, false, false},
			{false, false, false, false},
			{true, true, true, false},
			{false, false, true, false},
		},
		{
			{false, false, false, false},
			{false, true, false, false},
			{false, true, false, false},
			{true, true, false, false},
		},
	},
	L: {
		{
			{false, false, false, false},
			{false, false, true, false},
			{true, true, true, false},
			{false, false, false, false},
		},
		{
			{false, false, false, false},
			{false, true, false, false},
			{false, true, false, false},
			{false, true, true, false},
		},
		{
			{false, false, false, false},
			{false, false, false, false},
			{true, true, true, false},
			{true, false, false, false},
		},
		{
			{false, false, false, false},
			{true, true, false, false},
			{false, true, false, false},
			{false, true, false, false},
		},
	},
}

func NewTetromino(tetrominoType TetrominoType, position Point) (*Tetromino, error) {
	if tetrominoType < I || tetrominoType > L {
		return nil, fmt.Errorf("%w: %d", ErrInvalidTetrominoType, tetrominoType)
	}

	shapes, exists := tetrominoShapes[tetrominoType]
	if !exists || len(shapes) == 0 {
		return nil, fmt.Errorf("%w: テトロミノ形状が見つかりません", ErrInvalidTetrominoType)
	}

	shape := make([][]bool, len(shapes[0]))
	for i := range shape {
		shape[i] = make([]bool, len(shapes[0][i]))
		copy(shape[i], shapes[0][i])
	}

	return &Tetromino{
		Type:     tetrominoType,
		Shape:    shape,
		Position: position,
		size:     4,
	}, nil
}

func (t *Tetromino) GetBlocks() []Point {
	var blocks []Point
	for y := 0; y < t.size; y++ {
		for x := 0; x < t.size; x++ {
			if y < len(t.Shape) && x < len(t.Shape[y]) && t.Shape[y][x] {
				blocks = append(blocks, Point{
					X: t.Position.X + x,
					Y: t.Position.Y + y,
				})
			}
		}
	}
	return blocks
}

func (t *Tetromino) Move(delta Point) error {
	newPosition := t.Position.Add(delta)
	t.Position = newPosition
	return nil
}

func (t *Tetromino) Rotate() error {
	shapes, exists := tetrominoShapes[t.Type]
	if !exists {
		return fmt.Errorf("%w: テトロミノタイプが無効です", ErrRotationFailed)
	}

	if len(shapes) <= 1 {
		return nil
	}

	for i, shape := range shapes {
		if t.shapeEquals(shape) {
			nextIndex := (i + 1) % len(shapes)
			newShape := make([][]bool, len(shapes[nextIndex]))
			for j := range newShape {
				newShape[j] = make([]bool, len(shapes[nextIndex][j]))
				copy(newShape[j], shapes[nextIndex][j])
			}
			t.Shape = newShape
			return nil
		}
	}

	return fmt.Errorf("%w: 現在の形状が見つかりません", ErrRotationFailed)
}

func (t *Tetromino) shapeEquals(shape [][]bool) bool {
	if len(t.Shape) != len(shape) {
		return false
	}
	for i := range t.Shape {
		if len(t.Shape[i]) != len(shape[i]) {
			return false
		}
		for j := range t.Shape[i] {
			if t.Shape[i][j] != shape[i][j] {
				return false
			}
		}
	}
	return true
}
