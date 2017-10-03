package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const templateRoot = "templates/"
const outputFileName = "./templates.go"

var outFile *os.File

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	cleanup()

	openOutput()
	defer closeOutput()

	log.Println("Processing templates...")
	filepath.Walk(templateRoot, processFile)
}

func cleanup() {
	if err := os.Remove(outputFileName); err != nil {
		if os.IsNotExist(err) {
			// all good!  (it wasn't there in the first place)
		} else {

			log.Fatal("Error cleaning up old template output file", err)
		}
	}
}

func openOutput() {
	var err error
	if outFile, err = os.Create(outputFileName); err != nil {
		log.Fatal("Error creating templates.go", err)
	}
	outFile.Write([]byte(`
		package main

		var templates map[string]string

		func init() {
			templates = make(map[string]string)

	`))
}

func closeOutput() {

	outFile.Write([]byte("}\n"))

	if err := outFile.Close(); err != nil {
		log.Fatal("Error closing template output file", err)
	}
}

func processFile(path string, info os.FileInfo, err error) error {

	// dump out early for directories
	if info.IsDir() {
		return nil
	}

	mapKey := strings.TrimPrefix(path, templateRoot)
	log.Println("\t", mapKey)

	outFile.Write([]byte("templates[\"" + mapKey + "\"] = `"))

	if file, err := os.Open(path); err != nil {
		fmt.Println("Error opening file:", err)
	} else {
		if _, copyErr := io.Copy(outFile, file); copyErr != nil {
			log.Fatal("Error coyping template file", copyErr)
		}
	}

	outFile.Write([]byte("`\n"))

	return nil
}
