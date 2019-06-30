package builder

import (
	"go/ast"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/util"
)

func NewAstStructBuilder(structDecl *ast.GenDecl) *AstStructBuilder {
	return &AstStructBuilder{
		structDecl: structDecl,
	}
}

type AstStructBuilder struct {
	structDecl *ast.GenDecl
}

func (asb AstStructBuilder) Get() *ast.GenDecl {
	return asb.structDecl
}

func (asb AstStructBuilder) AddProperty(name, typeName string) {
	propertyNew := &ast.Field{
		Names: []*ast.Ident{
			{Name: name},
		},
		Type: util.StringToAstType(typeName),
	}
	propsList := asb.structDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
	propsList = append(propsList, propertyNew)
	asb.structDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List = propsList
}

// --------------------------------------------------

func NewAstCompositeLitBuilder(lit *ast.CompositeLit) *AstCompositeLitBuilder {
	return &AstCompositeLitBuilder{
		lit: lit,
	}
}

type AstCompositeLitBuilder struct {
	lit *ast.CompositeLit
}

func (aclb AstCompositeLitBuilder) AddElement(propName string, propValue ast.Expr) {
	expr := &ast.KeyValueExpr{
		Key:   ast.NewIdent(propName),
		Value: propValue,
	}
	aclb.lit.Elts = append(aclb.lit.Elts, expr)
}

// --------------------------------------------------

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
