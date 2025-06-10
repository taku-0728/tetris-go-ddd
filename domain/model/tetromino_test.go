package model

import (
	"errors"
	"testing"
)

func TestNewTetromino(t *testing.T) {
	tests := []struct {
		name         string
		tetrominoType TetrominoType
		position     Point
		expectError  bool
		errorType    error
	}{
		{
			name:         "正常なIピース作成",
			tetrominoType: I,
			position:     Point{X: 5, Y: 0},
			expectError:  false,
		},
		{
			name:         "正常なOピース作成",
			tetrominoType: O,
			position:     Point{X: 3, Y: 2},
			expectError:  false,
		},
		{
			name:         "正常なTピース作成",
			tetrominoType: T,
			position:     Point{X: 0, Y: 0},
			expectError:  false,
		},
		{
			name:         "無効なテトロミノタイプ（負の値）",
			tetrominoType: TetrominoType(-1),
			position:     Point{X: 0, Y: 0},
			expectError:  true,
			errorType:    ErrInvalidTetrominoType,
		},
		{
			name:         "無効なテトロミノタイプ（範囲外の値）",
			tetrominoType: TetrominoType(10),
			position:     Point{X: 0, Y: 0},
			expectError:  true,
			errorType:    ErrInvalidTetrominoType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tetromino, err := NewTetromino(tt.tetrominoType, tt.position)

			if tt.expectError {
				if err == nil {
					t.Errorf("NewTetromino() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("NewTetromino() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewTetromino() unexpected error = %v", err)
				return
			}

			if tetromino == nil {
				t.Error("NewTetromino() returned nil tetromino")
				return
			}

			if tetromino.Type != tt.tetrominoType {
				t.Errorf("NewTetromino() Type = %v, want %v", tetromino.Type, tt.tetrominoType)
			}

			if tetromino.Position != tt.position {
				t.Errorf("NewTetromino() Position = %v, want %v", tetromino.Position, tt.position)
			}
		})
	}
}

func TestTetromino_Move(t *testing.T) {
	tests := []struct {
		name             string
		initialPosition  Point
		delta           Point
		expectedPosition Point
	}{
		{
			name:             "右に移動",
			initialPosition:  Point{X: 5, Y: 10},
			delta:           Point{X: 1, Y: 0},
			expectedPosition: Point{X: 6, Y: 10},
		},
		{
			name:             "左に移動",
			initialPosition:  Point{X: 5, Y: 10},
			delta:           Point{X: -1, Y: 0},
			expectedPosition: Point{X: 4, Y: 10},
		},
		{
			name:             "下に移動",
			initialPosition:  Point{X: 5, Y: 10},
			delta:           Point{X: 0, Y: 1},
			expectedPosition: Point{X: 5, Y: 11},
		},
		{
			name:             "斜め移動",
			initialPosition:  Point{X: 3, Y: 7},
			delta:           Point{X: 2, Y: -3},
			expectedPosition: Point{X: 5, Y: 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tetromino, err := NewTetromino(I, tt.initialPosition)
			if err != nil {
				t.Fatalf("NewTetromino() error = %v", err)
			}

			err = tetromino.Move(tt.delta)
			if err != nil {
				t.Errorf("Tetromino.Move() error = %v", err)
			}

			if tetromino.Position != tt.expectedPosition {
				t.Errorf("Tetromino.Move() Position = %v, want %v", tetromino.Position, tt.expectedPosition)
			}
		})
	}
}

func TestTetromino_GetBlocks(t *testing.T) {
	tests := []struct {
		name         string
		tetrominoType TetrominoType
		position     Point
		minBlocks    int
		maxBlocks    int
	}{
		{
			name:         "Iピースのブロック数",
			tetrominoType: I,
			position:     Point{X: 0, Y: 0},
			minBlocks:    4,
			maxBlocks:    4,
		},
		{
			name:         "Oピースのブロック数",
			tetrominoType: O,
			position:     Point{X: 0, Y: 0},
			minBlocks:    4,
			maxBlocks:    4,
		},
		{
			name:         "Tピースのブロック数",
			tetrominoType: T,
			position:     Point{X: 0, Y: 0},
			minBlocks:    4,
			maxBlocks:    4,
		},
		{
			name:         "位置オフセットありのブロック",
			tetrominoType: I,
			position:     Point{X: 5, Y: 10},
			minBlocks:    4,
			maxBlocks:    4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tetromino, err := NewTetromino(tt.tetrominoType, tt.position)
			if err != nil {
				t.Fatalf("NewTetromino() error = %v", err)
			}

			blocks := tetromino.GetBlocks()
			
			if len(blocks) < tt.minBlocks || len(blocks) > tt.maxBlocks {
				t.Errorf("Tetromino.GetBlocks() blocks count = %d, want between %d and %d", len(blocks), tt.minBlocks, tt.maxBlocks)
			}

			for _, block := range blocks {
				if block.X < tt.position.X || block.Y < tt.position.Y {
					t.Errorf("Tetromino.GetBlocks() block %v is outside expected bounds from position %v", block, tt.position)
				}
			}
		})
	}
}

func TestTetromino_Rotate(t *testing.T) {
	tests := []struct {
		name         string
		tetrominoType TetrominoType
		position     Point
		expectError  bool
		errorType    error
	}{
		{
			name:         "Iピースの回転",
			tetrominoType: I,
			position:     Point{X: 5, Y: 5},
			expectError:  false,
		},
		{
			name:         "Tピースの回転",
			tetrominoType: T,
			position:     Point{X: 3, Y: 3},
			expectError:  false,
		},
		{
			name:         "Oピースの回転（変化なし）",
			tetrominoType: O,
			position:     Point{X: 2, Y: 2},
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tetromino, err := NewTetromino(tt.tetrominoType, tt.position)
			if err != nil {
				t.Fatalf("NewTetromino() error = %v", err)
			}

			originalBlocks := tetromino.GetBlocks()
			
			err = tetromino.Rotate()
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Tetromino.Rotate() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("Tetromino.Rotate() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("Tetromino.Rotate() unexpected error = %v", err)
				return
			}

			newBlocks := tetromino.GetBlocks()
			
			if len(newBlocks) != len(originalBlocks) {
				t.Errorf("Tetromino.Rotate() block count changed from %d to %d", len(originalBlocks), len(newBlocks))
			}

			if tetromino.Position != tt.position {
				t.Errorf("Tetromino.Rotate() Position changed from %v to %v", tt.position, tetromino.Position)
			}
		})
	}
}