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
	"strconv"
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
	"io"
	"log"
	"os"
	"path"
)

type TemplateFileInfo struct {
	mode       os.FileMode
	contents   string
}

var templates = map[string]TemplateFileInfo {
	`)
}

func closeOutput() {
	outFile.WriteString(`
}

func writeTemplateToFile(templateKey string, outputPathRoot string) {
	fileInfo := templates[templateKey]
	gzBytes, _ := base64.StdEncoding.DecodeString(fileInfo.contents)
	gz, _ := gzip.NewReader(bytes.NewBuffer(gzBytes))

	outputFileName := path.Join(outputPathRoot, templateKey)

	if err := os.MkdirAll(path.Dir(outputFileName), 0755); err != nil {
		log.Fatal("Unable to create directory:", err)
	}

	if outFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.mode); err != nil {
		log.Fatal("Error creating template output file:", err)
	} else {
		io.Copy(outFile, gz)
	}
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

	outFile.WriteString("\"" + mapKey + "\": {\n")

	outFile.WriteString(fmt.Sprintf("mode: 0%v,\n", strconv.FormatUint(uint64(info.Mode()), 8)))

	outFile.WriteString("contents: \"")
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

	outFile.WriteString("\",\n},\n")

	return nil
}
