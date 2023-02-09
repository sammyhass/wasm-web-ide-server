package wasm

import (
	"go/ast"
	"strings"
	"testing"
)

func TestParse_FindsAllExportedFunctions(t *testing.T) {
	var src = `package main

//export add
func add(a, b int) int {
	return a + b
}

//export sub
func sub(a, b int) int {
	return a - b
}

//export mul
func mul(a, b int) int {
	return a * b
}

func main() {}`

	r := strings.NewReader((src))

	exports, err := Parse(r)

	if err != nil {
		t.Error(err)
	}

	if len(exports) != 3 {
		t.Errorf("Expected 3 exported functions, got %d", len(exports))
	}

	for _, export := range exports {
		if !strings.HasPrefix(export.comment.Text, "//export ") {
			t.Errorf("Expected comment to start with '//export ', got '%s'", export.comment.Text)
		}

		expDeclName := strings.TrimPrefix(export.comment.Text, "//export ")
		declName := export.decl.(*ast.FuncDecl).Name.Name

		if expDeclName != declName {
			t.Errorf("Expected exported function name to be '%s', got '%s'", expDeclName, declName)
		}
	}
}

func TestParse_CorrectlyParsesSignature(t *testing.T) {
	input := `package main

	//export add
	func add(a, b int) int {
		return a + b
	}

	//export stringify
	func stringify(a, b int) (string, error) {
		return "", nil
	}
	`
	cases := []struct {
		name     string
		nParams  int
		nReturns int
	}{
		{
			name:     "add",
			nParams:  2,
			nReturns: 2,
		},
		{
			name:     "stringify",
			nParams:  2,
			nReturns: 2,
		},
	}

	r := strings.NewReader(input)

	exports, err := Parse(r)

	if err != nil {
		t.Error(err)
	}

	for _, c := range cases {
		for _, export := range exports {
			if export.decl.(*ast.FuncDecl).Name.Name == c.name {
				if len(export.decl.(*ast.FuncDecl).Type.Params.List) != c.nParams {
					t.Errorf("Expected %d parameters, got %d", c.nParams, len(export.decl.(*ast.FuncDecl).Type.Params.List))
				}

				if len(export.decl.(*ast.FuncDecl).Type.Results.List) != c.nReturns {
					t.Errorf("Expected %d returns, got %d", c.nReturns, len(export.decl.(*ast.FuncDecl).Type.Results.List))
				}
			}
		}
	}
}
