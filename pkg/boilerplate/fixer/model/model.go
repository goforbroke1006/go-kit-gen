package model

import (
	fixer2 "github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
)

func NewModelFixed(filename, serviceName string, serviceActions []string) *ModelFixed {
	return &ModelFixed{
		filename:        filename,
		serviceName:     serviceName,
		serviceActions:  serviceActions,
		templatesRelDir: "./",
	}
}

type ModelFixed struct {
	filename        string
	serviceName     string
	serviceActions  []string
	templatesRelDir string
}

func (mf ModelFixed) Fix() {
	fs, astFile := fixer2.OpenGolangSourceFile(mf.filename)
	// TODO:
	fixer2.WriteSourceFile(mf.filename, astFile, fs)
}
