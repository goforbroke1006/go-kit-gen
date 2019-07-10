package naming

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/string_util"
)

func GetServiceInterfaceName(serviceName string) string {
	return string_util.FirstLetterToUpperCase(serviceName) + "Service"
}
