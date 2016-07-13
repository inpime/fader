package config

import (
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppCofigTomlDecoder(t *testing.T) {
	cfg := make(configs)
	cfg.AddConfig(sectionName, &Settings{})

	type S struct {
		F string `species:"gopher" color:"blue"`
	}

	for _, addonSettings := range cfg {
		if _, err := toml.Decode(testconfig, addonSettings); err != nil {
			t.Error(err)
		}
	}

	assert.Equal(t, cfg.Config("main").(*Settings).Include(), []string{"a", "b"})
}

var testconfig = `
[main]

include = ["a", "b"]

[notmain]

include = ["c", "d"]

`
