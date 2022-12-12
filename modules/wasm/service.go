package wasm

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/sammyhass/web-ide/server/modules/file_server"
)

type WasmService struct{}

func NewWasmService() *WasmService {
	return &WasmService{}
}

/*
Compile takes a string of Go code and compiles it to web assembly using tinygo, returning the route to the compiled wasm file that is served from the file server.
*/
func (ws *WasmService) Compile(code string) (string, error) {

	tmpFile, err := os.CreateTemp("", "wasm-*.go")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write([]byte(code)); err != nil {
		return "", err
	}

	uniqueId := uuid.New().String()
	osPath := fmt.Sprintf("%s/%s.wasm", file_server.STATIC_DIR, uniqueId)
	routePath := fmt.Sprintf("%s/%s.wasm", file_server.CONTROLLER_ROUTE, uniqueId)

	cmd := exec.Command("tinygo", "build", "-o", osPath, "-target", "wasm", tmpFile.Name())

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	errs := []string{}

	if err := cmd.Run(); err != nil {
		for _, line := range strings.Split(stderr.String(), "\n") {

			hasFileName := strings.Contains(line, tmpFile.Name())
			if hasFileName {
				trimTill := strings.Index(line, tmpFile.Name())
				line = line[trimTill:]
				line = strings.ReplaceAll(line, tmpFile.Name(), "main.go")
				errs = append(errs, line)
			}
		}

		if len(errs) > 0 {
			return "", errors.New(strings.Join(errs, "\n"))
		}

		if err != nil {
			return "", fmt.Errorf("%v", stderr.String())
		}

	}

	return routePath, nil

}
