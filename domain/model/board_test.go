package model

import (
	"errors"
	"testing"
)

func TestNewBoard(t *testing.T) {
	tests := []struct {
		name        string
		width       int
		height      int
		expectError bool
		errorType   error
	}{
		{
			name:        "正常なボード作成",
			width:       10,
			height:      20,
			expectError: false,
		},
		{
			name:        "小さなボード作成",
			width:       5,
			height:      5,
			expectError: false,
		},
		{
			name:        "幅が0のボード",
			width:       0,
			height:      10,
			expectError: true,
			errorType:   ErrInvalidBoardSize,
		},
		{
			name:        "高さが0のボード",
			width:       10,
			height:      0,
			expectError: true,
			errorType:   ErrInvalidBoardSize,
		},
		{
			name:        "負の幅のボード",
			width:       -5,
			height:      10,
			expectError: true,
			errorType:   ErrInvalidBoardSize,
		},
		{
			name:        "負の高さのボード",
			width:       10,
			height:      -5,
			expectError: true,
			errorType:   ErrInvalidBoardSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := NewBoard(tt.width, tt.height)

			if tt.expectError {
				if err == nil {
					t.Errorf("NewBoard() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("NewBoard() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewBoard() unexpected error = %v", err)
				return
			}

			if board == nil {
				t.Error("NewBoard() returned nil board")
				return
			}

			if board.Width != tt.width {
				t.Errorf("NewBoard() Width = %v, want %v", board.Width, tt.width)
			}

			if board.Height != tt.height {
				t.Errorf("NewBoard() Height = %v, want %v", board.Height, tt.height)
			}

			if len(board.Grid) != tt.height {
				t.Errorf("NewBoard() Grid height = %v, want %v", len(board.Grid), tt.height)
			}

			for i, row := range board.Grid {
				if len(row) != tt.width {
					t.Errorf("NewBoard() Grid[%d] width = %v, want %v", i, len(row), tt.width)
				}
			}
		})
	}
}

func TestBoard_IsValidPosition(t *testing.T) {
	board, err := NewBoard(10, 20)
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	tests := []struct {
		name     string
		point    Point
		expected bool
	}{
		{
			name:     "左上角",
			point:    Point{X: 0, Y: 0},
			expected: true,
		},
		{
			name:     "右下角",
			point:    Point{X: 9, Y: 19},
			expected: true,
		},
		{
			name:     "中央",
			point:    Point{X: 5, Y: 10},
			expected: true,
		},
		{
			name:     "X座標が負",
			point:    Point{X: -1, Y: 5},
			expected: false,
		},
		{
			name:     "Y座標が負",
			point:    Point{X: 5, Y: -1},
			expected: false,
		},
		{
			name:     "X座標が範囲外",
			point:    Point{X: 10, Y: 5},
			expected: false,
		},
		{
			name:     "Y座標が範囲外",
			point:    Point{X: 5, Y: 20},
			expected: false,
		},
		{
			name:     "両座標が範囲外",
			point:    Point{X: 15, Y: 25},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := board.IsValidPosition(tt.point)
			if result != tt.expected {
				t.Errorf("Board.IsValidPosition() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBoard_SetBlock_and_IsOccupied(t *testing.T) {
	board, err := NewBoard(10, 20)
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	tests := []struct {
		name        string
		point       Point
		occupied    bool
		expectError bool
		errorType   error
	}{
		{
			name:        "有効な位置にブロック設置",
			point:       Point{X: 5, Y: 10},
			occupied:    true,
			expectError: false,
		},
		{
			name:        "有効な位置のブロック削除",
			point:       Point{X: 3, Y: 7},
			occupied:    false,
			expectError: false,
		},
		{
			name:        "無効な位置（範囲外）",
			point:       Point{X: 15, Y: 10},
			occupied:    true,
			expectError: true,
			errorType:   ErrOutOfBounds,
		},
		{
			name:        "無効な位置（負の座標）",
			point:       Point{X: -1, Y: 5},
			occupied:    true,
			expectError: true,
			errorType:   ErrOutOfBounds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := board.SetBlock(tt.point, tt.occupied)

			if tt.expectError {
				if err == nil {
					t.Errorf("Board.SetBlock() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("Board.SetBlock() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("Board.SetBlock() unexpected error = %v", err)
				return
			}

			occupied, err := board.IsOccupied(tt.point)
			if err != nil {
				t.Errorf("Board.IsOccupied() error = %v", err)
				return
			}

			if occupied != tt.occupied {
				t.Errorf("Board.IsOccupied() = %v, want %v", occupied, tt.occupied)
			}
		})
	}
}

func TestBoard_CanPlaceTetromino(t *testing.T) {
	board, err := NewBoard(10, 20)
	if err != nil {
		t.Fatalf("NewBoard() error = %v", err)
	}

	board.SetBlock(Point{X: 5, Y: 10}, true)

	tests := []struct {
		name      string
		tetromino *Tetromino
		expected  bool
	}{
		{
			name:      "nilテトロミノ",
			tetromino: nil,
			expected:  false,
		},
		{
			name: "有効な位置のテトロミノ",
			tetromino: func() *Tetromino {
				t, _ := NewTetromino(I, Point{X: 0, Y: 0})
				return t
			}(),
			expected: true,
		},
		{
			name: "占有済みブロックと重複するテトロミノ",
			tetromino: func() *Tetromino {
				t, _ := NewTetromino(I, Point{X: 4, Y: 9})
				return t
			}(),
			expected: false,
		},
		{
			name: "ボード範囲外のテトロミノ",
			tetromino: func() *Tetromino {
				t, _ := NewTetromino(I, Point{X: -2, Y: 0})
				return t
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := board.CanPlaceTetromino(tt.tetromino)
			if result != tt.expected {
				t.Errorf("Board.CanPlaceTetromino() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBoard_GetCompletedLines(t *testing.T) {
	tests := []struct {
		name          string
		setupBoard    func(*Board)
		expectedLines []int
	}{
		{
			name: "完成ラインなし",
			setupBoard: func(b *Board) {
				for x := 0; x < 5; x++ {
					b.SetBlock(Point{X: x, Y: 19}, true)
				}
			},
			expectedLines: []int{},
		},
		{
			name: "1つの完成ライン",
			setupBoard: func(b *Board) {
				for x := 0; x < b.Width; x++ {
					b.SetBlock(Point{X: x, Y: 19}, true)
				}
			},
			expectedLines: []int{19},
		},
		{
			name: "複数の完成ライン",
			setupBoard: func(b *Board) {
				for x := 0; x < b.Width; x++ {
					b.SetBlock(Point{X: x, Y: 18}, true)
					b.SetBlock(Point{X: x, Y: 19}, true)
				}
			},
			expectedLines: []int{18, 19},
		},
		{
			name: "離れた完成ライン",
			setupBoard: func(b *Board) {
				for x := 0; x < b.Width; x++ {
					b.SetBlock(Point{X: x, Y: 15}, true)
					b.SetBlock(Point{X: x, Y: 19}, true)
				}
			},
			expectedLines: []int{15, 19},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := NewBoard(10, 20)
			if err != nil {
				t.Fatalf("NewBoard() error = %v", err)
			}

			tt.setupBoard(board)

			result := board.GetCompletedLines()

			if len(result) != len(tt.expectedLines) {
				t.Errorf("Board.GetCompletedLines() length = %v, want %v", len(result), len(tt.expectedLines))
				return
			}

			for i, line := range result {
				if line != tt.expectedLines[i] {
					t.Errorf("Board.GetCompletedLines()[%d] = %v, want %v", i, line, tt.expectedLines[i])
				}
			}
		})
	}
}

func TestBoard_ClearLines(t *testing.T) {
	tests := []struct {
		name         string
		setupBoard   func(*Board)
		linesToClear []int
		expectError  bool
		errorType    error
		checkResult  func(*Board) bool
	}{
		{
			name: "1つのライン消去",
			setupBoard: func(b *Board) {
				for x := 0; x < b.Width; x++ {
					b.SetBlock(Point{X: x, Y: 19}, true)
				}
				b.SetBlock(Point{X: 0, Y: 18}, true)
			},
			linesToClear: []int{19},
			expectError:  false,
			checkResult: func(b *Board) bool {
				occupied, _ := b.IsOccupied(Point{X: 0, Y: 19})
				return occupied
			},
		},
		{
			name:         "範囲外のライン消去",
			setupBoard:   func(b *Board) {},
			linesToClear: []int{25},
			expectError:  true,
			errorType:    ErrOutOfBounds,
		},
		{
			name:         "空のライン配列",
			setupBoard:   func(b *Board) {},
			linesToClear: []int{},
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := NewBoard(10, 20)
			if err != nil {
				t.Fatalf("NewBoard() error = %v", err)
			}

			tt.setupBoard(board)

			err = board.ClearLines(tt.linesToClear)

			if tt.expectError {
				if err == nil {
					t.Errorf("Board.ClearLines() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("Board.ClearLines() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("Board.ClearLines() unexpected error = %v", err)
				return
			}

			if tt.checkResult != nil && !tt.checkResult(board) {
				t.Error("Board.ClearLines() result check failed")
			}
		})
	}
}

func TestBoard_IsGameOver(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func(*Board)
		expected   bool
	}{
		{
			name:       "空のボード",
			setupBoard: func(b *Board) {},
			expected:   false,
		},
		{
			name: "最上行にブロックあり",
			setupBoard: func(b *Board) {
				b.SetBlock(Point{X: 5, Y: 0}, true)
			},
			expected: true,
		},
		{
			name: "最上行以外にブロックあり",
			setupBoard: func(b *Board) {
				b.SetBlock(Point{X: 5, Y: 1}, true)
				b.SetBlock(Point{X: 3, Y: 10}, true)
			},
			expected: false,
		},
		{
			name: "最上行の複数位置にブロックあり",
			setupBoard: func(b *Board) {
				b.SetBlock(Point{X: 0, Y: 0}, true)
				b.SetBlock(Point{X: 9, Y: 0}, true)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := NewBoard(10, 20)
			if err != nil {
				t.Fatalf("NewBoard() error = %v", err)
			}

			tt.setupBoard(board)

			result := board.IsGameOver()
			if result != tt.expected {
				t.Errorf("Board.IsGameOver() = %v, want %v", result, tt.expected)
			}
		})
	}
}
