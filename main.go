package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"path"
)

//go:generate go run scripts/includeText.go

var appName string
var clobber bool
var rootDir string

func init() {
	log.SetOutput(os.Stdout)

	flag.StringVar(&appName, "appName", "", "(required) ESP32 App Name")
	flag.StringVar(&rootDir, "rootDir", "", "(optional) New root directory (defaults to appName)")
	flag.BoolVar(&clobber, "clobber", false, "(optional) if -clobber, existing app directory will be deleted")
	flag.Parse()

	if rootDir == "" {
		rootDir = appName
	}
}

func main() {
	if appName == "" {
		flag.Usage()
		os.Exit(1)
	}

	cleanupExisting()
	createTemplateApp()
}

func cleanupExisting() {
	if clobber {
		log.Println("Removing old installation at", rootDir)
		os.RemoveAll(rootDir)
	}
}

func createTemplateApp() {


	if err := os.MkdirAll(fmt.Sprintf("%v/src", rootDir), 0755); err != nil {
		log.Fatal("Unable to create src directory:", err)
	}

	log.Println("Writing template files...")
	for templateFileName, _ := range templates {
		fileName := fmt.Sprintf("%v/src/%v", rootDir, templateFileName)
		log.Println("\t", fileName)

		os.MkdirAll(path.Dir(fileName), 0755)

		if file, openErr := os.Create(fileName); openErr != nil {
			log.Fatal("Unable to open", fileName, "for writing:", openErr);
		} else {
			file.WriteString(getTemplate(templateFileName))
		}
	}
}