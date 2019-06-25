package source

import "go/ast"

func NewStructBuilder(decl *ast.GenDecl) *StructBuilder {
	return &StructBuilder{
		decl: decl,
	}
}

type StructBuilder struct {
	decl *ast.GenDecl
}

func (sb StructBuilder) GetStructProperties() []*ast.Field {
	return sb.decl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
}

func (sb StructBuilder) StructAddProperty(prop *ast.Field) {
	fields := sb.decl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
	fields = append(fields, prop)
	sb.decl.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List = fields
}

func NewInterfaceBuilder(decl *ast.GenDecl) *InterfaceBuilder {
	return &InterfaceBuilder{
		decl: decl,
	}
}

type InterfaceBuilder struct {
	decl *ast.GenDecl
}

func (ib InterfaceBuilder) GetMethods() []*ast.Field {
	return ib.decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
}

func (ib InterfaceBuilder) AppendMethod(f *ast.Field) {
	fields := ib.decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
	fields = append(fields, f)
	ib.decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List = fields
}
