package fixer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func FixMissingReqRespForEndpoint(filename string, expectedMethods []string) {
	var res []string

	fs := token.NewFileSet()
	fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if nil != err {
		log.Fatal(err)
	}

	for _, m := range expectedMethods {
		findReq := false
		findRes := false
		for _, d := range fileGo.Decls {
			if gen, ok := d.(*ast.GenDecl); ok {
				if len(gen.Specs) > 0 {
					if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
						if "type" != gen.Tok.String() {
							continue
						}
						if m+"Request" == typeSpec.Name.Name {
							findReq = true
							continue
						}
						if m+"Response" == typeSpec.Name.Name {
							findRes = true
							continue
						}
					}
				}
			}
		}
		if !findReq {
			res = append(res, m+"Request")
		}
		if !findRes {
			res = append(res, m+"Response")
		}
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if nil != err {
		log.Fatal(err)
	}
	for _, typeName := range res {
		_, err := f.WriteString("type " + typeName + " struct {\n    // TODO: \n}\n")
		if nil != err {
			fmt.Println(err)
		}
	}
}

func FixMissingEndpoints(filename, serviceName string, expectedMethods []string) {
	fs := token.NewFileSet()
	fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if nil != err {
		log.Fatal(err)
	}

	var endpointsStruct *ast.GenDecl
	for _, d := range fileGo.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			if len(gen.Specs) > 0 {
				if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
					if "type" != gen.Tok.String() {
						continue
					}
					if serviceName+"Endpoints" == typeSpec.Name.Name {
						endpointsStruct = gen
						continue
					}
				}
			}
		}
	}

	if nil == endpointsStruct {
		return
	}

	fieldsList := endpointsStruct.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List

	for _, m := range expectedMethods {
		exists := false
		for _, field := range fieldsList {
			if m+"Endpoint" == field.Names[0].Name {
				exists = true
				break
			}
		}

		if !exists {
			missing := &ast.Field{
				Names: []*ast.Ident{
					{Name: m + "Endpoint"},
				},
				Type: fieldsList[0].Type,
			}
			fieldsList = append(fieldsList, missing)
		}
	}

	endpointsStruct.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List = fieldsList

	out, err := os.OpenFile(filename, os.O_RDWR, 066)
	if nil != err {
		log.Fatal(err)
	}
	if err = out.Truncate(0); nil != err {
		log.Fatal(err)
	}
	if _, err = out.Seek(0, 0); nil != err {
		log.Fatal(err)
	}
	err = printer.Fprint(out, fs, fileGo)
	if nil != err {
		log.Fatal(err)
	}

}
