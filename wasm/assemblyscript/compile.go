package assemblyscript

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	"github.com/sammyhass/web-ide/server/wasm/util"
)

func Compile(assemblyScriptCode string) (util.CompileResult, error) {
	return compileWithOpts(assemblyScriptCode)
}

func compileWithOpts(code string) (util.CompileResult, error) {

	dir, delete, err := util.CreateTempCodeDir("main.ts", code)
	if err != nil {
		return util.CompileResult{}, err
	}
	defer delete()

	wasmF, err := os.CreateTemp(dir, "assemblyscript-*.wasm")
	if err != nil {
		return util.CompileResult{}, err
	}
	defer os.Remove(wasmF.Name())

	watF, err := os.CreateTemp(dir, "assemblyscript-*.wat")
	if err != nil {
		return util.CompileResult{}, err
	}
	defer os.Remove(watF.Name())

	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command("asc", "main.ts", "--outFile", wasmF.Name(), "--textFile", watF.Name())
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return util.CompileResult{}, errors.New(stderr.String())
		}
	}

	wasmBytes, err := os.ReadFile(wasmF.Name())
	if err != nil {
		return util.CompileResult{}, err
	}

	watBytes, err := os.ReadFile(watF.Name())
	if err != nil {
		return util.CompileResult{}, err
	}

	return util.CompileResult{
		Wasm: wasmBytes,
		Wat:  string(watBytes),
	}, nil

}
