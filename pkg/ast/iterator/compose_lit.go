package iterator

import "go/ast"

func NewAstCompositeLitIterator(lit *ast.CompositeLit) *AstCompositeLitIterator {
	return &AstCompositeLitIterator{
		lit: lit,
	}
}

type AstCompositeLitIterator struct {
	lit *ast.CompositeLit
}

func (acli AstCompositeLitIterator) GetElements() []ast.Expr {
	return acli.lit.Elts
}
