package builder

import (
	"go/ast"
	"go/token"
)

func NewAstFileBuilder(file *ast.File) *AstFileBuilder {
	return &AstFileBuilder{
		file: file,
	}
}

type AstFileBuilder struct {
	file *ast.File
}

func (afb *AstFileBuilder) SetPackage(packageName string) {
	afb.file.Name = ast.NewIdent(packageName)
}

func (afb *AstFileBuilder) AddImport(path string) {
	//importStr := &ast.ImportSpec{
	//	Path: &ast.BasicLit{
	//		Kind:     token.STRING,
	//		Value:    "\"" + path + "\"",
	//		ValuePos: 27,
	//	},
	//}
	//afb.file.Imports = append(afb.file.Imports, importStr)

	importDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"" + path + "\"",
				},
			},
		},
	}
	var decls []ast.Decl // TODO: optimize alloc
	decls = append(decls, importDecl)
	decls = append(decls, afb.file.Decls...)
	afb.file.Decls = decls
}
