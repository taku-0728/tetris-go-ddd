package service

import (
	"errors"
	"testing"
	"tetris/domain/model"
)

func TestNewGameService(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "正常なゲームサービス作成",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameService, err := NewGameService()

			if tt.expectError {
				if err == nil {
					t.Error("NewGameService() error = nil, wantErr")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("NewGameService() unexpected error = %v", err)
				return
			}

			if gameService == nil {
				t.Error("NewGameService() returned nil")
				return
			}

			if gameService.GetBoard() == nil {
				t.Error("NewGameService() board is nil")
			}

			if gameService.GetCurrentPiece() == nil {
				t.Error("NewGameService() current piece is nil")
			}

			if gameService.GetNextPiece() == nil {
				t.Error("NewGameService() next piece is nil")
			}

			if gameService.GetScore() != 0 {
				t.Errorf("NewGameService() initial score = %d, want 0", gameService.GetScore())
			}

			if gameService.GetLevel() != 1 {
				t.Errorf("NewGameService() initial level = %d, want 1", gameService.GetLevel())
			}

			if gameService.IsGameOver() {
				t.Error("NewGameService() game over = true, want false")
			}
		})
	}
}

func TestGameService_MovePiece(t *testing.T) {
	tests := []struct {
		name        string
		delta       model.Point
		expectError bool
		errorType   error
	}{
		{
			name:        "右に移動",
			delta:       model.Point{X: 1, Y: 0},
			expectError: false,
		},
		{
			name:        "左に移動",
			delta:       model.Point{X: -1, Y: 0},
			expectError: false,
		},
		{
			name:        "下に移動",
			delta:       model.Point{X: 0, Y: 1},
			expectError: false,
		},
		{
			name:        "大幅に左に移動（無効）",
			delta:       model.Point{X: -10, Y: 0},
			expectError: true,
			errorType:   ErrInvalidMove,
		},
		{
			name:        "大幅に右に移動（無効）",
			delta:       model.Point{X: 10, Y: 0},
			expectError: true,
			errorType:   ErrInvalidMove,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameService, err := NewGameService()
			if err != nil {
				t.Fatalf("NewGameService() error = %v", err)
			}

			originalPosition := gameService.GetCurrentPiece().Position

			err = gameService.MovePiece(tt.delta)

			if tt.expectError {
				if err == nil {
					t.Errorf("GameService.MovePiece() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("GameService.MovePiece() error = %v, wantErr %v", err, tt.errorType)
				}

				if gameService.GetCurrentPiece().Position != originalPosition {
					t.Error("GameService.MovePiece() position changed after invalid move")
				}
				return
			}

			if err != nil {
				t.Errorf("GameService.MovePiece() unexpected error = %v", err)
				return
			}

			expectedPosition := originalPosition.Add(tt.delta)
			if gameService.GetCurrentPiece().Position != expectedPosition {
				t.Errorf("GameService.MovePiece() position = %v, want %v", gameService.GetCurrentPiece().Position, expectedPosition)
			}
		})
	}
}

