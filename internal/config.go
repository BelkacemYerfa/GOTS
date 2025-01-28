package internal

import (
	"fmt"
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
	ContentType   string
	EmbeddedTypes map[string]any
}

func NewConfig(filePath string) *Config {
	return &Config{
		FilePath: filePath,
	}
}

func (c *Config) LoadContent() {
	file, err := os.OpenFile(c.FilePath, os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	content := make([]byte, stat.Size())

	_, err = file.Read(content)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	c.ContentType = string(content)
}

func (c *Config) ParseContent() {
	// TODO : parse the content
	content := make(map[string]any)

	err := yaml.Unmarshal([]byte(c.ContentType), &content)

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	for key, value := range content {
		switch value.(type) {
		case map[string]any:
			fmt.Println("Key: ", key)
			c.EmbeddedTypes = make(map[string]any, len(value.(map[string]any)))
			for k, v := range value.(map[string]any) {
				c.EmbeddedTypes[k] = v
			}
		default:
			fmt.Printf("the value is %v\n", value)
			c.EmbeddedTypes[key] = value.(map[string]any)
		}
	}
}
