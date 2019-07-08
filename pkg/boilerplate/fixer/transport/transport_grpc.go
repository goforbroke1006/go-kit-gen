package transport

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/builder"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
)

func NewGRPCTransportFixer(file *ast.File, serviceName string) GRPCTransportFixer {
	afi := iterator.NewAstFileIterator(file)
	structDecl := afi.GetStructDecl(naming.GetTransportStructName(serviceName, "grpc"))
	return &grpcTransportFixer{
		file: file,
		decl: structDecl,
		asdi: iterator.NewAstStructDeclIterator(structDecl),
		asb:  builder.NewAstStructBuilder(structDecl),
		apf:  &factory.AstPrimitiveFactory{},
	}
}

type GRPCTransportFixer interface {
	FixServerImplStructField(actionName string) error
	FixServerImplStructMethod(actionName string, argsNameToType map[string]string) error
}

type grpcTransportFixer struct {
	file *ast.File
	decl *ast.GenDecl
	asdi *iterator.AstStructDeclIterator
	asb  *builder.AstStructBuilder
	apf  *factory.AstPrimitiveFactory
}

func (tf grpcTransportFixer) FixServerImplStructField(actionName string) error {
	if tf.asdi.HasProperty("handle" + actionName) {
		return errors.New(fmt.Sprintf("field '%s' already exists", actionName))
	}
	tf.asb.AddProperty("handle"+actionName, "grpc.Handler")
	return nil
}

func (tf grpcTransportFixer) FixServerImplStructMethod(
	actionName string,
	argsNameToType map[string]string,
) error {
	name := tf.asdi.GetName()
	funcDecl := tf.apf.CreateFuncDecl(actionName, argsNameToType, nil, []ast.Expr{}, nil, &name)
	tf.file.Decls = append(tf.file.Decls, funcDecl) // TODO: create AsFileBuilder
	return nil
}
