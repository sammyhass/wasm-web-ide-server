package wasm

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sammyhass/web-ide/server/wasm/tinygo"
)

func TestTinyGoWasm2Wat(t *testing.T) {

	wasm, err := tinygo.Compile(`
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
`)

	if err != nil {
		t.Error(err)
	}

	wat, err := WasmToWat(bytes.NewReader(wasm))

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
