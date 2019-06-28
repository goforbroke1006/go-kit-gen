package iterator

import "go/ast"

type AstIterator struct {
	file *ast.File
}

func (iter AstIterator) GetFuncDecl(funcName string) *ast.FuncDecl {
	return nil // TODO: implement me
}

func (iter AstIterator) GetStructFuncDecl(funcName string, structReceiver string) *ast.FuncDecl {
	return nil // TODO: implement me
}

func (iter AstIterator) GetStructDecl(structName string) *ast.GenDecl {
	return nil // TODO: implement me
}

func (iter AstIterator) GetInterfaceTypeSpec(interfaceName string) *ast.TypeSpec {
	return nil // TODO: implement me
}
