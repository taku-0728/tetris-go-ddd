package console

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"tetris/application"
	"tetris/domain/model"
)

const (
	EmptyBlock  = "  "
	FilledBlock = "██"
	WallBlock   = "│"
	FloorBlock  = "─"
	CornerBlock = "└"
)

type Display struct {
	width  int
	height int
}

func NewDisplay() *Display {
	return &Display{
		width:  model.BoardWidth,
		height: model.BoardHeight,
	}
}

func (d *Display) ClearScreen() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (d *Display) Render(gameState application.GameState) error {
	if err := d.ClearScreen(); err != nil {
		return fmt.Errorf("画面クリアエラー: %w", err)
	}

	d.printHeader()
	d.printGameInfo(gameState)
	d.printBoard(gameState)
	d.printControls()

	if gameState.GameOver {
		d.printGameOver(gameState)
	}

	return nil
}

func (d *Display) printHeader() {
	fmt.Println("┌" + strings.Repeat("─", 40) + "┐")
	fmt.Println("│" + centerText("テトリス", 40) + "│")
	fmt.Println("├" + strings.Repeat("─", 40) + "┤")
}

func (d *Display) printGameInfo(gameState application.GameState) {
	fmt.Printf("│ スコア: %-10d ライン: %-10d │\n", gameState.Score, gameState.Lines)
	fmt.Printf("│ レベル: %-10d                    │\n", gameState.Level)
	fmt.Println("├" + strings.Repeat("─", 40) + "┤")
}

func (d *Display) printBoard(gameState application.GameState) {
	board := gameState.Board
	currentPiece := gameState.CurrentPiece

	gameBoard := make([][]bool, d.height)
	for i := range gameBoard {
		gameBoard[i] = make([]bool, d.width)
		copy(gameBoard[i], board.Grid[i])
	}

	if currentPiece != nil {
		blocks := currentPiece.GetBlocks()
		for _, block := range blocks {
			if block.Y >= 0 && block.Y < d.height && block.X >= 0 && block.X < d.width {
				gameBoard[block.Y][block.X] = true
			}
		}
	}

	for y := 0; y < d.height; y++ {
		fmt.Print("│")
		for x := 0; x < d.width; x++ {
			if gameBoard[y][x] {
				fmt.Print(FilledBlock)
			} else {
				fmt.Print(EmptyBlock)
			}
		}
		fmt.Println("│")
	}

	fmt.Println("└" + strings.Repeat("─", d.width*2) + "┘")
}

func (d *Display) printControls() {
	fmt.Println()
	fmt.Println("操作方法:")
	fmt.Println("  A/D: 左右移動")
	fmt.Println("  S: 下移動")
	fmt.Println("  W: 回転")
	fmt.Println("  Space: 一気に落下")
	fmt.Println("  P: 一時停止")
	fmt.Println("  Q: 終了")
}

func (d *Display) printGameOver(gameState application.GameState) {
	fmt.Println()
	fmt.Println("┌" + strings.Repeat("─", 30) + "┐")
	fmt.Println("│" + centerText("ゲームオーバー！", 30) + "│")
	fmt.Printf("│" + centerText(fmt.Sprintf("最終スコア: %d", gameState.Score), 30) + "│\n")
	fmt.Printf("│" + centerText(fmt.Sprintf("消去ライン: %d", gameState.Lines), 30) + "│\n")
	fmt.Println("│" + centerText("Rでリスタート、Qで終了", 30) + "│")
	fmt.Println("└" + strings.Repeat("─", 30) + "┘")
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}

	padding := width - len(text)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
}
