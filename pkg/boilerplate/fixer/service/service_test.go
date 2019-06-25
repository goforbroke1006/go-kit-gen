package service

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/filesystem"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

const serviceUnfinishedSampleFilename = "testdata/service.go.sample.txt"
const templatesDirRelateToTestDir = "../../../../"
const testServiceName = "SomeAwesomeHub"

var testServiceActions = []string{
	"MethodOne",
	"MethodTwo",
	"MethodThree",
	"SayHello",
}

func buildTestServiceFixer(subjectFilename string) *ServiceFixer {
	return &ServiceFixer{
		filename:        subjectFilename,
		serviceName:     testServiceName,
		serviceActions:  testServiceActions,
		templatesRelDir: templatesDirRelateToTestDir,
	}
}

func TestServiceFixer_addMissedMethodSignaturesInServiceInterface(t *testing.T) {
	subjectFilename := "testdata/TestServiceFixer_addMissedMethodSignaturesInServiceInterface.tmp"

	if err := filesystem.CopyFileToFile(serviceUnfinishedSampleFilename, subjectFilename); nil != err {
		log.Fatal(err)
	}

	serviceFixer := buildTestServiceFixer(subjectFilename)
	serviceFixer.addMissedMethodSignaturesInServiceInterface()

	result, err := ioutil.ReadFile(subjectFilename)
	if nil != err {
		t.Fatal(err)
	}

	for _, action := range testServiceActions {
		if !strings.Contains(string(result), action+"(ctx context.Context)") {
			t.Fatal("Interface method " + action + " is missed")
		}
	}
}

func TestServiceFixer_addMissedMethodImplementationsInPrivateServiceStruct(t *testing.T) {
	subjectFilename := "testdata/TestServiceFixer_addMissedMethodImplementationsInPrivateServiceStruct.tmp"

	if err := filesystem.CopyFileToFile(serviceUnfinishedSampleFilename, subjectFilename); nil != err {
		log.Fatal(err)
	}

	sf := buildTestServiceFixer(subjectFilename)
	sf.addMissedMethodImplementationsInPrivateServiceStruct()

	result, err := ioutil.ReadFile(subjectFilename)
	if nil != err {
		t.Fatal(err)
	}

	svcNameLow := naming.GetServicePrivateImplStructName(sf.serviceName)
	for _, action := range testServiceActions {
		if !strings.Contains(string(result), "func (svc "+svcNameLow+") "+action+"(") {
			t.Fatal("Struct " + svcNameLow + " method " + action + " is missed")
		}
	}
}
