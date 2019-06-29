package old

import (
	"fmt"
	"os"
	"text/template"
)

func CreateNewFromTemplate(targetFilename, tmplFilename string, data interface{}) {
	filePrefix := os.Getenv("GOPATH") + "/src/github.com/goforbroke1006/go-kit-gen/"
	transportTmpl := template.Must(template.ParseFiles(filePrefix + tmplFilename))
	file, _ := os.Create(targetFilename)
	err := transportTmpl.Execute(file, data)
	if nil != err {
		fmt.Println(err.Error())
	}
}
