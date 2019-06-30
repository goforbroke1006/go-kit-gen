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

// --------------------------------------------------

func NewAstStructDeclIterator(structDecl *ast.GenDecl) *AstStructDeclIterator {
	return &AstStructDeclIterator{
		structDecl: structDecl,
	}
}

type AstStructDeclIterator struct {
	structDecl *ast.GenDecl
}

func (asdi AstStructDeclIterator) GetProperties() []*ast.Field {
	return asdi.structDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
}

// --------------------------------------------------

func NewAstInterfaceTypeIterator(interfaceType *ast.GenDecl) *AstInterfaceTypeIterator {
	return &AstInterfaceTypeIterator{
		infcDecl: interfaceType,
	}
}

type AstInterfaceTypeIterator struct {
	infcDecl *ast.GenDecl
}

func (aiti AstInterfaceTypeIterator) GetMethodsFieldList() *ast.FieldList {
	return aiti.infcDecl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods
}

func (aiti AstInterfaceTypeIterator) GetMethod(methodName string) *ast.Field {
	for _, m := range aiti.GetMethodsFieldList().List {
		if methodName == m.Names[0].Name {
			return m
		}
	}
	return nil
}

// --------------------------------------------------

func NewAstCompositeLitIterator(lit *ast.CompositeLit) *AstCompositeLitIterator {
	return &AstCompositeLitIterator{
		lit: lit,
	}
}

type AstCompositeLitIterator struct {
	lit *ast.CompositeLit
}

func (acli AstCompositeLitIterator) GetElements() []ast.Expr {
	return acli.lit.Elts
}

// --------------------------------------------------

func NewAstFuncDeclIterator(funcDecl *ast.FuncDecl) *AstFuncDeclIterator {
	return &AstFuncDeclIterator{
		funcDecl: funcDecl,
	}
}

type AstFuncDeclIterator struct {
	funcDecl *ast.FuncDecl
}

func (afdi AstFuncDeclIterator) GetReturnSmtm() *ast.ReturnStmt {
	return afdi.funcDecl.Body.List[0].(*ast.ReturnStmt)
}

// --------------------------------------------------
