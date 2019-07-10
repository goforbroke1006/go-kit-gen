package service

type ServiceInterfaceGenerator interface {
	CreateInterfaceIfNotExists(protoServiceName string)
	CreateMethodSignatureIfNotExists(protoActionName string)
}

type serviceInterfaceGenerator struct {
}

func (s serviceInterfaceGenerator) CreateInterfaceIfNotExists(protoServiceName string) {
	panic("implement me")
}

func (s serviceInterfaceGenerator) CreateMethodSignatureIfNotExists(protoActionName string) {
	panic("implement me")
}

func NewServiceInterfaceGenerator() ServiceInterfaceGenerator {
	return &serviceInterfaceGenerator{}
}
