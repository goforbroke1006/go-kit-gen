package generator

import (
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/util"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
	"github.com/goforbroke1006/go-kit-gen/pkg/string_util"
	"go/ast"
	"go/token"
)

type TransportGRPCGenerator struct {
	crawler *source.FileCrawler
}

func (g TransportGRPCGenerator) CreateServerImpleStructIfNotExists(serviceName string) {
	structName := string_util.FirstLetterToLowerCase(serviceName) + "GRPCServer"
	if nil != g.crawler.GetStruct(structName) {
		return
	}
	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(structName),
				Type: &ast.StructType{Fields: &ast.FieldList{}},
			},
		},
	}
	g.crawler.PushBack(structDecl)
}

func (g TransportGRPCGenerator) CreateField(serviceName, actionName string) error {
	structName := string_util.FirstLetterToLowerCase(serviceName) + "GRPCServer"
	eStruct := g.crawler.GetStruct(structName)
	if nil == eStruct {
		return fmt.Errorf("struct %s does not exist", structName)
	}

	fieldName := string_util.FirstLetterToLowerCase(actionName)
	if eStruct.HasProperty(fieldName) {
		return nil
	}

	eStruct.AddProperty(fieldName, "grpc.Handler")

	return nil
}

func (g TransportGRPCGenerator) CreateMethodDecl(serviceName, pbGoPackage, actionName string) {
	structName := string_util.FirstLetterToLowerCase(serviceName) + "GRPCServer"

	serverStruct := g.crawler.GetStruct(structName)

	if nil == serverStruct {
		return
	}

	fieldName := string_util.FirstLetterToLowerCase(actionName)
	methodName := string_util.FirstLetterToUpperCase(actionName)

	if serverStruct.HasMethod(methodName) {
		return
	}

	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = ast.NewIdent(methodName)

	funcDecl.Recv = &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("s")},
				Type:  ast.NewIdent(structName),
			},
		},
	}

	funcDecl.Type = &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("ctx", "context.Context"),
				util.CreateField("req", "*"+pbGoPackage+"."+methodName+"Request"),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("", "*"+pbGoPackage+"."+methodName+"Response"),
				util.CreateField("", "error"),
			},
		},
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("_"),
					ast.NewIdent("resp"),
					ast.NewIdent("err"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("s"),
								Sel: ast.NewIdent(fieldName),
							},
							Sel: ast.NewIdent("ServeGRPC"),
						},
						Args: []ast.Expr{
							ast.NewIdent("ctx"),
							ast.NewIdent("req"),
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  ast.NewIdent("err"),
					Op: token.NEQ,
					Y:  ast.NewIdent("nil"),
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("nil"),
								ast.NewIdent("err"),
							},
						},
					},
				},
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.TypeAssertExpr{
						X: ast.NewIdent("resp"),
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(pbGoPackage),
								Sel: ast.NewIdent(methodName + "Response"),
							},
						},
					},
					ast.NewIdent("nil"), // TODO: improve it
				},
			},
		},
	}

	g.crawler.PushBack(funcDecl)
}

func (g TransportGRPCGenerator) CreateServerConstructorIfNotExists(serviceName, pbGoPackage string) {
	// TODO: implement me

	constructorName := "New" + string_util.FirstLetterToUpperCase(serviceName) + "GRPCServer"
	if nil != g.crawler.GetFunc(constructorName) {
		return
	}

	srvInterfaceName := string_util.FirstLetterToUpperCase(serviceName) + "Server"
	structName := string_util.FirstLetterToLowerCase(serviceName) + "GRPCServer"

	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = ast.NewIdent(constructorName)
	funcDecl.Type = &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("ctx", "context.Context"),
				util.CreateField("endpoint", "endpoint."+serviceName+"Endpoints"),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("", pbGoPackage+"."+srvInterfaceName),
			},
		},
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CompositeLit{
						Type: util.StringToAstType(structName),
						Elts: []ast.Expr{
							// TODO: improve it
						},
					},
				},
			},
		},
	}

	g.crawler.PushBack(funcDecl)
}

