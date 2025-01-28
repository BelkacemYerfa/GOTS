package internal

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type EmbType struct {
	Pkg   string
	Value string
}

type Config struct {
	FilePath      string
	EmbeddedTypes map[string]any
}

func NewConfig(filePath string) *Config {
	return &Config{
		FilePath: filePath,
	}
}

func (c *Config) LoadContent() string {
	file, err := os.OpenFile(c.FilePath, os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return ""
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		log.Fatalf("Error: %v", err)
		return ""
	}

	content := make([]byte, stat.Size())

	_, err = file.Read(content)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return ""
	}

	return string(content)
}

func (c *Config) ParseContent() {
	// TODO : parse the content
	parseContent := make(map[string]any)
	contentToParse := c.LoadContent()
	if contentToParse == "" {
		return // No content to parse
	}
	err := yaml.Unmarshal([]byte(contentToParse), &parseContent)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	for key, value := range parseContent {
		switch value.(type) {
		case map[string]any:
			c.EmbeddedTypes = make(map[string]any, len(value.(map[string]any)))
			for k, v := range value.(map[string]any) {
				c.EmbeddedTypes[k] = v
			}
		default:
			c.EmbeddedTypes[key] = value.(map[string]any)
		}
	}
}
