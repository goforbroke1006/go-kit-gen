package iterator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func NewAstFileIterator(filename string) *AstFileIterator {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if nil != err {
		log.Fatal(err)
	}
	return &AstFileIterator{file: file}
}

type AstFileIterator struct {
	file *ast.File
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

func (afi AstFileIterator) GetStructFuncDecl(funcName string, recvStructType string) *ast.FuncDecl {
	for _, decl := range afi.file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if nil == decl.(*ast.FuncDecl).Recv {
				continue
			}
			if recvStructType != decl.(*ast.FuncDecl).Recv.List[0].Type.(*ast.Ident).Name {
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

func (afi AstFileIterator) GetInterfaceTypeSpec(interfaceName string) *ast.TypeSpec {
	for _, d := range afi.file.Decls {
		switch d.(type) {
		case *ast.GenDecl:
			if "type" != d.(*ast.GenDecl).Tok.String() {
				continue
			}
			if interfaceName == d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name {
				return d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
			}
		}
	}
	return nil
}

// --------------------------------------------------

func NewAstInterfaceDeclIterator(interfaceType *ast.TypeSpec) *AstInterfaceDeclIterator {
	return &AstInterfaceDeclIterator{
		interfaceType: interfaceType,
	}
}

type AstInterfaceDeclIterator struct {
	interfaceType *ast.TypeSpec
}

func (aii AstInterfaceDeclIterator) GetMethodsFieldList() *ast.FieldList {
	return aii.interfaceType.Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods
}

// --------------------------------------------------