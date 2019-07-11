package generator

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/util"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
	"github.com/goforbroke1006/go-kit-gen/pkg/string_util"
	"go/ast"
	"go/token"
)

type ServicePrivateImplGenerator interface {
	CreateEmptyStructIfNotExists(pbGoServiceName string) (name string, err error)
	CreateMethodDeclIfNotExists(pbGoServiceName, pbGoActionName string)
}

type servicePrivateImplGenerator struct {
	crawler *source.FileCrawler
}

func (s servicePrivateImplGenerator) CreateEmptyStructIfNotExists(pbGoServiceName string) (name string, err error) {
	name = string_util.FirstLetterToLowerCase(pbGoServiceName) + "Service"

	if nil != s.crawler.GetStruct(name) {
		return
	}

	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.StructType{Fields: &ast.FieldList{}},
			},
		},
	}
	s.crawler.PushBack(structDecl)

	return
}

func (s servicePrivateImplGenerator) CreateMethodDeclIfNotExists(pbGoServiceName, pbGoActionName string) {
	name := string_util.FirstLetterToLowerCase(pbGoServiceName) + "Service"

	funcDecl := &ast.FuncDecl{}
	funcDecl.Name = &ast.Ident{
		Name:    pbGoActionName,
		NamePos: token.NoPos,
	}

	funcDecl.Recv = &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("svc")},
				Type:  ast.NewIdent(name),
			},
		},
	}

	funcDecl.Type = &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("ctx", "context.Context"),
				util.CreateField("arg1", "interface{}"),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				util.CreateField("", "interface{}"),
				util.CreateField("", "error"),
			},
		},
	}

	funcDecl.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("nil"), // TODO: improve it
					ast.NewIdent("nil"), // TODO: improve it
				},
			},
		},
	}

	s.crawler.PushBack(funcDecl)
}

func NewServicePrivateImplGenerator(crawler *source.FileCrawler) ServicePrivateImplGenerator {
	return &servicePrivateImplGenerator{
		crawler: crawler,
	}
}
