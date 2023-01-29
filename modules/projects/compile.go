package projects

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func createTempCodeDir(code string, goMod string) (string, error, func()) {
	tmpDir, err := os.MkdirTemp("", "project-dir-*")
	if err != nil {
		return "", err, nil
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
		return "", err, nil
	}

	if _, err := codeFile.Write([]byte(code)); err != nil {
		deleteDir()
		return "", err, nil
	}

	goModFile, err := createInTemp("go.mod")
	if err != nil {
		deleteDir()
		return "", err, nil
	}
	if _, err := goModFile.Write([]byte(goMod)); err != nil {
		deleteDir()
		return "", err, nil
	}

	return tmpDir, nil, deleteDir
}

// installDeps runs go get -d ./... in the given directory to install dependencies
func installDeps(dir string) error {
	cmd := exec.Command("go", "get", "-d", "./...")
	cmd.Dir = dir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v", stderr.String())
	}

	return nil
}

/*
compileProject takes a string of Go code and a string containing a go.mod file and compiles it to web assembly using tinygo, returning a
reader to the compiled wasm file
*/
func compileProject(code string, goMod string) (io.Reader, error) {

	dir, err, deleteDir := createTempCodeDir(code, goMod)
	if err != nil {
		return nil, err
	}
	defer deleteDir()

	filename := "main.go"
	out := "main.wasm"

	if err := installDeps(dir); err != nil {
		return nil, err
	}

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
			return nil, errors.New(strings.Join(errs, "\n"))
		}

		if err != nil {
			return nil, fmt.Errorf("%v", stderr.String())
		}

	}

	return os.Open(path.Join(dir, out))

}
