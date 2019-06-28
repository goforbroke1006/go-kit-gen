package boilerplate

import (
	"log"
	"os"
)

func InitProjectDirs(workingDirPath string) {
	if _, err := os.Stat(workingDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(workingDirPath, os.ModePerm); nil != err {
			log.Fatal(err)
		}
	}

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
