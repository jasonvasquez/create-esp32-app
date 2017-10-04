package main

import (
	"flag"
	//"fmt"
	"os"
	"os/exec"
	"log"
	//"path"
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
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		log.Fatal("Unable to create src directory:", err)
	}

	log.Println("Writing template files...")
	for templateFile, _ := range templates {
		log.Println("\t", templateFile)
		writeTemplateToFile(templateFile, rootDir)
	}

	log.Println("Initializing git...")
	cmd := exec.Command("git", "init")
	cmd.Dir = rootDir
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)

	cmd2 := exec.Command("git", "add", "-A",  ".")
	cmd2.Dir = rootDir
	log.Printf("Running command and waiting for it to finish...")
	err2 := cmd2.Run()
	log.Printf("Command finished with error: %v", err2)

	cmd3 := exec.Command("git", "commit", "-m", "Initial import")
	cmd3.Dir = rootDir
	log.Printf("Running command and waiting for it to finish...")
	err3 := cmd3.Run()
	log.Printf("Command finished with error: %v", err3)
}