package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	MappingAtomicTypes = map[string]string{
		"string":    "string",
		"int":       "number",
		"float64":   "number",
		"bool":      "boolean",
		"[]string":  "string[]",
		"[]int":     "number[]",
		"[]float64": "number[]",
		"[]bool":    "boolean[]",
		"struct":    "interface",
		"map":       "Record<string, any>",
	}
)

func main() {
	content := ReadFile("types.go")
	tsContent := ContentToInterface(content)
	CreateFile(nil, tsContent)
}

func ReadFile(filePath string) string {
	fmt.Println("Reading file...")
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	path := filepath.Join(pwd, filePath)
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return string(content)
}

func CreateFile(fileName *string, content string) {
	name := "types.ts"
	fmt.Println("Creating file...")
	path, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

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

func ContentToInterface(content string) string {
	fmt.Println("Converting JSON to interface...")

	tsContent := make([]string, 0)

	for _, line := range strings.Split(content, "\n") {
		var setWords = strings.Split(line, " ")
		for idx, word := range setWords {
			var peek = PeekSlice(setWords, idx+1)
			var val string
			if peek != nil {
				val = *peek
			}
			if strings.Contains(word, "type") && strings.Contains(setWords[idx+2], "struct") {
				structName := setWords[idx+1]
				tsContent = append(tsContent, fmt.Sprintf("export interface %s {\n", structName))
				setWords = ChopSlice(setWords, idx)
			}
			if valType, exists := MappingAtomicTypes[val]; exists {
				tsContent = append(tsContent, fmt.Sprintf("\t%s: %s;\n", setWords[0], valType))
				setWords = ChopSlice(setWords, 2)
			} else if !exists {
				// * this means that this is a custom struct
				// * create it before the current struct
				// newStruct := fmt.Sprintf("export interface %s {\n", val)
				// * search about it in the file
				// * if it doesn't exist, create it
				/* _, line := SearchStruct(content, val)
				if line == "" {
					log.Fatalf("Error: %s struct not found", val)
				}
				tsContent = append(tsContent, fmt.Sprintf("\t%s: %s;\n", setWords[0], val))
				tsContent = append(tsContent, newStruct)
				tsContent = append(tsContent, "}\n") */
			}

			if strings.Contains(word, "}") {
				tsContent = append(tsContent, "}\n")
			}
		}
	}
	return strings.Join(tsContent, "")
}

func SearchStruct(content string, structName string) (int, string) {
	for i, line := range strings.Split(content, "\n") {
		if strings.Contains(line, structName) {
			return i, line
		}
	}
	return -1, ""
}

func PeekSlice(slice []string, idx int) *string {
	if idx < len(slice) {
		return &slice[idx]
	}
	return nil
}

func ChopSlice(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}
