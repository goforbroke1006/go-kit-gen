package generator

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/util"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
	"github.com/pkg/errors"
	"go/ast"
	"go/token"
)

type ServiceInterfaceGenerator interface {
	CreateInterfaceIfNotExists(protoServiceName string) (name string, err error)
	CreateMethodSignatureIfNotExists(protoServiceName string, protoActionName string) error
}

type serviceInterfaceGenerator struct {
	crawler *source.FileCrawler
}

func (s serviceInterfaceGenerator) CreateInterfaceIfNotExists(protoServiceName string) (name string, err error) {
	name = protoServiceName + "Service"

	if nil != s.crawler.GetInterface(name) {
		return "", nil
	}

	intfDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						//
					},
				},
			},
		},
	}
	s.crawler.PushBack(intfDecl)

	return "", nil
}

func (s serviceInterfaceGenerator) CreateMethodSignatureIfNotExists(protoServiceName string, protoActionName string) error {
	// TODO: implement me

	name := protoServiceName + "Service"

	intfc := s.crawler.GetInterface(name)
	if nil == intfc {
		return errors.New("interface " + name + " does not exists")
	}

	if intfc.HasMethod(protoActionName) {
		return nil
	}

	intfc.AddMethod(protoActionName,
		[]*ast.Field{
			util.CreateField("ctx", "context.Context"),
			util.CreateField("arg1", "interface{}"),
		},
		[]*ast.Field{
			util.CreateField("", "interface{}"),
			util.CreateField("", "error"),
		},
	)
	return nil
}

func NewServiceInterfaceGenerator(crawler *source.FileCrawler) ServiceInterfaceGenerator {
	return &serviceInterfaceGenerator{
		crawler: crawler,
	}
}
