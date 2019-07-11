package generator

type TransportGRPCGenerator struct {
}

func (g TransportGRPCGenerator) CreateServerImpleStructIfNotExists(serviceName string) {
	// TODO: implement me
}

func (g TransportGRPCGenerator) CreateField(actionName string) {
	// TODO: implement me
}

func (g TransportGRPCGenerator) CreateMethodDecl(actionName string) {
	// TODO: implement me
}

func (g TransportGRPCGenerator) CreateServerContructorIfNotExists() {
	// TODO: implement me
}

func (g TransportGRPCGenerator) AddFieldInitInConstrucor(actionName string) {
	// TODO: implement me
}

func NewTransportGRPCGenerator() *TransportGRPCGenerator {
	return &TransportGRPCGenerator{}
}
