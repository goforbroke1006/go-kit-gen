package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/naming"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/project"
)

var (
	argWorkingDir    = flag.String("working-dir", "./", "Define project root dir")
	argProtoResFile  = flag.String("proto-res-file", "", "*.pb.go file location")
	argServiceName   = flag.String("service-name", "", "Service name")
	argTransportType = flag.String("transport-type", "", "Select transport type (grpc, http)")

	argEndpointFile  = flag.String("file-endpoint", "./endpoint/endpoint.go", "Endpoint file location related to selected working dir")
	argServiceFile   = flag.String("file-service", "./service/service.go", "Service file location related to selected working dir")
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

	//if len(*argProtoPath) == 0 {
	//	if *argProtoPath, err = filepath.Abs(*argProtoPath); nil != err {
	//		log.Fatal(err)
	//	}
	//} else {
	//	*argProtoPath = *argWorkingDir
	//}

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

	//output, err := shell.Execute("protoc", "--proto_path=%s"+*argProtoPath, "--go_out=plugins=grpc:%s"+*argWorkingDir, *argProtoFile, )
	//if nil != err {
	//	log.Fatalln(err.Error())
	//}
	//log.Println(output)

	project.InitProjectDirs(*argWorkingDir)

	//protoGenFilename := *argWorkingDir + "/" + strings.TrimSuffix(*argProtoFile, "proto") + "pb.go"
	_, pbGoFile := fixer.OpenGolangSourceFile(*argProtoResFile)
	protoPbResAfi := iterator.NewAstFileIterator(pbGoFile)

	serviceInterfaceName := naming.GetServiceInterfaceInPb(*argServiceName)
	serverInterface := protoPbResAfi.GetInterfaceTypeSpec(serviceInterfaceName)

	serverInterfaceAiti := iterator.NewAstInterfaceTypeIterator(serverInterface)

	//fmt.Println(itf.Name.Name)
	//var methodNames = map[string]map[string]string{}
	//methods := old.ExtractMethodsFromType(interfaces)
	//for _, f := range methods {
	//	methodNames[f.Names[0].Name] = nil
	//}

	//_, pbGoFile := fixer.OpenGolangSourceFile(protoGenFilename)
	//protoPbResAFI := iterator.NewAstFileIterator(pbGoFile)

	//{
	//	serviceFilename := *argWorkingDir + "/service/service.go"
	//	if _, err := os.Stat(serviceFilename); os.IsNotExist(err) {
	//		old.CreateNewFromTemplate(serviceFilename, "template/service.tmpl", data)
	//	} else {
	//		//fmt.Println("File", serviceFilename, "already exists! Please edit it manually!")
	//		serviceFixer := service.NewServiceFixer(nil, *argServiceName, methodNames)
	//		serviceFixer.Fix()
	//	}
	//}
	//
	//{
	//	endpointFilename := *argWorkingDir + "/endpoint/endpoint.go"
	//	if _, err := os.Stat(endpointFilename); os.IsNotExist(err) {
	//		old.CreateNewFromTemplate(endpointFilename, "template/endpoint.tmpl", data)
	//	} else {
	//		endpointFixer := endpoint.NewEndpointFixer(nil, *argServiceName, methodNames)
	//		endpointFixer.Fix()
	//	}
	//}
	//
	//{
	//	modelFilename := *argWorkingDir + "/model/model.go"
	//	if _, err := os.Stat(modelFilename); os.IsNotExist(err) {
	//		old.CreateNewFromTemplate(modelFilename, "template/model.tmpl", data)
	//	} else {
	//		fmt.Println("File", modelFilename, "already exists! Please edit it manually!")
	//	}
	//}

	//{
	//	transportFilename := *argWorkingDir + "/transport/transport.go"
	//	if _, err := os.Stat(transportFilename); os.IsNotExist(err) {
	//		old.CreateNewFromTemplate(transportFilename, "template/transport.tmpl", data)
	//	} else {
	//		fmt.Println("File", transportFilename, "already exists! Please edit it manually!")
	//	}
	//}

	/*{
		if _, err := os.Stat(*argTransportFile); os.IsNotExist(err) {
			if _, err := os.OpenFile(*argTransportFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660); nil != err {
				log.Fatal(err.Error())
			}
			file := &ast.File{}
			builder.NewAstFileBuilder(file).SetPackage("transport") // TODO: hard-code
			fileSet := token.NewFileSet()
			fixer.WriteSourceFile(*argTransportFile, file, fileSet)
		}

		fileSet, file := fixer.OpenGolangSourceFile(*argTransportFile)
		switch *argTransportType {
		case TransportTypeGRPC:

			transportFileAfi := iterator.NewAstFileIterator(file)
			transportFileAfb := builder.NewAstFileBuilder(file)

			if !transportFileAfi.HasImport("context") {
				transportFileAfb.AddImport("context") // TODO: hard-code
			}

			if !transportFileAfi.HasImport("github.com/go-kit/kit/transport/grpc") {
				transportFileAfb.AddImport("github.com/go-kit/kit/transport/grpc") // TODO: hard-code
			}

			transportFixer := transport.NewGRPCTransportFixer(file, *argServiceName)
			for _, mArgs := range serverInterfaceAiti.GetMethodsFieldList().List {
				actionName := mArgs.Names[0].Name
				if name, err := transportFixer.FixServerImplStructField(actionName); nil != err {
					log.Println("warn", err.Error())
				} else {
					log.Println("info", "field ", name, "was created")
				}

				argFields := mArgs.Type.(*ast.FuncType).Params.List
				argFields[1].Type.(*ast.StarExpr).X.(*ast.Ident).Name = protoPbResAfi.GetPackageName() + "." + argFields[1].Type.(*ast.StarExpr).X.(*ast.Ident).Name

				retFields := mArgs.Type.(*ast.FuncType).Results.List
				retFields[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name = protoPbResAfi.GetPackageName() + "." + retFields[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name

				if err := transportFixer.FixServerImplStructMethod(actionName, argFields, retFields); nil != err {
					log.Println("warn", err.Error())
				}

				if err := transportFixer.FixDecodeMethod(actionName, argFields, retFields); nil != err {
					log.Println("warn", err.Error())
				}
			}
			// TODO;
			break
		case TransportTypeHTTP:
			// TODO:
			break
		default:
			log.Fatalln("unexpected transport type, can't find fixer for", *argTransportType)
		}
		fixer.WriteSourceFile(*argTransportFile, file, fileSet)
	}*/

}
