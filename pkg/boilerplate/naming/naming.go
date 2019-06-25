package naming

import "strings"

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
