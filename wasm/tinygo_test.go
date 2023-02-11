package wasm

import (
	"os"
	"testing"
)

var src = `package main

import (
	"syscall/js"
	"math/rand"
)


var buffer []byte

const WIDTH = 600
const HEIGHT = 600
const CELL_SIZE = 8

//export getIndex
func getIndex(row int, col int) int {
	return row * WIDTH + col
}

//export getWidth
func getWidth() int { return WIDTH }

//export getHeight
func getHeight() int { return HEIGHT }

//export getCellSize
func getCellSize() int { return CELL_SIZE }

//export initGame
func initGame() {
	buffer = make([]byte, WIDTH*HEIGHT)
	for i, _ := range buffer {
		if rand.Float64() > 0.5 {
			buffer[i] = 1
		} else {
			buffer[i] = 0
		}
	}
}

func getLiveNeighbours(row, col int) int {
	liveNeighbours := 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			if row+i < 0 || row+i >= HEIGHT || col+j < 0 || col+j >= WIDTH { // out of bounds
				continue
			}

			if buffer[getIndex(row+i, col+j)] == 1 {
				liveNeighbours++
			}
		}
	}

	return liveNeighbours
}

//export tick
func tick() {
	newBuffer := make([]byte, len(buffer))
	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			liveNeighbours := getLiveNeighbours(i, j)
			idx := getIndex(i, j)

			if buffer[idx] == 1 {
				if liveNeighbours < 2 || liveNeighbours > 3{
					newBuffer[idx] = 0
				} else {
					newBuffer[idx] = 1
				}
			} else {
				if liveNeighbours == 3 {
					newBuffer[idx] = 1
				} else {
					newBuffer[idx] = 0
				}
			}
		}
	}

	buffer = newBuffer
}

func updateBuffer(this js.Value, args []js.Value) any {
	tick()

	return js.CopyBytesToJS(js.Global().Get("buffer"), buffer)
}

func main() {
	c := make(chan int, 0)

	js.Global().Set("updateBuffer", js.FuncOf(updateBuffer))

	<-c
}
`

func TestCompile_WorksWithValidGoFile(t *testing.T) {

	res, err := compileTinyGo(src, CompileOpts{
		GenWat: true,
	})

	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
	}

	if len(res.Wasm) == 0 {
		t.Error("Expected wasm to be non-empty")
	}

	if res.Wat == "" {
		t.Error("Expected wat to be non-empty")
	}

}

func TestCompile_ReturnsErrorWithInvalidGoFile(t *testing.T) {
	_, err := compileTinyGo("package main", CompileOpts{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCompile_BeforeDelete(t *testing.T) {
	_, err := compileTinyGo(src, CompileOpts{
		BeforeDelete: func(f *os.File) error {
			if _, err := os.Stat(f.Name()); err != nil {
				t.Errorf("Expected file to exist, got %s", err)
				return err
			}
			return nil
		},
	})

	if err != nil {
		t.Error(err)
	}
}
