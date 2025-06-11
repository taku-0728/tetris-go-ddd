package input

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	ErrInputCancelled = errors.New("入力がキャンセルされました")
	ErrInvalidInput   = errors.New("無効な入力です")
)

type KeyboardInput struct {
	inputChan chan string
	ctx       context.Context
	cancel    context.CancelFunc
	once      sync.Once
}

func NewKeyboardInput() *KeyboardInput {
	ctx, cancel := context.WithCancel(context.Background())
	return &KeyboardInput{
		inputChan: make(chan string, 10),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (k *KeyboardInput) Start() error {
	k.setupSignalHandler()
	go k.readInput()
	return nil
}

func (k *KeyboardInput) Stop() {
	k.once.Do(func() {
		k.cancel()
		close(k.inputChan)
	})
}

func (k *KeyboardInput) GetInput() (string, error) {
	select {
	case input, ok := <-k.inputChan:
		if !ok {
			return "", ErrInputCancelled
		}
		return input, nil
	case <-k.ctx.Done():
		return "", ErrInputCancelled
	}
}

func (k *KeyboardInput) HasInput() bool {
	select {
	case <-k.inputChan:
		return true
	default:
		return false
	}
}

func (k *KeyboardInput) readInput() {
	defer k.cancel()

	for {
		select {
		case <-k.ctx.Done():
			return
		default:
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				continue
			}

			select {
			case k.inputChan <- input:
			case <-k.ctx.Done():
				return
			}
		}
	}
}

func (k *KeyboardInput) setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			k.cancel()
		case <-k.ctx.Done():
		}
	}()
}

func MapInputToCommand(input string) (string, error) {
	commandMap := map[string]string{
		"a":       "left",
		"A":       "left",
		"d":       "right",
		"D":       "right",
		"s":       "down",
		"S":       "down",
		"w":       "rotate",
		"W":       "rotate",
		" ":       "drop",
		"p":       "pause",
		"P":       "pause",
		"q":       "quit",
		"Q":       "quit",
		"r":       "restart",
		"R":       "restart",
		"left":    "left",
		"right":   "right",
		"down":    "down",
		"rotate":  "rotate",
		"drop":    "drop",
		"pause":   "pause",
		"quit":    "quit",
		"restart": "restart",
	}

	if command, exists := commandMap[input]; exists {
		return command, nil
	}

	return "", fmt.Errorf("%w: %s", ErrInvalidInput, input)
}
