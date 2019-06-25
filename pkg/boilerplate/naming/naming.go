package naming

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
