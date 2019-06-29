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

func MapToFieldList(decls map[string]string) []*ast.Field {
	fields := make([]*ast.Field, len(decls))

	counter := 0

	for varName, typeName := range decls {
		paramDecl := &ast.Field{}

		if len(varName) > 0 {
			paramDecl.Names = []*ast.Ident{
				ast.NewIdent(varName),
			}
		} else {
			paramDecl.Names = []*ast.Ident{
				ast.NewIdent("_"),
			}
		}

		paramDecl.Type = StringToAstType(typeName)

		fields[counter] = paramDecl
		counter++
	}

	return fields
}
