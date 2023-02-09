package wasm

import (
	"bytes"
	"strings"
	"testing"
)

func TestWasm2Wat(t *testing.T) {

	wasm, err := Compile(src)

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
