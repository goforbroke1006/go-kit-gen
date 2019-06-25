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
