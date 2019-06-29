package fixer

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func OpenGolangSourceFile(filename string) (*token.FileSet, *ast.File) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if nil != err {
		log.Fatal(err)
	}
	return fs, file
}

func WriteSourceFile(filename string, file *ast.File, fset *token.FileSet) {
	out, err := os.OpenFile(filename, os.O_RDWR, 066)
	if nil != err {
		log.Fatal(err)
	}
	if err = out.Truncate(0); nil != err {
		log.Fatal(err)
	}
	if _, err = out.Seek(0, 0); nil != err {
		log.Fatal(err)
	}

	//fs := token.NewFileSet()
	err = printer.Fprint(out, fset, file)
	if nil != err {
		log.Fatal(err)
	}
}
