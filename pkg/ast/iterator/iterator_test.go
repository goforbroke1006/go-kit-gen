package iterator

import (
	"go/ast"
	"testing"
)

func TestAstIterator_GetFuncDecl(t *testing.T) {
	var file *ast.File

	// TODO: init file from testdata

	iter := AstIterator{
		file: file,
	}
	if got := iter.GetFuncDecl("SayHello"); nil == got {
		t.Errorf("AstIterator.GetFuncDecl() = nil, want funcDecl")
	}

}
