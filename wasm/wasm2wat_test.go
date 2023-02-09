package wasm

import "testing"

func TestWasm2Wat(t *testing.T) {

	wasm, err := Compile(src)

	if err != nil {
		t.Error(err)
	}

	wat, err := WasmToWat(wasm)

	if err != nil {
		t.Error(err)
	}

	if wat == "" {
		t.Error("Expected wat to be non-empty")
	}

}
