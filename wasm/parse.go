package wasm

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
)

type exportedFunc struct {
	decl    ast.Decl
	comment *ast.Comment
}

func parseExportedComments(
	comments []*ast.CommentGroup,
	out chan<- *ast.Comment,
	done chan<- bool,
) {
	for _, group := range comments {
		for _, comment := range group.List {
			if strings.HasPrefix(comment.Text, "//export ") {
				out <- comment
			}
		}
	}

	done <- true

}

func parseDecls(
	comments <-chan *ast.Comment,
	decls []ast.Decl,
	foundAllComments <-chan bool,
	done chan<- bool,
	out chan<- exportedFunc,
) error {
	for {
		select {
		case comment := <-comments:
			for _, decl := range decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					// ensure the comment and function are one line apart

					if fn.Name.Name == strings.TrimPrefix(comment.Text, "//export ") && fn.Pos() == comment.End()+1 {
						out <- exportedFunc{
							decl:    decl,
							comment: comment,
						}
					}
				}
			}
		case <-foundAllComments:
			done <- true
			return nil
		}
	}

}

func parseExports(
	src io.Reader,
	done chan<- bool,
	out chan<- exportedFunc,
) {
	fset := token.NewFileSet()

	parsed, err := parser.ParseFile(fset, "main.go", src, parser.ParseComments)

	if err != nil {
		panic(err)
	}

	foundAllComments := make(chan bool, 1)
	foundAllExports := make(chan bool, 1)

	comments := make(chan *ast.Comment)

	go parseExportedComments(parsed.Comments, comments, foundAllComments)
	go parseDecls(comments, parsed.Decls, foundAllComments, foundAllExports, out)

	<-foundAllExports

	done <- true

}

func Parse(src io.Reader) (exports []exportedFunc, err error) {
	done := make(chan bool, 1)
	out := make(chan exportedFunc)

	go parseExports(src, done, out)

	for {
		select {
		case export := <-out:
			exports = append(exports, export)
		case <-done:
			return exports, nil
		}
	}
}
