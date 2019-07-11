package source

import (
	"go/ast"
)

func NewFileCrawler(file *ast.File) *FileCrawler {
	return &FileCrawler{
		file: file,
	}
}

type FileCrawler struct {
	file *ast.File
}

func (c *FileCrawler) SetPackageIfNotDefined(packageName string) {
	if nil == c.file.Name {
		c.file.Name = ast.NewIdent(packageName)
	}
}

func (c *FileCrawler) GetStruct(name string) *StructCrawler {
	return nil // TODO: implement me
}

func (c *FileCrawler) GetInterface(name string) *InterfaceCrawler {
	for _, d := range c.file.Decls {
		if _, ok := d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec); !ok {
			continue
		}

		if _, ok := d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType); !ok {
			continue
		}

		if name == d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name {
			return NewInterfaceCrawler(d.(*ast.GenDecl))
		}
	}
	return nil // TODO: implement me
}

func (c *FileCrawler) PushBack(decl ast.Decl) {
	// TODO: implement me
	c.file.Decls = append(c.file.Decls, decl)
}
