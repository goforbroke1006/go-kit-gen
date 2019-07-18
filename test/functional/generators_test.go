package functional

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/goforbroke1006/go-kit-gen/pkg/generator"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
)

const serviceName = "SomeAwesomeHub"

func TestServiceFixer(t *testing.T) {
	serviceFilename, _ := filepath.Abs("../testdata/pkg/service/service.go")

	if err := os.Remove(serviceFilename); nil != err {
		t.Log(err)
	}

	fileSet := token.NewFileSet()
	serviceFileNode, err := parser.ParseFile(fileSet, serviceFilename, nil, parser.ParseComments)
	if err != nil {
		serviceFileNode = &ast.File{}
	}

	crawler := source.NewFileCrawler(serviceFileNode)
	crawler.SetPackageIfNotDefined("service")

	srvGen := generator.NewServiceInterfaceGenerator(crawler)
	_, err = srvGen.CreateInterfaceIfNotExists(serviceName)

	_ = srvGen.CreateMethodSignatureIfNotExists(serviceName, "MethodOne")
	_ = srvGen.CreateMethodSignatureIfNotExists(serviceName, "MethodTwo")
	_ = srvGen.CreateMethodSignatureIfNotExists(serviceName, "MethodThree")
	_ = srvGen.CreateMethodSignatureIfNotExists(serviceName, "SayHello")

	srvImplGen := generator.NewServicePrivateImplGenerator(crawler)
	_, err = srvImplGen.CreateEmptyStructIfNotExists(serviceName)

	srvImplGen.CreateMethodDeclIfNotExists(serviceName, "MethodOne")
	srvImplGen.CreateMethodDeclIfNotExists(serviceName, "MethodTwo")
	srvImplGen.CreateMethodDeclIfNotExists(serviceName, "MethodThree")
	srvImplGen.CreateMethodDeclIfNotExists(serviceName, "SayHello")

	generator.CreateServiceConstructor(crawler, serviceName)

	if file, err := os.OpenFile(serviceFilename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		t.Fatal(err.Error())
	} else {
		if err = printer.Fprint(file, fileSet, serviceFileNode); nil != err {
			t.Fatal(err)
		}
	}

	result, err := ioutil.ReadFile(serviceFilename)
	if nil != err {
		t.Fatal(err)
	}

	matches := []string{
		`type SomeAwesomeHubService interface`,
		`MethodOne\([\w]+ context.Context`,
		`MethodTwo\([\w]+ context.Context`,
		`MethodThree\([\w]+ context.Context`,
		`SayHello\([\w]+ context.Context`,

		`type someAwesomeHubService struct`,
		`func \([\w]+ someAwesomeHubService\) MethodOne\(ctx context.Context`,
		`func \([\w]+ someAwesomeHubService\) MethodTwo\(ctx context.Context`,
		`func \([\w]+ someAwesomeHubService\) MethodThree\(ctx context.Context`,
		`func \([\w]+ someAwesomeHubService\) SayHello\(ctx context.Context`,

		`func NewSomeAwesomeHubService\(\) SomeAwesomeHubService {`,
		`return &someAwesomeHubService{}`,
	}

	for _, m := range matches {
		re := regexp.MustCompile(m)
		match := re.FindStringSubmatch(string(result))
		if len(match) > 0 {
			t.Log(fmt.Sprintf("find substring : %s", match[0]))
		} else {
			t.Fatalf("Can't find any row by regex '%s'", m)
		}
	}

}

