package main

import (
	"api"
	"interfaces"
	"io/ioutil"
	"os"
	"store/boltdb"
	"time"

	"log"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

var (
	db            *bolt.DB
	fileManager   interfaces.FileManager
	bucketManager interfaces.BucketManager
	settings      *api.Settings
)

const (
	FADER_INITFILE = "FADER_INITFILE"
	FADER_DBPATH   = "FADER_DBPATH"
)

func main() {
	log.Println("generate config")

	log.Println("init settings...")
	settings = api.SettingsOrDefault(&api.Settings{
		DatabasePath: os.Getenv(FADER_DBPATH),
	})

	defer func() {
		// IMPORT

		importManager := interfaces.NewImportManager(
			bucketManager,
			fileManager,
		)

		data, err := importManager.Export("v1", "fader2", `Fader2. Fader console v1.`)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(os.Getenv(FADER_INITFILE), data, 0600)
		if err != nil {
			panic(err)
		}

		os.RemoveAll(settings.DatabasePath)
	}()

	log.Println("remove old DB...")
	os.RemoveAll(settings.DatabasePath)

	var err error
	log.Println("open DB...")
	db, err = bolt.Open(settings.DatabasePath, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	log.Println("setup managers...")
	bucketManager = boltdb.NewBucketManager(db)
	fileManager = boltdb.NewFileManager(db)

	log.Println("gen...")
	gen()
}

func gen() {
	var err error
	////////////////////////////////////////////////////////////////////////////
	// SETTINGS
	////////////////////////////////////////////////////////////////////////////

	settingBucketID := uuid.NewV4()
	faderConsoleBucketID := uuid.NewV4()

	bucketFile := interfaces.NewBucket()
	bucketFile.BucketID = settingBucketID
	bucketFile.BucketName = api.ConfigBucketName
	err = bucketManager.CreateBucket(bucketFile)
	log.Printf("create bucket %q", bucketFile.BucketName)
	if err != nil {
		panic(err)
	}

	bucketFile = interfaces.NewBucket()
	bucketFile.BucketID = faderConsoleBucketID
	bucketFile.BucketName = "fader.consolev1"
	err = bucketManager.CreateBucket(bucketFile)
	log.Printf("create bucket %q", bucketFile.BucketName)
	if err != nil {
		panic(err)
	}

	// main.toml

	file := interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = settingBucketID
	file.FileName = api.MainConfigFileName
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
	log.Printf("create file %q", file.FileName)
	if err != nil {
		panic(err)
	}

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
	log.Printf("create file %q", file.FileName)
	if err != nil {
		panic(err)
	}

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

[[routing.routs]]
name = "homepage"
path = "/"
bucket = "fader.consolev1"
file = "index.html"
licenses = ["guest"]
methods = ["get"]

[[routing.routs]]
name = "fileView"
path = "/files/{file_id}"
bucket = "fader.consolev1"
file = "view_file.html"
licenses = ["guest"]
methods = ["get"]
    `)

	err = fileManager.CreateFile(file)
	log.Printf("create file %q", file.FileName)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// index.html
	////////////////////////////////////////////////////////////////////////////

	// fc = fader console
	file = interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = faderConsoleBucketID
	file.FileName = "index.html"
	file.LuaScript = []byte(`
local basic = require("basic")

c = ctx()
c:Set("YourName", c:QueryParam("name"))

print("=====")
print(basic.name, basic.author)
print(basic.name, basic.author)
print(basic.name, basic.author)
print("has faderConfoleIndex", c:Route():Has("faderConfoleIndex"))
print("has empty", c:Route():Has(""))
print("has qwe", c:Route():Has("qwe"))
print("has qwe", c:Route("qwe"):Has())

print("current, имя", c:Route():Name())
print("current, путь", c:Route():Path())
print("current, бакет", c:Route():Bucket())
print("current, файл", c:Route():File())
print("current, аргументы", c:Route():Args())

print("faderConfoleIndex", c:Route("faderConfoleIndex"):URL())
print("homepage", c:Route("homepage"):URL())
print("current", c:Route():URL())
print("file view", c:Route("fileView"):URL("file_id", "ID", "qwd", "qwdqwdqwdqwd"))
vvv = basic.PrimaryIDsData
vvv:Add(basic.PrimaryNamesData)
vvv.Add(basic.AccessStatusData)
basic.check(vvv)
basic.check(vvv.Add(basic.AccessStatusData))

print("")
print("=====")
print("")

founrRoute = c:Route("qwe")
print(founrRoute:Has())
if founrRoute:Has() then
	print("qwe найден!!!")
end

founrRoute = c:Route("homepage")
print("homepage", founrRoute:URL())
print(founrRoute:Has())
if founrRoute:Has() then
	print("homepage найден!!!")
end


c:Set("baskets", basic:ListBuckets())


print("=====")
`)
	file.ContentType = "text/html"
	file.RawData = []byte(`
{% if ctx.Get("YourName") == "" %}
<p>You have no name? Can i name you <a href="?name=Super Star">Super Star</a>?</p>
<p>Don't like the name?</p>
<form>
	<fieldset>
		<legend>What is your name?</legend>
		<input type="text" name="name" placeholder="What is your name?"/>
		<button>Set</button>
	</fieldset>
</form>
{% else %}
<h1>Welcome {{ ctx.Get("YourName") }}!</h1>
{% endif %}
<p><small>Fader2. Fader console v1.</small></p>
{# current route #}
<p><small><a href="?">clear</a></small></p>

<ul>
{% for i in ctx.Get("baskets") %}
<li>
	Name: {{ i.BucketName }}
	<ul>
		{% for ii in ListFilesByBucketID(i.BucketID) %}
		<li>
			<a href="{{ Route("fileView").URLPath("file_id", ii.FileName) }}">Name: {{ ii.FileName }}</a>
		</li>
		{% endfor %}
	</ul>
</li>
{% endfor %}
<ul>
    `)
	// TODO: CSRF код для формы
	// TODO: текущий роут

	err = fileManager.CreateFile(file)
	log.Printf("create file %q", file.FileName)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// view_file.html
	////////////////////////////////////////////////////////////////////////////

	// fc = fader console
	file = interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = faderConsoleBucketID
	file.FileName = "view_file.html"
	file.LuaScript = []byte(`
local basic = require("basic")

c = ctx()
c:Set("YourName", c:QueryParam("name"))

print("=====")
print(basic.name, basic.author)
print(basic.name, basic.author)
print(basic.name, basic.author)
print("has faderConfoleIndex", c:Route():Has("faderConfoleIndex"))
print("has empty", c:Route():Has(""))
print("has qwe", c:Route():Has("qwe"))
print("has qwe", c:Route("qwe"):Has())

print("current, имя", c:Route():Name())
print("current, путь", c:Route():Path())
print("current, бакет", c:Route():Bucket())
print("current, файл", c:Route():File())
print("current, аргументы", c:Route():Args())

print("faderConfoleIndex", c:Route("faderConfoleIndex"):URL())
print("homepage", c:Route("homepage"):URL())
print("current", c:Route():URL())

print("")
print("=====")
print("")

founrRoute = c:Route("qwe")
print(founrRoute:Has())
if founrRoute:Has() then
	print("qwe найден!!!")
end

founrRoute = c:Route("homepage")
print("homepage", founrRoute:URL())
print(founrRoute:Has())
if founrRoute:Has() then
	print("homepage найден!!!")
end


c:Set("baskets", basic:ListBuckets())


print("=====")
`)
	file.ContentType = "text/html"
	file.RawData = []byte(`
{% if ctx.Get("YourName") == "" %}
<p>You have no name? Can i name you <a href="?name=Super Star">Super Star</a>?</p>
<p>Don't like the name?</p>
<form>
	<fieldset>
		<legend>What is your name?</legend>
		<input type="text" name="name" placeholder="What is your name?"/>
		<button>Set</button>
	</fieldset>
</form>
{% else %}
<h1>Welcome {{ ctx.Get("YourName") }}!</h1>
{% endif %}
<p><small>Fader2. Fader console v1.</small></p>
{# current route #}
<p><small><a href="?">clear</a></small></p>

<ul>
{% for i in ctx.Get("baskets") %}
<li>
	Name: {{ i.BucketName }}
	<ul>
		{% for ii in ListFilesByBucketID(i.BucketID) %}
		<li>
			<a href="/files/view">Name: {{ ii.FileName }}</a>
		</li>
		{% endfor %}
	</ul>
</li>
{% endfor %}
<ul>
    `)
	// TODO: CSRF код для формы
	// TODO: текущий роут

	err = fileManager.CreateFile(file)
	log.Printf("create file %q", file.FileName)
	if err != nil {
		panic(err)
	}
}
