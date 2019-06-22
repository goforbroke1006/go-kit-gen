package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	workingDir = flag.String("working-dir", "./", "Define project root dir")
	protoPath  = flag.String("proto-path", "", "Services and messages blueprints *.proto file location")
	protoFile  = flag.String("proto-file", "", "Services and messages blueprints *.proto file location")
)

func init() {
	flag.Parse()
}

func main() {

	command := exec.Command(
		"protoc",
		fmt.Sprintf("--proto_path=%s", *protoPath),
		fmt.Sprintf("--go_out=plugins=grpc:%s", *workingDir),
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

	{
		command = exec.Command("mkdir", (*workingDir)+"/endpoint/")
		command.Start()
		command.Wait()
	}

	{
		command = exec.Command("mkdir", (*workingDir)+"/service/")
		command.Start()
		command.Wait()
	}

	endpointTmpl, err := template.ParseFiles((*workingDir) + "/../template/endpoint.tmpl")
	if nil != err {
		log.Fatal(err)
	}

	serviceTmpl, err := template.ParseFiles((*workingDir) + "/../template/service.tmpl")
	if nil != err {
		log.Fatal(err)
	}

	interfaces := parseServerInterfaces(
		"/home/goforbroke/go/src/github.com/goforbroke1006/go-kit-gen/debug-wd/pb/api/v1/fake-data-provider-service.pb.go",
	)
	for _, itf := range interfaces {
		//fmt.Println(itf.Name.Name)
		var methodNames []string
		methods := extractMethodsFromType(itf)
		for _, f := range methods {
			//fmt.Println("", ">", f.Names[0].Name)
			if nil != err {
				fmt.Println(err.Error())
				continue
			}
			methodNames = append(methodNames, f.Names[0].Name)
		}

		serviceName := strings.TrimSuffix(itf.Name.Name, "Server")
		serviceNameLow := strings.ToLower(serviceName[0:1]) + serviceName[1:]

		data := struct {
			ServiceName    string
			ServiceNameLow string
			MethodNames    []string
		}{
			ServiceName:    serviceName,
			ServiceNameLow: serviceNameLow,
			MethodNames:    methodNames,
		}

		endpointFile, _ := os.Create((*workingDir) + "/endpoint/" + serviceNameLow + ".go")
		err = endpointTmpl.Execute(endpointFile, data)
		if nil != err {
			fmt.Println(err.Error())
		}

		serviceFile, _ := os.Create((*workingDir) + "/service/" + serviceNameLow + ".go")
		err = serviceTmpl.Execute(serviceFile, data)
		if nil != err {
			fmt.Println(err.Error())
		}

	}

}

func parseServerInterfaces(filename string) []*ast.TypeSpec {
	var result []*ast.TypeSpec

	fs := token.NewFileSet()
	fileGo, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if nil != err {
		log.Fatal(err)
	}

	for _, d := range fileGo.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			if len(gen.Specs) > 0 {
				if typeSpec, ok := gen.Specs[0].(*ast.TypeSpec); ok {
					if "type" != gen.Tok.String() {
						continue
					}
					if strings.HasPrefix(typeSpec.Name.Name, "Unimplemented") {
						continue
					}
					if !strings.HasSuffix(typeSpec.Name.Name, "Server") {
						continue
					}

					{
						s := typeSpec.Name.Name[0:1]
						upper := strings.ToUpper(s)
						if s != upper {
							continue
						}
					}

					result = append(result, typeSpec)
				}
			}
			continue
		}
	}

	return result
}

func extractMethodsFromType(t *ast.TypeSpec) []*ast.Field {
	return t.Name.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
}