func TestEndpointFixer(t *testing.T) {
	endpointFilename, _ := filepath.Abs("../testdata/pkg/endpoint/endpoint.go")

	if err := os.Remove(endpointFilename); nil != err {
		t.Log(err)
	}

	fileSet := token.NewFileSet()
	fileNode, err := parser.ParseFile(fileSet, endpointFilename, nil, parser.ParseComments)
	if err != nil {
		//log.Fatal(err)
		fileNode = &ast.File{}
	}

	crawler := source.NewFileCrawler(fileNode)
	crawler.SetPackageIfNotDefined("endpoint")
	crawler.AddImportIfNotExists("../service", "")
	crawler.AddImportIfNotExists("github.com/go-kit/kit/endpoint", "goKitEndpoint")

	// TODO: magic here
	eptsStructGen := generator.NewEndpointsStructGenerator(crawler)
	eptsStructGen.CreateEndpointStructIfNotExists(serviceName)

	_ = eptsStructGen.CreateEndpointStructField(serviceName, "MethodOne")
	_ = eptsStructGen.CreateEndpointStructField(serviceName, "MethodTwo")
	_ = eptsStructGen.CreateEndpointStructField(serviceName, "MethodThree")
	_ = eptsStructGen.CreateEndpointStructField(serviceName, "SayHello")

	eptsStructGen.CreateRequestStruct("MethodOne")
	eptsStructGen.CreateResponseStruct("MethodOne")
	eptsStructGen.CreateRequestStruct("MethodTwo")
	eptsStructGen.CreateResponseStruct("MethodTwo")
	eptsStructGen.CreateRequestStruct("MethodThree")
	eptsStructGen.CreateResponseStruct("MethodThree")
	eptsStructGen.CreateRequestStruct("SayHello")
	eptsStructGen.CreateResponseStruct("SayHello")

	eptsStructGen.CreateConstructorIfNotExists(serviceName)
	eptsStructGen.SetFieldInConstructor(serviceName, "MethodOne")
	eptsStructGen.SetFieldInConstructor(serviceName, "MethodTwo")
	eptsStructGen.SetFieldInConstructor(serviceName, "MethodThree")
	eptsStructGen.SetFieldInConstructor(serviceName, "SayHello")

	eptsStructGen.CreateMakeEndpointFunc(serviceName, "MethodOne")
	eptsStructGen.CreateMakeEndpointFunc(serviceName, "MethodTwo")
	eptsStructGen.CreateMakeEndpointFunc(serviceName, "MethodThree")
	eptsStructGen.CreateMakeEndpointFunc(serviceName, "SayHello")

	if file, err := os.OpenFile(endpointFilename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		t.Fatal(err.Error())
	} else {
		if err = printer.Fprint(file, fileSet, fileNode); nil != err {
			t.Fatal(err)
		}
	}

	result, err := ioutil.ReadFile(endpointFilename)
	if nil != err {
		t.Fatal(err)
	}

	matches := []string{
		`MethodOneEndpoint([\s]+)goKitEndpoint.Endpoint`,
		`MethodTwoEndpoint([\s]+)goKitEndpoint.Endpoint`,
		`MethodThreeEndpoint([\s]+)goKitEndpoint.Endpoint`,
		`SayHelloEndpoint([\s]+)goKitEndpoint.Endpoint`,

		`type MethodOneRequest struct`,
		`type MethodTwoRequest struct`,
		`type MethodThreeRequest struct`,
		`type SayHelloRequest struct`,

		`type MethodOneResponse struct`,
		`type MethodTwoResponse struct`,
		`type MethodThreeResponse struct`,
		`type SayHelloResponse struct`,

		`func NewSomeAwesomeHubEndpoints\(svc service.SomeAwesomeHubService\) SomeAwesomeHubEndpoints`,

		`MethodOneEndpoint\:([\t\s]+)makeMethodOneEndpoint\(svc\)`,
		`MethodTwoEndpoint\:([\t\s]+)makeMethodTwoEndpoint\(svc\)`,
		`MethodThreeEndpoint\:([\t\s]+)makeMethodThreeEndpoint\(svc\)`,
		`SayHelloEndpoint\:([\t\s]+)makeSayHelloEndpoint\(svc\)`,

		`func makeMethodOneEndpoint\(`,
		`func makeMethodTwoEndpoint\(`,
		`func makeMethodThreeEndpoint\(`,
		`func makeSayHelloEndpoint\(`,
	}

	for _, m := range matches {
		re := regexp.MustCompile(m)
		match := re.FindStringSubmatch(string(result))
		if len(match) > 0 {
			fmt.Println(match[0])
		} else {
			t.Fatalf("Can't find any row by regex '%s'", m)
		}
	}

}

