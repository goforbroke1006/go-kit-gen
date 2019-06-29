package iterator

import (
	"testing"

	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
)

func TestAstIterator_GetFuncDecl(t *testing.T) {
	_, file := fixer.OpenGolangSourceFile("testdata/enpoint.tmp.go")
	iter := NewAstFileIterator(file)
	if got := iter.GetFuncDecl("makeMethodOneEndpoint"); nil == got {
		t.Errorf("AstFileIterator.GetFuncDecl() = nil, want funcDecl")
	}
}

func TestAstIterator_GetStructFuncDecl(t *testing.T) {
	_, file := fixer.OpenGolangSourceFile("testdata/enpoint.tmp.go")
	iter := NewAstFileIterator(file)
	if got := iter.GetStructFuncDecl("MethodWithReceiver", "SomeAwesomeHubEndpoints"); nil == got {
		t.Errorf("AstFileIterator.GetStructFuncDecl() = nil, want structFuncDecl")
	}
}

func TestAstIterator_GetStructDecl(t *testing.T) {
	_, file := fixer.OpenGolangSourceFile("testdata/enpoint.tmp.go")
	iter := NewAstFileIterator(file)
	if got := iter.GetStructDecl("SomeAwesomeHubEndpoints"); nil == got {
		t.Errorf("AstFileIterator.GetStructDecl() = nil, want structDecl")
	}
}

func TestAstIterator_GetInterfaceTypeSpec(t *testing.T) {
	_, file := fixer.OpenGolangSourceFile("testdata/enpoint.tmp.go")
	iter := NewAstFileIterator(file)
	if got := iter.GetInterfaceTypeSpec("TestInterface"); nil == got {
		t.Errorf("AstFileIterator.GetStructDecl() = nil, want interfaceType")
	}
}
