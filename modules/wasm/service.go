package wasm

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/sammyhass/web-ide/server/modules/file_server"
)

type WasmService struct {
}

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

	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
		return "", errors.New(stderr.String())
	}

	return routePath, nil

}
