package builder

import "go/ast"

func NewAstInterfaceBuilder(infcDecl *ast.GenDecl) *AstInterfaceBuilder {
	return &AstInterfaceBuilder{
		infcDecl: infcDecl,
	}
}

type AstInterfaceBuilder struct {
	infcDecl *ast.GenDecl
}

func (aib AstInterfaceBuilder) AddFuncSignature(funcDecl *ast.Field) {
	list := aib.infcDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
	list = append(list, funcDecl)
	aib.infcDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List = list
}
