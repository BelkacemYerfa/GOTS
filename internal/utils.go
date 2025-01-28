package internal

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
)

// Mapping of Go types to TypeScript types
var typeMapping = map[string]string{
	"int":            "number",
	"int32":          "number",
	"int64":          "number",
	"uint":           "number",
	"uint32":         "number",
	"uint64":         "number",
	"float32":        "number",
	"float64":        "number",
	"string":         "string",
	"bool":           "boolean",
	"byte":           "string",
	"timeTime":       "Date",
	"sqlNullString":  "string | null",
	"sqlNullInt64":   "number | null",
	"sqlNullFloat64": "number | null",
	"sqlNullBool":    "boolean | null",
	"uuidUUID":       "string",
	"jsonRawMessage": "any",
	"bigInt":         "string",
	"netIP":          "string",
	"urlURL":         "string",
}

// Store custom types defined in the Go file
var customTypes = make(map[string]bool)
var embeddedTypes = make(map[string]any)

func Transpile(srcFile, outputFile string, config *Config) {
	// Parse the Go source file

	embeddedTypes = config.EmbeddedTypes

	file, err := os.Create("custom_types.go")

	file.WriteString("package main\n")

	file.WriteString("import (\n")

	for _, val := range embeddedTypes {
		frmStr := fmt.Sprintf("\"%s\"\n", val.(map[string]any)["pkg"])
		if content, err := os.ReadFile("custom_types.go"); err == nil {
			if !strings.Contains(string(content), frmStr) {
				file.WriteString(frmStr)
			}
		}
	}
	file.WriteString(")\n")

	for _, val := range embeddedTypes {
		file.WriteString(fmt.Sprintf("%v\n", val.(map[string]any)["value"]))
	}

	fset := token.NewFileSet()

	embNode, err := parser.ParseFile(fset, "custom_types.go", nil, parser.AllErrors)

	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	var tsInterfaces bytes.Buffer

	ast.Inspect(embNode, func(n ast.Node) bool {
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
		fields, embedded := parseFields(structType.Fields.List)
		tsInterface := generateTSInterface(structName, fields, embedded)
		tsInterfaces.WriteString(tsInterface + "\n")
		// * we close the file to kill the process
		file.Close()

		return true
	})

	tokenSet := token.NewFileSet()

	node, err := parser.ParseFile(tokenSet, srcFile, nil, parser.AllErrors)

	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

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
		fields, embedded := parseFields(structType.Fields.List)
		tsInterface := generateTSInterface(structName, fields, embedded)
		tsInterfaces.WriteString(tsInterface + "\n")
		return true
	})

	// Print the TypeScript interfaces
	createFile(&outputFile, tsInterfaces.String())

	err = os.Remove("custom_types.go")

	if err != nil {
		// ! file doesn't get removed
		log.Fatalf("Error: %v", err)
	}

}

// Parse fields from a Go struct and return a map of field names to TypeScript types
// and a list of embedded structs
func parseFields(fieldList []*ast.Field) (map[string]string, []any) {
	fields := make(map[string]string)
	embedded := []any{}

	for _, field := range fieldList {

		if len(field.Names) == 0 {
			// Handle embedded structs
			embeddedName := fieldTypeToString(field.Type)
			fmt.Println(embeddedName)
			if mapped, ok := embeddedTypes[embeddedName]; ok {
				embedded = append(embedded, mapped)
			} else {
				embedded = append(embedded, embeddedName)
			}
			continue
		}

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
	return fields, embedded
}

// Extract embedded types from a Go struct
func getEmbeddedTypes(fieldList []*ast.Field) []string {
	var embedded []string
	for _, field := range fieldList {
		if field.Names == nil { // Embedded fields have no names
			embedded = append(embedded, fieldTypeToString(field.Type))
		}
	}

	fmt.Println(embedded)
	return embedded
}

// Convert Go AST field type to string representation
func fieldTypeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.SelectorExpr: // Qualified identifiers like "gorm.Model"
		pkg := fieldTypeToString(t.X) // Package name
		return pkg + t.Sel.Name
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
	fmt.Println(goType)
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
func generateTSInterface(structName string, fields map[string]string, embedded []any) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("export interface %s", structName))
	if len(embedded) > 0 {
		extends := make([]string, 0, len(embedded))
		for _, embed := range embedded {
			fmt.Println(embeddedTypes)
			if tsEmbed, ok := embeddedTypes[embed.(string)]; ok {
				extends = append(extends, tsEmbed.(EmbType).Value)
			} else {
				extends = append(extends, embed.(string))
			}
		}
		builder.WriteString(fmt.Sprintf(" extends %s", strings.Join(extends, ", ")))
	}
	builder.WriteString(" {\n")
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
