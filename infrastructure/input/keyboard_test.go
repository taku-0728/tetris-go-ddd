package input

import (
	"errors"
	"testing"
)

func TestMapInputToCommand(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
		errorType   error
	}{
		// 移動コマンド
		{
			name:        "小文字a - 左移動",
			input:       "a",
			expected:    "left",
			expectError: false,
		},
		{
			name:        "大文字A - 左移動",
			input:       "A",
			expected:    "left",
			expectError: false,
		},
		{
			name:        "left - 左移動",
			input:       "left",
			expected:    "left",
			expectError: false,
		},
		{
			name:        "小文字d - 右移動",
			input:       "d",
			expected:    "right",
			expectError: false,
		},
		{
			name:        "大文字D - 右移動",
			input:       "D",
			expected:    "right",
			expectError: false,
		},
		{
			name:        "right - 右移動",
			input:       "right",
			expected:    "right",
			expectError: false,
		},
		{
			name:        "小文字s - 下移動",
			input:       "s",
			expected:    "down",
			expectError: false,
		},
		{
			name:        "大文字S - 下移動",
			input:       "S",
			expected:    "down",
			expectError: false,
		},
		{
			name:        "down - 下移動",
			input:       "down",
			expected:    "down",
			expectError: false,
		},
		// 回転コマンド
		{
			name:        "小文字w - 回転",
			input:       "w",
			expected:    "rotate",
			expectError: false,
		},
		{
			name:        "大文字W - 回転",
			input:       "W",
			expected:    "rotate",
			expectError: false,
		},
		{
			name:        "rotate - 回転",
			input:       "rotate",
			expected:    "rotate",
			expectError: false,
		},
		// ドロップコマンド
		{
			name:        "スペース - ドロップ",
			input:       " ",
			expected:    "drop",
			expectError: false,
		},
		{
			name:        "drop - ドロップ",
			input:       "drop",
			expected:    "drop",
			expectError: false,
		},
		// 一時停止コマンド
		{
			name:        "小文字p - 一時停止",
			input:       "p",
			expected:    "pause",
			expectError: false,
		},
		{
			name:        "大文字P - 一時停止",
			input:       "P",
			expected:    "pause",
			expectError: false,
		},
		{
			name:        "pause - 一時停止",
			input:       "pause",
			expected:    "pause",
			expectError: false,
		},
		// 終了コマンド
		{
			name:        "小文字q - 終了",
			input:       "q",
			expected:    "quit",
			expectError: false,
		},
		{
			name:        "大文字Q - 終了",
			input:       "Q",
			expected:    "quit",
			expectError: false,
		},
		{
			name:        "quit - 終了",
			input:       "quit",
			expected:    "quit",
			expectError: false,
		},
		// リスタートコマンド
		{
			name:        "小文字r - リスタート",
			input:       "r",
			expected:    "restart",
			expectError: false,
		},
		{
			name:        "大文字R - リスタート",
			input:       "R",
			expected:    "restart",
			expectError: false,
		},
		{
			name:        "restart - リスタート",
			input:       "restart",
			expected:    "restart",
			expectError: false,
		},
		// 無効な入力
		{
			name:        "無効な文字",
			input:       "x",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
		{
			name:        "無効な文字列",
			input:       "invalid",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
		{
			name:        "空文字",
			input:       "",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
		{
			name:        "数字",
			input:       "1",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
		{
			name:        "特殊文字",
			input:       "@",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
		{
			name:        "複数文字の無効入力",
			input:       "abc",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
		{
			name:        "タブ文字",
			input:       "\t",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
		{
			name:        "改行文字",
			input:       "\n",
			expected:    "",
			expectError: true,
			errorType:   ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MapInputToCommand(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("MapInputToCommand() error = nil, wantErr %v", tt.errorType)
					return
				}
				if !errors.Is(err, tt.errorType) {
					t.Errorf("MapInputToCommand() error = %v, wantErr %v", err, tt.errorType)
				}
				return
			}

			if err != nil {
				t.Errorf("MapInputToCommand() unexpected error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("MapInputToCommand() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewKeyboardInput(t *testing.T) {
	keyboardInput := NewKeyboardInput()

	tests := []struct {
		name    string
		check   func() bool
		message string
	}{
		{
			name:    "KeyboardInputが正常に作成される",
			check:   func() bool { return keyboardInput != nil },
			message: "NewKeyboardInput() returned nil",
		},
		{
			name:    "inputChanが初期化される",
			check:   func() bool { return keyboardInput.inputChan != nil },
			message: "NewKeyboardInput() inputChan is nil",
		},
		{
			name:    "contextが初期化される",
			check:   func() bool { return keyboardInput.ctx != nil },
			message: "NewKeyboardInput() context is nil",
		},
		{
			name:    "cancelが初期化される",
			check:   func() bool { return keyboardInput.cancel != nil },
			message: "NewKeyboardInput() cancel function is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.message)
			}
		})
	}

	keyboardInput.Stop()
}

func TestKeyboardInput_StartAndStop(t *testing.T) {
	tests := []struct {
		name      string
		operation func(*KeyboardInput) error
	}{
		{
			name: "Start操作",
			operation: func(ki *KeyboardInput) error {
				return ki.Start()
			},
		},
		{
			name: "Stop操作",
			operation: func(ki *KeyboardInput) error {
				ki.Stop()
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyboardInput := NewKeyboardInput()
			defer keyboardInput.Stop()

			err := tt.operation(keyboardInput)
			if err != nil {
				t.Errorf("Operation failed with error: %v", err)
			}
		})
	}
}

func TestKeyboardInput_ContextCancellation(t *testing.T) {
	keyboardInput := NewKeyboardInput()

	if keyboardInput.ctx.Err() != nil {
		t.Error("Context should not be canceled initially")
	}

	keyboardInput.Stop()

	select {
	case <-keyboardInput.ctx.Done():
	default:
		t.Error("Context should be canceled after Stop()")
	}
}

func TestInputCommandMapping_Completeness(t *testing.T) {
	expectedCommands := []string{
		"left", "right", "down", "rotate", "drop", "pause", "quit", "restart",
	}

	tests := []struct {
		name   string
		inputs []string
	}{
		{
			name:   "左移動の全バリエーション",
			inputs: []string{"a", "A", "left"},
		},
		{
			name:   "右移動の全バリエーション",
			inputs: []string{"d", "D", "right"},
		},
		{
			name:   "下移動の全バリエーション",
			inputs: []string{"s", "S", "down"},
		},
		{
			name:   "回転の全バリエーション",
			inputs: []string{"w", "W", "rotate"},
		},
		{
			name:   "ドロップの全バリエーション",
			inputs: []string{" ", "drop"},
		},
		{
			name:   "一時停止の全バリエーション",
			inputs: []string{"p", "P", "pause"},
		},
		{
			name:   "終了の全バリエーション",
			inputs: []string{"q", "Q", "quit"},
		},
		{
			name:   "リスタートの全バリエーション",
			inputs: []string{"r", "R", "restart"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var commands []string
			for _, input := range tt.inputs {
				command, err := MapInputToCommand(input)
				if err != nil {
					t.Errorf("MapInputToCommand(%q) error = %v", input, err)
					continue
				}
				commands = append(commands, command)
			}

			if len(commands) == 0 {
				t.Error("No valid commands found")
				return
			}

			expectedCommand := commands[0]
			for _, cmd := range commands {
				if cmd != expectedCommand {
					t.Errorf("Inconsistent command mapping: got %q, expected %q", cmd, expectedCommand)
				}
			}

			found := false
			for _, expected := range expectedCommands {
				if expectedCommand == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Command %q is not in expected commands list", expectedCommand)
			}
		})
	}
}
