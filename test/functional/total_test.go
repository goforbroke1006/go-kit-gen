package functional

import (
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/generator/service"
	"io/ioutil"
	"regexp"
	"testing"
)

const serviceName = "SomeAwesomeHub"
const pbGoFilename = "../testdata/some-awesome-hub.pb.go"

func TestServiceFixer(t *testing.T) {

	const serviceFilename = "../testdata/pkg/service/service.go"

	iterator.NewAstInterfaceTypeIterator()

	// TODO: call all fixers
	serviceInterfaceGenerator := service.NewServiceInterfaceGenerator()
	serviceInterfaceGenerator.CreateInterfaceIfNotExists(serviceName)
	serviceInterfaceGenerator.CreateMethodSignatureIfNotExists(serviceName)

	result, err := ioutil.ReadFile(serviceFilename)
	if nil != err {
		t.Fatal(err)
	}

	matches := []string{
		"type SomeAwesomeHub interface",
		"MethodOne\\([\\w]+ context.Context",
	}

	for _, m := range matches {
		re := regexp.MustCompile(m)
		match := re.FindStringSubmatch(string(result))
		if len(match) > 0 {
			fmt.Println(match[1])
		} else {
			t.Fatalf("Can't find any row by regex '%s'", m)
		}
	}

}

//func TestEndpointFixer(t *testing.T) {
//
//	const endpointFilename = "../testdata/pkg/endpoint/endpoint.go"
//
//	// TODO: call all fixers
//
//	result, err := ioutil.ReadFile(endpointFilename)
//	if nil != err {
//		t.Fatal(err)
//	}
//
//	if !strings.Contains(string(result), "type "+action+"Request struct") {
//		t.Fatal("" + action + "Request struct is missed")
//	}
//
//}
