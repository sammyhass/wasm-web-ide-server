package wasm

import (
	"bytes"
	"io"
	"os"
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

	bytes, err := Compile(src)

	if err != nil {
		t.Error(err)
	}

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

func TestStripWASM(
	t *testing.T,
) {
	wasm, err := Compile(src)
	if err != nil {
		t.Error(err)
	}

	f, err := os.CreateTemp("", "wasm-strip-*.wasm")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, bytes.NewReader(wasm))
	if err != nil {
		t.Error(err)
	}

	if err = StripWasm(f); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	stat, err := f.Stat()
	if err != nil {
		t.Error(err)
	}

	if stat.Size() > int64(len(wasm)) {
		t.Errorf("Expected new file size to be smaller or equal to original, got %d", stat.Size())
	}

	t.Logf("Original size: %d, new size: %d", len(wasm), stat.Size())

}
