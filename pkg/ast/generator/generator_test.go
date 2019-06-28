package generator

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"testing"
)

func TestSourceFileBuilder_CreateFunc(t *testing.T) {
	file := &ast.File{
		Name: ast.NewIdent("testdata"),
	}

	apb := AstPrimitiveBuilder{}
	funcDecl := apb.CreateFuncDecl(
		"MakeSomeAwesomeHubEndpoints",
		map[string]string{
			"ctx": "context.Context",
			"svc": "service.SomeAwesomeHubService",
		},
		map[string]string{
			"resp": "",
			"err":  "error",
		},
		[]interface{}{
			apb.CreateCompositeLiteralExpr(
				"SomeAwesomeHubEndpoints",
				map[string]ast.Expr{
					"MethodOneEndpoint": apb.CreateMethodCallExpr(
						"makeMethodOneEndpoint",
						[]string{"svc"},
					),
				},
			),
			nil,
		},
	)

	file.Decls = append(file.Decls, funcDecl)

	fs := token.NewFileSet()
	printer.Fprint(os.Stdout, fs, file)

	// TODO: add assertion check
}
