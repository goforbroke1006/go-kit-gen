package factory

import (
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"strings"
	"testing"
)

func TestSourceFileBuilder_CreateStructDecl(t *testing.T) {
	file := &ast.File{
		Name: ast.NewIdent("testdata"),
	}

	apb := AstPrimitiveFactory{}
	structDecl := apb.CreateStructDecl(
		"SomeAwesomeHubEndpoints",
		map[string]string{
			"MethodOneEndpoint": "endpoint.Endpoint",
			"MethodTwoEndpoint": "endpoint.Endpoint",
		},
	)

	file.Decls = append(file.Decls, structDecl)

	fs := token.NewFileSet()
	sb := &strings.Builder{}
	printer.Fprint(sb, fs, file)

	log.Println(sb.String())

	if !strings.Contains(sb.String(), "type SomeAwesomeHubEndpoints struct") {
		t.Errorf("AstPrimitiveFactory.CreateStructDecl() works wrong, can't find struct declaration")
	}

	if !strings.Contains(sb.String(), "MethodOneEndpoint	endpoint.Endpoint") {
		t.Errorf("AstPrimitiveFactory.CreateStructDecl() works wrong, can't find struct's prop declaration")
	}

	// TODO: add assertion check
}

func TestSourceFileBuilder_CreateFunc(t *testing.T) {
	file := &ast.File{
		Name: ast.NewIdent("testdata"),
	}

	apb := AstPrimitiveFactory{}
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
					"MethodOneEndpoint": apb.CreateMethodCallExpr("makeMethodOneEndpoint", []string{"svc"}),
					"MethodTwoEndpoint": apb.CreateMethodCallExpr("makeMethodTwoEndpoint", []string{"svc"}),
				},
			),
			nil,
		},
	)

	file.Decls = append(file.Decls, funcDecl)

	fs := token.NewFileSet()
	sb := &strings.Builder{}

	ast.Print(fs, funcDecl)

	printer.Fprint(sb, fs, file)

	log.Println(sb.String())

	if !strings.Contains(sb.String(), "func MakeSomeAwesomeHubEndpoints") {
		t.Errorf("AstPrimitiveFactory.CreateFuncDecl() works wrong, can't find func declaration")
	}

	if !strings.Contains(sb.String(), "MethodOneEndpoint: makeMethodOneEndpoint(svc)") {
		t.Errorf("AstPrimitiveFactory.CreateCompositeLiteralExpr() works wrong, funct return literal")
	}
}
