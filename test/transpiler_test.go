package test

import (
	"gots/internal"
	"strings"
	"testing"
)

func TestRecusiveTypes(t *testing.T) {
	srcFile := "test_files/type1.go"
	configFile := "./gots.yaml"
	outputFile := "types/types.ts"

	expContent := `export interface Node {
  Value: string;
  LeftNode?: Node;
  RightNode?: Node;
}

export interface Tree {
  Value: string;
  LeftTree?: Node;
  RightTree?: Node;
}
`

	config := internal.NewConfig(configFile)
	config.ParseContent()

	transpiler := internal.NewTranspiler(srcFile, outputFile, config)

	content := transpiler.Transpile()

	if !strings.Contains(content, expContent) {
		t.Errorf("expected an exact match between the files, got an error")
	}
}
