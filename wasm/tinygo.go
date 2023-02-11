package wasm

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

/*
compileProject takes a string of Go code  and compiles it to WASM
*/
func compileTinyGo(code string, opts CompileOpts) (CompileResult, error) {
	result := CompileResult{}

	dir, deleteDir, err := createTempCodeDir("main.go", code)
	if err != nil {
		return result, err
	}
	defer deleteDir()

	filename := "main.go"
	out := "main.wasm"

	cmd := exec.Command("tinygo", "build", "-o", out, "-target", "wasm", filename)
	cmd.Dir = dir

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	errs := []string{}
	if err := cmd.Run(); err != nil {
		if err != nil {
			for _, line := range strings.Split(stderr.String(), "\n") {
				if strings.Contains(line, filename) {
					errs = append(errs, line)
				}
			}

			return result, errors.New(strings.Join(errs, "\n"))
		}
	}

	f, err := os.Open(path.Join(dir, out))
	if err != nil {
		return result, err
	}

	defer f.Close()

	if opts.BeforeDelete != nil {
		if err := opts.BeforeDelete(f); err != nil {
			return result, err
		}
	}

	wasmBytes, err := io.ReadAll(f)
	if err != nil {
		return result, err
	}

	result.Wasm = wasmBytes
	if opts.GenWat {
		wat, err := WasmToWat(bytes.NewReader(wasmBytes))
		if err != nil {
			return result, err
		}

		result.Wat = wat
	}

	return result, nil

}
