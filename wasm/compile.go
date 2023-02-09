package wasm

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func createTempCodeDir(code string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "project-dir-*")
	if err != nil {
		return "", nil, err
	}

	deleteDir := func() {
		os.RemoveAll(tmpDir)
	}

	createInTemp := func(filename string) (*os.File, error) {
		return os.Create(path.Join(tmpDir, filename))
	}

	codeFile, err := createInTemp("main.go")
	if err != nil {
		deleteDir()
		return "", nil, err
	}

	if _, err := codeFile.Write([]byte(code)); err != nil {
		deleteDir()
		return "", nil, err
	}

	return tmpDir, deleteDir, nil
}

/*
compileProject takes a string of Go code and a string containing a go.mod file and compiles it to web assembly using tinygo, returning a reader to the compiled wasm file and a function to delete the temporary directory containing the compiled code.
*/
func Compile(code string) (*os.File, func(), error) {
	dir, deleteDir, err := createTempCodeDir(code)
	if err != nil {
		return nil, func() {}, err
	}

	filename := "main.go"
	out := "main.wasm"

	cmd := exec.Command("tinygo", "build", "-o", out, "-target", "wasm", filename)
	cmd.Dir = dir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	var errs []string

	if err := cmd.Run(); err != nil {
		for _, line := range strings.Split(stderr.String(), "\n") {
			hasFileName := strings.Contains(line, filename)
			if hasFileName {
				errs = append(errs, line)
			}
		}

		if len(errs) > 0 {
			return nil, deleteDir, errors.New(strings.Join(errs, "\n"))
		}

		if err != nil {
			return nil, deleteDir, fmt.Errorf("%v", stderr.String())
		}
	}

	stripWasm(path.Join(dir, out))

	wasmFile, err := os.Open(path.Join(dir, out))

	if err != nil {
		return nil, deleteDir, err
	}

	return wasmFile, deleteDir, nil

}

func stripWasm(
	fname string,
) error {

	cmd := exec.Command("wasm-strip", fname)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil

}
