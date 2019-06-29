package endpoint

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var ef *EndpointFixer

const endpointUnfinishedSampleFilename = "testdata/endpoint.go.sample.txt"
const experimentalSampleFilename = "testdata/endpoint-result.tmp"
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

func recreateTestSourceFile(sampleFilename, filename string) {
	initContent, err := ioutil.ReadFile(sampleFilename)
	if nil != err {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, []byte(initContent), os.ModePerm)
	if nil != err {
		log.Fatal(err)
	}
}

//func TestMain(m *testing.M) {
//
//	_, file := fixer.OpenGolangSourceFile(experimentalSampleFilename)
//	ef = NewEndpointFixer(file, testServiceName, testServiceActions)
//	recreateTestSourceFile(endpointUnfinishedSampleFilename, experimentalSampleFilename)
//
//	os.Exit(m.Run())
//}

func TestEndpointFixer_addMissedRequestModels(t *testing.T) {
	recreateTestSourceFile(endpointUnfinishedSampleFilename, experimentalSampleFilename)

	fset, file := fixer.OpenGolangSourceFile(experimentalSampleFilename)

	ef = NewEndpointFixer(file, testServiceName, testServiceActions)
	ef.addMissedRequestModels()

	fixer.WriteSourceFile(experimentalSampleFilename, file, fset)

	result, err := ioutil.ReadFile(experimentalSampleFilename)
	if nil != err {
		t.Fatal(err)
	}

	for action := range ef.serviceActions {
		if !strings.Contains(string(result), "type "+action+"Request struct") {
			t.Fatal("" + action + "Request struct is missed")
		}
	}
}

func TestEndpointFixer_addMissedResponseModels(t *testing.T) {
	ef.addMissedResponseModels()

	result, err := ioutil.ReadFile(experimentalSampleFilename)
	if nil != err {
		t.Fatal(err)
	}

	for action := range ef.serviceActions {
		if !strings.Contains(string(result), "type "+action+"Response struct") {
			t.Fatal("" + action + "Response struct is missed")
		}
	}
}

func TestEndpointFixer_addMissedPropertiesInEndpointsStruct(t *testing.T) {
	ef.addMissedPropertiesInEndpointsStruct()

	result, err := ioutil.ReadFile(experimentalSampleFilename)
	if nil != err {
		t.Fatal(err)
	}

	for action := range ef.serviceActions {
		if !strings.Contains(string(result), action+"Endpoint	endpoint.Endpoint") {
			t.Fatal("" + action + "Endpoint field is missed")
		}
	}
}

func TestEndpointFixer_addMissedPropertyInitializationInMakeEndpointsFunc(t *testing.T) {
	ef.addMissedPropertyInitializationInMakeEndpointsFunc()

	result, err := ioutil.ReadFile(experimentalSampleFilename)
	if nil != err {
		t.Fatal(err)
	}

	if !strings.Contains(string(result), "SayHelloEndpoint:") {
		t.Fatal("Endpoint prop init is missed")
	}
}

func TestEndpointFixer_addMissedEndpointBuilderFunc(t *testing.T) {
	ef.addMissedEndpointBuilderFunc()

	result, err := ioutil.ReadFile(experimentalSampleFilename)
	if nil != err {
		t.Fatal(err)
	}

	if !strings.Contains(string(result), "func makeSayHelloEndpoint(") {
		t.Fatal("makeSayHelloEndpoint endpoint builder func is missed")
	}
}