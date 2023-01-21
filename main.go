package main

import (
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	directoryPath := os.Args[1:]

	if len(directoryPath) < 1 {
		log.Println("Please provide a path of where the files are located")
		os.Exit(1)
	}

	fileList, err := findFilesInPath(directoryPath[0], "*.template.yaml")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if len(fileList) == 0 {
		log.Println("No files found")
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
	start := time.Now()
	for _, file := range fileList {
		wg.Add(1)
		file := file
		go validate(file, jsonSchema, &wg)

	}
	wg.Wait()
	duration := time.Since(start)
	log.Println("Total time taken to validate", len(fileList), "files:", duration.Seconds(), "seconds")
}

// validate validates a yaml file against a json schema
func validate(file string, jsonSchema *jsonschema.Schema, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Validating file", file)
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
	if err := jsonSchema.Validate(data); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// findFilesInPath finds all files in a path that match a pattern
func findFilesInPath(path string, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, "./"+path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

//The above code is a simple example of how to use the jsonschema package to validate a yaml file against a json schema.
//The code is not optimized for performance and is just a simple example.