func TestGameService_RotatePiece(t *testing.T) {
	tests := []struct {
		name         string
		setupGame    func(*GameService)
		expectError  bool
		errorType    error
	}{
		{
			name:        "正常な回転",
			setupGame:   func(g *GameService) {},
			expectError: false,
		},
		{
			name: "ゲームオーバー時の回転",
			setupGame: func(g *GameService) {
				g.gameOver = true
			},
			expectError: true,
			errorType:   ErrGameOver,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameService, err := NewGameService()
			if err != nil {
				t.Fatalf("NewGameService() error = %v", err)
			}

			tt.setupGame(gameService)

			originalShape := make([][]bool, len(gameService.GetCurrentPiece().Shape))
			for i := range originalShape {
				originalShape[i] = make([]bool, len(gameService.GetCurrentPiece().Shape[i]))
				copy(originalShape[i], gameService.GetCurrentPiece().Shape[i])
			}

			err = gameService.RotatePiece()

			if tt.expectError {
				if err == nil {
					t.Errorf("GameService.RotatePiece() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("GameService.RotatePiece() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("GameService.RotatePiece() unexpected error = %v", err)
				return
			}
		})
	}
}

func TestGameService_DropPiece(t *testing.T) {
	tests := []struct {
		name        string
		setupGame   func(*GameService)
		expectError bool
		errorType   error
	}{
		{
			name:        "正常なドロップ",
			setupGame:   func(g *GameService) {},
			expectError: false,
		},
		{
			name: "ゲームオーバー時のドロップ",
			setupGame: func(g *GameService) {
				g.gameOver = true
			},
			expectError: true,
			errorType:   ErrGameOver,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameService, err := NewGameService()
			if err != nil {
				t.Fatalf("NewGameService() error = %v", err)
			}

			tt.setupGame(gameService)

			err = gameService.DropPiece()

			if tt.expectError {
				if err == nil {
					t.Errorf("GameService.DropPiece() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("GameService.DropPiece() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("GameService.DropPiece() unexpected error = %v", err)
			}

			// ドロップ後はピースがロックされて新しいピースが生成されるため、
			// 単純な位置の比較ではなく、ドロップが正常に実行されたかをチェック
		})
	}
}

func TestGameService_Update(t *testing.T) {
	tests := []struct {
		name        string
		setupGame   func(*GameService)
		expectError bool
		errorType   error
	}{
		{
			name:        "正常な更新",
			setupGame:   func(g *GameService) {},
			expectError: false,
		},
		{
			name: "ゲームオーバー時の更新",
			setupGame: func(g *GameService) {
				g.gameOver = true
			},
			expectError: true,
			errorType:   ErrGameOver,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameService, err := NewGameService()
			if err != nil {
				t.Fatalf("NewGameService() error = %v", err)
			}

			tt.setupGame(gameService)

			err = gameService.Update()

			if tt.expectError {
				if err == nil {
					t.Errorf("GameService.Update() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("GameService.Update() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("GameService.Update() unexpected error = %v", err)
			}
		})
	}
}

func TestGameService_InitialState(t *testing.T) {
	gameService, err := NewGameService()
	if err != nil {
		t.Fatalf("NewGameService() error = %v", err)
	}

	tests := []struct {
		name     string
		check    func() bool
		message  string
	}{
		{
			name:    "初期スコアが0",
			check:   func() bool { return gameService.GetScore() == 0 },
			message: "Initial score should be 0",
		},
		{
			name:    "初期ラインが0",
			check:   func() bool { return gameService.GetLines() == 0 },
			message: "Initial lines should be 0",
		},
		{
			name:    "初期レベルが1",
			check:   func() bool { return gameService.GetLevel() == 1 },
			message: "Initial level should be 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.message)
			}
		})
	}
}

func TestGameService_Getters(t *testing.T) {
	gameService, err := NewGameService()
	if err != nil {
		t.Fatalf("NewGameService() error = %v", err)
	}

	tests := []struct {
		name     string
		getter   func() interface{}
		expected interface{}
	}{
		{
			name:     "GetScore初期値",
			getter:   func() interface{} { return gameService.GetScore() },
			expected: 0,
		},
		{
			name:     "GetLines初期値",
			getter:   func() interface{} { return gameService.GetLines() },
			expected: 0,
		},
		{
			name:     "GetLevel初期値",
			getter:   func() interface{} { return gameService.GetLevel() },
			expected: 1,
		},
		{
			name:     "IsGameOver初期値",
			getter:   func() interface{} { return gameService.IsGameOver() },
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			if result != tt.expected {
				t.Errorf("Getter() = %v, want %v", result, tt.expected)
			}
		})
	}

	if gameService.GetBoard() == nil {
		t.Error("GetBoard() returned nil")
	}

	if gameService.GetCurrentPiece() == nil {
		t.Error("GetCurrentPiece() returned nil")
	}

	if gameService.GetNextPiece() == nil {
		t.Error("GetNextPiece() returned nil")
	}
}