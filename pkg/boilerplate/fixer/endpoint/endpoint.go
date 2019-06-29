package endpoint

import (
	"go/ast"
	"log"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/builder"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
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
	endpointsStructName := naming.GetEndpointsStructName(ef.serviceName)

	afi := iterator.NewAstFileIterator(ef.file)
	endpointsStructure := afi.GetStructDecl(endpointsStructName)

	if nil == endpointsStructure {
		log.Println("Cant find endpoint struct", endpointsStructName)
		return
	}

	asdi := iterator.NewAstStructDeclIterator(endpointsStructure)
	asb := builder.NewAstStructBuilder(endpointsStructure)

	properties := asdi.GetProperties()
	for action := range ef.serviceActions {
		found := false
		endpointFieldName := naming.GetEndpointFieldName(action)
		for _, prop := range properties {
			if endpointFieldName == prop.Names[0].Name {
				found = true
				break
			}
		}
		if !found {
			asb.AddProperty(endpointFieldName, "endpoint.Endpoint")
		}
	}
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

	for action := range ef.serviceActions {
		builderFuncName := naming.GetEndpointBuilderFuncName(action)
		builderFuncDecl := source.FindFuncDeclByName(ef.file, builderFuncName)
		if nil == builderFuncDecl {

			// TODO: do magic
			//ef.file.Decls = append(ef.file.Decls, builderFuncDecl)

			// TODO: move with comments!!!
		}
	}
	//fixer.WriteSourceFile(ef.filename, file, fs)
}
