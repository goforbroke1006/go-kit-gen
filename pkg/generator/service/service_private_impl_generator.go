package service

type ServicePrivateImplGenerator interface {
	CreateEmptyStructIfNotExists(pbGoServiceName string) (name string, err error)
	CreateMethodDeclIfNotExists(pbGoActionName string, params [][2]string, returns [][2]string)
}

type servicePrivateImplGenerator struct {
}

func (s servicePrivateImplGenerator) CreateEmptyStructIfNotExists(pbGoServiceName string) (name string, err error) {
	// TODO: implement me
	return "", nil
}

func (s servicePrivateImplGenerator) CreateMethodDeclIfNotExists(pbGoActionName string, params [][2]string, returns [][2]string) {
	// TODO: implement me
}

func NewServicePrivateImplGenerator() ServicePrivateImplGenerator {
	return &servicePrivateImplGenerator{}
}
