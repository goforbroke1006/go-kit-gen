package builder

import "go/ast"

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
