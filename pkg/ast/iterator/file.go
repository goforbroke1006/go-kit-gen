package iterator

import (
	"go/ast"
)

func NewAstFileIterator(file *ast.File) *AstFileIterator {
	return &AstFileIterator{file: file}
}

type AstFileIterator struct {
	file *ast.File
}

func (afi AstFileIterator) GetPackageName() string {
	return afi.file.Name.Name
}

func (afi AstFileIterator) HasImport(relPath string) bool {
	for _, imp := range afi.file.Imports {
		if "\""+relPath+"\"" == imp.Path.Value {
			return true
		}
	}
	return false
}

func (afi AstFileIterator) GetFuncDecl(funcName string) *ast.FuncDecl {
	for _, d := range afi.file.Decls {
		switch d.(type) {
		case *ast.FuncDecl:
			if funcName == d.(*ast.FuncDecl).Name.Name {
				return d.(*ast.FuncDecl)
			}
		}
	}
	return nil
}

func (afi AstFileIterator) GetStructFuncDecl(funcName string, recvStructTypeName string) *ast.FuncDecl {
	for _, decl := range afi.file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if nil == decl.(*ast.FuncDecl).Recv {
				continue
			}
			if nil == decl.(*ast.FuncDecl).Recv.List[0].Type {
				continue
			}
			if recvStructTypeName != decl.(*ast.FuncDecl).Recv.List[0].Type.(*ast.Ident).Name {
				continue
			}
			if funcName == decl.(*ast.FuncDecl).Name.Name {
				return decl.(*ast.FuncDecl)
			}
		}
	}
	return nil
}

func (afi AstFileIterator) GetStructDecl(structName string) *ast.GenDecl {
	for _, d := range afi.file.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			if len(gen.Specs) > 0 {
				if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
					if "type" != gen.Tok.String() {
						continue
					}
					if structName == typeSpec.Name.Name {
						return gen
					}
				}
			}
		}
	}
	return nil
}

func (afi AstFileIterator) GetStructDeclFull(structName string) (*ast.GenDecl, []*ast.FuncDecl) {
	structDecl := afi.GetStructDecl(structName)

	var structFuncs []*ast.FuncDecl
	for _, decl := range afi.file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if nil == decl.(*ast.FuncDecl).Recv {
				continue
			}
			if structName == decl.(*ast.FuncDecl).Recv.List[0].Type.(*ast.Ident).Name {
				structFuncs = append(structFuncs, decl.(*ast.FuncDecl))
			}
		}
	}

	return structDecl, structFuncs
}

func (afi AstFileIterator) GetInterfaceTypeSpec(interfaceName string) *ast.GenDecl {
	for _, d := range afi.file.Decls {
		switch d.(type) {
		case *ast.GenDecl:
			if "type" != d.(*ast.GenDecl).Tok.String() {
				continue
			}
			if interfaceName == d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name {
				return d.(*ast.GenDecl)
			}
		}
	}
	return nil
}
