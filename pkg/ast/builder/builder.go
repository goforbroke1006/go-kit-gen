package builder

import (
	"go/ast"
	"go/token"
	"strings"
)

type AstPrimitiveBuilder struct {
}

func (apb AstPrimitiveBuilder) CreateStructDecl(
	structName string,
	properties map[string]string,
) *ast.GenDecl {
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
				Fields: &ast.FieldList{
					List: declsMapToFieldList(properties),
				},
			},
		},
	)
	return structDecl
}

func (apb AstPrimitiveBuilder) CreateFuncDecl(
	name string,
	params map[string]string,
	returns map[string]string,
	returnStmtVals []interface{},
) *ast.FuncDecl {
	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = ast.NewIdent(name)

	funcDecl.Type = &ast.FuncType{
		Params:  &ast.FieldList{List: declsMapToFieldList(params)},
		Results: &ast.FieldList{List: declsMapToFieldList(returns)},
	}

	var returnEexprs []ast.Expr
	for _, resVal := range returnStmtVals {
		if nil == resVal {
			returnEexprs = append(returnEexprs, ast.NewIdent("nil"))
		} else {
			returnEexprs = append(returnEexprs, resVal.(ast.Expr))
		}
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
	//apb.file.Comments = append(apb.file.Comments, cg)

	return funcDecl
}

func (apb AstPrimitiveBuilder) CreateCompositeLiteralExpr(
	structName string,
	namesToValues map[string]ast.Expr,
) *ast.CompositeLit {
	ident := &ast.CompositeLit{}
	ident.Type = ast.NewIdent(structName)
	ident.Elts = mapToStructFieldsInitialization(namesToValues)
	return ident
}

func (apb AstPrimitiveBuilder) CreateMethodCallExpr(
	funcName string,
	args []string,
) *ast.CallExpr {

	var exprs []ast.Expr
	for _, argName := range args {
		exprs = append(exprs, ast.NewIdent(argName))
	}

	return &ast.CallExpr{
		Fun:  ast.NewIdent(funcName),
		Args: exprs,
	}
}

func declsMapToFieldList(decls map[string]string) []*ast.Field {
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

		fields[counter] = paramDecl
		counter++
	}

	return fields
}

func mapToStructFieldsInitialization(m map[string]ast.Expr) []ast.Expr {
	var result []ast.Expr
	for name, value := range m {
		initRow := &ast.KeyValueExpr{
			Key:   ast.NewIdent(name),
			Value: value,
		}
		result = append(result, initRow)
	}
	return result
}
