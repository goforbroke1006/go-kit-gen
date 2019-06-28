package iterator

import (
	"testing"
)

func TestAstIterator_GetFuncDecl(t *testing.T) {
	iter := NewAstFileIterator("testdata/enpoint.tmp.go")
	if got := iter.GetFuncDecl("makeMethodOneEndpoint"); nil == got {
		t.Errorf("AstFileIterator.GetFuncDecl() = nil, want funcDecl")
	}
}

func TestAstIterator_GetStructFuncDecl(t *testing.T) {
	iter := NewAstFileIterator("testdata/enpoint.tmp.go")
	if got := iter.GetStructFuncDecl("MethodWithReceiver", "SomeAwesomeHubEndpoints"); nil == got {
		t.Errorf("AstFileIterator.GetStructFuncDecl() = nil, want structFuncDecl")
	}
}

func TestAstIterator_GetStructDecl(t *testing.T) {
	iter := NewAstFileIterator("testdata/enpoint.tmp.go")
	if got := iter.GetStructDecl("SomeAwesomeHubEndpoints"); nil == got {
		t.Errorf("AstFileIterator.GetStructDecl() = nil, want structDecl")
	}
}

func TestAstIterator_GetInterfaceTypeSpec(t *testing.T) {
	iter := NewAstFileIterator("testdata/enpoint.tmp.go")
	if got := iter.GetInterfaceTypeSpec("TestInterface"); nil == got {
		t.Errorf("AstFileIterator.GetStructDecl() = nil, want interfaceType")
	}
}
