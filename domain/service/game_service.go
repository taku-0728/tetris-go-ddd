package service

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"tetris/domain/model"
)

var (
	ErrGameOver    = errors.New("ゲームが終了しています")
	ErrInvalidMove = errors.New("無効な移動です")
	ErrNoPiece     = errors.New("アクティブなピースがありません")
)

type GameService struct {
	board        *model.Board
	currentPiece *model.Tetromino
	nextPiece    *model.Tetromino
	score        int
	lines        int
	level        int
	gameOver     bool
}

func NewGameService() (*GameService, error) {
	board, err := model.NewBoard(model.BoardWidth, model.BoardHeight)
	if err != nil {
		return nil, fmt.Errorf("ボード作成エラー: %w", err)
	}

	service := &GameService{
		board:    board,
		score:    0,
		lines:    0,
		level:    1,
		gameOver: false,
	}

	if err := service.spawnNewPiece(); err != nil {
		return nil, fmt.Errorf("初期ピース生成エラー: %w", err)
	}

	if err := service.generateNextPiece(); err != nil {
		return nil, fmt.Errorf("次ピース生成エラー: %w", err)
	}

	return service, nil
}

func (g *GameService) GetBoard() *model.Board {
	return g.board
}

func (g *GameService) GetCurrentPiece() *model.Tetromino {
	return g.currentPiece
}

func (g *GameService) GetNextPiece() *model.Tetromino {
	return g.nextPiece
}

func (g *GameService) GetScore() int {
	return g.score
}

func (g *GameService) GetLines() int {
	return g.lines
}

func (g *GameService) GetLevel() int {
	return g.level
}

func (g *GameService) IsGameOver() bool {
	return g.gameOver
}

func (g *GameService) MovePiece(delta model.Point) error {
	if g.gameOver {
		return ErrGameOver
	}
	if g.currentPiece == nil {
		return ErrNoPiece
	}

	originalPosition := g.currentPiece.Position
	if err := g.currentPiece.Move(delta); err != nil {
		return fmt.Errorf("ピース移動エラー: %w", err)
	}

	if !g.board.CanPlaceTetromino(g.currentPiece) {
		g.currentPiece.Position = originalPosition
		return ErrInvalidMove
	}

	return nil
}

func (g *GameService) RotatePiece() error {
	if g.gameOver {
		return ErrGameOver
	}
	if g.currentPiece == nil {
		return ErrNoPiece
	}

	originalShape := make([][]bool, len(g.currentPiece.Shape))
	for i := range originalShape {
		originalShape[i] = make([]bool, len(g.currentPiece.Shape[i]))
		copy(originalShape[i], g.currentPiece.Shape[i])
	}

	if err := g.currentPiece.Rotate(); err != nil {
		return fmt.Errorf("ピース回転エラー: %w", err)
	}

	if !g.board.CanPlaceTetromino(g.currentPiece) {
		g.currentPiece.Shape = originalShape
		return ErrInvalidMove
	}

	return nil
}

func (g *GameService) DropPiece() error {
	if g.gameOver {
		return ErrGameOver
	}

	for {
		err := g.MovePiece(model.Point{X: 0, Y: 1})
		if err != nil {
			if errors.Is(err, ErrInvalidMove) {
				break
			}
			return err
		}
	}

	return g.lockPiece()
}

func (g *GameService) Update() error {
	if g.gameOver {
		return ErrGameOver
	}

	err := g.MovePiece(model.Point{X: 0, Y: 1})
	if err != nil {
		if errors.Is(err, ErrInvalidMove) {
			return g.lockPiece()
		}
		return err
	}

	return nil
}

func (g *GameService) lockPiece() error {
	if g.currentPiece == nil {
		return ErrNoPiece
	}

	if err := g.board.PlaceTetromino(g.currentPiece); err != nil {
		return fmt.Errorf("ピース配置エラー: %w", err)
	}

	completedLines := g.board.GetCompletedLines()
	if len(completedLines) > 0 {
		if err := g.board.ClearLines(completedLines); err != nil {
			return fmt.Errorf("ライン消去エラー: %w", err)
		}
		g.updateScore(len(completedLines))
	}

	if g.board.IsGameOver() {
		g.gameOver = true
		return nil
	}

	g.currentPiece = g.nextPiece
	if err := g.generateNextPiece(); err != nil {
		return fmt.Errorf("次ピース生成エラー: %w", err)
	}

	if !g.board.CanPlaceTetromino(g.currentPiece) {
		g.gameOver = true
	}

	return nil
}

func (g *GameService) spawnNewPiece() error {
	tetrominoType := model.TetrominoType(rand.IntN(7))
	position := model.Point{X: model.BoardWidth/2 - 2, Y: 0}

	piece, err := model.NewTetromino(tetrominoType, position)
	if err != nil {
		return fmt.Errorf("テトロミノ生成エラー: %w", err)
	}

	g.currentPiece = piece
	return nil
}

func (g *GameService) generateNextPiece() error {
	tetrominoType := model.TetrominoType(rand.IntN(7))
	position := model.Point{X: model.BoardWidth/2 - 2, Y: 0}

	piece, err := model.NewTetromino(tetrominoType, position)
	if err != nil {
		return fmt.Errorf("次テトロミノ生成エラー: %w", err)
	}

	g.nextPiece = piece
	return nil
}

func (g *GameService) updateScore(linesCleared int) {
	g.lines += linesCleared
	g.level = (g.lines / 10) + 1

	baseScore := map[int]int{
		1: 100,
		2: 300,
		3: 500,
		4: 800,
	}

	if score, exists := baseScore[linesCleared]; exists {
		g.score += score * g.level
	}
}
