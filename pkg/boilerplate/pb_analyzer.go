package boilerplate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

func ExtractServerInterface(filename, serviceName string) *ast.TypeSpec {
	fs := token.NewFileSet()
	fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if nil != err {
		log.Fatal(err)
	}

	for _, d := range fileGo.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			if len(gen.Specs) > 0 {
				if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
					if "type" != gen.Tok.String() {
						continue
					}
					if strings.HasPrefix(typeSpec.Name.Name, "Unimplemented") {
						continue
					}
					if typeSpec.Name.Name == serviceName+"Server" {
						return typeSpec
					}
				}
			}
			continue
		}
	}

	return nil
}
