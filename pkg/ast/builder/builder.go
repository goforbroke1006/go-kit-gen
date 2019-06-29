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
