package service

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
)

type ServiceInterfaceGenerator interface {
	CreateInterfaceIfNotExists(protoServiceName string) (name string, err error)
	CreateMethodSignatureIfNotExists(protoActionName string, params [][2]string, returns [][2]string)
}

type serviceInterfaceGenerator struct {
}

func (s serviceInterfaceGenerator) CreateInterfaceIfNotExists(protoServiceName string) (name string, err error) {
	return "", nil
}

func (s serviceInterfaceGenerator) CreateMethodSignatureIfNotExists(protoActionName string, params [][2]string, returns [][2]string) {
	// TODO: implement me
}

func NewServiceInterfaceGenerator(file *source.FileCrawler) ServiceInterfaceGenerator {
	return &serviceInterfaceGenerator{}
}
