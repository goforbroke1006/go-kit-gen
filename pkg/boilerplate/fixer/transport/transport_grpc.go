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
	structName := naming.GetTransportStructName(serviceName, "grpc")
	structDecl := afi.GetStructDecl(structName)

	apf := &factory.AstPrimitiveFactory{}

	if nil == structDecl {
		structDecl = apf.CreateStructDecl(structName, map[string]string{})
		file.Decls = append(file.Decls, structDecl)
	}

	return &grpcTransportFixer{
		file: file,
		afi:  afi,
		decl: structDecl,
		asdi: iterator.NewAstStructDeclIterator(structDecl),
		asb:  builder.NewAstStructBuilder(structDecl),
		apf:  apf,
	}
}

type GRPCTransportFixer interface {
	FixServerImplStructField(actionName string) (fieldName string, err error)
	FixServerImplStructMethod(actionName string, argsNameToType []*ast.Field, returnsNameToType []*ast.Field) error
	FixDecodeMethod(actionName string, argsNameToType []*ast.Field, returnsNameToType []*ast.Field) error
}

type grpcTransportFixer struct {
	file *ast.File
	afi  *iterator.AstFileIterator
	decl *ast.GenDecl
	asdi *iterator.AstStructDeclIterator
	asb  *builder.AstStructBuilder
	apf  *factory.AstPrimitiveFactory
}

func (tf grpcTransportFixer) FixServerImplStructField(actionName string) (fieldName string, err error) {
	fieldName = "handle" + actionName

	if tf.asdi.HasProperty(fieldName) {
		return fieldName, errors.New(fmt.Sprintf("field '%s' already exists", fieldName))
	}
	tf.asb.AddProperty(fieldName, "grpc.Handler")
	return
}

func (tf grpcTransportFixer) FixServerImplStructMethod(
	actionName string,
	argsDecls []*ast.Field,
	returnsDecls []*ast.Field,
) error {
	if nil != tf.afi.GetFuncDecl(actionName) {
		return errors.New("method '" + actionName + "' already exists")
	}
	name := tf.asdi.GetName()
	funcDecl := tf.apf.CreateFuncDecl2(actionName, argsDecls, returnsDecls, []ast.Expr{nil, nil}, nil, &name)
	tf.file.Decls = append(tf.file.Decls, funcDecl) // TODO: create AsFileBuilder
	return nil
}

func (tf grpcTransportFixer) FixDecodeMethod(
	actionName string,
	argsDecls []*ast.Field,
	returnsDecls []*ast.Field,
) error {
	methodName := naming.GetDecodePbToEndpRequestMethodName(actionName)
	if nil != tf.afi.GetFuncDecl(methodName) {
		return errors.New("method '" + methodName + "' already exists")
	}

	//funcDecl := tf.apf.CreateFuncDecl2(methodName, argsDecls, returnsDecls, []ast.Expr{nil, nil}, nil, nil)
	//tf.file.Decls = append(tf.file.Decls, funcDecl) // TODO: create AsFileBuilder
	return nil
}
