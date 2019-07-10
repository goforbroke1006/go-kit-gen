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
	"github.com/goforbroke1006/go-kit-gen/pkg/generator/service"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
)

const serviceName = "SomeAwesomeHub"

func TestServiceFixer(t *testing.T) {
	serviceFilename, _ := filepath.Abs("../testdata/pkg/service/service.go")

	if err := os.Remove(serviceFilename); nil != err {
		t.Fatal(err)
	}

	fileSet := token.NewFileSet()
	serviceFileNode, err := parser.ParseFile(fileSet, serviceFilename, nil, parser.ParseComments)
	if err != nil {
		//log.Fatal(err)
		serviceFileNode = &ast.File{}
	}

	endpointSourceCrawler := source.NewFileCrawlerByAF(serviceFileNode)
	endpointSourceCrawler.SetPackageIfNotDefined("service")

	serviceInterfaceGenerator := service.NewServiceInterfaceGenerator(endpointSourceCrawler)
	serviceInterfaceName, err := serviceInterfaceGenerator.CreateInterfaceIfNotExists(serviceName)
	{
		serviceInterfaceGenerator.CreateMethodSignatureIfNotExists(
			"MethodOne",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
		serviceInterfaceGenerator.CreateMethodSignatureIfNotExists(
			"MethodTwo",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
		serviceInterfaceGenerator.CreateMethodSignatureIfNotExists(
			"MethodThree",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
		serviceInterfaceGenerator.CreateMethodSignatureIfNotExists(
			"SayHello",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
	}

	servicePrivateImplGenerator := service.NewServicePrivateImplGenerator()
	servImplName, err := servicePrivateImplGenerator.CreateEmptyStructIfNotExists(serviceName)
	{
		servicePrivateImplGenerator.CreateMethodDeclIfNotExists(
			"MethodOne",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
		servicePrivateImplGenerator.CreateMethodDeclIfNotExists(
			"MethodTwo",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
		servicePrivateImplGenerator.CreateMethodDeclIfNotExists(
			"MethodThree",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
		servicePrivateImplGenerator.CreateMethodDeclIfNotExists(
			"SayHello",
			[][2]string{{"ctx", "context.Context"}, {"agr1", "interface{}"}},
			[][2]string{{"", "interface{}"}, {"", "error"}},
		)
	}

	constructorGenerator := generator.NewConstructorGenerator()
	constructorGenerator.CreateServiceConstructorMethod(serviceInterfaceName, servImplName)

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
		"type SomeAwesomeHub interface",
		`MethodOne\([\w]+ context.Context`,
		`MethodTwo\([\w]+ context.Context`,
		`MethodThree\([\w]+ context.Context`,
		`SayHello\([\w]+ context.Context`,

		`type someAwesomeHub struct`,
		`func ([\w]+ someAwesomeHub) MethodOne(ctx context.Context`,
		`func ([\w]+ someAwesomeHub) MethodTwo(ctx context.Context`,
		`func ([\w]+ someAwesomeHub) MethodThree(ctx context.Context`,
		`func ([\w]+ someAwesomeHub) SayHello(ctx context.Context`,

		`func NewSomeAwesomeHub\(\) SomeAwesomeHub {`,
		`return &someAwesomeHub{}`,
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

func TestEndpointFixer(t *testing.T) {
	endpointFilename, _ := filepath.Abs("../testdata/pkg/endpoint/endpoint.go")

	if err := os.Remove(endpointFilename); nil != err {
		t.Log(err)
	}

	fileSet := token.NewFileSet()
	serviceFileNode, err := parser.ParseFile(fileSet, endpointFilename, nil, parser.ParseComments)
	if err != nil {
		//log.Fatal(err)
		serviceFileNode = &ast.File{}
	}

	endpointSourceCrawler := source.NewFileCrawlerByAF(serviceFileNode)
	endpointSourceCrawler.SetPackageIfNotDefined("endpoint")

	// TODO: magic here

	if file, err := os.OpenFile(endpointFilename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		t.Fatal(err.Error())
	} else {
		if err = printer.Fprint(file, fileSet, serviceFileNode); nil != err {
			t.Fatal(err)
		}
	}

	//result, err := ioutil.ReadFile(serviceFilename)
	//if nil != err {
	//	t.Fatal(err)
	//}

	// TODO: asserts here
}
