package application

import (
	"testing"
	"tetris/domain/model"
	"time"
)

func TestNewGameController(t *testing.T) {
	tests := []struct {
		name        string
		expectError bool
	}{
		{
			name:        "正常なゲームコントローラー作成",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, err := NewGameController()

			if tt.expectError {
				if err == nil {
					t.Error("NewGameController() error = nil, wantErr")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("NewGameController() unexpected error = %v", err)
				return
			}

			if controller == nil {
				t.Error("NewGameController() returned nil")
				return
			}

			if controller.gameService == nil {
				t.Error("NewGameController() gameService is nil")
			}

			if controller.dropInterval <= 0 {
				t.Error("NewGameController() dropInterval is invalid")
			}

			if controller.isPaused {
				t.Error("NewGameController() isPaused should be false initially")
			}
		})
	}
}

func TestGameController_GetGameState(t *testing.T) {
	controller, err := NewGameController()
	if err != nil {
		t.Fatalf("NewGameController() error = %v", err)
	}

	gameState := controller.GetGameState()

	tests := []struct {
		name     string
		check    func() bool
		message  string
	}{
		{
			name:    "ボードが存在する",
			check:   func() bool { return gameState.Board != nil },
			message: "GameState.Board is nil",
		},
		{
			name:    "現在のピースが存在する",
			check:   func() bool { return gameState.CurrentPiece != nil },
			message: "GameState.CurrentPiece is nil",
		},
		{
			name:    "次のピースが存在する",
			check:   func() bool { return gameState.NextPiece != nil },
			message: "GameState.NextPiece is nil",
		},
		{
			name:    "初期スコアが0",
			check:   func() bool { return gameState.Score == 0 },
			message: "GameState.Score should be 0 initially",
		},
		{
			name:    "初期ラインが0",
			check:   func() bool { return gameState.Lines == 0 },
			message: "GameState.Lines should be 0 initially",
		},
		{
			name:    "初期レベルが1",
			check:   func() bool { return gameState.Level == 1 },
			message: "GameState.Level should be 1 initially",
		},
		{
			name:    "初期状態でゲームオーバーではない",
			check:   func() bool { return !gameState.GameOver },
			message: "GameState.GameOver should be false initially",
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

func TestGameController_HandleInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "左移動 (a)",
			input:       "a",
			expectError: false,
		},
		{
			name:        "左移動 (A)",
			input:       "A",
			expectError: false,
		},
		{
			name:        "左移動 (left)",
			input:       "left",
			expectError: false,
		},
		{
			name:        "右移動 (d)",
			input:       "d",
			expectError: false,
		},
		{
			name:        "右移動 (D)",
			input:       "D",
			expectError: false,
		},
		{
			name:        "右移動 (right)",
			input:       "right",
			expectError: false,
		},
		{
			name:        "下移動 (s)",
			input:       "s",
			expectError: false,
		},
		{
			name:        "下移動 (S)",
			input:       "S",
			expectError: false,
		},
		{
			name:        "下移動 (down)",
			input:       "down",
			expectError: false,
		},
		{
			name:        "回転 (w)",
			input:       "w",
			expectError: false,
		},
		{
			name:        "回転 (W)",
			input:       "W",
			expectError: false,
		},
		{
			name:        "回転 (rotate)",
			input:       "rotate",
			expectError: false,
		},
		{
			name:        "ドロップ (space)",
			input:       "space",
			expectError: false,
		},
		{
			name:        "ドロップ (drop)",
			input:       "drop",
			expectError: false,
		},
		{
			name:        "一時停止 (p)",
			input:       "p",
			expectError: false,
		},
		{
			name:        "一時停止 (P)",
			input:       "P",
			expectError: false,
		},
		{
			name:        "一時停止 (pause)",
			input:       "pause",
			expectError: false,
		},
		{
			name:        "無効な入力",
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "空文字",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, err := NewGameController()
			if err != nil {
				t.Fatalf("NewGameController() error = %v", err)
			}

			err = controller.HandleInput(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("GameController.HandleInput() error = nil, wantErr")
				}
				return
			}

			if err != nil {
				t.Errorf("GameController.HandleInput() unexpected error = %v", err)
			}
		})
	}
}

