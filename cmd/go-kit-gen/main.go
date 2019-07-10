package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"
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
	//
}
