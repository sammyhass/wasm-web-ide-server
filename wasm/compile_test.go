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

	wasm, err := Compile(src)

	if err != nil {
		t.Error(err)
	}

	bytes, err := io.ReadAll(wasm)

	if err != nil {
		t.Error(err)
	}

	if len(bytes) == 0 {
		t.Error("Expected compiled wasm to be non-empty")
	}
}

func TestCompile_ReturnsErrorWithInvalidGoFile(t *testing.T) {
	_, err := Compile("package main")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
