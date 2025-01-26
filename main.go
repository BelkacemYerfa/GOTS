package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")

	content := ReadFile("types.go")
	jsonContent := StructToJSON(content)
	CreateFile(nil, jsonContent)
}

func ReadFile(filePath string) string {
	fmt.Println("Reading file...")
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Current directory: ", pwd)
	path := filepath.Join(pwd, filePath)
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("File content: ", string(content))
	return string(content)
}

func CreateFile(fileName *string, content string) {
	name := "types.ts"
	fmt.Println("Creating file...")
	path, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Current directory: ", path)

	if fileName != nil {
		name = *fileName
	}

	file, err := os.Create(filepath.Join(path, name))

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer file.Close()

	_, err = file.WriteString(content)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func StructToJSON(s interface{}) string {
	fmt.Println("Converting struct to JSON...")
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return string(bytes)
}

func JSONToInterface(jsonString string) string {
	fmt.Println("Converting JSON to interface...")
	for _, line := range strings.Split(jsonString, "\n") {
		for _, word := range strings.Split(line, " ") {
			println(line, word)
		}
	}
	return ""
}