func (g TransportGRPCGenerator) AddFieldInitInConstructor(serviceName, actionName string) error {
	constructorName := "New" + string_util.FirstLetterToUpperCase(serviceName) + "GRPCServer"
	funcDecl := g.crawler.GetFunc(constructorName)
	if nil == funcDecl {
		return fmt.Errorf("method '%s' does not exist", constructorName)
	}

	list := funcDecl.Body.List[len(funcDecl.Body.List)-1].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts

	fieldName := string_util.FirstLetterToLowerCase(actionName)

	for _, el := range list {
		if fieldName == el.(*ast.KeyValueExpr).Key.(*ast.Ident).Name {
			return nil
		}
	}

	// TODO: implement me

	assignment := &ast.KeyValueExpr{
		Key: ast.NewIdent(fieldName),
		Value: &ast.CallExpr{
			Fun: ast.NewIdent("grpc.NewServer"),
			Args: []ast.Expr{
				ast.NewIdent("endpoint." + actionName + "Endpoint"),
				ast.NewIdent("decode" + actionName + "Request"),
				ast.NewIdent("encode" + actionName + "Response"),
			},
		},
	}
	list = append(list, assignment)

	funcDecl.Body.List[len(funcDecl.Body.List)-1].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts = list

	return nil
}

func (g TransportGRPCGenerator) CreateDecodeRequestMethod(pbGoPackage, actionName string) {
	funcName := "decode" + string_util.FirstLetterToUpperCase(actionName) + "Request"
	funcDecl := g.crawler.GetFunc(funcName)
	if nil != funcDecl {
		return
	}

	funcDecl = &ast.FuncDecl{}
	funcDecl.Name = ast.NewIdent(funcName)
	funcDecl.Type = &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("ctx", "context.Context"),
				util.CreateField("r", "interface{}"),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("", "interface{}"),
				util.CreateField("", "error"),
			},
		},
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("req"),
				},
				Rhs: []ast.Expr{
					&ast.TypeAssertExpr{
						X: ast.NewIdent("r"),
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(pbGoPackage),
								Sel: ast.NewIdent(actionName + "Request"),
							},
						},
					},
				},
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CompositeLit{
						Type: util.StringToAstType("endpoint." + actionName + "Request"),
						Elts: []ast.Expr{},
					},
					ast.NewIdent("nil"),
				},
			},
		},
	}

	g.crawler.PushBack(funcDecl)
}

func (g TransportGRPCGenerator) CreateEncodeResponseMethod(pbGoPackage, actionName string) {
	funcName := "encode" + string_util.FirstLetterToUpperCase(actionName) + "Response"
	funcDecl := g.crawler.GetFunc(funcName)
	if nil != funcDecl {
		return
	}

	funcDecl = &ast.FuncDecl{}
	funcDecl.Name = ast.NewIdent(funcName)
	funcDecl.Type = &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("ctx", "context.Context"),
				util.CreateField("r", "interface{}"),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("", "interface{}"),
				util.CreateField("", "error"),
			},
		},
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{
					ast.NewIdent("resp"),
				},
				Rhs: []ast.Expr{
					&ast.TypeAssertExpr{
						X: ast.NewIdent("r"),
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("endpoint"),
								Sel: ast.NewIdent(actionName + "Response"),
							},
						},
					},
				},
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CompositeLit{
						Type: util.StringToAstType(pbGoPackage + "." + actionName + "Response"),
						Elts: []ast.Expr{},
					},
					ast.NewIdent("nil"),
				},
			},
		},
	}

	g.crawler.PushBack(funcDecl)
}

func NewTransportGRPCGenerator(crawler *source.FileCrawler) *TransportGRPCGenerator {
	return &TransportGRPCGenerator{
		crawler: crawler,
	}
}
