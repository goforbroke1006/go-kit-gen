package source

import (
	"go/ast"
)

type StructCrawler struct {
	file *ast.File
	decl *ast.GenDecl
}

func (s StructCrawler) GetName() string {
	return s.decl.Specs[0].(*ast.TypeSpec).Name.Name
}

func (s StructCrawler) HasMethod(name string) bool {
	for _, d := range s.file.Decls {
		if _, ok := d.(*ast.FuncDecl); !ok {
			continue
		}
		if name == d.(*ast.FuncDecl).Name.Name && s.GetName() == d.(*ast.FuncDecl).Recv.List[0].Type.(*ast.Ident).Name {
			return true
		}
	}
	return false
}

func NewStructCrawler(file *ast.File, decl *ast.GenDecl) *StructCrawler {
	return &StructCrawler{
		file: file,
		decl: decl,
	}
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
