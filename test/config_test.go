package test

import (
	"fmt"
	"gots/internal"
	"testing"
)

// test for loading the default config file without
func TestConfigGlobalLoading(t *testing.T) {
	config := internal.NewConfig("./gots.yaml")
	config.LoadContent()

	if len(config.EmbeddedTypes) != 0 {
		t.Errorf("expected len to be %v, got 0", len(config.EmbeddedTypes))
	}
}

func TestConfigParse(t *testing.T) {
	config := internal.NewConfig("./gots.yaml")
	config.LoadContent()
	config.ParseContent()

	if len(config.EmbeddedTypes) == 0 {
		t.Errorf("expected len to be %v, got 0", len(config.EmbeddedTypes))
	}

	expected := map[string]map[string]string{
		"gorm.Model": {
			"pkg":   "time",
			"value": "type gormModel struct { ID uint `gorm:'primarykey'`; CreatedAt time.Time; UpdatedAt time.Time; DeletedAt time.Time `gorm:'index'`}",
		},
	}

	for k, v := range expected {

		vKey := config.EmbeddedTypes[k]

		fmt.Printf("%v\n", vKey)
		fmt.Printf("%v\n", expected)

		if vKey == nil {
			t.Errorf("expected %v, got %v", k, nil)
			return
		}

		for k := range v {
			if vKey.(map[string]any)[k] == "" {
				t.Errorf("expected %v, got %v", k, nil)
			}
		}
	}
}
