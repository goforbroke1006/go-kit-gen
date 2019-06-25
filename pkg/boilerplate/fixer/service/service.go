package service

import (
	"go/ast"
	"go/token"
	"log"

	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/source"
)

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
	// TODO:
	fixer.WriteSourceFile(sf.filename, file, fs)
}
