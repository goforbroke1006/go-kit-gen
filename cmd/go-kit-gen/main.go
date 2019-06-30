package main

import (
	"flag"
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/ast/iterator"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer/endpoint"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer/service"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/project"
	"github.com/goforbroke1006/go-kit-gen/pkg/old"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	workingDir  = flag.String("working-dir", "./", "Define project root dir")
	protoPath   = flag.String("proto-path", "", "Services and messages blueprints *.proto file location")
	protoFile   = flag.String("proto-file", "", "Services and messages blueprints *.proto file location")
	serviceName = flag.String("service-name", "", "Service name")
)

func init() {
	flag.Parse()
}

func main() {

	// TODO: validate args/opts

	var workingDirPath string
	if strings.HasPrefix(*workingDir, "/") {
		workingDirPath = *workingDir
	} else {
		appWorkingDir, err := os.Getwd()
		if nil != err {
			log.Fatal(err)
		}
		workingDirPath = appWorkingDir + "/" + *workingDir
	}

	if len(*protoPath) == 0 {
		dir, err := os.Getwd()
		if nil != err {
			log.Fatal(err)
		}
		*protoPath = dir
	}

	if len(*protoFile) == 0 {
		log.Fatal("--proto-file is required")
	}
	if len(*serviceName) == 0 {
		log.Fatal("--service-name is required")
	}

	command := exec.Command(
		"protoc",
		fmt.Sprintf("--proto_path=%s", *protoPath),
		fmt.Sprintf("--go_out=plugins=grpc:%s", workingDirPath),
		*protoFile,
	)
	err := command.Start()
	if nil != err {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = command.Wait()
	if nil != err {
		log.Printf("Command finished with error: %v", err)
	}

	project.InitProjectDirs(workingDirPath)

	protoGenFilename := workingDirPath + "/" + strings.TrimSuffix(*protoFile, "proto") + "pb.go"
	interfaces := old.ExtractServerInterface(protoGenFilename, *serviceName)

	//fmt.Println(itf.Name.Name)
	var methodNames = map[string]map[string]string{}
	methods := old.ExtractMethodsFromType(interfaces)
	for _, f := range methods {
		methodNames[f.Names[0].Name] = nil
	}

	//serviceName := strings.TrimSuffix(interfaces.Name.Name, "Server")
	serviceNameLow := strings.ToLower((*serviceName)[0:1]) + (*serviceName)[1:]

	var protoFileRelDir string
	{
		protoFileRelDirParts := strings.Split(*protoFile, "/")
		protoFileRelDirParts = protoFileRelDirParts[:len(protoFileRelDirParts)-1]
		protoFileRelDir = strings.Join(protoFileRelDirParts, "/")
	}

	_, file := fixer.OpenGolangSourceFile(protoGenFilename)
	protoPbResAFI := iterator.NewAstFileIterator(file)

	data := struct {
		ProtoPackageName string
		ServiceName      string
		ServiceNameLow   string
		ProtoFileRelDir  string
		MethodNames      map[string]map[string]string
	}{
		ProtoPackageName: protoPbResAFI.GetPackageName(),
		ServiceName:      *serviceName,
		ServiceNameLow:   serviceNameLow,
		ProtoFileRelDir:  protoFileRelDir,
		MethodNames:      methodNames,
	}

	{
		serviceFilename := workingDirPath + "/service/service.go"
		if _, err := os.Stat(serviceFilename); os.IsNotExist(err) {
			old.CreateNewFromTemplate(serviceFilename, "template/service.tmpl", data)
		} else {
			//fmt.Println("File", serviceFilename, "already exists! Please edit it manually!")
			serviceFixer := service.NewServiceFixer(nil, *serviceName, methodNames)
			serviceFixer.Fix()
		}
	}

	{
		endpointFilename := workingDirPath + "/endpoint/endpoint.go"
		if _, err := os.Stat(endpointFilename); os.IsNotExist(err) {
			old.CreateNewFromTemplate(endpointFilename, "template/endpoint.tmpl", data)
		} else {
			endpointFixer := endpoint.NewEndpointFixer(nil, *serviceName, methodNames)
			endpointFixer.Fix()
		}
	}

	{
		modelFilename := workingDirPath + "/model/model.go"
		if _, err := os.Stat(modelFilename); os.IsNotExist(err) {
			old.CreateNewFromTemplate(modelFilename, "template/model.tmpl", data)
		} else {
			fmt.Println("File", modelFilename, "already exists! Please edit it manually!")
		}
	}

	{
		transportFilename := workingDirPath + "/transport/transport.go"
		if _, err := os.Stat(transportFilename); os.IsNotExist(err) {
			old.CreateNewFromTemplate(transportFilename, "template/transport.tmpl", data)
		} else {
			fmt.Println("File", transportFilename, "already exists! Please edit it manually!")
		}
	}

}
