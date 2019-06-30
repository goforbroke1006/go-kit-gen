package endpoint

import (
	"go/ast"
	"log"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/builder"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
)

func NewEndpointFixer(file *ast.File, serviceName string, serviceActions map[string]map[string]string) *EndpointFixer {
	return &EndpointFixer{
		file:           file,
		serviceName:    serviceName,
		serviceActions: serviceActions,
	}
}

type EndpointFixer struct {
	file           *ast.File
	serviceName    string
	serviceActions map[string]map[string]string
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
	afi := iterator.NewAstFileIterator(ef.file)

	endpointsFuncName := naming.GetMakeEndpointsFuncName(ef.serviceName)
	funcDecl := afi.GetFuncDecl(endpointsFuncName)

	if nil == funcDecl {
		log.Println("Can't find make endpoints func", endpointsFuncName)
		return
	}

	afdi := iterator.NewAstFuncDeclIterator(funcDecl)
	litValueOfEndpointStruct := afdi.GetReturnSmtm().Results[0].(*ast.CompositeLit)
	aclb := builder.NewAstCompositeLitBuilder(litValueOfEndpointStruct)

	apf := factory.AstPrimitiveFactory{}

	for action := range ef.serviceActions {
		endpointFieldName := naming.GetEndpointFieldName(action)

		found := false
		for _, initVal := range funcDecl.Body.List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts {
			if endpointFieldName == initVal.(*ast.KeyValueExpr).Key.(*ast.Ident).Name {
				found = true
				break
			}
		}

		if !found {
			aclb.AddElement(
				endpointFieldName,
				apf.CreateMethodCallExpr(
					naming.GetEndpointBuilderFuncName(action),
					[]string{"svc"},
				),
			)
		}
	}
}

func (ef EndpointFixer) addMissedEndpointBuilderFunc() {
	apf := factory.AstPrimitiveFactory{}
	afi := iterator.NewAstFileIterator(ef.file)

	for action := range ef.serviceActions {
		builderFuncName := naming.GetEndpointBuilderFuncName(action)
		builderFuncDecl := afi.GetFuncDecl(builderFuncName)
		if nil == builderFuncDecl {
			builderFuncDecl := apf.CreateFuncDecl(
				builderFuncName,
				map[string]string{"svc": ""},
				map[string]string{"": "endpoint.Endpoint"},
				[]ast.Expr{
					apf.CreateAnonFuncObjectDecl(
						map[string]string{"ctx": "context.Context", "req": ""},
						map[string]string{"res": "", "err": "error"},
						[]ast.Expr{nil, nil},
					),
				},
			)
			ef.file.Decls = append(ef.file.Decls, builderFuncDecl)
			log.Println("Create", builderFuncName, "function")
		}
	}
}
