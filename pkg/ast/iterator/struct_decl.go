package iterator

import "go/ast"

func NewAstStructDeclIterator(structDecl *ast.GenDecl) *AstStructDeclIterator {
	return &AstStructDeclIterator{
		structDecl: structDecl,
	}
}

type AstStructDeclIterator struct {
	structDecl *ast.GenDecl
}

func (asdi AstStructDeclIterator) GetName() string {
	return asdi.structDecl.Specs[0].(*ast.TypeSpec).Name.Name
}

func (asdi AstStructDeclIterator) GetProperties() []*ast.Field {
	return asdi.structDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
}

func (asdi AstStructDeclIterator) HasProperty(name string) bool {
	for _, field := range asdi.GetProperties() {
		if name == field.Names[0].Name {
			return true
		}
	}
	return false
}
