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

	sfb := AstPrimitiveBuilder{}
	funcDecl := sfb.CreateFuncDecl(
		"MakeSomeAwesomeHubEndpoints",
		map[string]string{
			"ctx":    "context.Context",
			"svc":    "service.SomeAwesomeHub",
			"config": "",
		},
		map[string]string{
			"resp": "",
			"err":  "error",
		},
		[]interface{}{
			sfb.CreateCompositeLit(
				"SomeEndpoints",
				map[string]ast.Expr{
					"MethodOneEndpoint": sfb.CreateMethodCallExpr("makeMethodOneEndpoint", []string{
						"svc",
					}),
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
