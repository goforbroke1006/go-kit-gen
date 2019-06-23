package boilerplate

import (
	"fmt"
	"os"
	"text/template"
)

func CreateNewFromTemplate(targetFilename, tmplFilename string, data interface{}) {
	transportTmpl := template.Must(template.ParseFiles(tmplFilename))
	file, _ := os.Create(targetFilename)
	err := transportTmpl.Execute(file, data)
	if nil != err {
		fmt.Println(err.Error())
	}
}
