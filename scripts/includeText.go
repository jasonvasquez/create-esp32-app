package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
	outFile.WriteString(`
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
)

var templates = map[string]string {
	`)
}

func closeOutput() {
	outFile.WriteString(`
}

func getTemplate(key string) string {
	gzBytes,_ := base64.StdEncoding.DecodeString(templates[key])
	gz, _:= gzip.NewReader(bytes.NewBuffer(gzBytes))

	var b bytes.Buffer
	b.ReadFrom(gz)
	return b.String()
}
	`)

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

	outFile.WriteString("\"" + mapKey + "\": \"")

	if fileBytes, err := ioutil.ReadFile(path); err != nil {
		fmt.Println("Error opening file:", err)
	} else {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, gzErr := gz.Write(fileBytes); gzErr != nil {
			log.Fatal("Error compressing file contents:", gzErr)
		} else {
			gz.Flush()
			gz.Close()

			outFile.WriteString(base64.StdEncoding.EncodeToString(b.Bytes()))
		}
	}

	outFile.WriteString("\",\n")

	return nil
}
