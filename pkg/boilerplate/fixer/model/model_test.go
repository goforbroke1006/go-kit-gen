package model

import (
	"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer"
	"github.com/goforbroke1006/go-kit-gen/pkg/filesystem"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const modelUnfinishedSampleFilename = "testdata/model.sample.tmp.go"
const testServiceName = "SomeAwesomeHub"

var testServiceActions = map[string]map[string]string{
	"MethodOne": {
		"FieldOne": "string",
	},
	"MethodTwo": {
		"FieldOne": "string",
		"FieldTwo": "uint64",
	},
	"MethodThree": {
		"FieldOne":   "string",
		"FieldTwo":   "string",
		"FieldThree": "",
	},
	"SayHello": {},
}

func TestModelFixer_Fix(t *testing.T) {
	subjectFilename := "testdata/result/TestModelFixer_Fix.tmp.go"

	if err := filesystem.CopyFileToFile(modelUnfinishedSampleFilename, subjectFilename); nil != err {
		log.Fatal(err)
	}

	fset, file := fixer.OpenGolangSourceFile(subjectFilename)

	mf := NewModelFixer(
		os.Getenv("GOPATH")+"github.com/goforbroke1006/go-kit-gen/pkg/boilerplate/fixer",
		file, "pb", testServiceName, testServiceActions)
	mf.Fix()

	fixer.WriteSourceFile(subjectFilename, file, fset)

	result, err := ioutil.ReadFile(subjectFilename)
	if nil != err {
		t.Fatal(err)
	}

	for action := range testServiceActions {
		if !strings.Contains(string(result), "DecodeGRPC"+action+"Request") {
			t.Fatal("Function DecodeGRPC" + action + "Request is missed")
		}
		if !strings.Contains(string(result), "EncodeGRPC"+action+"Response") {
			t.Fatal("Function EncodeGRPC" + action + "Response is missed")
		}
	}
}
