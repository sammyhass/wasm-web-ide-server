package model

import "strings"

type FileLanguage int

const (
	GO = iota
	HTML
	CSS
	JS
)

func languageFromFileName(fileName string) string {
	return strings.Split(fileName, ".")[1]
}

/*
Files are represented as a map of file name to file content
*/
type ProjectFiles map[string]string

func ProjectFilesToFileViews(files ProjectFiles) []FileView {
	var fileViews []FileView
	for path, content := range files {
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

var DefaultFiles = ProjectFiles{
	"main.go":    DefaultGo,
	"index.html": DefaultHtml,
	"styles.css": DefaultCss,
	"app.js":     DefaultJs,
}
