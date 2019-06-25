package naming

func GetServiceInterfaceName(serviceName string) string {
	return serviceName + "Service"
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
