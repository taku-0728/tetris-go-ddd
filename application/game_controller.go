package application

import (
	"errors"
	"fmt"
	"tetris/domain/model"
	"tetris/domain/service"
	"time"
)

type GameState struct {
	Board        *model.Board
	CurrentPiece *model.Tetromino
	NextPiece    *model.Tetromino
	Score        int
	Lines        int
	Level        int
	GameOver     bool
}

type GameController struct {
	gameService  *service.GameService
	dropTimer    time.Time
	dropInterval time.Duration
	isPaused     bool
}

func NewGameController() (*GameController, error) {
	gameService, err := service.NewGameService()
	if err != nil {
		return nil, fmt.Errorf("ゲームサービス初期化エラー: %w", err)
	}

	return &GameController{
		gameService:  gameService,
		dropTimer:    time.Now(),
		dropInterval: time.Second,
		isPaused:     false,
	}, nil
}

func (gc *GameController) GetGameState() GameState {
	return GameState{
		Board:        gc.gameService.GetBoard(),
		CurrentPiece: gc.gameService.GetCurrentPiece(),
		NextPiece:    gc.gameService.GetNextPiece(),
		Score:        gc.gameService.GetScore(),
		Lines:        gc.gameService.GetLines(),
		Level:        gc.gameService.GetLevel(),
		GameOver:     gc.gameService.IsGameOver(),
	}
}

func (gc *GameController) Update() error {
	if gc.isPaused || gc.gameService.IsGameOver() {
		return nil
	}

	if time.Since(gc.dropTimer) >= gc.dropInterval {
		if err := gc.gameService.Update(); err != nil {
			if errors.Is(err, service.ErrGameOver) {
				return nil
			}
			return fmt.Errorf("ゲーム更新エラー: %w", err)
		}
		gc.dropTimer = time.Now()
		gc.updateDropInterval()
	}

	return nil
}

func (gc *GameController) HandleInput(input string) error {
	if gc.gameService.IsGameOver() {
		return nil
	}

	switch input {
	case "left", "a", "A":
		return gc.movePieceLeft()
	case "right", "d", "D":
		return gc.movePieceRight()
	case "down", "s", "S":
		return gc.movePieceDown()
	case "rotate", "w", "W":
		return gc.rotatePiece()
	case "drop", "space":
		return gc.dropPiece()
	case "pause", "p", "P":
		gc.togglePause()
		return nil
	default:
		return fmt.Errorf("不明な入力: %s", input)
	}
}

func (gc *GameController) movePieceLeft() error {
	err := gc.gameService.MovePiece(model.Point{X: -1, Y: 0})
	if err != nil && !errors.Is(err, service.ErrInvalidMove) {
		return fmt.Errorf("左移動エラー: %w", err)
	}
	return nil
}

func (gc *GameController) movePieceRight() error {
	err := gc.gameService.MovePiece(model.Point{X: 1, Y: 0})
	if err != nil && !errors.Is(err, service.ErrInvalidMove) {
		return fmt.Errorf("右移動エラー: %w", err)
	}
	return nil
}

func (gc *GameController) movePieceDown() error {
	err := gc.gameService.MovePiece(model.Point{X: 0, Y: 1})
	if err != nil && !errors.Is(err, service.ErrInvalidMove) {
		return fmt.Errorf("下移動エラー: %w", err)
	}
	if err == nil {
		gc.dropTimer = time.Now()
	}
	return nil
}

func (gc *GameController) rotatePiece() error {
	err := gc.gameService.RotatePiece()
	if err != nil && !errors.Is(err, service.ErrInvalidMove) {
		return fmt.Errorf("回転エラー: %w", err)
	}
	return nil
}

func (gc *GameController) dropPiece() error {
	err := gc.gameService.DropPiece()
	if err != nil {
		return fmt.Errorf("ドロップエラー: %w", err)
	}
	gc.dropTimer = time.Now()
	return nil
}

func (gc *GameController) togglePause() {
	gc.isPaused = !gc.isPaused
}

func (gc *GameController) IsPaused() bool {
	return gc.isPaused
}

func (gc *GameController) updateDropInterval() {
	level := gc.gameService.GetLevel()
	baseInterval := 1000 * time.Millisecond
	reduction := time.Duration(level-1) * 100 * time.Millisecond

	gc.dropInterval = baseInterval - reduction
	if gc.dropInterval < 100*time.Millisecond {
		gc.dropInterval = 100 * time.Millisecond
	}
}

func (gc *GameController) Reset() error {
	gameService, err := service.NewGameService()
	if err != nil {
		return fmt.Errorf("ゲームリセットエラー: %w", err)
	}

	gc.gameService = gameService
	gc.dropTimer = time.Now()
	gc.dropInterval = time.Second
	gc.isPaused = false

	return nil
}
