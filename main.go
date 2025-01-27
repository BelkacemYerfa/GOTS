package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Mapping of Go types to TypeScript types
var typeMapping = map[string]string{
	"int":     "number",
	"int32":   "number",
	"int64":   "number",
	"uint":    "number",
	"uint32":  "number",
	"uint64":  "number",
	"float32": "number",
	"float64": "number",
	"string":  "string",
	"bool":    "boolean",
}

// Command line flags
var srcFile string
var srcFolder string
var outputFile string

// Store custom types defined in the Go file
var customTypes = make(map[string]bool)

var root = &cobra.Command{
	Use:   "gots",
	Short: "pipe the files to convert types from",
	Long:  `pipe the files to convert types from`,
	// TODO : make a function to that the file can't be inserted with the folder
	Run: func(cmd *cobra.Command, args []string) {
		if srcFile == "" && srcFolder == "" {
			fmt.Println("No source file or folder provided, check the help command")
			return
		}
		if srcFile != "" && srcFolder != "" {
			fmt.Println("You can't provide both a file and a folder, just one")
			return
		}
		switch true {
		case srcFile != "":
			transpile(srcFile, outputFile)
		case srcFolder != "":
			fmt.Println("Folder functionality, Coming soon...")
		default:
			fmt.Println("No source file or folder provided")
		}
	},
}

func Execute() error {
	root.Flags().StringVarP(&srcFile, "src", "s", "", "source file to convert types from")
	root.Flags().StringVarP(&srcFolder, "folder", "f", "", "source folder to convert types from")
	root.Flags().StringVarP(&outputFile, "output", "o", "", "output file name")

	// TODO : Support multi files through folders (input, output)
	return root.Execute()
}

func main() {
	Execute()
}

func transpile(srcFile, outputFile string) {
	// Parse the Go source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, srcFile, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	var tsInterfaces bytes.Buffer

	// Traverse the AST to collect custom types and structs
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Track custom types
		customTypes[typeSpec.Name.Name] = true

		// Process structs
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Generate the TypeScript interface
		structName := typeSpec.Name.Name
		fields := parseFields(structType.Fields.List)
		tsInterface := generateTSInterface(structName, fields)
		tsInterfaces.WriteString(tsInterface + "\n")
		return true
	})

	// Print the TypeScript interfaces
	createFile(&outputFile, tsInterfaces.String())
}

// Parse fields from a Go struct and return a map of field names to TypeScript types
func parseFields(fieldList []*ast.Field) map[string]string {
	fields := make(map[string]string)
	for _, field := range fieldList {
		goType := fieldTypeToString(field.Type)

		// Check if the type is a pointer
		isPointer := strings.HasPrefix(goType, "*")
		if isPointer {
			goType = goType[1:] // Strip the pointer symbol for mapping
		}

		tsType := mapGoTypeToTSType(goType)

		// Handle field names (can be multiple names for the same type)
		for _, name := range field.Names {
			fieldName := name.Name
			if isPointer {
				fieldName += "?" // Mark as optional in TypeScript
			}
			fields[fieldName] = tsType
		}
	}
	return fields
}

// Convert Go AST field type to string representation
func fieldTypeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + fieldTypeToString(t.X)
	case *ast.ArrayType:
		return "[]" + fieldTypeToString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", fieldTypeToString(t.Key), fieldTypeToString(t.Value))
	default:
		return "unknown"
	}
}

// Map Go type to TypeScript type, handling custom types
func mapGoTypeToTSType(goType string) string {
	if tsType, ok := typeMapping[goType]; ok {
		return tsType
	}
	if strings.HasPrefix(goType, "[]") {
		elementType := goType[2:]
		return mapGoTypeToTSType(elementType) + "[]"
	}
	// Check if the type is a custom type
	if customTypes[goType] {
		return goType // Use the same name as the TypeScript interface
	}
	return goType // Fallback for unknown types
}

// Generate TypeScript interface string
func generateTSInterface(structName string, fields map[string]string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("export interface %s {\n", structName))
	for fieldName, tsType := range fields {
		builder.WriteString(fmt.Sprintf("  %s: %s;\n", fieldName, tsType))
	}
	builder.WriteString("}\n")
	return builder.String()
}

func createFile(fileName *string, content string) {
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
