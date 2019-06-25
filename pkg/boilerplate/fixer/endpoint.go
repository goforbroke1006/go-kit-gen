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
