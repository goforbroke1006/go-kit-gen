package generator

import (
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
	"github.com/goforbroke1006/go-kit-gen/pkg/string_util"
	"go/ast"
	"go/token"
)

type EndpointsStructGenerator struct {
	crawler *source.FileCrawler
}

func (g EndpointsStructGenerator) CreateEndpointStructIfNotExists(serviceName string) {
	name := string_util.FirstLetterToUpperCase(serviceName) + "Endpoints"

	if nil != g.crawler.GetStruct(name) {
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
	g.crawler.PushBack(structDecl)

	return
}

func (g EndpointsStructGenerator) CreateEndpointStructField(serviceName, actionName string) error {
	name := string_util.FirstLetterToUpperCase(serviceName) + "Endpoints"
	eStruct := g.crawler.GetStruct(name)
	if nil == eStruct {
		return fmt.Errorf("struct %s does not exist", name)
	}

	eStruct.AddProperty(actionName, "endpoint.Endpoint")

	return nil
}

func (g EndpointsStructGenerator) CreateRequestStruct(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateResponseStruct(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateMakeEndpointFunc(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateConstructorIfNotExists() {
	// TODO: implement me
}

func (g EndpointsStructGenerator) SetFieldInContructor(actionName string) {
	// TODO: implement me
}

func NewEndpointsStructGenerator(crawler *source.FileCrawler) *EndpointsStructGenerator {
	return &EndpointsStructGenerator{
		crawler: crawler,
	}
}
