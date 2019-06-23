package boilerplate

import (
	"log"
	"os"
)

func InitProjectDirs(workingDirPath string) {
	dirs := []string{
		"service",
		"endpoint",
		"model",
		"transport",
	}
	for _, d := range dirs {
		if _, err := os.Stat(workingDirPath + "/" + d); os.IsNotExist(err) {
			if err := os.Mkdir(workingDirPath+"/"+d, os.ModePerm); nil != err {
				log.Fatal(err)
			}
		}
	}
}
