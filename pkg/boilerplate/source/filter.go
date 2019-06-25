package source

import "go/ast"

func FindFuncDeclByName(file *ast.File, name string) *ast.FuncDecl {
	for _, d := range file.Decls {
		switch d.(type) {
		case *ast.FuncDecl:
			if name == d.(*ast.FuncDecl).Name.Name {
				return d.(*ast.FuncDecl)
			}
		}
	}
	return nil
}

func FindStructDeclByName(file *ast.File, name string) *ast.GenDecl {
	for _, d := range file.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			if len(gen.Specs) > 0 {
				if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
					if "type" != gen.Tok.String() {
						continue
					}
					if name == typeSpec.Name.Name {
						return gen
					}
				}
			}
		}
	}
	return nil
}

func FindStructMethods(file *ast.File, structName string) []*ast.FuncDecl {
	var result []*ast.FuncDecl
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			if nil == decl.(*ast.FuncDecl).Recv {
				continue
			}
			if structName == decl.(*ast.FuncDecl).Recv.List[0].Type.(*ast.Ident).Name {
				result = append(result, decl.(*ast.FuncDecl))
			}
		}
	}
	return result
}

func FindInterfaceByName(file *ast.File, name string) *ast.GenDecl {
	for _, d := range file.Decls {
		switch d.(type) {
		case *ast.GenDecl:
			if "type" != d.(*ast.GenDecl).Tok.String() {
				continue
			}
			if name == d.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.Name {
				return d.(*ast.GenDecl)
			}
		}
	}
	return nil
}
