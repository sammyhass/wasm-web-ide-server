package model

import (
	"errors"
	"fmt"
	"strings"
)

func languageFromFileName(fileName string) string {
	return strings.Split(fileName, ".")[1]
}

/*
ProjectFiles are represented as a map of file name to file content
*/
type ProjectFiles map[string]string

func ProjectFilesToFileViews(files ProjectFiles) []FileView {
	var fileViews []FileView
	for path, content := range files {
		fmt.Println(path, strings.Split(path, "."))
		fileViews = append(fileViews, FileView{
			Name:     path,
			Content:  content,
			Language: languageFromFileName(path),
		})
	}
	return fileViews
}

type FileView struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Language string `json:"language"`
}

var DefaultGo = `package main

import (
	"syscall/js"
)

func main() {
	js.Global().Get("alert").Invoke("Hello WASM!")
}`

var DefaultHtml = `<h1>Hello World</h1>`

var DefaultCss = `h1 {
	color: red;
}`

var DefaultJs = `console.log("Hello World")`

var goModTemplate = `module %s
go 1.19`

func DefaultGoMod(projName string) string {
	// replace spaces with dashes and all special characters with nothing
	slug := strings.Map(func(r rune) rune {
		if r == ' ' {
			return '-'
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, projName)
	return fmt.Sprintf(goModTemplate, slug)
}

var DefaultFiles = ProjectFiles{
	"main.go":    DefaultGo,
	"index.html": DefaultHtml,
	"styles.css": DefaultCss,
	"app.js":     DefaultJs,
}

func GetFileContent(files []FileView, filename string) (string, error) {
	for _, file := range files {
		if file.Name == filename {
			return file.Content, nil
		}
	}

	return "", errors.New(filename + " not found")
}
