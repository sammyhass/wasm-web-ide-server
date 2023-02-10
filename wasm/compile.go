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

func Compile(code string) ([]byte, error) {
	return CompileWithOpts(code, CompileOpts{})
}

type CompileOpts struct {
	BeforeDelete func(wasm *os.File) error // BeforeDelete is called before the temp directory is deleted, it is passed the compiled WASM file
}

/*
compileProject takes a string of Go code  and compiles it to WASM
*/
func CompileWithOpts(code string, opts CompileOpts) ([]byte, error) {

	dir, deleteDir, err := createTempCodeDir(code)
	if err != nil {
		return nil, err
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

			return nil, errors.New(strings.Join(errs, "\n"))
		}
	}

	f, err := os.Open(path.Join(dir, out))
	if err != nil {
		return nil, err
	}

	defer f.Close()

	if opts.BeforeDelete != nil {
		if err := opts.BeforeDelete(f); err != nil {
			return nil, err
		}
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}

func StripWasm(
	f *os.File,
) error {
	if f == nil {
		return errors.New("nil file")
	}

	cmd := exec.Command("wasm-strip", f.Name())
	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}

		return err
	}

	return nil

}
