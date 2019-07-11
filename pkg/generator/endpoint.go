package generator

type EndpointsStructGenerator struct {
}

func (g EndpointsStructGenerator) CreateEndpointStructIfNotExists(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateEndpointStructField(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateRequestStruct(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateResponseStruct(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateMakeEndpointFunc(actionName string) {
	// TODO: implement me
}

func (g EndpointsStructGenerator) CreateConstructorIfNotExists() {
	// TODO: implement me
}

func (g EndpointsStructGenerator) SetFieldInContructor(actionName string) {
	// TODO: implement me
}

func NewEndpointsStructGenerator() *EndpointsStructGenerator {
	return &EndpointsStructGenerator{}
}
