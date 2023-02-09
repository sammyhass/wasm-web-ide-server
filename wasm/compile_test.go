package wasm

import (
	"io"
	"testing"
)

var src = `package main

//export add
func add(a, b int) int {
	return a + b
}

//export sub
func sub(a, b int) int {
	return a - b
}

func main() {}`

func TestCompile_WorksWithValidGoFile(t *testing.T) {

	wasm, del, err := Compile(src)
	if err != nil {
		t.Error(err)
	}
	defer del()

	bytes, err := io.ReadAll(wasm)

	if err != nil {
		t.Error(err)
	}

	if len(bytes) == 0 {
		t.Error("Expected compiled wasm to be non-empty")
	}
}

func TestCompile_ReturnsErrorWithInvalidGoFile(t *testing.T) {
	_, delete, err := Compile("package main")
	defer delete()

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
