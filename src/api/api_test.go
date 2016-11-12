package api

import (
	"net/http"
	"os"
	"testing"
	"interfaces"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/satori/go.uuid"
	"io/ioutil"
)

func TestApiGlobal_simple(t *testing.T) {
	err := Setup(e, &Settings{})
	defer func() {
		os.RemoveAll(settings.DatabasePath)
	}()
	assert.NoError(t, err)

	setupSysConfigFilesCase2()

	err = appConfigUpdateFn()
	assert.NoError(t, err)

	s, b := request(echo.GET, "/route22/abc-def.123_456?a=b&c=d'd'd;", e)
	assert.Equal(t, http.StatusOK, s)
	assert.Equal(t, []byte(`Hello check extension &#39;example&#39; fader abc-def.123_456 d&#39;d&#39;d !`), b)
	//                            | addons function               | |lua| |route param|   |url param  |
}

func TestSettingsINitFile_simple(t *testing.T) {
	err := Setup(e, &Settings{})
	defer func() {
		// IMPORT

		importManager := interfaces.NewImportManager(
			bucketManager,
			fileManager,
		)

		data, err := importManager.Export("v1", "fader2", `Fader2. Fader console v1.`)
		assert.NoError(t, err, "create archive")
		err = ioutil.WriteFile("../../fader2.setup.txt", data, 0600)
		assert.NoError(t, err, "write file")

		os.RemoveAll(settings.DatabasePath)
	}()
	assert.NoError(t, err)

	////////////////////////////////////////////////////////////////////////////
	// SETTINGS
	////////////////////////////////////////////////////////////////////////////

	settingBucketID := uuid.NewV4()
	faderConsoleBucketID := uuid.NewV4()

	bucketFile := interfaces.NewBucket()
	bucketFile.BucketID = settingBucketID
	bucketFile.BucketName = configBucketName
	err = bucketManager.CreateBucket(bucketFile)
	assert.NoError(t, err, "create bucket %q", bucketFile.BucketName)

	

	bucketFile = interfaces.NewBucket()
	bucketFile.BucketID = faderConsoleBucketID
	bucketFile.BucketName = "fader.consolev1"
	err = bucketManager.CreateBucket(bucketFile)
	assert.NoError(t, err, "create bucket %q", bucketFile.BucketName)


	// main.toml

	file := interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = settingBucketID
	file.FileName = mainConfigFileName
	file.LuaScript = []byte{}
	file.ContentType = "text/toml"
	file.RawData = []byte(`
[main]

# include only works in 'main.toml'
include = [
    "fader.console.v1.toml",
	"fader.console.v1.routing.toml",
    
	# your application
]

tplcache = false

[routing.csrf]
enabled = true

# REQUIRED: after the first start to please change a secret value (once)
secret = "secret" 
tokenlookup = "form:csrf"

[routing.csrf.cookie]
name = "csrf" # cookie name
path = "/" # cookie path
age = 86400 # 24H
    `)

	err = fileManager.CreateFile(file)
	assert.NoError(t, err, "create file %q", file.FileName)

	////////////////////////////////////////////////////////////////////////////
	// FADER.CONSOLE.v1
	////////////////////////////////////////////////////////////////////////////

	file = interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = settingBucketID
	file.FileName = "fader.console.v1.toml"
	file.LuaScript = []byte{}
	file.ContentType = "text/toml"
	file.RawData = []byte(`
#################################
# replace main config
#################################
# [main]
# tplcache = false
# 
# [routing.csrf]
#
# [routing.csrf.cookie]
#
# [[routing.routs]]
# 

#################################
# custom config
#################################
#
# [addons.a.b]
# c = "d"
    `)

	err = fileManager.CreateFile(file)
	assert.NoError(t, err, "create file %q", file.FileName)

	file = interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = settingBucketID
	file.FileName = "fader.console.v1.routing.toml"
	file.LuaScript = []byte{}
	file.ContentType = "text/toml"
	file.RawData = []byte(`
[[routing.routs]]
name = "faderConfoleIndex"
path = "/fader2/console"
bucket = "fader.consolev1"
file = "index.html"
licenses = ["guest"]
methods = ["get"]
    `)

	err = fileManager.CreateFile(file)
	assert.NoError(t, err, "create file %q", file.FileName)


	// fc = fader console
	file = interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = faderConsoleBucketID
	file.FileName = "index.html"
	file.LuaScript = []byte(``)
	file.ContentType = "text/html"
	file.RawData = []byte(`
<h1>Welcome!</h1>
<smal>Fader2. Fader console v1.</smal>
    `)

	err = fileManager.CreateFile(file)
	assert.NoError(t, err, "create file %q", file.FileName)
}