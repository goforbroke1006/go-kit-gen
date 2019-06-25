package endpoint

import (
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
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
	fs, file := fixer.OpenGolangSourceFile(ef.filename)

	tmpl := template.Must(template.ParseFiles(ef.templatesRelDir + "template/empty-struct.tmpl"))

	for _, action := range ef.serviceActions {
		structName := naming.GetEndpointRequestStructName(action)
		reqModelDecl := source.FindStructDeclByName(file, structName)
		if nil == reqModelDecl {
			filename := action + ".tmp"
			fileTmp, err := os.Create(filename)
			if nil != err {
				fmt.Println(err.Error())
			}
			err = tmpl.Execute(fileTmp, struct{ Name string }{Name: structName})
			if nil != err {
				fmt.Println(err.Error())
			}

			_, fileAstTmp := fixer.OpenGolangSourceFile(filename)
			structDecl := source.FindStructDeclByName(fileAstTmp, structName)

			file.Decls = append(file.Decls, structDecl)
		}
	}

	fixer.WriteSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedResponseModels() {
	fs, file := fixer.OpenGolangSourceFile(ef.filename)

	tmpl := template.Must(template.ParseFiles(ef.templatesRelDir + "template/empty-struct.tmpl"))

	for _, action := range ef.serviceActions {
		structName := naming.GetEndpointResponseStructName(action)
		respModelDecl := source.FindStructDeclByName(file, structName)
		if nil == respModelDecl {
			filename := action + ".tmp"
			fileTmp, err := os.Create(filename)
			if nil != err {
				fmt.Println(err.Error())
			}
			err = tmpl.Execute(fileTmp, struct{ Name string }{Name: structName})
			if nil != err {
				fmt.Println(err.Error())
			}

			_, fileAstTmp := fixer.OpenGolangSourceFile(filename)
			structDecl := source.FindStructDeclByName(fileAstTmp, structName)

			file.Decls = append(file.Decls, structDecl)
		}
	}

	fixer.WriteSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedPropertiesInEndpointsStruct() {
	fs, file := fixer.OpenGolangSourceFile(ef.filename)

	endpointsStructName := naming.GetEndpointsStructName(ef.serviceName)
	endpointsStructure := source.FindStructDeclByName(file, endpointsStructName)
	if nil == endpointsStructure {
		log.Println("Cant find endpoint struct", endpointsStructName)
		return
	}

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

	fixer.WriteSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedPropertyInitializationInMakeEndpointsFunc() {
	fs, file := fixer.OpenGolangSourceFile(ef.filename)
	endpointsFuncName := naming.GetMakeEndpointsFuncName(ef.serviceName)
	funcDecl := source.FindFuncDeclByName(file, endpointsFuncName)
	if nil == funcDecl {
		log.Println("Can't find make endpoints func", endpointsFuncName)
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

	fixer.WriteSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedEndpointBuilderFunc() {
	fs, file := fixer.OpenGolangSourceFile(ef.filename)

	tmpl := template.Must(template.ParseFiles(ef.templatesRelDir + "template/endpoint/endpoint-builder-func.tmpl"))

	for _, act := range ef.serviceActions {
		builderFuncName := naming.GetEndpointBuilderFuncName(act)
		builderFuncDecl := source.FindFuncDeclByName(file, builderFuncName)
		if nil == builderFuncDecl {

			//buffer := &bytes.Buffer{}
			filename := act + ".tmp"
			fileTmp, err := os.Create(filename)
			if nil != err {
				fmt.Println(err.Error())
			}
			err = tmpl.Execute(fileTmp, struct{ ActionName string }{ActionName: act})
			if nil != err {
				fmt.Println(err.Error())
			}

			_, fileAstTmp := fixer.OpenGolangSourceFile(filename)
			funcDecl := source.FindFuncDeclByName(fileAstTmp, builderFuncName)

			file.Decls = append(file.Decls, funcDecl)
			// TODO: move with comments!!!
		}
	}
	fixer.WriteSourceFile(ef.filename, file, fs)
}