func TestGRPCTransportFixer(t *testing.T) {
	transportFilename, _ := filepath.Abs("../testdata/pkg/transport/transport_grpc.go")

	//if err := os.Remove(transportFilename); nil != err {
	//	t.Log(err)
	//}

	fileSet := token.NewFileSet()
	serviceFileNode, err := parser.ParseFile(fileSet, transportFilename, nil, parser.ParseComments)
	if err != nil {
		serviceFileNode = &ast.File{}
	}

	crawler := source.NewFileCrawler(serviceFileNode)
	crawler.SetPackageIfNotDefined("transport")
	crawler.AddImportIfNotExists("context", "")
	crawler.AddImportIfNotExists("github.com/go-kit/kit/transport/grpc", "")

	transportGRPCGen := generator.NewTransportGRPCGenerator(crawler)
	transportGRPCGen.CreateServerImpleStructIfNotExists(serviceName)

	_ = transportGRPCGen.CreateField(serviceName, "MethodOne")
	_ = transportGRPCGen.CreateField(serviceName, "MethodTwo")
	_ = transportGRPCGen.CreateField(serviceName, "MethodThree")
	_ = transportGRPCGen.CreateField(serviceName, "SayHello")

	transportGRPCGen.CreateMethodDecl(serviceName, "pb", "MethodOne")
	transportGRPCGen.CreateMethodDecl(serviceName, "pb", "MethodTwo")
	transportGRPCGen.CreateMethodDecl(serviceName, "pb", "MethodThree")
	transportGRPCGen.CreateMethodDecl(serviceName, "pb", "SayHello")

	transportGRPCGen.CreateServerConstructorIfNotExists()

	transportGRPCGen.AddFieldInitInConstrucor("MethodOne")
	transportGRPCGen.AddFieldInitInConstrucor("MethodTwo")
	transportGRPCGen.AddFieldInitInConstrucor("MethodThree")
	transportGRPCGen.AddFieldInitInConstrucor("SayHello")

	if file, err := os.OpenFile(transportFilename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		t.Fatal(err.Error())
	} else {
		if err = printer.Fprint(file, fileSet, serviceFileNode); nil != err {
			t.Fatal(err)
		}
	}

	result, err := ioutil.ReadFile(transportFilename)
	if nil != err {
		t.Fatal(err)
	}

	matches := []string{
		`type someAwesomeHubGRPCServer struct`,

		`methodOne[\s]+grpc.Handler`,
		`methodTwo[\s]+grpc.Handler`,
		`methodThree[\s]+grpc.Handler`,
		`sayHello[\s]+grpc.Handler`,

		`func \(s someAwesomeHubGRPCServer\) MethodOne\(ctx context.Context, req \*pb.MethodOneRequest\) \(\*pb.MethodOneResponse, error\)`,
		`func \(s someAwesomeHubGRPCServer\) MethodTwo\(ctx context.Context, req \*pb.MethodTwoRequest\) \(\*pb.MethodTwoResponse, error\)`,
		`func \(s someAwesomeHubGRPCServer\) MethodThree\(ctx context.Context, req \*pb.MethodThreeRequest\) \(\*pb.MethodThreeResponse, error\)`,
		`func \(s someAwesomeHubGRPCServer\) SayHello\(ctx context.Context, req \*pb.SayHelloRequest\) \(\*pb.SayHelloResponse, error\)`,

		`func NewSomeAwesomeHubGRPCServer\([\w]+ context.Context, endpoint endpoint.Endpoints\) [\w]+.FakeDataProviderServer`,

		`methodOne: grpc.NewServer\(`,
		`methodTwo: grpc.NewServer\(`,
		`methodThree: grpc.NewServer\(`,
		`sayHello: grpc.NewServer\(`,
	}

	for _, m := range matches {
		re := regexp.MustCompile(m)
		match := re.FindStringSubmatch(string(result))
		if len(match) > 0 {
			fmt.Println(match[0])
		} else {
			t.Fatalf("Can't find any row by regex '%s'", m)
		}
	}
}
