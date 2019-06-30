package model

import (
	"go/ast"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
)

func NewModelFixer(file *ast.File, serviceName string, serviceActions map[string]map[string]string) *ModelFixer {
	return &ModelFixer{
		file:           file,
		serviceName:    serviceName,
		serviceActions: serviceActions,
	}
}

type ModelFixer struct {
	file           *ast.File
	serviceName    string
	serviceActions map[string]map[string]string
}

func (mf ModelFixer) Fix() {
	// TODO: import "pb" 		package from proto file
	// TODO: import "endpoint" 	package from target project

	afi := iterator.NewAstFileIterator(mf.file)
	apf := factory.AstPrimitiveFactory{}

	for action := range mf.serviceActions {
		{
			decoderFuncName := naming.GetDecodePbToEndpRequestMethodName(action)
			decoderFuncDecl := afi.GetFuncDecl(decoderFuncName)

			if nil == decoderFuncDecl {
				decoderFuncDecl = apf.CreateFuncDecl(
					decoderFuncName,
					map[string]string{"_": "context.Context", "pbReq": ""}, // TODO: set real request type like pb.<ACTION>Request
					map[string]string{"endpReq": "", "err": "error"},       // TODO: set real request type like endpoint.<ACTION>Response
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
					map[string]string{"_": "context.Context", "endpResp": ""}, // TODO: set real request type like endpoint.<ACTION>Response
					map[string]string{"pbResp": "", "err": "error"},           // TODO: set real request type like pb.<ACTION>Response
					[]ast.Expr{nil, nil},
					nil, nil,
				)
				mf.file.Decls = append(mf.file.Decls, encoderFuncDecl)
			}
		}
	}

}
