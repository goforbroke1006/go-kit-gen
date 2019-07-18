package source

import (
	"go/ast"
	"go/token"
)

func NewFileCrawler(file *ast.File) *FileCrawler {
	return &FileCrawler{
		file: file,
	}
}

type FileCrawler struct {
	file *ast.File
}

func (c FileCrawler) HasImport(packagePath string) bool {
	for _, d := range c.file.Decls {
		//fmt.Printf("%v", d)
		if _, ok := d.(*ast.GenDecl); !ok {
			continue
		}
		if token.IMPORT != d.(*ast.GenDecl).Tok {
			continue
		}

		for _, spec := range d.(*ast.GenDecl).Specs {
			if "\""+packagePath+"\"" == spec.(*ast.ImportSpec).Path.Value {
				return true
			}
		}
	}
	return false
}

func (c FileCrawler) AddImportIfNotExists(packagePath, alias string) {
	if c.HasImport(packagePath) {
		return
	}

	importDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					//Kind:  token.STRING,
					Value: "\"" + packagePath + "\"",
				},
			},
		},
	}
	if len(alias) > 0 {
		importDecl.Specs[0].(*ast.ImportSpec).Name = ast.NewIdent(alias)
	}

	var decls []ast.Decl
	decls = append(decls, importDecl)
	decls = append(decls, c.file.Decls...)
	c.file.Decls = decls
}

func (c *FileCrawler) SetPackageIfNotDefined(packageName string) {
	if nil == c.file.Name {
		c.file.Name = ast.NewIdent(packageName)
	}
}

func (c *FileCrawler) GetStruct(name string) *StructCrawler {
	for _, d := range c.file.Decls {
		if _, ok := d.(*ast.GenDecl); !ok {
			continue
		}

		if _, ok := d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec); !ok {
			continue
		}

		if _, ok := d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType); !ok {
			continue
		}

		if name == d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name {
			return NewStructCrawler(c.file, d.(*ast.GenDecl))
		}
	}

	return nil // TODO: implement me
}

func (c *FileCrawler) GetInterface(name string) *InterfaceCrawler {
	for _, d := range c.file.Decls {
		if _, ok := d.(*ast.GenDecl); !ok {
			continue
		}

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

func (c *FileCrawler) GetFunc(name string) *ast.FuncDecl {
	for _, d := range c.file.Decls {
		if _, ok := d.(*ast.FuncDecl); !ok {
			continue
		}
		if name == d.(*ast.FuncDecl).Name.Name {
			return d.(*ast.FuncDecl)
		}

	}
	return nil
}
