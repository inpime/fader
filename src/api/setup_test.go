package api

import (
	"interfaces"
	"os"
	"testing"

	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestConfiguratorUpdate_simple(t *testing.T) {
	err := Setup(e, &Settings{})
	defer func() {
		os.RemoveAll(settings.DatabasePath)
	}()
	assert.NoError(t, err)
	assert.NoError(t, setupSysConfigFilesCase1())

	// first setup
	// err = appConfigUpdateFn()
	// assert.NoError(t, err)
	// err = appRoutesUpdateFn()
	// assert.NoError(t, err)

	configOnce := sync.Once{}
	routesOnce := sync.Once{}

	firstInitWaiting := sync.WaitGroup{}
	firstInitWaiting.Add(1) // config
	firstInitWaiting.Add(1) // routing
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				err = appConfigUpdateFn()
				assert.NoError(t, err)
				configOnce.Do(func() {
					firstInitWaiting.Done()
				})
			}

			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				err = appRoutesUpdateFn()
				assert.NoError(t, err)
				routesOnce.Do(func() {
					firstInitWaiting.Done()
				})
			}

			wg.Done()
		}()
	}

	firstInitWaiting.Wait() // wait the first initialization

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 100; i++ {
				if config.Config().Addons["faderRoutesToml"] != "yes" {
					assert.Fail(t, "not expented value")
				}

				if config.Config().Addons["faderAppToml"] != "yes" {
					assert.Fail(t, "not expented value")
				}
			}

			wg.Done()
		}()

		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				r := vrouter.Get("route2")
				assert.NotNil(t, r, "route")
				assert.Equal(t, r.GetName(), "route2")

				url, err := r.URLPath("param", "123")
				assert.NoError(t, err, "get url")
				assert.Equal(t, "/route2/123", url.String())
			}

			wg.Done()
		}()
	}

	wg.Wait()

	setupSysConfigFilesCase2()

	// update config setup
	err = appConfigUpdateFn()
	assert.NoError(t, err)
	err = appRoutesUpdateFn()
	assert.NoError(t, err)

	assert.Len(t, config.Config().Routing.Routs, 3)
	assert.Len(t, config.Config().Addons, 3)
	assert.EqualValues(t, config.Config().Addons["faderAppToml"], "yes2")
	assert.EqualValues(t, config.Config().Addons["faderRoutesToml"], "yes2")
	assert.EqualValues(t, config.Config().Routing.CSRF.Cookie.Name, "csrf")
	assert.EqualValues(t, config.Config().Routing.CSRF.Cookie.Path, "/")

	r := vrouter.Get("route2")
	assert.Nil(t, r, "route")

	r = vrouter.Get("route22")
	assert.NotNil(t, r, "route")
	assert.Equal(t, r.GetName(), "route22")

	url, err := r.URLPath("param", "123")
	assert.NoError(t, err, "get url")
	assert.Equal(t, "/route22/123", url.String())
}

func setupSysConfigFilesCase1() error {
	bucketSettings := interfaces.NewBucket()
	bucketSettings.BucketID = uuid.NewV4()
	bucketSettings.BucketName = configBucketName
	if err := bucketManager.CreateBucket(bucketSettings); err != nil {
		logger.Panicln("[FAIL] create bucket", err)
		return err
	}

	mainConfigFile := interfaces.NewFile()
	mainConfigFile.FileID = uuid.NewV4()
	mainConfigFile.BucketID = bucketSettings.BucketID
	mainConfigFile.FileName = mainConfigFileName
	mainConfigFile.LuaScript = []byte{}
	mainConfigFile.ContentType = "text/toml"
	mainConfigFile.RawData = []byte(`
[main]

include = [
    "fader.routes.toml",
    "fader.app.toml",
]
tplcache = false

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
path = "/route1/{param:[a-zA-Z0-9._-]+}"
bucket = "actions"
file = "route1"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[addons.a.b]
c = "d"
    `)

	if err := fileManager.CreateFile(mainConfigFile); err != nil {
		logger.Panicln("[FAIL] create mainConfigFile", err)
		return err
	}

	// includes files ----------------------------------------------------

	// fader.routes.toml
	faderRoutesToml := interfaces.NewFile()
	faderRoutesToml.FileID = uuid.NewV4()
	faderRoutesToml.BucketID = bucketSettings.BucketID
	faderRoutesToml.FileName = "fader.routes.toml"
	faderRoutesToml.LuaScript = []byte{}
	faderRoutesToml.ContentType = "text/toml"
	faderRoutesToml.RawData = []byte(`
[[routing.routs]]
name = "route2"
path = "/route2/{param:[a-zA-Z0-9._-]+}"
bucket = "actions"
file = "route2"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[addons]
faderRoutesToml = "yes"

[addons.a.b]
faderRoutesToml = "yes"
    `)

	if err := fileManager.CreateFile(faderRoutesToml); err != nil {
		logger.Panicln("[FAIL] create faderRoutesToml", err)
		return err
	}

	// fader.app.toml
	faderAppToml := interfaces.NewFile()
	faderAppToml.FileID = uuid.NewV4()
	faderAppToml.BucketID = bucketSettings.BucketID
	faderAppToml.FileName = "fader.app.toml"
	faderAppToml.LuaScript = []byte{}
	faderAppToml.ContentType = "text/toml"
	faderAppToml.RawData = []byte(`
[[routing.routs]]
name = "route3"
path = "/route3/{param:[a-zA-Z0-9._-]+}"
bucket = "actions"
file = "route3"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[addons]
faderAppToml = "yes"

[addons.a.b]
faderAppToml = "yes"
    `)

	if err := fileManager.CreateFile(faderAppToml); err != nil {
		logger.Panicln("[FAIL] create faderAppToml", err)
		return err
	}

	return nil
}

