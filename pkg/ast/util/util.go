package util

import (
	"go/ast"
	"strings"
)

func StringToAstType(typeName string) ast.Expr {
	if len(typeName) > 0 {
		if strings.Contains(typeName, ".") {
			typeParts := strings.Split(typeName, ".")
			return &ast.SelectorExpr{
				X:   ast.NewIdent(typeParts[0]),
				Sel: ast.NewIdent(typeParts[1]),
			}
		} else {
			return ast.NewIdent(typeName)
		}
	} else {
		return &ast.InterfaceType{
			Methods:    &ast.FieldList{Opening: 1, Closing: 2},
			Interface:  0,
			Incomplete: false,
		}
	}
}
