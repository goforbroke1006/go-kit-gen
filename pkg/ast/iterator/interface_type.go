package iterator

import "go/ast"

func NewAstInterfaceTypeIterator(interfaceType *ast.GenDecl) *AstInterfaceTypeIterator {
	return &AstInterfaceTypeIterator{
		infcDecl: interfaceType,
	}
}

type AstInterfaceTypeIterator struct {
	infcDecl *ast.GenDecl
}

func (aiti AstInterfaceTypeIterator) GetMethodsFieldList() *ast.FieldList {
	return aiti.infcDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods
}

func (aiti AstInterfaceTypeIterator) GetMethod(methodName string) *ast.Field {
	for _, m := range aiti.GetMethodsFieldList().List {
		if methodName == m.Names[0].Name {
			return m
		}
	}
	return nil
}
