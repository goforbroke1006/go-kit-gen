package main

import (
	"flag"
	"github.com/goforbroke1006/go-kit-gen/pkg/generator"
	"github.com/goforbroke1006/go-kit-gen/pkg/source"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	argWorkingDir    = flag.String("working-dir", "./", "Define project root dir")
	argProtoResFile  = flag.String("proto-res-file", "", "*.pb.go file location")
	argServiceName   = flag.String("service-name", "", "Service name")
	argTransportType = flag.String("transport-type", "", "Select transport type (grpc, http)")

	argServiceFile   = flag.String("file-service", "./service/service.go", "Service file location related to selected working dir")
	argEndpointFile  = flag.String("file-endpoint", "./endpoint/endpoint.go", "Endpoint file location related to selected working dir")
	argTransportFile = flag.String("file-transport", "./transport/transport_<argTransportType>.go", "Transport file location related to selected working dir")
)

const (
	TransportTypeGRPC = "grpc"
	TransportTypeHTTP = "http"
)

func init() {
	flag.Parse()

	var err error

	if *argWorkingDir, err = filepath.Abs(*argWorkingDir); nil != err {
		log.Fatal(err)
	}

	if len(*argProtoResFile) == 0 {
		// TODO: check file exists
		log.Fatal("--proto-file is required")
	}
	if !strings.HasPrefix(*argProtoResFile, "/") {
		if *argProtoResFile, err = filepath.Abs(*argWorkingDir + "/" + *argProtoResFile); nil != err {
			log.Fatal(err.Error())
		}
	}

	if len(*argServiceName) == 0 {
		log.Fatal("--service-name is required")
	}

	if *argTransportType != TransportTypeGRPC && *argTransportType != TransportTypeHTTP {
		log.Fatal("unexpected --file-transport value")
	}

	if !strings.HasPrefix(*argEndpointFile, "/") {
		if *argEndpointFile, err = filepath.Abs(*argWorkingDir + "/" + *argEndpointFile); nil != err {
			log.Fatal(err.Error())
		}
	}

	if !strings.HasPrefix(*argServiceFile, "/") {
		if *argServiceFile, err = filepath.Abs(*argWorkingDir + "/" + *argServiceFile); nil != err {
			log.Fatal(err.Error())
		}
	}

	*argTransportFile = strings.Replace(*argTransportFile, "<argTransportType>", *argTransportType, -1)
	if !strings.HasPrefix(*argTransportFile, "/") {
		if *argTransportFile, err = filepath.Abs(*argWorkingDir + "/" + *argTransportFile); nil != err {
			log.Fatal(err.Error())
		}
	}
}

func main() {
	fileSet := token.NewFileSet()
	protoResFileNode, err := parser.ParseFile(fileSet, *argProtoResFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err.Error())
	}
	protoCrawler := source.NewFileCrawler(protoResFileNode)
	actions := protoCrawler.GetInterface(*argServiceName + "Server").GetMethods()

	var actionsNames []string
	for _, a := range actions {
		actionsNames = append(actionsNames, a.Names[0].Name)
	}

	fixServiceFile(*argServiceFile, *argServiceName, actionsNames)
	fixEndpointFile(*argEndpointFile, *argServiceName, actionsNames)

	if TransportTypeGRPC == *argTransportType {
		fixTransportGRPCFile(protoCrawler.GetPackage(),
			*argTransportFile, *argServiceName, actionsNames)
	}
}

