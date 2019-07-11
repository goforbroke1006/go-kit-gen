package util

import (
	"go/ast"
	"go/token"
	"strings"
)

func CreateStructDecl(structName string, properties []*ast.Field) *ast.GenDecl {
	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
	}
	structDecl.Specs = append(
		structDecl.Specs,
		&ast.TypeSpec{
			Name: &ast.Ident{
				Name: structName,
			},
			Type: &ast.StructType{
				Fields: &ast.FieldList{List: properties},
			},
		},
	)
	return structDecl
}

func CreateField(name, typ string) *ast.Field {
	paramDecl := &ast.Field{}

	if len(name) > 0 {
		paramDecl.Names = []*ast.Ident{
			ast.NewIdent(name),
		}
	}

	paramDecl.Type = StringToAstType(typ)
	return paramDecl
}

func MapToFieldList(decls [][2]string) []*ast.Field {
	fields := make([]*ast.Field, len(decls))

	counter := 0

	for _, typeName := range decls {
		paramDecl := &ast.Field{}

		if len(typeName[0]) > 0 {
			paramDecl.Names = []*ast.Ident{
				ast.NewIdent(typeName[0]),
			}
		}

		paramDecl.Type = StringToAstType(typeName[1])

		fields[counter] = paramDecl
		counter++
	}

	return fields
}

func StringToAstType(typeName string) ast.Expr {
	if "interface{}" == typeName || 0 == len(typeName) {
		return &ast.InterfaceType{
			Methods:    &ast.FieldList{Opening: 1, Closing: 2},
			Interface:  0,
			Incomplete: false,
		}
	}

	if strings.Contains(typeName, ".") {
		typeParts := strings.Split(typeName, ".")
		return &ast.SelectorExpr{
			X:   ast.NewIdent(typeParts[0]),
			Sel: ast.NewIdent(typeParts[1]),
		}
	}

	return ast.NewIdent(typeName)

}
