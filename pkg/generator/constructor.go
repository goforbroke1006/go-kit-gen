package generator

type ConstructorGenerator interface {
	CreateServiceConstructorMethod(interfaceName string, implStructName string)
}

type constructorGenerator struct {
}

func (c constructorGenerator) CreateServiceConstructorMethod(interfaceName string, implStructName string) {
	// TODO: implement me
}

func NewConstructorGenerator() ConstructorGenerator {
	return &constructorGenerator{}
}
