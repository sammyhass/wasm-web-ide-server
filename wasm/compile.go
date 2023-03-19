package wasm

import (
	"errors"
	"os"
	"path"

	"github.com/sammyhass/web-ide/server/model"
)

type CompileOpts struct {
	GenWat       bool                      // whether or not to generate a wat file along with the wasm file
	BeforeDelete func(wasm *os.File) error // BeforeDelete is called before the temp directory is deleted, it is passed the compiled WASM file
}

type CompileResult struct {
	Wasm []byte
	Wat  string
}

func Compile(language model.ProjectLanguage, code string, options CompileOpts) (CompileResult, error) {
	switch language {
	case model.LanguageAssemblyScript:
		return compileAssemblyScript(code, options)
	case model.LanguageGo:
		return compileTinyGo(code, options)
	default:
		return CompileResult{}, errors.New("unknown language")
	}
}

func createTempCodeDir(fname string, code string) (string, func(), error) {
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
