package generator

import (
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/util"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
	"github.com/goforbroke1006/go-kit-gen/pkg/string_util"
	"go/ast"
	"go/token"
)

type EndpointsStructGenerator struct {
	crawler *source.FileCrawler
}

func (g EndpointsStructGenerator) CreateEndpointStructIfNotExists(serviceName string) {
	name := string_util.FirstLetterToUpperCase(serviceName) + "Endpoints"

	if nil != g.crawler.GetStruct(name) {
		return
	}

	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{Fields: &ast.FieldList{}},
			},
		},
	}
	g.crawler.PushBack(structDecl)

	return
}

func (g EndpointsStructGenerator) CreateEndpointStructField(serviceName, actionName string) error {
	name := string_util.FirstLetterToUpperCase(serviceName) + "Endpoints"
	eStruct := g.crawler.GetStruct(name)
	if nil == eStruct {
		return fmt.Errorf("struct %s does not exist", name)
	}

	if eStruct.HasProperty(actionName + "Endpoint") {
		return nil
	}

	eStruct.AddProperty(actionName+"Endpoint", "goKitEndpoint.Endpoint")

	return nil
}

func (g EndpointsStructGenerator) CreateRequestStruct(actionName string) {
	name := string_util.FirstLetterToUpperCase(actionName) + "Request"

	if nil != g.crawler.GetStruct(name) {
		return
	}

	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{Fields: &ast.FieldList{
					// TODO: specify struct fields/parameters
				}},
			},
		},
	}
	g.crawler.PushBack(structDecl)
}

func (g EndpointsStructGenerator) CreateResponseStruct(actionName string) {
	name := string_util.FirstLetterToUpperCase(actionName) + "Response"

	if nil != g.crawler.GetStruct(name) {
		return
	}

	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{Fields: &ast.FieldList{
					// TODO: specify struct fields/parameters
				}},
			},
		},
	}
	g.crawler.PushBack(structDecl)
}

func (g EndpointsStructGenerator) CreateConstructorIfNotExists(serviceName string) {
	serviceName = string_util.FirstLetterToUpperCase(serviceName)
	constructorName := "New" + serviceName + "Endpoints"
	if nil != g.crawler.GetFunc(constructorName) {
		return
	}

	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = ast.NewIdent(constructorName)

	structName := string_util.FirstLetterToUpperCase(serviceName) + "Endpoints"
	funcDecl.Type = &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("svc", "service."+serviceName+"Service"),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("", structName),
			},
		},
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CompositeLit{
						Type: util.StringToAstType(structName),
					}, // TODO: improve it
				},
			},
		},
	}

	g.crawler.PushBack(funcDecl)
}

func (g EndpointsStructGenerator) SetFieldInConstructor(serviceName, actionName string) {
	// TODO: implement me

	constructorName := "New" + serviceName + "Endpoints"
	funcDecl := g.crawler.GetFunc(constructorName)
	list := funcDecl.Body.List[len(funcDecl.Body.List)-1].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts

	for _, el := range list {
		if actionName+"Endpoint" == el.(*ast.KeyValueExpr).Key.(*ast.Ident).Name {
			return
		}
	}

	makeFuncName := "make" + string_util.FirstLetterToUpperCase(actionName) + "Endpoint"

	list = append(list, &ast.KeyValueExpr{
		Key: ast.NewIdent(actionName + "Endpoint"),
		Value: &ast.CallExpr{
			Fun: ast.NewIdent(makeFuncName),
			Args: []ast.Expr{
				ast.NewIdent("svc"),
			},
		},
	})

	funcDecl.Body.List[len(funcDecl.Body.List)-1].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts = list
}

func (g EndpointsStructGenerator) CreateMakeEndpointFunc(serviceName, actionName string) {
	name := "make" + string_util.FirstLetterToUpperCase(actionName) + "Endpoint"

	if nil != g.crawler.GetFunc(name) {
		return
	}

	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = &ast.Ident{
		Name:    name,
		NamePos: token.NoPos,
	}

	funcDecl.Type = &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("svc", "service."+serviceName+"Service"),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("", "goKitEndpoint.Endpoint"),
			},
		},
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("nil"), // TODO: improve it
				},
			},
		},
	}

	g.crawler.PushBack(funcDecl)
}

func NewEndpointsStructGenerator(crawler *source.FileCrawler) *EndpointsStructGenerator {
	return &EndpointsStructGenerator{
		crawler: crawler,
	}
}
