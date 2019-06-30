package factory

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/util"
	"go/ast"
	"go/token"
)

type AstPrimitiveFactory struct {
}

func (apb AstPrimitiveFactory) CreateStructDecl(
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
				Fields: &ast.FieldList{List: util.MapToFieldList(properties)},
			},
		},
	)
	return structDecl
}

func (apb AstPrimitiveFactory) CreateFuncDecl(
	name string,
	params map[string]string,
	returns map[string]string,
	returnStmtVals []ast.Expr,
	receiverName *string,
	receiverTypeName *string,
) *ast.FuncDecl {
	if nil != returnStmtVals && len(returnStmtVals) != len(returns) {
		panic("return expr list must have same size like return declaration list")
	}

	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = &ast.Ident{
		Name:    name,
		NamePos: token.NoPos,
	}

	if nil != receiverTypeName {
		if nil == receiverName || len(*receiverName) == 0 {
			receiverName = new(string)
			*receiverName = "self"
		}

		funcDecl.Recv = &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(*receiverName)},
					Type:  ast.NewIdent(*receiverTypeName),
				},
			},
		}
	}

	funcDecl.Type = &ast.FuncType{
		Params:  &ast.FieldList{List: util.MapToFieldList(params)},
		Results: &ast.FieldList{List: util.MapToFieldList(returns)},
	}

	if len(returnStmtVals) > 0 {
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
	}

	return funcDecl
}

func (apb AstPrimitiveFactory) CreateFuncSignatureExpr(
	name string,
	params map[string]string,
	returns map[string]string,
) *ast.Field {
	funcDecl := &ast.Field{}
	funcDecl.Names = []*ast.Ident{
		{
			Name:    name,
			NamePos: token.NoPos,
		},
	}

	funcDecl.Type = &ast.FuncType{
		Params:  &ast.FieldList{List: util.MapToFieldList(params)},
		Results: &ast.FieldList{List: util.MapToFieldList(returns)},
	}

	return funcDecl
}

func (apb AstPrimitiveFactory) CreateAnonFuncObjectDecl(
	params map[string]string,
	returns map[string]string,
	returnStmtVals []ast.Expr,
) *ast.FuncLit {
	if nil != returnStmtVals && len(returnStmtVals) != len(returns) {
		panic("return expr list must have same size like return declaration list")
	}

	funcLit := &ast.FuncLit{}

	funcLit.Type = &ast.FuncType{
		Params:  &ast.FieldList{List: util.MapToFieldList(params)},
		Results: &ast.FieldList{List: util.MapToFieldList(returns)},
	}

	var returnEexprs []ast.Expr
	for _, resVal := range returnStmtVals {
		if nil == resVal {
			returnEexprs = append(returnEexprs, ast.NewIdent("nil"))
		} else {
			returnEexprs = append(returnEexprs, resVal.(ast.Expr))
		}
	}

	funcLit.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: returnEexprs,
			},
		},
	}

	return funcLit
}

func (apb AstPrimitiveFactory) CreateCompositeLiteralExpr(
	structName string,
	namesToValues map[string]ast.Expr,
) *ast.CompositeLit {
	ident := &ast.CompositeLit{}
	ident.Type = ast.NewIdent(structName)

	for name, value := range namesToValues {
		keyIdent := ast.NewIdent(name)
		keyIdent.NamePos = 0
		initRow := &ast.KeyValueExpr{
			Key:   keyIdent,
			Value: value,
		}
		ident.Elts = append(ident.Elts, initRow)
	}

	return ident
}

func (apb AstPrimitiveFactory) CreateMethodCallExpr(
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
