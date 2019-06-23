package main

import (
	"flag"
	"fmt"
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate"
	"go/ast"
	"go/parser"
	"go/token"
	//"io/ioutil"
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

	boilerplate.InitProjectDirs(workingDirPath)

	protoGenFilename := workingDirPath + "/" + strings.TrimSuffix(*protoFile, "proto") + "pb.go"
	interfaces := boilerplate.ExtractServerInterface(protoGenFilename, *serviceName)

	//fmt.Println(itf.Name.Name)
	var methodNames []string
	methods := extractMethodsFromType(interfaces)
	for _, f := range methods {
		methodNames = append(methodNames, f.Names[0].Name)
	}

	serviceName := strings.TrimSuffix(interfaces.Name.Name, "Server")
	serviceNameLow := strings.ToLower(serviceName[0:1]) + serviceName[1:]

	var protoFileRelDir string
	{
		protoFileRelDirParts := strings.Split(*protoFile, "/")
		protoFileRelDirParts = protoFileRelDirParts[:len(protoFileRelDirParts)-1]
		protoFileRelDir = strings.Join(protoFileRelDirParts, "/")
	}

	data := struct {
		ProtoPackageName string
		ServiceName      string
		ServiceNameLow   string
		ProtoFileRelDir  string
		MethodNames      []string
	}{
		ProtoPackageName: "pb", // FIXME: hardcode
		ServiceName:      serviceName,
		ServiceNameLow:   serviceNameLow,
		ProtoFileRelDir:  protoFileRelDir,
		MethodNames:      methodNames,
	}

	{
		endpointFilename := (*workingDir) + "/endpoint/endpoint.go"
		if _, err := os.Stat(endpointFilename); os.IsNotExist(err) {
			boilerplate.CreateNewFromTemplate(endpointFilename, "template/endpoint.tmpl", data)
		} else {
			//fmt.Println("File", endpointFilename, "already exists! Please edit it manually!")
			//fmt.Println(findMissingReqRespForEndpoint(endpointFilename, methodNames))

			f, err := os.OpenFile(endpointFilename, os.O_APPEND|os.O_WRONLY, 0644)
			if nil != err {
				log.Fatal(err)
			}
			for _, typeName := range findMissingReqRespForEndpoint(endpointFilename, methodNames) {
				_, err := f.WriteString("type " + typeName + " struct {\n    // TODO: \n}\n")
				if nil != err {
					fmt.Println(err)
				}
			}
			fixMissingEndpoints(endpointFilename, serviceName, methodNames)
		}
	}

	{
		serviceFilename := workingDirPath + "/service/service.go"
		if _, err := os.Stat(serviceFilename); os.IsNotExist(err) {
			boilerplate.CreateNewFromTemplate(serviceFilename, "template/service.tmpl", data)
		} else {
			fmt.Println("File", serviceFilename, "already exists! Please edit it manually!")
		}
	}

	{
		modelFilename := workingDirPath + "/model/model.go"
		if _, err := os.Stat(modelFilename); os.IsNotExist(err) {
			boilerplate.CreateNewFromTemplate(modelFilename, "template/model.tmpl", data)
		} else {
			fmt.Println("File", modelFilename, "already exists! Please edit it manually!")
		}
	}

	{
		transportFilename := workingDirPath + "/transport/transport.go"
		if _, err := os.Stat(transportFilename); os.IsNotExist(err) {
			boilerplate.CreateNewFromTemplate(transportFilename, "template/transport.tmpl", data)
		} else {
			fmt.Println("File", transportFilename, "already exists! Please edit it manually!")
		}
	}

}

func extractMethodsFromType(t *ast.TypeSpec) []*ast.Field {
	return t.Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
}

func findMissingReqRespForEndpoint(filename string, expectedMethods []string) []string {
	var res []string

	fs := token.NewFileSet()
	fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if nil != err {
		log.Fatal(err)
	}

	for _, m := range expectedMethods {
		findReq := false
		findRes := false
		for _, d := range fileGo.Decls {
			if gen, ok := d.(*ast.GenDecl); ok {
				if len(gen.Specs) > 0 {
					if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
						if "type" != gen.Tok.String() {
							continue
						}
						if m+"Request" == typeSpec.Name.Name {
							findReq = true
							continue
						}
						if m+"Response" == typeSpec.Name.Name {
							findRes = true
							continue
						}
					}
				}
			}
		}
		if !findReq {
			res = append(res, m+"Request")
		}
		if !findRes {
			res = append(res, m+"Response")
		}
	}
	return res
}

func fixMissingEndpoints(filename, serviceName string, expectedMethods []string) {
	fs := token.NewFileSet()
	fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if nil != err {
		log.Fatal(err)
	}

	var endpointsStruct *ast.GenDecl
	for _, d := range fileGo.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			if len(gen.Specs) > 0 {
				if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
					if "type" != gen.Tok.String() {
						continue
					}
					if serviceName+"Endpoints" == typeSpec.Name.Name {
						endpointsStruct = gen
						continue
					}
				}
			}
		}
	}

	if nil == endpointsStruct {
		return
	}

	//for _, m := range expectedMethods {
	//	exists := false
	//	for _, field := range endpointsStruct.Specs[0].(*ast.TypeSpec).Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
	//		if m+"Endpoint" == field.Names[0].Name {
	//			exists = true
	//			break
	//		}
	//	}
	//
	//	if !exists {
	//		fs.
	//		endpointsStruct.Rparen
	//	}
	//}

}
