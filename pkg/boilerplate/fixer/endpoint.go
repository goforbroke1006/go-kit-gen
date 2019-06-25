package fixer

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"text/template"

	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/source"
)

func NewEndpointFixer(filename, serviceName string, serviceActions []string) *EndpointFixer {
	return &EndpointFixer{
		filename:        filename,
		serviceName:     serviceName,
		serviceActions:  serviceActions,
		templatesRelDir: "./",
	}
}

type EndpointFixer struct {
	filename        string
	serviceName     string
	serviceActions  []string
	templatesRelDir string
}

func (ef EndpointFixer) Fix() {
	ef.addMissedRequestModels()
	ef.addMissedResponseModels()
	ef.addMissedPropertiesInEndpointsStruct()
	ef.addMissedPropertyInitializationInMakeEndpointsFunc()
	ef.addMissedEndpointBuilderFunc()
}

func (ef EndpointFixer) addMissedRequestModels() {
	// TODO:
}

func (ef EndpointFixer) addMissedResponseModels() {
	// TODO:
}

func (ef EndpointFixer) addMissedPropertiesInEndpointsStruct() {
	fs, file := openGolangSourceFile(ef.filename)

	endpointsStructure := source.FindStructDeclByName(file, naming.GetEndpointsStructName(ef.serviceName))

	propBuilder := source.NewStructBuilder(endpointsStructure)

	properties := propBuilder.GetStructProperties()
	for _, action := range ef.serviceActions {
		found := false
		for _, prop := range properties {
			if prop.Names[0].Name == naming.GetEndpointFieldName(action) {
				found = true
				break
			}
		}

		if !found {
			missing := &ast.Field{
				Names: []*ast.Ident{
					{Name: naming.GetEndpointFieldName(action)},
				},
				Type: properties[0].Type,
			}
			propBuilder.StructAddProperty(missing)
		}
	}

	writeSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedPropertyInitializationInMakeEndpointsFunc() {
	fs, file := openGolangSourceFile(ef.filename)
	funcDecl := source.FindFuncDeclByName(file, naming.GetMakeEndpointsFuncName(ef.serviceName))
	if nil == funcDecl {
		log.Println("Can't find make endpoints func")
		return
	}

	var missedActionsList []string

	{
		for _, action := range ef.serviceActions {
			found := false
			for _, initVal := range funcDecl.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts {
				if initVal.(*ast.KeyValueExpr).Key.(*ast.Ident).Name == naming.GetEndpointFieldName(action) {
					found = true
					break
				}
			}
			if !found {
				missedActionsList = append(missedActionsList, action)
			}
		}
	}

	{
		exprs := funcDecl.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts
		for _, missedAction := range missedActionsList {
			expr := &ast.KeyValueExpr{
				Key: ast.NewIdent(naming.GetEndpointFieldName(missedAction)),
				//Colon: exprs[0].(*ast.KeyValueExpr).Colon,
				Value: &ast.CallExpr{
					Fun: ast.NewIdent(naming.GetEndpointBuilderFuncName(missedAction)),
					Args: []ast.Expr{
						ast.NewIdent("svc"), // FIXME: hardcode
					},
				},
			}
			exprs = append(exprs, expr)
		}
		funcDecl.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts = exprs
	}

	writeSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedEndpointBuilderFunc() {
	fs, file := openGolangSourceFile(ef.filename)
	for _, act := range ef.serviceActions {
		builderFuncName := naming.GetEndpointBuilderFuncName(act)
		builderFuncDecl := source.FindFuncDeclByName(file, builderFuncName)
		if nil == builderFuncDecl {

			transportTmpl := template.Must(template.ParseFiles(ef.templatesRelDir + "template/endpoint/endpoint-builder-func.tmpl"))
			//buffer := &bytes.Buffer{}
			filename := act + ".tmp"
			fileTmp, err := os.Create(filename)
			if nil != err {
				fmt.Println(err.Error())
			}
			err = transportTmpl.Execute(fileTmp, struct{ ActionName string }{ActionName: act})
			if nil != err {
				fmt.Println(err.Error())
			}

			_, fileAstTmp := openGolangSourceFile(filename)
			funcDecl := source.FindFuncDeclByName(fileAstTmp, builderFuncName)

			file.Decls = append(file.Decls, funcDecl)
			// TODO: move with comments!!!
		}
	}
	writeSourceFile(ef.filename, file, fs)
}

//func FixMissingReqRespForEndpoint(filename string, expectedMethods []string) {
//	var res []string
//
//	fs := token.NewFileSet()
//	fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
//	if nil != err {
//		log.Fatal(err)
//	}
//
//	for _, m := range expectedMethods {
//		findReq := false
//		findRes := false
//		for _, d := range fileGo.Decls {
//			if gen, ok := d.(*ast.GenDecl); ok {
//				if len(gen.Specs) > 0 {
//					if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
//						if "type" != gen.Tok.String() {
//							continue
//						}
//						if m+"Request" == typeSpec.Name.Name {
//							findReq = true
//							continue
//						}
//						if m+"Response" == typeSpec.Name.Name {
//							findRes = true
//							continue
//						}
//					}
//				}
//			}
//		}
//		if !findReq {
//			res = append(res, m+"Request")
//		}
//		if !findRes {
//			res = append(res, m+"Response")
//		}
//	}
//
//	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
//	if nil != err {
//		log.Fatal(err)
//	}
//	for _, typeName := range res {
//		_, err := f.WriteString("type " + typeName + " struct {\n    // TODO: \n}\n")
//		if nil != err {
//			fmt.Println(err)
//		}
//	}
//}
//
//func FixMissingEndpoints(filename, serviceName string, expectedMethods []string) {
//	fs, fileGo := openGolangSourceFile(filename)
//
//	endpointsStruct := extractEndpointsStruct(fileGo, serviceName)
//	if nil != endpointsStruct {
//		fillEndpointsStructWithMissedProperties(endpointsStruct, expectedMethods)
//	} else {
//		log.Println("WARN", "Can't find endpoints struct!")
//	}
//
//	method := extractEndpointMakeMethod(fileGo, serviceName)
//	if nil != method {
//		missedInits := findMissedEndpointsInitializations(method, expectedMethods)
//		fillMissedEndpointInits(method, missedInits)
//	}
//
//	writeSourceFile(filename, fileGo, fs)
//}
//
//func extractEndpointsStruct(src *ast.File, serviceName string) *ast.GenDecl {
//	var structDecl *ast.GenDecl
//	for _, d := range src.Decls {
//		if gen, ok := d.(*ast.GenDecl); ok {
//			if len(gen.Specs) > 0 {
//				if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
//					if "type" != gen.Tok.String() {
//						continue
//					}
//					if serviceName+"Endpoints" == typeSpec.Name.Name {
//						structDecl = gen
//						continue
//					}
//				}
//			}
//		}
//	}
//	return structDecl
//}
//
//func fillEndpointsStructWithMissedProperties(endpointsStruct *ast.GenDecl, serviceMethods []string) {
//	fieldsList := endpointsStruct.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
//	for _, m := range serviceMethods {
//		exists := false
//		for _, field := range fieldsList {
//			if m+"Endpoint" == field.Names[0].Name {
//				exists = true
//				break
//			}
//		}
//
//		if !exists {
//			missing := &ast.Field{
//				Names: []*ast.Ident{
//					{Name: m + "Endpoint"},
//				},
//				Type: fieldsList[0].Type,
//			}
//			fieldsList = append(fieldsList, missing)
//		}
//	}
//	endpointsStruct.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List = fieldsList
//}
//
//func extractEndpointMakeMethod(src *ast.File, serviceName string) *ast.FuncDecl {
//	return source.FindFuncDeclByName(src, naming.GetMakeEndpointsFuncName(serviceName))
//}
//
//func findMissedEndpointsInitializations(makeEndpointFunc *ast.FuncDecl, serviceMethods []string) []string {
//	var missed []string
//	for _, sm := range serviceMethods {
//		find := false
//		for _, initVal := range makeEndpointFunc.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts {
//			if initVal.(*ast.KeyValueExpr).Key.(*ast.Ident).Name == sm+"Endpoint" {
//				find = true
//				break
//			}
//		}
//		if !find {
//			missed = append(missed, sm)
//		}
//	}
//	return missed
//}
//
//func fillMissedEndpointInits(makeEndpointFunc *ast.FuncDecl, missedInits []string) {
//	exprs := makeEndpointFunc.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts
//
//	for _, missed := range missedInits {
//		expr := &ast.KeyValueExpr{
//			Key:   ast.NewIdent(missed + "Endpoint"),
//			Colon: exprs[0].(*ast.KeyValueExpr).Colon,
//			Value: &ast.CallExpr{
//				Fun: ast.NewIdent("make" + missed + "Endpoint"),
//				Args: []ast.Expr{
//					ast.NewIdent("svc"), // FIXME: hardcode
//				},
//			},
//		}
//		exprs = append(exprs, expr)
//	}
//
//	makeEndpointFunc.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts = exprs
//}
