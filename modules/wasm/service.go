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

type Service struct{}

func NewWasmService() *Service {
	return &Service{}
}

func createTempCodeDir(code string, goMod string) (string, error, func()) {
	tmpDir, err := os.MkdirTemp("", "project-dir-*")
	if err != nil {
		return "", err, nil
	}

	deleteDir := func() {
		os.RemoveAll(tmpDir)
	}

	codeFile, err := os.CreateTemp(tmpDir, "main.go")
	if err != nil {
		deleteDir()
		return "", err, nil
	}

	if _, err := codeFile.Write([]byte(code)); err != nil {
		deleteDir()
		return "", err, nil
	}

	goModFile, err := os.CreateTemp(tmpDir, "go.mod")
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

/*
Compile takes a string of Go code and a string containing a go.mod file and compiles it to web assembly using tinygo, returning the route to the compiled wasm file that is served from the file server.
*/
func (ws *Service) Compile(code string, goMod string) (string, error) {

	dir, err, delete := createTempCodeDir(code, goMod)
	if err != nil {
		return "", err
	}
	defer delete()

	uniqueId := uuid.New().String()
	osPath := fmt.Sprintf("%s/%s.wasm", file_server.STATIC_DIR, uniqueId)
	routePath := fmt.Sprintf("%s/%s.wasm", file_server.CONTROLLER_ROUTE, uniqueId)

	filename := "main.go"

	cmd := exec.Command("tinygo", "build", "-o", osPath, "-target", "wasm", filename)
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
			return "", errors.New(strings.Join(errs, "\n"))
		}

		if err != nil {
			return "", fmt.Errorf("%v", stderr.String())
		}

	}

	return routePath, nil

}
