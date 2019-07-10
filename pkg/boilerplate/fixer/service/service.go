package service

import (
	"go/ast"
	"log"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/builder"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/old/source"
)

func NewServiceFixer(file *ast.File, serviceName string, serviceActions map[string]map[string]string) *ServiceFixer {
	return &ServiceFixer{
		file:           file,
		serviceName:    serviceName,
		serviceActions: serviceActions,
	}
}

type ServiceFixer struct {
	file           *ast.File
	serviceName    string
	serviceActions map[string]map[string]string
}

func (sf ServiceFixer) Fix() {
	sf.addMissedMethodSignaturesInServiceInterface()
	sf.addMissedMethodImplementationsInPrivateServiceStruct()
}

func (sf ServiceFixer) addMissedMethodSignaturesInServiceInterface() {

	afi := iterator.NewAstFileIterator(sf.file)

	serviceInterfaceName := naming.GetServiceInterfaceName(sf.serviceName)
	intfc := afi.GetInterfaceTypeSpec(serviceInterfaceName)
	if nil == intfc {
		log.Println("Cant find service interface", serviceInterfaceName)
		return
	}

	aiti := iterator.NewAstInterfaceTypeIterator(intfc)
	aib := builder.NewAstInterfaceBuilder(intfc)
	apf := factory.AstPrimitiveFactory{}

	for action := range sf.serviceActions {
		methodDecl := aiti.GetMethod(action)
		if nil == methodDecl {
			funcSign := apf.CreateFuncSignatureExpr(
				action,
				map[string]string{"ctx": "context.Context"},
				map[string]string{},
			)
			aib.AddFuncSignature(funcSign)
			log.Println("Add", action, "method signature to", serviceInterfaceName, "interface")
		}
	}
}

func (sf ServiceFixer) addMissedMethodImplementationsInPrivateServiceStruct() {

	implStructName := naming.GetServicePrivateImplStructName(sf.serviceName)
	svcImplStructDecl := source.FindStructDeclByName(sf.file, implStructName)
	if nil == svcImplStructDecl {
		log.Println("Cant find service impl struct", implStructName)
		return
	}

	afi := iterator.NewAstFileIterator(sf.file)
	apf := factory.AstPrimitiveFactory{}

	//tmpl := template.Must(template.ParseFiles(sf.templatesRelDir + "template/service/service-impl-struct-sample.tmpl"))

	recvName := "svc"
	for action := range sf.serviceActions {
		svcMethod := afi.GetStructFuncDecl(action, implStructName)
		if nil == svcMethod {
			svcMethod := apf.CreateFuncDecl(
				action,
				map[string]string{"ctx": "context.Context"},
				map[string]string{"res": "", "err": "error"},
				[]ast.Expr{nil, nil},
				&recvName, &implStructName,
			)
			sf.file.Decls = append(sf.file.Decls, svcMethod)

			log.Println("Create", "(", implStructName, ")", action, "func")

			//tmpFilename := action + ".tmp"
			//fileTmp, err := os.Create(tmpFilename)
			//if nil != err {
			//	fmt.Println(err.Error())
			//}
			//err = tmpl.Execute(fileTmp, struct {
			//	StructName string
			//	MethodName string
			//}{StructName: implStructName, MethodName: action})
			//if nil != err {
			//	fmt.Println(err.Error())
			//}
			//
			//_, fileAstTmp := fixer2.OpenGolangSourceFile(tmpFilename)
			//missedFuncDecl := source.FindFuncDeclByName(fileAstTmp, action)
			//
			//file.Decls = append(file.Decls, missedFuncDecl)
			//err = os.Remove(tmpFilename)
			//if nil != err {
			//	log.Println("Can't remove file", tmpFilename)
			//}

			//missedFuncDecl := &ast.FuncDecl{
			//	Recv: &ast.FieldList{
			//		List: []*ast.Field{
			//			{
			//				Names: []*ast.Ident{
			//					{
			//						Name: "svc",
			//						Obj: &ast.Object{
			//							Name: "svc",
			//							Decl: &ast.Field{
			//								Type: ast.NewIdent(implStructName),
			//							},
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//}

			// TODO:
		}
	}
}
