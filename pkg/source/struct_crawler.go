package source

import (
	"go/ast"
)

type StructCrawler struct {
}

func (s StructCrawler) GetName() string {
	return "" // TODO: implement me
}

func NewStructCrawler() *StructCrawler {
	return &StructCrawler{}
}

type InterfaceCrawler struct {
	intfDecl *ast.GenDecl
}

func (crawler InterfaceCrawler) HasMethod(name string) bool {
	list := crawler.intfDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
	for _, m := range list {
		if name == m.Names[0].Name {
			return true
		}
	}
	return false
}

func (crawler InterfaceCrawler) AddMethod(name string, args []*ast.Field, rets []*ast.Field) {
	// TODO: implement me

	field := &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent(name),
		},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{List: args},
			Results: &ast.FieldList{List: rets},
		},
	}

	list := crawler.intfDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
	list = append(list, field)
	crawler.intfDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List = list
}

func NewInterfaceCrawler(intfDecl *ast.GenDecl) *InterfaceCrawler {
	return &InterfaceCrawler{
		intfDecl: intfDecl,
	}
}
