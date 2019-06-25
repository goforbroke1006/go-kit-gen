package filesystem

import (
	"io/ioutil"
	"os"
)

func CopyFileToFile(fromFilename, toFilename string) error {
	initContent, err := ioutil.ReadFile(fromFilename)
	if nil != err {
		return err
	}

	err = ioutil.WriteFile(toFilename, []byte(initContent), os.ModePerm)
	if nil != err {
		return err
	}

	return nil
}
