package model

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
)

func NewModelFixer(
	workingDir string,
	file *ast.File,
	pbPackage string,
	serviceName string,
	serviceActions map[string]map[string]string,
) *ModelFixer {
	endpointDir, _ := filepath.Abs(workingDir + "/endpoint")
	endpointFullPackage := strings.TrimPrefix(endpointDir, os.Getenv("GOPATH"))

	return &ModelFixer{
		endpointPkg:    endpointFullPackage,
		file:           file,
		pbPackage:      pbPackage,
		serviceName:    serviceName,
		serviceActions: serviceActions,
	}
}

type ModelFixer struct {
	endpointPkg    string
	file           *ast.File
	pbPackage      string
	serviceName    string
	serviceActions map[string]map[string]string
}

func (mf ModelFixer) Fix() {
	afi := iterator.NewAstFileIterator(mf.file)

	// TODO: import "pb" 		package from proto file

	// TODO: import "endpoint" 	package from target project
	//if !afi.HasImport(mf.pbPackage) {
	//	mf.file.Imports = append(mf.file.Imports, &ast.ImportSpec{
	//		Path: &ast.BasicLit{
	//			ValuePos: 1000,
	//			Kind:     token.STRING,
	//			Value:    "\"" + mf.endpointPkg + "\"",
	//		},
	//	})
	//}

	apf := factory.AstPrimitiveFactory{}

	for action := range mf.serviceActions {
		{
			decoderFuncName := naming.GetDecodePbToEndpRequestMethodName(action)
			decoderFuncDecl := afi.GetFuncDecl(decoderFuncName)

			if nil == decoderFuncDecl {
				decoderFuncDecl = apf.CreateFuncDecl(
					decoderFuncName,
					map[string]string{
						"_":     "context.Context",
						"pbReq": mf.pbPackage + "." + action + "Request",
					},
					map[string]string{
						"endpReq": "endpoint." + action + "Request",
						"err":     "error",
					},
					[]ast.Expr{nil, nil},
					nil, nil,
				)
				mf.file.Decls = append(mf.file.Decls, decoderFuncDecl)
			}
		}

		{
			encoderFuncName := naming.GetEncodeEndpToPbResponseMethodName(action)
			encoderFuncDecl := afi.GetFuncDecl(encoderFuncName)

			if nil == encoderFuncDecl {
				encoderFuncDecl = apf.CreateFuncDecl(
					encoderFuncName,
					map[string]string{
						"_":        "context.Context",
						"endpResp": "endpoint." + action + "Response",
					},
					map[string]string{
						"pbResp": mf.pbPackage + "." + action + "Response",
						"err":    "error",
					},
					[]ast.Expr{nil, nil},
					nil, nil,
				)
				mf.file.Decls = append(mf.file.Decls, encoderFuncDecl)
			}
		}
	}

}
