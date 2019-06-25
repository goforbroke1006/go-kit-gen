package fixer

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var ef *EndpointFixer

func TestMain(m *testing.M) {
	ef = &EndpointFixer{
		filename:    "testdata/result.tmp",
		serviceName: "SomeAwesomeHub",
		serviceActions: []string{
			"MethodOne",
			"MethodTwo",
			"MethodThree",
			"SayHello",
		},
		templatesRelDir: "../../../",
	}

	initContent, err := ioutil.ReadFile("testdata/endpoint.go.sample.txt")
	if nil != err {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(ef.filename, []byte(initContent), os.ModePerm)
	if nil != err {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestEndpointFixer_addMissedRequestModels(t *testing.T) {
	ef.addMissedRequestModels()

	result, err := ioutil.ReadFile(ef.filename)
	if nil != err {
		t.Fatal(err)
	}

	for _, action := range ef.serviceActions {
		if !strings.Contains(string(result), "type "+action+"Request struct") {
			t.Fatal("" + action + "Request struct is missed")
		}
	}
}

func TestEndpointFixer_addMissedResponseModels(t *testing.T) {
	ef.addMissedResponseModels()

	result, err := ioutil.ReadFile(ef.filename)
	if nil != err {
		t.Fatal(err)
	}

	for _, action := range ef.serviceActions {
		if !strings.Contains(string(result), "type "+action+"Response struct") {
			t.Fatal("" + action + "Response struct is missed")
		}
	}
}

func TestEndpointFixer_addMissedPropertiesInEndpointsStruct(t *testing.T) {
	ef.addMissedPropertiesInEndpointsStruct()

	result, err := ioutil.ReadFile(ef.filename)
	if nil != err {
		t.Fatal(err)
	}

	for _, action := range ef.serviceActions {
		if !strings.Contains(string(result), action+"Endpoint	endpoint.Endpoint") {
			t.Fatal("" + action + "Endpoint field is missed")
		}
	}
}

func TestEndpointFixer_addMissedPropertyInitializationInMakeEndpointsFunc(t *testing.T) {
	ef.addMissedPropertyInitializationInMakeEndpointsFunc()

	result, err := ioutil.ReadFile(ef.filename)
	if nil != err {
		t.Fatal(err)
	}

	if !strings.Contains(string(result), "SayHelloEndpoint:") {
		t.Fatal("Endpoint prop init is missed")
	}
}

func TestEndpointFixer_addMissedEndpointBuilderFunc(t *testing.T) {
	ef.addMissedEndpointBuilderFunc()

	result, err := ioutil.ReadFile(ef.filename)
	if nil != err {
		t.Fatal(err)
	}

	if !strings.Contains(string(result), "func makeSayHelloEndpoint(") {
		t.Fatal("makeSayHelloEndpoint endpoint builder func is missed")
	}
}
