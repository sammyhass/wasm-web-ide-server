package wasm

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path"
)

func compileAssemblyScript(assemblyScriptCode string, options CompileOpts) (CompileResult, error) {
	codeFileName := "main.ts"
	dir, delete, err := createTempCodeDir(codeFileName, assemblyScriptCode)
	if err != nil {
		return CompileResult{}, err
	}
	defer delete()

	wasmF, err := os.CreateTemp(dir, "assemblyscript-*.wasm")
	if err != nil {
		return CompileResult{}, err
	}
	defer os.Remove(wasmF.Name())

	command := []string{codeFileName, "--outFile", path.Base(wasmF.Name())}
	var watF *os.File
	if options.GenWat {
		watF, err = os.CreateTemp(dir, "assemblyscript-*.wat")
		if err != nil {
			return CompileResult{}, err
		}
		defer os.Remove(watF.Name())

		command = append(command, "--textFile", path.Base(watF.Name()))
	}

	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command("asc", command...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return CompileResult{}, errors.New(stderr.String())
		}
	}

	wasmBytes, err := os.ReadFile(wasmF.Name())
	if err != nil {
		return CompileResult{}, err
	}

	if options.BeforeDelete != nil {
		if err := options.BeforeDelete(wasmF); err != nil {
			return CompileResult{}, err
		}
	}

	if options.GenWat && watF != nil {
		watBytes, err := os.ReadFile(watF.Name())
		if err != nil {
			return CompileResult{}, err
		}

		return CompileResult{
			Wasm: wasmBytes,
			Wat:  string(watBytes),
		}, nil
	}

	return CompileResult{
		Wasm: wasmBytes,
	}, nil

}
