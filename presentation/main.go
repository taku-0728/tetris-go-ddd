package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"tetris/application"
	"tetris/infrastructure/console"
	"tetris/infrastructure/input"
	"time"
)

func main() {
	if err := runGame(); err != nil {
		log.Fatalf("ゲーム実行エラー: %v", err)
	}
}

func runGame() error {
	rand.Seed(time.Now().UnixNano())

	gameController, err := application.NewGameController()
	if err != nil {
		return fmt.Errorf("ゲームコントローラー初期化エラー: %w", err)
	}

	display := console.NewDisplay()
	keyboardInput := input.NewKeyboardInput()

	if err := keyboardInput.Start(); err != nil {
		return fmt.Errorf("キーボード入力初期化エラー: %w", err)
	}
	defer keyboardInput.Stop()

	fmt.Println("テトリスゲームを開始します！")
	fmt.Println("何かキーを押してゲームを開始してください...")
	
	_, err = keyboardInput.GetInput()
	if err != nil && !errors.Is(err, input.ErrInputCancelled) {
		return fmt.Errorf("入力待機エラー: %w", err)
	}

	gameLoop := &GameLoop{
		controller: gameController,
		display:    display,
		input:      keyboardInput,
	}

	return gameLoop.Run()
}

type GameLoop struct {
	controller *application.GameController
	display    *console.Display
	input      *input.KeyboardInput
}

func (gl *GameLoop) Run() error {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := gl.update(); err != nil {
				return fmt.Errorf("ゲーム更新エラー: %w", err)
			}

			gameState := gl.controller.GetGameState()
			if err := gl.display.Render(gameState); err != nil {
				return fmt.Errorf("描画エラー: %w", err)
			}

		default:
			inputStr, err := gl.input.GetInput()
			if err != nil {
				if errors.Is(err, input.ErrInputCancelled) {
					return nil
				}
				continue
			}

			command, err := input.MapInputToCommand(inputStr)
			if err != nil {
				continue
			}

			if err := gl.handleCommand(command); err != nil {
				return err
			}
		}
	}
}

func (gl *GameLoop) update() error {
	return gl.controller.Update()
}

func (gl *GameLoop) handleCommand(command string) error {
	switch command {
	case "quit":
		return fmt.Errorf("ゲーム終了")
	case "restart":
		return gl.controller.Reset()
	default:
		return gl.controller.HandleInput(command)
	}
}