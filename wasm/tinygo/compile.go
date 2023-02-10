package tinygo

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/sammyhass/web-ide/server/wasm/util"
)

func Compile(code string) ([]byte, error) {
	return compileWithOpts(code, compileOpts{})
}

type compileOpts struct {
	BeforeDelete func(wasm *os.File) error // BeforeDelete is called before the temp directory is deleted, it is passed the compiled WASM file
}

/*
compileProject takes a string of Go code  and compiles it to WASM
*/
func compileWithOpts(code string, opts compileOpts) ([]byte, error) {

	dir, deleteDir, err := util.CreateTempCodeDir("main.go", code)
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
