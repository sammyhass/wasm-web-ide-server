package wasm

import (
	"io"
	"os"
	"os/exec"
)

// Convert a WASM reader to a WebAssembly Text Format (WAT) string
func WasmToWat(wasmReader io.Reader) (string, error) {
	f, err := os.CreateTemp("", "wasm2wat-*.wasm")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, wasmReader)
	if err != nil {
		return "", err
	}

	watF, err := os.CreateTemp("", "wasm2wat-*.wat")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("wasm2wat", "--enable-all", f.Name(), "-o", watF.Name())

	if err := cmd.Run(); err != nil {
		return "", err
	}

	watBytes, err := io.ReadAll(watF)
	if err != nil {
		return "", err
	}

	return string(watBytes), nil
}
