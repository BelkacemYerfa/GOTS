package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

type Transpiler struct {
	SourceFile    string
	OutputFile    string
	Config        *Config
	CustomTypes   map[string]bool
	EmbeddedTypes map[string]any
}

func NewTranspiler(srcFile, outputFile string, config *Config) *Transpiler {
	return &Transpiler{
		SourceFile:    srcFile,
		OutputFile:    outputFile,
		Config:        config,
		CustomTypes:   make(map[string]bool),
		EmbeddedTypes: config.EmbeddedTypes,
	}
}

// * Today's work
// TODO : write tests for different types (Recursive, Union , Generics, ...etc)

// * Upcoming features
// TODO : support multi packages in yaml file
// TODO : support comments transpiler (flag -tc (--transpile-comments))
// TODO : support multi output target	(flag -mo (--multi-output))
// TODO : support file versioning
// TODO : support types in addition to interfaces
// TODO : support better error handling
// TODO : Enhanced recursive mode
// 		Filtering files
// ?	exclude certain files (-ex | --exclude testdata/*)
// ? 	support in config file

// * Types to support
// TODO : Enum Support from go to ts
//	Example : type UserRole int
// 		const (
//     	Admin UserRole = iota
//     	User
//     	Guest
// 		)
//	Ts : type UserRole = "Admin" | "User" | "Guest";

// Transpile orchestrates the entire transpilation process
func (t *Transpiler) Transpile() string {

	f := File{
		filename: t.OutputFile,
	}
	f.ReadFile()

	tsInterfaces := &strings.Builder{}

	// Process external types
	externalTS := t.handleExternalTypes()
	if !strings.Contains(f.content, strings.Split(externalTS, "\n")[0]) {
		tsInterfaces.WriteString(externalTS)
	}

	// Parse and transpile Go file
	goFileTS := t.transpileFile(t.SourceFile)
	tsInterfaces.WriteString(goFileTS)

	return tsInterfaces.String()
}

// handleExternalTypes processes embedded types and generates TypeScript interfaces
func (t *Transpiler) handleExternalTypes() string {
	tempFile := ".custom_types.go"
	defer os.Remove(tempFile) // Ensure cleanup of the temporary file

	file, err := os.Create(tempFile)
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer file.Close()

	file.WriteString("package main\nimport (\n")
	for _, val := range t.EmbeddedTypes {
		pkg := fmt.Sprintf("\"%s\"\n", val.(map[string]any)["pkg"])
		if content, err := os.ReadFile(tempFile); err == nil && !strings.Contains(string(content), pkg) {
			file.WriteString(pkg)
		}
	}
	file.WriteString(")\n")

	for _, val := range t.EmbeddedTypes {
		file.WriteString(fmt.Sprintf("%v\n", val.(map[string]any)["value"]))
	}

	fset := token.NewFileSet()
	embNode, err := parser.ParseFile(fset, tempFile, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse temp file: %v", err)
	}

	var tsInterfaces bytes.Buffer
	ast.Inspect(embNode, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		t.CustomTypes[typeSpec.Name.Name] = true
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		fields, embedded := t.parseFields(structType.Fields.List)
		tsInterface := t.generateTSInterface(typeSpec.Name.Name, fields, embedded)
		tsInterfaces.WriteString(tsInterface + "\n")
		return true
	})

	return tsInterfaces.String()
}

// transpileFile parses and transpiles the given Go source file
func (t *Transpiler) transpileFile(srcFile string) string {
	var tsInterfaces strings.Builder
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, srcFile, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		t.CustomTypes[typeSpec.Name.Name] = true
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		fields, embedded := t.parseFields(structType.Fields.List)
		tsInterface := t.generateTSInterface(typeSpec.Name.Name, fields, embedded)
		tsInterfaces.WriteString(tsInterface + "\n")
		return true
	})

	return tsInterfaces.String()
}

// parseFields extracts fields and embedded types from a struct
func (t *Transpiler) parseFields(fieldList []*ast.Field) (map[string]string, []any) {
	fields := make(map[string]string)
	embedded := []any{}

	for _, field := range fieldList {
		if len(field.Names) == 0 {
			embeddedName := t.fieldTypeToString(field.Type)
			if mapped, ok := t.EmbeddedTypes[embeddedName]; ok {
				embedded = append(embedded, mapped)
			} else {
				embedded = append(embedded, embeddedName)
			}
			continue
		}

		goType := t.fieldTypeToString(field.Type)
		isPointer := strings.HasPrefix(goType, "*")
		if isPointer {
			goType = goType[1:]
		}
		tsType := MapGoTypeToTSType(goType, t.CustomTypes)

		for _, name := range field.Names {
			fieldName := name.Name
			if isPointer {
				fieldName += "?" // Optional in TypeScript
			}
			fields[fieldName] = tsType
		}
	}
	return fields, embedded
}

// fieldTypeToString converts a Go field type to its string representation
func (t *Transpiler) fieldTypeToString(expr ast.Expr) string {
	switch typ := expr.(type) {
	case *ast.SelectorExpr:
		pkg := t.fieldTypeToString(typ.X)
		if t.EmbeddedTypes[pkg+"."+typ.Sel.Name] != nil {
			return pkg + typ.Sel.Name
		}
		return typ.Sel.Name
	case *ast.Ident:
		return typ.Name
	case *ast.StarExpr:
		return "*" + t.fieldTypeToString(typ.X)
	case *ast.ArrayType:
		return "[]" + t.fieldTypeToString(typ.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", t.fieldTypeToString(typ.Key), t.fieldTypeToString(typ.Value))
	default:
		return "unknown"
	}
}

// generateTSInterface creates a TypeScript interface from Go struct data
func (t *Transpiler) generateTSInterface(name string, fields map[string]string, embedded []any) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("export interface %s", name))
	if len(embedded) > 0 {
		extends := make([]string, len(embedded))
		for i, embed := range embedded {
			if tsEmbed, ok := t.EmbeddedTypes[embed.(string)]; ok {
				extends[i] = tsEmbed.(EmbType).Value
			} else {
				extends[i] = embed.(string)
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
