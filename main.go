package main

import (
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	var fileNames []string

	// RUN_MODE can be either "generate" or "cli"
	if os.Getenv("RUN_MODE") == "generate" {
		fileNames = getGeneratedFileList()
	} else {
		fileNames = os.Args[1:]
	}

	if len(fileNames) < 1 {
		log.Println("Please provide a schema file")
		os.Exit(1)
	}

	schema, err := os.ReadFile("./schemas/schema.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(string(schema))); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	jsonSchema, err := compiler.Compile("schema.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Starting wait group
	wg.Add(1)
	validateSchema(jsonSchema, fileNames)
	// Waiting for all go routines to finish
	wg.Wait()
}

func validateSchema(schema *jsonschema.Schema, files []string) {

	start := time.Now()
	// Closing wait group once all go routines are finished
	defer wg.Done()
	for index, file := range files {
		log.Println("Validating file", index+1, "of", len(files))
		content, err := os.ReadFile(file)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		var data interface{}
		if err := yaml.Unmarshal(content, &data); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if err := schema.Validate(data); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("File", index+1, "of", len(files), "is valid")
	}
	log.Println("All files are valid")
	duration := time.Since(start)
	log.Println("Total time taken to validate files:", duration.Seconds(), "seconds")
	os.Exit(0)
}

func getGeneratedFileList() []string {
	filePath := "./files/generated/"
	var fileList []string
	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		fileList = append(fileList, filePath+file.Name())
	}
	return fileList
}
