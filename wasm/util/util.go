package util

import (
	"os"
	"path"
)

type CompileResult struct {
	Wasm []byte
	Wat  string
}

func CreateTempCodeDir(fname string, code string) (string, func(), error) {
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

	codeFile, err := createInTemp(fname)
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
