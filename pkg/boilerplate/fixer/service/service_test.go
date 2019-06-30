package service

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/filesystem"
)

const serviceUnfinishedSampleFilename = "testdata/service.go.sample.txt"
const testServiceName = "SomeAwesomeHub"

var testServiceActions = map[string]map[string]string{
	"MethodOne": {
		"FieldOne": "string",
	},
	"MethodTwo": {
		"FieldOne": "string",
		"FieldTwo": "uint64",
	},
	"MethodThree": {
		"FieldOne":   "string",
		"FieldTwo":   "string",
		"FieldThree": "",
	},
	"SayHello": {},
}

func TestServiceFixer_addMissedMethodSignaturesInServiceInterface(t *testing.T) {
	subjectFilename := "testdata/addMissedMethodSignaturesInServiceInterface.tmp.go"

	if err := filesystem.CopyFileToFile(serviceUnfinishedSampleFilename, subjectFilename); nil != err {
		log.Fatal(err)
	}

	fset, file := fixer.OpenGolangSourceFile(subjectFilename)

	sf := NewServiceFixer(file, testServiceName, testServiceActions)
	sf.addMissedMethodSignaturesInServiceInterface()

	fixer.WriteSourceFile(subjectFilename, file, fset)

	result, err := ioutil.ReadFile(subjectFilename)
	if nil != err {
		t.Fatal(err)
	}

	for action := range testServiceActions {
		if !strings.Contains(string(result), action+"(ctx context.Context)") {
			t.Fatal("Interface method " + action + " is missed")
		}
	}
}

func TestServiceFixer_addMissedMethodImplementationsInPrivateServiceStruct(t *testing.T) {
	subjectFilename := "testdata/addMissedMethodImplementationsInPrivateServiceStruct.tmp.go"

	if err := filesystem.CopyFileToFile(serviceUnfinishedSampleFilename, subjectFilename); nil != err {
		log.Fatal(err)
	}

	fset, file := fixer.OpenGolangSourceFile(subjectFilename)

	sf := NewServiceFixer(file, testServiceName, testServiceActions)
	sf.addMissedMethodImplementationsInPrivateServiceStruct()

	fixer.WriteSourceFile(subjectFilename, file, fset)

	result, err := ioutil.ReadFile(subjectFilename)
	if nil != err {
		t.Fatal(err)
	}

	svcNameLow := naming.GetServicePrivateImplStructName(sf.serviceName)
	for action := range testServiceActions {
		if !strings.Contains(string(result), "func (svc "+svcNameLow+") "+action+"(") {
			t.Fatal("Struct " + svcNameLow + " method " + action + " is missed")
		}
	}
}
