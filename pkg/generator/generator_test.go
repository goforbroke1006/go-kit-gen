package generator

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"testing"
)

func TestSourceFileBuilder_CreateFunc(t *testing.T) {
	//filename := "testdata/TestSourceFileBuilder_CreateFunc.test.tmp.go"

	//fs := token.NewFileSet()
	//fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	//if nil != err {
	//	log.Fatal(err)
	//}

	sfb := SourceFileBuilder{
		file: &ast.File{
			Name: ast.NewIdent("testdata"),
		},
	}
	sfb.CreateFunc(
		"SomeFunc",
		map[string]string{
			"ctx": "context.Context",
			"req": "",
		},
		map[string]string{
			"resp": "",
			"err":  "error",
		},
	)

	fs := token.NewFileSet()
	printer.Fprint(os.Stdout, fs, sfb.file)
}
