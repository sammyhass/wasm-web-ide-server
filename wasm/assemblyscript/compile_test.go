package assemblyscript

import (
	"strings"
	"testing"
)

func TestCompile_ValidAssemblyScript(t *testing.T) {
	src := `
		export function add(a: i32, b: i32): i32 {
			return a + b;
		}

		export function sub(a: i32, b: i32): i32 {
			return a - b;
		}
		`

	res, err := Compile(src)

	if err != nil {
		t.Error(err)
	}

	if len(res.Wasm) == 0 {
		t.Error("Expected wasm to be non-empty")
	}

	if res.Wat == "" {
		t.Error("Expected wat to be non-empty")
	}

	if !strings.Contains(res.Wat, "add") {
		t.Error("Expected wat to contain 'add'")
	}

	if !strings.Contains(res.Wat, "sub") {
		t.Error("Expected wat to contain 'sub'")
	}

}
