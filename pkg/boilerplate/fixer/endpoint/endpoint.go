package endpoint

import (
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"go/ast"
	"log"
	"os"
	"text/template"

	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/old/source"
)

func NewEndpointFixer(file *ast.File, serviceName string, serviceActions map[string]map[string]string) *EndpointFixer {
	return &EndpointFixer{
		file:           file,
		serviceName:    serviceName,
		serviceActions: serviceActions,
	}
}

type EndpointFixer struct {
	file            *ast.File
	serviceName     string
	serviceActions  map[string]map[string]string
	templatesRelDir string
}

func (ef EndpointFixer) GetFile() *ast.File {
	return ef.file
}

func (ef EndpointFixer) Fix() {
	ef.addMissedRequestModels()
	ef.addMissedResponseModels()
	ef.addMissedPropertiesInEndpointsStruct()
	ef.addMissedPropertyInitializationInMakeEndpointsFunc()
	ef.addMissedEndpointBuilderFunc()
}

func (ef EndpointFixer) addMissedRequestModels() {
	afi := iterator.NewAstFileIterator(ef.file)
	apf := factory.AstPrimitiveFactory{}

	for action, propsDescrList := range ef.serviceActions {
		structName := naming.GetEndpointRequestStructName(action)
		reqModelDecl := afi.GetStructDecl(structName)
		if nil == reqModelDecl {
			reqModelDecl = apf.CreateStructDecl(structName, propsDescrList)
			ef.file.Decls = append(ef.file.Decls, reqModelDecl)
			log.Println("Create", structName, "struct")
		}
	}
}

func (ef EndpointFixer) addMissedResponseModels() {
	afi := iterator.NewAstFileIterator(ef.file)
	apf := factory.AstPrimitiveFactory{}

	for action := range ef.serviceActions {
		structName := naming.GetEndpointResponseStructName(action)
		respModelDecl := afi.GetStructDecl(structName)
		if nil == respModelDecl {
			respModelDecl = apf.CreateStructDecl(
				structName,
				map[string]string{
					"Err": "string",
				},
			)
			ef.file.Decls = append(ef.file.Decls, respModelDecl)
			log.Println("Create", structName, "struct")
		}
	}
}

func (ef EndpointFixer) addMissedPropertiesInEndpointsStruct() {
	//fs, file := fixer.OpenGolangSourceFile(ef.filename)

	endpointsStructName := naming.GetEndpointsStructName(ef.serviceName)
	endpointsStructure := source.FindStructDeclByName(ef.file, endpointsStructName)
	if nil == endpointsStructure {
		log.Println("Cant find endpoint struct", endpointsStructName)
		return
	}

	propBuilder := source.NewStructBuilder(endpointsStructure)

	properties := propBuilder.GetStructProperties()
	for action := range ef.serviceActions {
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

	//fixer.WriteSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedPropertyInitializationInMakeEndpointsFunc() {
	//fs, file := fixer.OpenGolangSourceFile(ef.filename)
	endpointsFuncName := naming.GetMakeEndpointsFuncName(ef.serviceName)
	funcDecl := source.FindFuncDeclByName(ef.file, endpointsFuncName)
	if nil == funcDecl {
		log.Println("Can't find make endpoints func", endpointsFuncName)
		return
	}

	var missedActionsList []string

	{
		for action := range ef.serviceActions {
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

	//fixer.WriteSourceFile(ef.filename, file, fs)
}

func (ef EndpointFixer) addMissedEndpointBuilderFunc() {
	//fs, file := fixer.OpenGolangSourceFile(ef.filename)

	tmpl := template.Must(template.ParseFiles(ef.templatesRelDir + "template/endpoint/endpoint-builder-func.tmpl"))

	for action := range ef.serviceActions {
		builderFuncName := naming.GetEndpointBuilderFuncName(action)
		builderFuncDecl := source.FindFuncDeclByName(ef.file, builderFuncName)
		if nil == builderFuncDecl {

			//buffer := &bytes.Buffer{}
			tmpFilename := action + ".tmp"
			fileTmp, err := os.Create(tmpFilename)
			if nil != err {
				fmt.Println(err.Error())
			}
			err = tmpl.Execute(fileTmp, struct{ ActionName string }{ActionName: action})
			if nil != err {
				fmt.Println(err.Error())
			}

			_, fileAstTmp := fixer.OpenGolangSourceFile(tmpFilename)
			funcDecl := source.FindFuncDeclByName(fileAstTmp, builderFuncName)

			err = os.Remove(tmpFilename)
			if nil != err {
				log.Println("Can't remove file", tmpFilename)
			}

			ef.file.Decls = append(ef.file.Decls, funcDecl)
			// TODO: move with comments!!!
		}
	}
	//fixer.WriteSourceFile(ef.filename, file, fs)
}
