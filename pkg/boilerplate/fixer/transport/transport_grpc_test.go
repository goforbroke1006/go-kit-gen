package transport

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/filesystem"
	"log"
	"testing"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/builder"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/factory"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
)

func Test_grpcTransportFixer_FixServerImplStructField(t *testing.T) {
	const serviceName = "SomeAwesomeHub"
	type args struct {
		actionName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "must not create existing field", args: args{actionName: "MethodOne"}, wantErr: true},
		{name: "must create field", args: args{actionName: "MethodThree"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subjectFilename := "./testdata/result/FixServerImplStructField-" + tt.name + "-.tmp.go"

			if err := filesystem.CopyFileToFile("./testdata/transport_sample.go", subjectFilename); nil != err {
				log.Fatal(err)
			}

			fileSet, file := fixer.OpenGolangSourceFile(subjectFilename)
			afi := iterator.NewAstFileIterator(file)
			structDecl := afi.GetStructDecl(naming.GetTransportStructName(serviceName, "grpc"))

			tf := grpcTransportFixer{
				file: file,
				decl: structDecl,
				asdi: iterator.NewAstStructDeclIterator(structDecl),
				asb:  builder.NewAstStructBuilder(structDecl),
				apf:  &factory.AstPrimitiveFactory{},
			}
			if err := tf.FixServerImplStructField(tt.args.actionName); (err != nil) != tt.wantErr {
				t.Errorf("grpcTransportFixer.FixServerImplStructField() error = %v, wantErr %v", err, tt.wantErr)
			}

			fixer.WriteSourceFile(subjectFilename, file, fileSet)
		})
	}
}
