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

//func TestFixMissingEndpoints(t *testing.T) {
//	initContent := `package endpoint
//
//import (
//	"context"
//	"github.com/go-kit/kit/endpoint"
//)
//
//// TODO:
//
//type SomeAwesomeHubEndpoints struct {
//	MethodOneEndpoint endpoint.Endpoint
//	MethodTwoEndpoint endpoint.Endpoint
//	MethodThreeEndpoint endpoint.Endpoint
//}
//
//func MakeSomeAwesomeHubEndpoints(svc interface{}) FakeDataProviderEndpoints {
//	return SomeAwesomeHubEndpoints {
//		MethodOneEndpoint: makeMethodOneEndpoint(svc),
//		MethodTwoEndpoint: makeMethodTwoEndpoint(svc),
//		MethodThreeEndpoint: makeMethodThreeEndpoint(svc),
//	}
//}`
//	filename := "testdata/TestFixMissingEndpoints__endpoint.tmp"
//	err := ioutil.WriteFile(filename, []byte(initContent), os.ModePerm)
//	if nil != err {
//		t.Fatal(err)
//	}
//
//	FixMissingEndpoints(
//		filename, "SomeAwesomeHub", []string{
//			"MethodOne",
//			"MethodTwo",
//			"MethodThree",
//			"SayHello",
//		},
//	)
//
//	result, err := ioutil.ReadFile(filename)
//	if nil != err {
//		t.Fatal(err)
//	}
//
//	if !strings.Contains(string(result), "TODO:") {
//		t.Fatal("Missing comments")
//	}
//
//	if !strings.Contains(string(result), "SayHelloEndpoint	endpoint.Endpoint") {
//		t.Fatal("Missing required endpoint property declaration")
//	}
//
//	if !strings.Contains(string(result), "SayHelloEndpoint:") {
//		t.Fatal("Missing required endpoint property initialization")
//	}
//}
//
//func Test_extractEndpointMakeMethod_findMissedEndpointsInitializations(t *testing.T) {
//	initContent, err := ioutil.ReadFile("testdata/endpoint.go.sample.txt")
//	if nil != err {
//		t.Fatal(err)
//	}
//
//	filename := "testdata/Test_extractEndpointMakeMethod.tmp"
//	err = ioutil.WriteFile(filename, []byte(initContent), os.ModePerm)
//	if nil != err {
//		t.Fatal(err)
//	}
//
//	fs := token.NewFileSet()
//	fileGo, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
//	if nil != err {
//		log.Fatal(err)
//	}
//
//	makeMethod := extractEndpointMakeMethod(fileGo, "SomeAwesomeHub")
//	if nil == makeMethod {
//		t.Fatal("Make endpoint method not found")
//	}
//
//	initializations := findMissedEndpointsInitializations(makeMethod, []string{
//		"MethodOne",
//		"MethodTwo",
//		"MethodThree",
//		"SayHello",
//	})
//	if nil == initializations || len(initializations) == 0 || initializations[0] != "SayHello" {
//		t.Fatal("Dont show missed service method endpoint init")
//	}
//
//	fillMissedEndpointInits(makeMethod, initializations)
//
//	writeSourceFile(filename, fileGo, fs)
//
//	result, err := ioutil.ReadFile(filename)
//	if nil != err {
//		t.Fatal(err)
//	}
//
//	if !strings.Contains(string(result), "SayHelloEndpoint:") {
//		t.Fatal("Endpoint prop init is missed")
//	}
//}

func TestEndpointFixer_addMissedRequestModels(t *testing.T) {
	//initContent, err := ioutil.ReadFile("testdata/endpoint.go.sample.txt")
	//if nil != err {
	//	t.Fatal(err)
	//}
	//
	//err = ioutil.WriteFile(ef.filename, []byte(initContent), os.ModePerm)
	//if nil != err {
	//	t.Fatal(err)
	//}

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
	//initContent, err := ioutil.ReadFile("testdata/endpoint.go.sample.txt")
	//if nil != err {
	//	t.Fatal(err)
	//}
	//
	//err = ioutil.WriteFile(ef.filename, []byte(initContent), os.ModePerm)
	//if nil != err {
	//	t.Fatal(err)
	//}

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
	//initContent, err := ioutil.ReadFile("testdata/endpoint.go.sample.txt")
	//if nil != err {
	//	t.Fatal(err)
	//}
	//
	//err = ioutil.WriteFile(ef.filename, []byte(initContent), os.ModePerm)
	//if nil != err {
	//	t.Fatal(err)
	//}

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
	//initContent, err := ioutil.ReadFile("testdata/endpoint.go.sample.txt")
	//if nil != err {
	//	t.Fatal(err)
	//}
	//
	//filename := "testdata/TestEndpointFixer_addMissedPropertyInitializationInMakeEndpointsFunc.tmp"
	//err = ioutil.WriteFile(filename, []byte(initContent), os.ModePerm)
	//if nil != err {
	//	t.Fatal(err)
	//}

	//ef := EndpointFixer{
	//	filename:    filename,
	//	serviceName: "SomeAwesomeHub",
	//	serviceActions: []string{
	//		"MethodOne",
	//		"MethodTwo",
	//		"MethodThree",
	//		"SayHello",
	//	},
	//}
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