func fixServiceFile(serviceFilename, serviceName string, actionNames []string) {
	fileSet := token.NewFileSet()
	serviceFileNode, err := parser.ParseFile(fileSet, serviceFilename, nil, parser.ParseComments)
	if err != nil {
		log.Println(err)
		serviceFileNode = &ast.File{}
	}

	crawler := source.NewFileCrawler(serviceFileNode)
	crawler.SetPackageIfNotDefined("service")

	srvGen := generator.NewServiceInterfaceGenerator(crawler)
	_, err = srvGen.CreateInterfaceIfNotExists(serviceName)
	if nil != err {
		log.Fatalln(err.Error())
	}

	for _, an := range actionNames {
		err = srvGen.CreateMethodSignatureIfNotExists(serviceName, an)
		if nil != err {
			log.Fatalln(err.Error())
		}
	}

	srvImplGen := generator.NewServicePrivateImplGenerator(crawler)
	_, err = srvImplGen.CreateEmptyStructIfNotExists(serviceName)
	if nil != err {
		log.Fatalln(err.Error())
	}

	for _, an := range actionNames {
		srvImplGen.CreateMethodDeclIfNotExists(serviceName, an)
	}

	generator.CreateServiceConstructor(crawler, serviceName)

	if file, err := os.OpenFile(serviceFilename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		log.Fatal(err.Error())
	} else {
		if err = printer.Fprint(file, fileSet, serviceFileNode); nil != err {
			log.Fatalln(err)
		}
	}
}

func fixEndpointFile(endpointFilename, serviceName string, actionNames []string) {
	fileSet := token.NewFileSet()
	fileNode, err := parser.ParseFile(fileSet, endpointFilename, nil, parser.ParseComments)
	if nil != err {
		log.Println(err)
		fileNode = &ast.File{}
	}

	crawler := source.NewFileCrawler(fileNode)
	crawler.SetPackageIfNotDefined("endpoint")
	crawler.AddImportIfNotExists("../service", "")
	crawler.AddImportIfNotExists("github.com/go-kit/kit/endpoint", "goKitEndpoint")

	structGen := generator.NewEndpointsStructGenerator(crawler)
	structGen.CreateEndpointStructIfNotExists(serviceName)

	for _, a := range actionNames {
		_ = structGen.CreateEndpointStructField(serviceName, a)
	}

	// create models for arguments
	for _, a := range actionNames {
		structGen.CreateRequestStruct(a)
		structGen.CreateResponseStruct(a)
	}

	structGen.CreateConstructorIfNotExists(serviceName)
	for _, a := range actionNames {
		structGen.SetFieldInConstructor(serviceName, a)
	}

	for _, a := range actionNames {
		structGen.CreateMakeEndpointFunc(serviceName, a)
	}

	if file, err := os.OpenFile(endpointFilename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		log.Fatal(err.Error())
	} else {
		if err = printer.Fprint(file, fileSet, fileNode); nil != err {
			log.Fatalln(err)
		}
	}
}

func fixTransportGRPCFile(pbGoPackage, filename, serviceName string, actionsNames []string) {
	fileSet := token.NewFileSet()
	fileNode, err := parser.ParseFile(fileSet, filename, nil, parser.ParseComments)
	if nil != err {
		log.Println(err)
		fileNode = &ast.File{}
	}

	crawler := source.NewFileCrawler(fileNode)
	crawler.SetPackageIfNotDefined("transport")
	crawler.AddImportIfNotExists("../endpoint", "")
	crawler.AddImportIfNotExists("context", "")
	crawler.AddImportIfNotExists("github.com/go-kit/kit/transport/grpc", "")

	transportGRPCGen := generator.NewTransportGRPCGenerator(crawler)
	transportGRPCGen.CreateServerImpleStructIfNotExists(serviceName)

	for _, a := range actionsNames {
		_ = transportGRPCGen.CreateField(serviceName, a)
	}

	for _, a := range actionsNames {
		transportGRPCGen.CreateMethodDecl(serviceName, pbGoPackage, a)
	}

	transportGRPCGen.CreateServerConstructorIfNotExists(serviceName, pbGoPackage)

	for _, a := range actionsNames {
		_ = transportGRPCGen.AddFieldInitInConstructor(serviceName, a)
	}

	for _, a := range actionsNames {
		transportGRPCGen.CreateDecodeRequestMethod(pbGoPackage, a)
		transportGRPCGen.CreateEncodeResponseMethod(pbGoPackage, a)
	}

	if file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666); nil != err {
		log.Fatal(err.Error())
	} else {
		if err = printer.Fprint(file, fileSet, fileNode); nil != err {
			log.Fatalln(err)
		}
	}
}
