package naming

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/string_util"
	"strings"
)

func GetServiceInterfaceName(serviceName string) string {
	return serviceName + "Service"
}

func GetServicePrivateImplStructName(serviceName string) string {
	return strings.ToLower(serviceName[0:1]) + serviceName[1:] + "Service"
}

func GetEndpointRequestStructName(actionName string) string {
	return actionName + "Request"
}

func GetEndpointResponseStructName(actionName string) string {
	return actionName + "Response"
}

func GetEndpointsStructName(serviceName string) string {
	return serviceName + "Endpoints"
}

func GetEndpointFieldName(actionName string) string {
	return actionName + "Endpoint"
}

func GetMakeEndpointsFuncName(serviceName string) string {
	return "Make" + serviceName + "Endpoints"
}

func GetEndpointBuilderFuncName(actionName string) string {
	return "make" + actionName + "Endpoint"
}

// model

func GetDecodePbToEndpRequestMethodName(actionName string) string {
	return "DecodeGRPC" + actionName + "Request"
}

func GetEncodeEndpToPbResponseMethodName(actionName string) string {
	return "EncodeGRPC" + actionName + "Response"
}

//

func GetTransportStructName(serviceName, transportType string) string {
	return strings.ToLower(transportType) + string_util.FirstLetterToUpperCase(serviceName) + "Server"
}
