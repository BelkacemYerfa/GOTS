package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Mapping of Go types to TypeScript types
var TypeMapping = map[string]string{
	"int":         "number",
	"int32":       "number",
	"int64":       "number",
	"uint":        "number",
	"uint32":      "number",
	"uint64":      "number",
	"float32":     "number",
	"float64":     "number",
	"string":      "string",
	"bool":        "boolean",
	"byte":        "string",
	"Time":        "Date",
	"NullString":  "string | null",
	"NullInt64":   "number | null",
	"NullFloat64": "number | null",
	"NullBool":    "boolean | null",
	"UUID":        "string",
	"RawMessage":  "any",
	"Int":         "string",
	"IP":          "string",
	"URL":         "string",
}

// Map Go type to TypeScript type, handling custom types
func MapGoTypeToTSType(goType string, customTypes map[string]bool) string {
	if tsType, ok := TypeMapping[goType]; ok {
		return tsType
	}
	if strings.HasPrefix(goType, "[]") {
		elementType := goType[2:]
		return MapGoTypeToTSType(elementType, customTypes) + "[]"
	}
	// Check if the type is a custom type
	if customTypes[goType] {
		return goType // Use the same name as the TypeScript interface
	}
	return goType // Fallback for unknown types
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

	if strings.Contains(*fileName, "/") {
		// * this means that the user has provided a path to the folder

		splitted := strings.Split(*fileName, "/")
		subPath := strings.Join(splitted[:len(splitted)-1], "/")
		name = splitted[len(splitted)-1]

		if _, err := os.Stat(subPath); os.IsNotExist(err) {
			os.MkdirAll(subPath, os.ModePerm)
		}

		file, err := os.Create(filepath.Join(path, subPath, name))

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		defer file.Close()

		_, err = file.WriteString(content)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	} else {

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

	fmt.Println("Done")
}

func UpdateFile(fileName *string, content string) {
	name := "types.ts"
	if fileName != nil {
		name = *fileName
	}
	file := NewFile(name, nil, content)

	file.AddContentToFile(content)
}

func MultiFile(folderPath string, recursive bool) []string {
	// * this means that the user has provided a path to the folder
	path, err := os.Getwd()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Fatalf("Error: %v", err)
	}

	files, err := os.ReadDir(filepath.Join(path, folderPath))

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	Files := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			Files = append(Files, filepath.Join(path, folderPath, file.Name()))
		} else if recursive {
			Files = append(Files, MultiFile(filepath.Join(folderPath, file.Name()), recursive)...)
		}
	}

	return Files
}
