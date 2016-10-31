package api

import (
	"testing"

	"sync"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

var testRawAppConfig = `# settings@main

[main]

include = [
    "console.conf", # application "fader console"
    "filestatic.conf", # addon "filestatic"
    "importexport.conf", # addon "importexport"
    
    "fader.conf"
    # your application
]
tplcache = true

[routing.csrf]
enabled = true
secret = "secret" # change to unique value
tokenlookup = "form:csrf"

[routing.csrf.cookie]
name = "csrf" # cookie name
path = "/" # cookie path
age = 86400 # 24H

[[routing.routs]]
name = "route1"
path = "/"
bucket = "first bucket"
file = "entrypoint"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[addons.a.b]
c = "d"

`

func TestConfigurator_simple(t *testing.T) {

	config := &Config{
		Addons: make(map[string]interface{}),
	}

	_, err := toml.Decode(testRawAppConfig, config)
	assert.NoError(t, err)

	assert.True(t, config.Main.TplCache)
	assert.EqualValues(t, config.Main.Include, []string{
		"console.conf",
		"filestatic.conf",
		"importexport.conf",
		"fader.conf",
	})
	assert.True(t, config.Routing.CSRF.Enabled)
	assert.EqualValues(t, config.Routing.CSRF.Secret, "secret")
	assert.EqualValues(t, config.Routing.CSRF.TokentLookup, "form:csrf")
	assert.EqualValues(t, config.Routing.CSRF.Cookie.Name, "csrf")
	assert.EqualValues(t, config.Routing.CSRF.Cookie.Path, "/")
	assert.EqualValues(t, config.Routing.CSRF.Cookie.Age, 86400)
	assert.EqualValues(t, config.Addons, map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "d",
			},
		},
	})

	assert.True(t, len(config.Routing.Routs) == 1)
	assert.Equal(t, config.Routing.Routs[0].Name, "route1")
	assert.Equal(t, config.Routing.Routs[0].Path, "/")
	assert.Equal(t, config.Routing.Routs[0].Bucket, "first bucket")
	assert.Equal(t, config.Routing.Routs[0].File, "entrypoint")
	assert.Equal(t, config.Routing.Routs[0].AllowedLicenses, []string{"guest"})
	assert.Equal(t, config.Routing.Routs[0].AllowedMethods, []string{"get"})
	assert.Equal(t, config.Routing.Routs[0].CSRF, false)
}

func TestConfigMerge_simple(t *testing.T) {
	config := &Config{
		Addons: make(map[string]interface{}),
	}

	_, err := toml.Decode(testRawAppConfig, config)
	assert.NoError(t, err)

	includeConfigs := []string{
		`
[[routing.routs]]
name = "route2"
path = "/route2"
bucket = "first bucket"
file = "entrypoint2"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get", "post"]

[[routing.routs]]
name = "route3"
path = "/route3"
bucket = "first bucket"
file = "entrypoint3"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["post"]

[addons.c.d]
k = "f"

[addons.a.b]
a = "b"
c = "b"
		`, `
[[routing.routs]]
name = "route4"
path = "/route4"
bucket = "first bucket"
file = "entrypoint4"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[[routing.routs]]
name = "route5"
path = "/route5"
bucket = "first bucket"
file = "entrypoint5"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]
CSRF = true


[addons.c.d]
k = "f"
f = "k"

[addons.a.b]
a = "b"
c = "b"
b = "ac"
		`,
	}

	wg := sync.WaitGroup{}

	for _, tomlConfig := range includeConfigs {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				_config := &Config{
					Addons: make(map[string]interface{}),
				}

				_, err := toml.Decode(tomlConfig, _config)
				assert.NoError(t, err)
				err = config.Merge(_config)
				assert.NoError(t, err, "merge config")

				wg.Done()
			}()
		}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				assert.True(t, config.Config().Main.TplCache)

				wg.Done()
			}()
		}
	}

	wg.Wait()

	//
	{
		// Check use
		_config := &Config{
			Addons: make(map[string]interface{}),
		}

		_, err := toml.Decode(`
		[main]
		use = true
		tplcache = false
		`, _config)
		assert.NoError(t, err)
		err = config.Merge(_config)
		assert.NoError(t, err, "merge config")
	}

	{
		// Check use merge
		_config := &Config{
			Addons: make(map[string]interface{}),
		}

		_, err := toml.Decode(`
		[main]
		tplcache = true
		`, _config)
		assert.NoError(t, err)
		err = config.Merge(_config)
		assert.NoError(t, err, "merge config")
	}
	//

	assert.False(t, config.Main.TplCache) // shold be false
	assert.EqualValues(t, config.Main.Include, []string{
		"console.conf",
		"filestatic.conf",
		"importexport.conf",
		"fader.conf",
	})
	assert.True(t, config.Routing.CSRF.Enabled)
	assert.EqualValues(t, config.Routing.CSRF.Secret, "secret")
	assert.EqualValues(t, config.Routing.CSRF.TokentLookup, "form:csrf")
	assert.EqualValues(t, config.Routing.CSRF.Cookie.Name, "csrf")
	assert.EqualValues(t, config.Routing.CSRF.Cookie.Path, "/")
	assert.EqualValues(t, config.Routing.CSRF.Cookie.Age, 86400)
	assert.EqualValues(t, config.Addons, map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"a": "b",
				"c": "b",
				"b": "ac",
			},
		},
		"c": map[string]interface{}{
			"d": map[string]interface{}{
				"k": "f",
				"f": "k",
			},
		},
	})

	assert.Len(t, config.Routing.Routs, 401) // 100 * 4 + 1
	assert.Equal(t, config.Routing.Routs[0].Name, "route1")
	assert.Equal(t, config.Routing.Routs[0].Path, "/")
	assert.Equal(t, config.Routing.Routs[0].Bucket, "first bucket")
	assert.Equal(t, config.Routing.Routs[0].File, "entrypoint")
	assert.Equal(t, config.Routing.Routs[0].AllowedLicenses, []string{"guest"})
	assert.Equal(t, config.Routing.Routs[0].AllowedMethods, []string{"get"})
	assert.Equal(t, config.Routing.Routs[0].CSRF, false)

	// TODO: bad way to check

	assert.Equal(t, config.Routing.Routs[400].Name, "route5")
	assert.Equal(t, config.Routing.Routs[400].Path, "/route5")
	assert.Equal(t, config.Routing.Routs[400].Bucket, "first bucket")
	assert.Equal(t, config.Routing.Routs[400].File, "entrypoint5")
	assert.Equal(t, config.Routing.Routs[400].AllowedLicenses, []string{"guest"})
	assert.Equal(t, config.Routing.Routs[400].AllowedMethods, []string{"get"})
	assert.Equal(t, config.Routing.Routs[400].CSRF, true)
}
