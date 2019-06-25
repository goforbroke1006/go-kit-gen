package service

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"text/template"

	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/source"
)

func NewServiceFixer(filename, serviceName string, serviceActions []string) *ServiceFixer {
	return &ServiceFixer{
		filename:        filename,
		serviceName:     serviceName,
		serviceActions:  serviceActions,
		templatesRelDir: "./",
	}
}

type ServiceFixer struct {
	filename        string
	serviceName     string
	serviceActions  []string
	templatesRelDir string
}

func (sf ServiceFixer) Fix() {
	sf.addMissedMethodSignaturesInServiceInterface()
	sf.addMissedMethodImplementationsInPrivateServiceStruct()
}

func (sf ServiceFixer) addMissedMethodSignaturesInServiceInterface() {
	fs, file := fixer.OpenGolangSourceFile(sf.filename)

	serviceInterfaceName := naming.GetServiceInterfaceName(sf.serviceName)
	ifc := source.FindInterfaceByName(file, serviceInterfaceName)
	if nil == ifc {
		log.Println("Cant find service interface", serviceInterfaceName)
		return
	}

	ib := source.NewInterfaceBuilder(ifc)
	actualMethodsList := ib.GetMethods()

	for _, action := range sf.serviceActions {
		found := false
		for _, methodDecl := range actualMethodsList {
			if action == methodDecl.Names[0].Name {
				found = true
				break
			}
		}
		if !found {
			missedMethodDecl := &ast.Field{
				Names: []*ast.Ident{
					{
						Name: action,
						Obj:  &ast.Object{Kind: ast.Fun, Name: action},
					},
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "ctx", NamePos: token.NoPos},
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("context"),
									Sel: ast.NewIdent("Context"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.InterfaceType{
									Methods: &ast.FieldList{},
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
			}

			ib.AppendMethod(missedMethodDecl)
		}
	}

	fixer.WriteSourceFile(sf.filename, file, fs)
}

func (sf ServiceFixer) addMissedMethodImplementationsInPrivateServiceStruct() {
	fs, file := fixer.OpenGolangSourceFile(sf.filename)
	implStructName := naming.GetServicePrivateImplStructName(sf.serviceName)
	svcImplStructDecl := source.FindStructDeclByName(file, implStructName)
	if nil == svcImplStructDecl {
		log.Println("Cant find service impl struct", implStructName)
		return
	}

	actualMethodsList := source.FindStructMethods(file, implStructName)

	tmpl := template.Must(template.ParseFiles(sf.templatesRelDir + "template/service/service-impl-struct-sample.tmpl"))

	for _, action := range sf.serviceActions {
		found := false
		for _, m := range actualMethodsList {
			if action == m.Name.Name {
				found = true
				break
			}
		}
		if !found {
			filename := action + ".tmp"
			fileTmp, err := os.Create(filename)
			if nil != err {
				fmt.Println(err.Error())
			}
			err = tmpl.Execute(fileTmp, struct {
				StructName string
				MethodName string
			}{StructName: implStructName, MethodName: action})
			if nil != err {
				fmt.Println(err.Error())
			}

			_, fileAstTmp := fixer.OpenGolangSourceFile(filename)
			missedFuncDecl := source.FindFuncDeclByName(fileAstTmp, action)

			file.Decls = append(file.Decls, missedFuncDecl)

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

	fixer.WriteSourceFile(sf.filename, file, fs)
}
