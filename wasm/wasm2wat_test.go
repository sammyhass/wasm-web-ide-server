package wasm

import (
	"bytes"
	"strings"
	"testing"
)

func TestTinyGoWasm2Wat(t *testing.T) {

	wasm, err := compileTinyGo(`
	package main

	//export add
	func add(a, b int) int {
		return a + b
	}


	//export sub
	func sub(a, b int) int {
		return a - b
	}

	func main() {}
`, CompileOpts{
		GenWat: false,
	})

	if err != nil {
		t.Error(err)
	}

	wat, err := WasmToWat(
		bytes.NewReader(wasm.Wasm),
	)

	if err != nil {
		t.Error(err)
	}

	if wat == "" {
		t.Error("Expected wat to be non-empty")
	}

	if !strings.Contains(wat, "add") {
		t.Error("Expected wat to contain 'add'")
	}

	if !strings.Contains(wat, "sub") {
		t.Error("Expected wat to contain 'sub'")
	}

}
