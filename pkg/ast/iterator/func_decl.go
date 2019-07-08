package iterator

import "go/ast"

func NewAstFuncDeclIterator(funcDecl *ast.FuncDecl) *AstFuncDeclIterator {
	return &AstFuncDeclIterator{
		funcDecl: funcDecl,
	}
}

type AstFuncDeclIterator struct {
	funcDecl *ast.FuncDecl
}

func (afdi AstFuncDeclIterator) GetReturnSmtm() *ast.ReturnStmt {
	return afdi.funcDecl.Body.List[0].(*ast.ReturnStmt)
}