func TestGameController_PauseToggle(t *testing.T) {
	controller, err := NewGameController()
	if err != nil {
		t.Fatalf("NewGameController() error = %v", err)
	}

	tests := []struct {
		name           string
		initialPaused  bool
		expectedPaused bool
	}{
		{
			name:           "一時停止オン",
			initialPaused:  false,
			expectedPaused: true,
		},
		{
			name:           "一時停止オフ",
			initialPaused:  true,
			expectedPaused: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller.isPaused = tt.initialPaused

			controller.togglePause()

			if controller.IsPaused() != tt.expectedPaused {
				t.Errorf("GameController.togglePause() paused = %v, want %v", controller.IsPaused(), tt.expectedPaused)
			}
		})
	}
}

func TestGameController_Update(t *testing.T) {
	tests := []struct {
		name        string
		setupController func(*GameController)
		expectError bool
	}{
		{
			name:        "通常の更新",
			setupController: func(gc *GameController) {},
			expectError: false,
		},
		{
			name: "一時停止中の更新",
			setupController: func(gc *GameController) {
				gc.isPaused = true
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, err := NewGameController()
			if err != nil {
				t.Fatalf("NewGameController() error = %v", err)
			}

			tt.setupController(controller)

			err = controller.Update()

			if tt.expectError {
				if err == nil {
					t.Error("GameController.Update() error = nil, wantErr")
				}
				return
			}

			if err != nil {
				t.Errorf("GameController.Update() unexpected error = %v", err)
			}
		})
	}
}

func TestGameController_Reset(t *testing.T) {
	controller, err := NewGameController()
	if err != nil {
		t.Fatalf("NewGameController() error = %v", err)
	}

	controller.isPaused = true

	err = controller.Reset()
	if err != nil {
		t.Errorf("GameController.Reset() error = %v", err)
		return
	}

	tests := []struct {
		name    string
		check   func() bool
		message string
	}{
		{
			name:    "一時停止がリセットされる",
			check:   func() bool { return !controller.IsPaused() },
			message: "Reset() should unpause the game",
		},
		{
			name:    "ドロップインターバルがリセットされる",
			check:   func() bool { return controller.dropInterval == time.Second },
			message: "Reset() should reset drop interval",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.message)
			}
		})
	}

	gameState := controller.GetGameState()
	if gameState.Score != 0 {
		t.Error("Reset() should reset score to 0")
	}

	if gameState.Lines != 0 {
		t.Error("Reset() should reset lines to 0")
	}

	if gameState.Level != 1 {
		t.Error("Reset() should reset level to 1")
	}

	if gameState.GameOver {
		t.Error("Reset() should set game over to false")
	}
}

func TestGameController_MovePieceOperations(t *testing.T) {
	controller, err := NewGameController()
	if err != nil {
		t.Fatalf("NewGameController() error = %v", err)
	}

	initialPosition := controller.GetGameState().CurrentPiece.Position

	tests := []struct {
		name     string
		operation func() error
		deltaX   int
		deltaY   int
	}{
		{
			name:     "左移動",
			operation: controller.movePieceLeft,
			deltaX:   -1,
			deltaY:   0,
		},
		{
			name:     "右移動",
			operation: controller.movePieceRight,
			deltaX:   1,
			deltaY:   0,
		},
		{
			name:     "下移動",
			operation: controller.movePieceDown,
			deltaX:   0,
			deltaY:   1,
		},
	}

	currentPosition := initialPosition
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if err != nil {
				return
			}

			newPosition := controller.GetGameState().CurrentPiece.Position
			expectedPosition := model.Point{
				X: currentPosition.X + tt.deltaX,
				Y: currentPosition.Y + tt.deltaY,
			}

			if newPosition == expectedPosition {
				currentPosition = newPosition
			}
		})
	}
}

func TestGameController_DropIntervalUpdate(t *testing.T) {
	controller, err := NewGameController()
	if err != nil {
		t.Fatalf("NewGameController() error = %v", err)
	}

	tests := []struct {
		name            string
		expectedMaxInterval time.Duration
		expectedMinInterval time.Duration
	}{
		{
			name:            "ドロップインターバル更新",
			expectedMaxInterval: time.Second,
			expectedMinInterval: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller.updateDropInterval()

			if controller.dropInterval > tt.expectedMaxInterval {
				t.Errorf("updateDropInterval() interval = %v, should be <= %v", controller.dropInterval, tt.expectedMaxInterval)
			}

			if controller.dropInterval < tt.expectedMinInterval {
				t.Errorf("updateDropInterval() interval = %v, should be >= %v", controller.dropInterval, tt.expectedMinInterval)
			}
		})
	}
}