package fixer

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFixMissingEndpoints(t *testing.T) {
	initContent := `package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

// TODO:

type SomeAwesomeHubEndpoints struct {
	MethodOneEndpoint endpoint.Endpoint
	MethodTwoEndpoint endpoint.Endpoint
	MethodThreeEndpoint endpoint.Endpoint
}

func MakeFakeDataProviderEndpoints(svc interface{}) FakeDataProviderEndpoints {
	return FakeDataProviderEndpoints{
		MethodOneEndpoint: makeMethodOneEndpoint(svc),
		MethodTwoEndpoint: makeMethodTwoEndpoint(svc),
		MethodThreeEndpoint: makeMethodThreeEndpoint(svc),
	}
}`
	filename := "testdata/TestFixMissingEndpoints__endpoint.tmp"
	err := ioutil.WriteFile(filename, []byte(initContent), os.ModePerm)
	if nil != err {
		t.Fatal(err)
	}

	FixMissingEndpoints(
		filename, "SomeAwesomeHub", []string{
			"MethodOne",
			"MethodTwo",
			"MethodThree",
			"SayHello",
		},
	)

	result, err := ioutil.ReadFile(filename)
	if nil != err {
		t.Fatal(err)
	}

	if !strings.Contains(string(result), "TODO:") {
		t.Fatal("Missing comments")
	}

	if !strings.Contains(string(result), "SayHelloEndpoint	endpoint.Endpoint") {
		t.Fatal("Missing required endpoint property declaration")
	}

	if !strings.Contains(string(result), "SayHelloEndpoint:") {
		t.Fatal("Missing required endpoint property initialization")
	}
}
