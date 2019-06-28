package generator

import (
	"go/ast"
	"strings"
)

type SourceFileBuilder struct {
	file *ast.File
}

func (sfb SourceFileBuilder) CreateFunc(name string, params map[string]string, returns map[string]string) {
	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = ast.NewIdent(name)
	funcDecl.Type = &ast.FuncType{
		Params:  declsMapToFieldList(params),
		Results: declsMapToFieldList(returns),
	}

	var returnEexprs []ast.Expr
	for range returns {
		returnEexprs = append(returnEexprs, ast.NewIdent("nil"))
	}
	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: returnEexprs,
			},
		},
	}
	funcDecl.Pos()

	// TODO:

	//cg := &ast.CommentGroup{
	//	List: []*ast.Comment{
	//		{
	//			Slash: funcDecl.Pos() + 100,
	//			Text: "// TODO: implement me",
	//		},
	//	},
	//}
	//sfb.file.Comments = append(sfb.file.Comments, cg)

	sfb.file.Decls = append(sfb.file.Decls, funcDecl)
}

func declsMapToFieldList(decls map[string]string) *ast.FieldList {
	fields := make([]*ast.Field, len(decls))

	counter := 0

	for varName, typeName := range decls {
		paramDecl := &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(varName),
			},
		}

		if len(typeName) > 0 {
			if strings.Contains(typeName, ".") {
				typeParts := strings.Split(typeName, ".")
				paramDecl.Type = &ast.SelectorExpr{
					X:   ast.NewIdent(typeParts[0]),
					Sel: ast.NewIdent(typeParts[1]),
				}
			} else {
				paramDecl.Type = ast.NewIdent(typeName)
			}
		} else {
			paramDecl.Type = &ast.InterfaceType{
				Methods:    &ast.FieldList{Opening: 1, Closing: 2},
				Interface:  0,
				Incomplete: false,
			}
		}

		//fields = append(fields, paramDecl)
		fields[counter] = paramDecl
		counter++
	}

	return &ast.FieldList{List: fields}
}