func setupSysConfigFilesCase2() error {
	bucketSettings := interfaces.NewBucket()
	bucketSettings.BucketID = uuid.NewV4()
	bucketSettings.BucketName = configBucketName
	if err := bucketManager.CreateBucket(bucketSettings); err != nil {
		logger.Panicln("[FAIL] create bucket", err)
		return err
	}

	mainConfigFile := interfaces.NewFile()
	mainConfigFile.FileID = uuid.NewV4()
	mainConfigFile.BucketID = bucketSettings.BucketID
	mainConfigFile.FileName = mainConfigFileName
	mainConfigFile.LuaScript = []byte{}
	mainConfigFile.ContentType = "text/toml"
	mainConfigFile.RawData = []byte(`
[main]

include = [
    "fader.routes.toml",
    "fader.app.toml",
]
tplcache = false

[routing.csrf]
enabled = true
secret = "secret" # change to unique value
tokenlookup = "form:csrf"

[routing.csrf.cookie]
name = "csrf" # cookie name
path = "/" # cookie path
age = 86400 # 24H

[[routing.routs]]
name = "route12"
path = "/route12/{param:[a-zA-Z0-9._-]+}"
bucket = "actions"
file = "route12"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[addons.a.b]
c = "d"
    `)

	if err := fileManager.CreateFile(mainConfigFile); err != nil {
		logger.Panicln("[FAIL] create mainConfigFile", err)
		return err
	}

	// includes files ----------------------------------------------------

	// fader.routes.toml
	faderRoutesToml := interfaces.NewFile()
	faderRoutesToml.FileID = uuid.NewV4()
	faderRoutesToml.BucketID = bucketSettings.BucketID
	faderRoutesToml.FileName = "fader.routes.toml"
	faderRoutesToml.LuaScript = []byte{}
	faderRoutesToml.ContentType = "text/toml"
	faderRoutesToml.RawData = []byte(`
[[routing.routs]]
name = "route22"
path = "/route22/{param:[a-zA-Z0-9._-]+}"
bucket = "actions"
file = "route22"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[addons]
faderRoutesToml = "yes2"

[addons.a.b]
faderRoutesToml = "yes2"
    `)

	if err := fileManager.CreateFile(faderRoutesToml); err != nil {
		logger.Panicln("[FAIL] create faderRoutesToml", err)
		return err
	}

	// fader.app.toml
	faderAppToml := interfaces.NewFile()
	faderAppToml.FileID = uuid.NewV4()
	faderAppToml.BucketID = bucketSettings.BucketID
	faderAppToml.FileName = "fader.app.toml"
	faderAppToml.LuaScript = []byte{}
	faderAppToml.ContentType = "text/toml"
	faderAppToml.RawData = []byte(`
[[routing.routs]]
name = "route23"
path = "/route23/{param:[a-zA-Z0-9._-]+}"
bucket = "actions"
file = "route23"
#specialhandler = ""
#args = "arg1 arg2"
licenses = ["guest"]
methods = ["get"]

[addons]
faderAppToml = "yes2"

[addons.a.b]
faderAppToml = "yes2"
    `)

	if err := fileManager.CreateFile(faderAppToml); err != nil {
		logger.Panicln("[FAIL] create faderAppToml", err)
		return err
	}

	// actions

	actionsBucket := interfaces.NewBucket()
	actionsBucket.BucketID = uuid.NewV4()
	actionsBucket.BucketName = "actions"
	if err := bucketManager.CreateBucket(actionsBucket); err != nil {
		logger.Panicln("[FAIL] create bucket", err)
		return err
	}

	router22File := interfaces.NewFile()
	router22File.FileID = uuid.NewV4()
	router22File.BucketID = actionsBucket.BucketID
	router22File.FileName = "route22"
	router22File.LuaScript = []byte(`
    c = ctx()
	c:Set("name", "fader")
`)
	router22File.ContentType = "text/html"
	router22File.RawData = []byte(`Hello {{ ctx.Get("name") }} {{ ctx.Get("param") }} {{ ctx.QueryParam("c") }} !`)

	if err := fileManager.CreateFile(router22File); err != nil {
		logger.Panicln("[FAIL] create router22File", err)
		return err
	}

	return nil
}
