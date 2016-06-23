package main

import (
	"api"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/inpime/dbox"
	"gopkg.in/olivere/elastic.v3"
	// "io/ioutil"
	"log"
	"net"
	// "net/http"
	"os"
	// "sort"
	"store"
	"time"
	"utils"
)

func initElasticSearch() error {
	dockerElasticSearchHost := os.Getenv("FADER_ES_ADDR_DOCKER")
	if len(dockerElasticSearchHost) > 0 {
		log.Printf("Elasticsearch via docker. Host: %q", dockerElasticSearchHost)
		ips, err := net.LookupIP(dockerElasticSearchHost)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Lookup for elasticsearch returns the following IPs:")
		for _, ip := range ips {
			api.Cfg.Search.Host = "http://" + ip.String() + ":9200"
			log.Printf("%v", ip)
			break
		}
	}

	db, err := elastic.NewClient(
		// elastic.SetSniff(false),
		elastic.SetURL(api.Cfg.Search.Host),
		elastic.SetInfoLog(logrus.New()),
		elastic.SetHealthcheckTimeoutStartup(time.Second*60),
		// elastic.SetTraceLog(logrus.New()),
	)

	if err != nil {
		panic(err)
	}

	// ------------------------
	// search
	// ------------------------

	// store.ElasticSearchIndexName defined in initConfig()

	exists, err := db.IndexExists(store.ElasticSearchIndexName).Do()

	if err != nil {
		logrus.WithError(err).Errorf("exsist index %q", store.ElasticSearchIndexName)
		return err
	}

	if !exists {
		// Create a new index.
		createIndex, err := db.CreateIndex(store.ElasticSearchIndexName).Do()

		if err != nil {
			// Handle error
			logrus.WithError(err).Errorf("create index %q", store.ElasticSearchIndexName)
			return err
		}

		if !createIndex.Acknowledged {
			// Not acknowledged

			logrus.Errorf("create index %q: not acknowledged", store.ElasticSearchIndexName)
		}

		// OK
	}

	store.ESClient = db

	return nil
}

func initStroes() {
	db, err := bolt.Open(api.Cfg.Store.BoltDBFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		panic(err)
	}

	// set default buckets store for buckets
	dbox.BucketStore = store.NewBoltDBStore(db, api.BucketsBucketName)

	// main boltdb client for stores
	store.BoltDBClient = db

	// dbox.RegistryStore("localfs.static.rawdata", store.NewLocalStore(nil, api.Cfg.Store.StaticPath))
	// dbox.RegistryStore("boltdb.static.structdata", store.NewBoltDBStore(db, "static.files.structdata"))
	// dbox.RegistryStore(api.SettingsStoreName, store.NewBoltDBStore(db, "settings"))
	// dbox.RegistryStore("boltdb.console", store.NewBoltDBStore(db, "console"))
	// dbox.RegistryStore("boltdb.pages", store.NewBoltDBStore(db, "pages"))
	// dbox.RegistryStore(api.UsersStoreName, store.NewBoltDBStore(db, "users"))
	// dbox.RegistryStore(api.BucketsStoreName, store.NewBoltDBStore(db, "buckets"))

	// // for develop
	// dbox.RegistryStore("fs.settings.rawdata", store.NewLocalStore(nil, api.Cfg.WorkspacePath+"/.settings/"))
	// dbox.RegistryStore("fs.pages.rawdata", store.NewLocalStore(nil, api.Cfg.WorkspacePath+"/.pages/"))
	// dbox.RegistryStore("fs.console.rawdata", store.NewLocalStore(nil, api.Cfg.WorkspacePath+"/.console/"))
	// dbox.RegistryStore("memory", dbox.NewMemoryStore())

	// ------------------
	// Setup buckets
	// ------------------

	bucket, err := store.BucketByName(api.BucketsBucketName)
	bucket.InitInOneStore(dbox.BoltDBStoreType) // store key - boltdb.buckets

	if err == dbox.ErrNotFound {
		fmt.Printf("buckets: create bucket %q\n", bucket.Name())

		// bucket.SetMetaDataStoreName(api.BucketsStoreName)
		// bucket.SetRawDataStoreName(api.BucketsStoreName)
		// bucket.SetMapDataStoreName(api.BucketsStoreName)

		// Mapping bucket

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		bucket.MMapDataFilesMapping().LoadFromM(store.BucketMapMappingDefault)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	bucket, err = store.BucketByName(api.SettingsBucketName)
	bucket.InitRawDataStore(dbox.LocalStoreType, false)  // store key - fs.settings.rawdata
	bucket.InitMetaDataStore(dbox.BoltDBStoreType, true) // store key - boltdb.settings
	bucket.InitMapDataStore(dbox.BoltDBStoreType, true)  // store key - boltdb.settings

	if err == dbox.ErrNotFound {
		fmt.Printf("settings: create bucket %q\n", bucket.Name())

		// bucket.SetMetaDataStoreName(api.SettingsStoreName)
		// bucket.SetRawDataStoreName("fs.settings.rawdata")
		// bucket.SetMapDataStoreName(api.SettingsStoreName)

		// Mapping bucket

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		// bucket.MMapDataFilesMapping().LoadFromM(store.FileEmptyMapDataMapping)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	bucket, err = store.BucketByName(api.StaticBucketName)
	bucket.InitRawDataStore(dbox.LocalStoreType, false)  // store key - fs.static.rawdata
	bucket.InitMetaDataStore(dbox.BoltDBStoreType, true) // store key - boltdb.static
	bucket.InitMapDataStore(dbox.BoltDBStoreType, true)  // store key - boltdb.static

	if err == dbox.ErrNotFound {
		fmt.Printf("settings: create bucket %q\n", bucket.Name())

		// bucket.SetMetaDataStoreName("boltdb.static.structdata")
		// bucket.SetRawDataStoreName("localfs.static.rawdata")
		// bucket.SetMapDataStoreName("boltdb.static.structdata")

		// Mapping bucket

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		// bucket.MMapDataFilesMapping().LoadFromM(store.FileEmptyMapDataMapping)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	bucket, err = store.BucketByName(api.PagesBucketName)
	bucket.InitRawDataStore(dbox.LocalStoreType, false)  // store key - fs.pages.rawdata
	bucket.InitMetaDataStore(dbox.BoltDBStoreType, true) // store key - boltdb.pages
	bucket.InitMapDataStore(dbox.BoltDBStoreType, true)  // store key - boltdb.pages

	if err == dbox.ErrNotFound {
		fmt.Printf("settings: create bucket %q\n", bucket.Name())

		// bucket.SetMetaDataStoreName("boltdb.pages")
		// bucket.SetRawDataStoreName("fs.pages.rawdata")
		// bucket.SetMapDataStoreName("boltdb.pages")

		// Mapping bucket

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		// bucket.MMapDataFilesMapping().LoadFromM(store.FileEmptyMapDataMapping)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	bucket, err = store.BucketByName(api.UsersBucketName)
	bucket.InitInOneStore(dbox.BoltDBStoreType) // store key - boltdb.users

	if err == dbox.ErrNotFound {
		fmt.Printf("settings: create bucket %q\n", bucket.Name())

		// bucket.SetMetaDataStoreName(api.UsersStoreName)
		// bucket.SetRawDataStoreName(api.UsersStoreName)
		// bucket.SetMapDataStoreName(api.UsersStoreName)

		// Mapping bucket

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		bucket.MMapDataFilesMapping().
			Set("licenses", utils.Map().
				Set("type", "string").
				Set("index", "not_analyzed"))

		bucket.UpdateMapping()
		bucket.Sync()
	}

	bucket, err = store.BucketByName(api.ConsoleBucketName)
	bucket.InitRawDataStore(dbox.LocalStoreType, false)  // store key - fs.console.rawdata
	bucket.InitMetaDataStore(dbox.BoltDBStoreType, true) // store key - boltdb.console
	bucket.InitMapDataStore(dbox.BoltDBStoreType, true)  // store key - boltdb.console

	if err == dbox.ErrNotFound {
		fmt.Printf("settings: create bucket %q\n", bucket.Name())

		// bucket.SetMetaDataStoreName("boltdb.console")
		// bucket.SetRawDataStoreName("fs.console.rawdata")
		// bucket.SetMapDataStoreName("boltdb.console")

		// Mapping bucket

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		// bucket.MMapDataFilesMapping().LoadFromM(store.FileEmptyMapDataMapping)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	// --------------------
	// Settings
	// --------------------

	file, err := store.LoadOrNewFile(api.SettingsBucketName, api.MainSettingsFileName)
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/toml")
		file.RawData().Write([]byte(`# main settings

routs = ["routing"]
pageCaching = false

[usercontent]
bucket="static"`))
		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.SettingsBucketName, api.RoutingSettingsFileName)
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/toml")
		file.RawData().Write([]byte(`#routing
# -----------------
# Web
# -----------------

[[routs]]
path = "/"
handler = "pages index"
methods = ["get"]
licenses = ["guest", "user", "admin"]

[[routs]]
path = "/usercontent/{fileid}"
handler = "usercontent"
special = true
methods = ["get"]
licenses = ["guest", "user", "admin"]

[[routs]]
path = "/console/settings/export"
handler = "exportapp"
special = true
methods = ["get"]
licenses = ["admin"]

[[routs]]
path = "/console/settings/import"
handler = "importapp"
special = true
methods = ["post"]
licenses = ["admin"]

# -----------------
# Sessions
# -----------------

[[routs]]
path = "/sessions/logout"
handler = "pages logout"
methods = ["post"]
licenses = ["user", "admin"]

[[routs]]
path = "/sessions/login"
handler = "pages login"
methods = ["post"]
licenses = ["guest"]

[[routs]]
path = "/sessions/current"
handler = "pages sessioninfo"
methods = ["get"]
licenses = ["user", "admin"]

# -----------------
# Console
# -----------------

[[routs]]
path = "/console"
handler = "console dashboard"
methods = ["get"]
licenses = ["admin"]

[[routs]]
path = "/console/dashboard"
handler = "console dashboard"
methods = ["get"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets"
handler = "console bucketslist"
methods = ["get"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/{bucket_id}/files"
handler = "console fileslist"
methods = ["get"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/{bucket_id}/files/{file_id}/view"
handler = "console viewfile"
methods = ["get"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/{bucket_id}/newfile"
handler = "console newfileform"
methods = ["get"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/newbucket"
handler = "console newbucketform"
methods = ["get"]
licenses = ["admin"]

# -----------------
# Manager files
# -----------------

[[routs]]
path = "/console/buckets/{bucket_id}/files/{file_id}/textdata/put"
handler = "console putfile_textdata"
methods = ["post"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/{bucket_id}/files/{file_id}/rawdata/put_via_uploader"
handler = "console putfile_rawdata_viauploader"
methods = ["post"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/{bucket_id}/files/{file_id}/structdata/put"
handler = "console putfile_structdata"
methods = ["post"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/{bucket_id}/files/{file_id}/properties/put"
handler = "console putfile_properties"
methods = ["post"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets/{bucket_id}/files"
handler = "console newfile"
methods = ["post"]
licenses = ["admin"]

[[routs]]
path = "/console/buckets"
handler = "console newbucket"
methods = ["post"]
licenses = ["admin"]



`))

		file.Sync()
	}

	// --------------------
	// Users
	// --------------------

	file, err = store.LoadOrNewFile(api.UsersBucketName, api.GuestUserFileName)
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		utils.M(file.MapData()).Set("fullName", "Guest")
		user := api.FileAsUser(file)
		user.AddLicense(api.GuestLicense)

		user.Sync()
	}

	file, err = store.LoadOrNewFile(api.UsersBucketName, "demo")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		utils.M(file.MapData()).Set("fullName", "Demo")
		user := api.FileAsUser(file)
		user.AddLicense(api.UserLicense)
		user.AddLicense(api.AdminLicense)

		user.Sync()
	}

	// --------------------
	// Pages
	// --------------------

	file, err = store.LoadOrNewFile(api.PagesBucketName, "index")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{#index page#}
    {% set session = ctx.Session() %}

<hr />
<p>Flash messages:</p>
{{ session.Get("counter") }} + 1 =
{{ session.Set("counter", session.Get("counter")+1) }}
{{ session.Get("counter") }}
{#{{ session.Save() }}#}
<hr />

<hr />
{% if session.IsAuth() %}
<form class="uk-form" action="/sessions/logout" method="post">
    <button>Logout</button>
</form>
{% else %}
<p>Login</p>
<form class="uk-form" action="/sessions/login" method="post">
   <input type="text" name="email" value="{{ctx.FormValue("email")}}"/>
   <input type="text" name="password" value="{{ctx.FormValue("password")}}"/>
   <button>Login</button>
</form>

{% endif %}
<hr />`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.PagesBucketName, "login")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{#login#}
{% set session = ctx.Session() %}

{% if ctx.IsPost() %}
{% set res = session.Signin(ctx.FormValue("email"), ctx.FormValue("password")) %}

{% if res|is_error %}
{{ session.AddFlash(res.Error()) }}
{{ ctx.Redirect302("/?error")}}
{% else %}
{{ ctx.Redirect302("/sessions/current")}}
{% endif %}

{% endif %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.PagesBucketName, "logout")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{#logout#}
{% set session = ctx.Session() %}

{% if ctx.IsPost() %}
{{ session.Logout() }}
{{ ctx.Redirect302("/")}}
{% endif %}
`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.PagesBucketName, "sessioninfo")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{#sessioninfo#}
{% set user = ctx.CurrentUser() %}

<pre>
ID: {{ user.ID() }}
Name: {{ user.Name() }}
Licenses: {{ user.Licenses()|stringformat:"%#v" }}
</pre>`))

		file.Sync()
	}

	// --------------------
	// Console
	// --------------------

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "layout")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`<html>
    <head>
        <title>{% block title %}{% endblock %}</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta charset="utf-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="stylesheet" href="//cdn.jsdelivr.net/uikit/2.26.2/css/uikit.almost-flat.min.css" />
        <link rel="stylesheet" href="//cdn.jsdelivr.net/uikit/2.26.2/css/components/notify.almost-flat.min.css" />
        <link rel="stylesheet" href="//cdn.jsdelivr.net/uikit/2.26.2/css/components/placeholder.almost-flat.min.css" />
        <link rel="stylesheet" href="//cdn.jsdelivr.net/uikit/2.26.3/css/components/autocomplete.almost-flat.min.css" />
        <link rel="stylesheet" href="//cdn.jsdelivr.net/uikit/2.26.2/css/components/upload.almost-flat.min.css" />
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/jsoneditor/5.5.5/jsoneditor.min.css" />
		
        <script src="//cdn.jsdelivr.net/mithril/0.2.5/mithril.min.js"></script>
        <script src="//cdn.jsdelivr.net/lodash/4.11.2/lodash.min.js"></script>
        <script id="mapContentTypeToACEType" type="application/json">{
    "text/css": "ace/mode/css",
    "text/plain": "ace/mode/django",
    "text/html": "ace/mode/django",
    "application/xhtml+xml": "ace/mode/django",
    "text/javascript": "ace/mode/javascript",
    "application/javascript": "ace/mode/javascript",
    "application/x-javascript": "ace/mode/javascript",
    "text/toml": "ace/mode/toml",
    "text/xml": "ace/mode/xml",
    "application/soap+xml": "ace/mode/xml",
    "application/json": "ace/mode/json",
    "text/json": "ace/mode/json",
    "text/csv": "ace/mode/json"
}</script>
        
		{% block head %}{% endblock %}
		{% ssi "console/head" %}
    </head>
    <body>
		{% include "console/navbar" %}
		
		<div class="uk-width-1-2 uk-container uk-container-center">
		{% block breadcrumb %}{% endblock %}
		</div>
		
        {% block content %}{% endblock %}
		
		{% include "console/footer" %}
        
        {% ssi "console/scripts" %}

        <script src="//cdn.jsdelivr.net/jquery/2.1.4/jquery.min.js"></script>
        <script src="//cdn.jsdelivr.net/uikit/2.26.2/js/uikit.min.js"></script>
        <script src="//cdn.jsdelivr.net/uikit/2.26.2/js/components/notify.min.js"></script>
        <script src="//cdn.jsdelivr.net/uikit/2.26.2/js/components/upload.min.js"></script>
        <script src="//cdn.jsdelivr.net/uikit/2.26.3/js/core/dropdown.min.js"></script>
        


        <script src="//cdnjs.cloudflare.com/ajax/libs/jsoneditor/5.5.6/jsoneditor.min.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/ace.js"></script>
        
        
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/mode-snippets.js"></script>

        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/mode-toml.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/snippets/toml.js"></script>
        
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/mode-django.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/snippets/django.js"></script>
        
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/mode-html.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/snippets/html.js"></script>
        
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/mode-css.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/snippets/css.js"></script>
        
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/mode-xml.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/snippets/xml.js"></script>
        
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/mode-json.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.3/snippets/json.js"></script>
        
		{% block scripts %}{% endblock %}
    </body>
</html>`))
		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "navbar")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`<nav class="uk-navbar">

	<a class="uk-navbar-brand" href="/console">Fader</a>
	<ul class="uk-navbar-nav">
		<li class="{# uk-active #}"><a href="/console/buckets">Buckets</a></li>
	</ul>
	
	
	<div class="uk-navbar-content"><a href="/">Go to site</a></div>

	<div class="uk-navbar-content uk-navbar-flip">
		<div class="uk-navbar-content uk-text-muted">{{ ctx.CurrentUser().Name() }}</div>
		<ul class="uk-navbar-nav">
			<li class="uk-navbar-content">
				<form class="uk-form uk-margin-remove uk-display-inline-block" action="/sessions/logout" method="post">
					<button class="uk-button uk-button-link"><a href="#">Logout</a></button>
				</form>
			</li>
		</ul>
	</div>

</nav>`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "footer")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`<footer style="padding: 120px 0;">
	<div class="uk-container uk-container-center uk-text-center">

		<ul class="uk-subnav uk-subnav-line uk-flex-center">
			<li><a href="http://github.com/inpime/fader">GitHub</a></li>
			<li><a href="http://github.com/inpime/fader/issues">Issues</a></li>
			<li><a href="http://github.com/inpime/fader/blob/master/CHANGELOG.md">Changelog</a></li>
			<li><a href="https://community.inpime.com">Community</a></li>
			<li><a href="https://twitter.com/fader">Twitter</a></li>
			<li><a href="https://facebook.com/fader">Facebook</a></li>
			<li><a href="https://plus.google.com/fader">Google+</a></li>
		</ul>

		<div class="uk-panel">
			<p>Licensed under <a href="http://opensource.org/licenses/MIT">MIT license</a>.</p>
			<a href="http://inpime.com"><img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAFMAAABTCAYAAADjsjsAAAAAAXNSR0IArs4c6QAAB/BJREFUeAHtXOluI0UQbt9XHOfyEpSwCAISSPBAPBRPw3vwAyEBuwIRbXaTTWI7vq/Yob4JLbUnM3FXX1lHKcmZid3Tx9fd9XVVV0/mp59/uRMv4gSBrJNcXjKJEHgB0+FAyDvMyyircjEndqplITJCjCZzMZzOxWK5mZrnScAk3MRevSKO9mqiXimtdMId4XjVH4l3130xmd2u/Pap/xMcTAD53fG+2Nui0ZggGUrwarsq9gnsP89a4mY4TUjF+wp57tbKVGZF5LIZMZnfiiHNgs5gIhboPUcSHMzj/XoqkGqbcoTA90cH4rfTy2jqq79x7uvlovj28x1RKRUePDadL8Rf79uiP5k9+M3ki6AEVCrkxfFBXbueWardyeGOdvp4whoB+OOXzUQgkbZUyIkfXjcF0rmQoGAe7dZEFnOOIfVKUWzRx0RODhtiXXG2HabWKyiYu/VkPalWKOk+Tb8mpZXf1Ur5B+Qmf4tf0WHbhh2m5hUQzAxNKzMVXabpyJVqmTd146sKbnlIHwzMjCDWNCROk8e46iTjAAkHWej1IQDBgtxEsJjnyvx2yXpkRsxuK8HAREVb/TG7vveLeP5z3dGULCk9QNHRSG8rQcG8uBmK24VeA2XDPnaHwmTUwCR91+rLbB69npG1hTWnrQQFE0C+Pe9oq05YKqdXXeM2fmgNxFXv8VF93h6Q6dozLkN9MCiYKLhNJhysjnWkMiA9+fvpNY3kdSnV5qze48m3523xx1lbtEnFqJbjmOz+Nx/a4p/L7tq6rOaa/p/ZWiU9P61foDt//fuC7O97G3yrUoDTiORO9MdzcdkbiY+kEtTGa2WckAh5tAfj6JMlu7yYz1EHLaw6KaGY6KsnARMlT0hHvW/3o08+lxH5XJZ041IsXSCY0tol6VGfnqgnA1NtL6YyRsumSzAwi2TFHDbINqepVqCRmINRnCLwDbdoqrdIv5oI3GxwYMyJ8NBJaQt46OUzTcbXqUcwMF81quILhscItnJ7cGFEDl9/tiO2NMxJrBZcSvrwcFkK5bVV4nl+4B6rGTgf4BRBx+mIqUWWlnc4MDVGSrySTfK2c6RAJPYNOYJ1ZTjZwJGJRmKkcQVbFxzvJxzJhZxeObA0xzO+zf9YG4KMzJrBqESl0QHwNepIs1GJ9o100iLNaDZzso5VywsCpqmnHBXd314/1Uu0ED8h0uEImNy1hAGTST5qIw/ISoo21dUvY/cnpCcfW2rFkkf/YnfStYQB03Cao7HFfFY0qulT/ZD2lbCNy5WBox1JtVzvYJqSj1rJgxRWRzTIV82GmlTrHhbrcOqWyVGwdzBNyUdFZZ824h6yeob2w/cii0pNq3M/moJ84FNyK97BtCEf2dQCEcx2bKojtMZ0R9EH+aCu/sFkks9ldyQxXLkitEVKpVgQr5vb8l/2dXPBZJAPGvnvZS/RHofexFTHB+Euac4LHWRdm5GyTK+ODi75YFNrTl6eHgVrNWqr0XHwOu3S6KyU9BfyspHqNSIfD8silOEVTC75SMctQgrjYKKyr5t1USnaVXlEJqQvB7RXncklH7n2a/UmiaYeAqxspjc6ZDB2E/GGvOLiF0wG+WD6YdRAbskLceNgHzveWPzvi3yQt18wGeSD3UI1ZuB6zRYtKm8ivsgHdfEGJpd85BSXAGEH07Vu80k+XsHkks8gFoeE0BYXIdiyc3D1ST5eweSSzzDB8XBNrO5SfJKPXzA55EM1SXKJtXvTFT1qC6xP8vELJoN8sL5MOvuzuFuKzvDxWCEOwD7JxxuYtuSjAuSK1X2TjzcwueSTNMUloK7O6vgmH29gcsknTZfB2sGhpw4FXtmKb/LxByaTfOJrTAmcdOC6mOppHSbLcnH1smjXCU2RlU8jH/wufeGdgX5Itcw3fvVNPijPOZguyUcCAksIQbKmEoJ8UDfnYLokHwkeHMI2YIYgHy9guiIfCSSumO4RqxueQw9BPqin85HJiXYDSGnkg8qpgkW96egMQT6oq3Mwa2V9T/g0xfJRQVTvrygA1kRCkA/q5RRMxKWXGecj456idUDBi6R7UErmFYp8nIPJJR+OLsMCHutO7im3UOTjHEyOvkThHF2GmCPIukNSUSLlD6fDlMeMbp1Oc46+BPkMKUxFV3DUBXK/Hax/ZJDTYbp1SUvnFEyO5QPyMTl9Bh3ImeqhyAcAOwMzRwe2OXvaXPJRR4M+mHes0a+WYXLvDMxaBUuih7FqaZWy0WWY6jonfXF0UN3xTKuLq+/1F4VrSuS+gcVm+mGqv6EDprs1CubKQPsmyF3G+FBWQm5aXzkDk6MvUTPbYyPd0YzISJ/AtNCwTORumjN8mAjOwue5iRMwsaCuMl60ZDsqP9VOcAImLJ91L2NSAXhsz0dNZ3OP1UVocVIil3x0PUU2YGCbOLS4AZOxR44G2jB5aIA45TkBc4veBKgr8EvKoFbdZzYlnTWYmYh89FdYmOIpK8NNwSy1ntZg4gV2nGjekI6H1FZ7+sEaTM4URxtCMLknrNZmaw0m2yGcEDq4tpYbksAaTM7IfM7kg/62AhML9SrpTF15zuRjDSZMSA75jDycpNXtyBDprEYm1/Lpje1fqxgCFNMy9OdoQgl44Sfe8gfB2hHbEEsy42T0WvTD/38QjIF3sj1nsQKzM5xQmLR5QNVzA9Zqmj83MGzb8wKmLYLK8y9gKmDY3r6AaYug8vx/uKcBD2xbzKwAAAAASUVORK5CYII=" width="90" height="30" title="Fader is a product inpime.com" alt="Fader is a product inpime.com"></a>
		</div>

	</div>
</footer>`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "dashboard")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{% extends "console/layout" %}

{% block breadcrumb %}
<ul class="uk-breadcrumb">
    <li><a href="/console/buckets">buckets</a></li>
</ul>
{% endblock %}

{% block content %}
<div class="uk-width-1-2 uk-container uk-container-center">
<h1>Dashboard</h1>
<div class="uk-panel uk-panel-box">
    <div class="uk-panel-title">Import/export settings</div>
    <a class="uk-button" href="/console/settings/export"><i class="uk-icon-download"></i> export</a>    
    <hr />
    <div> 
        <div id="upload-drop" class="uk-placeholder uk-text-center">
            <i class="uk-icon-cloud-upload uk-icon-medium uk-text-muted uk-margin-small-right"></i> Attach archive by dropping them here or <a class="uk-form-file">selecting one<input id="upload-select" type="file" name="BinData" /></a> for import settings.
        </div>
        <div id="progressbar" class="uk-progress uk-hidden">
            <div class="uk-progress-bar" style="width: 0%;">0%</div>
        </div>
    </div>
</div>

</div>
{% endblock %}

{% block scripts %}
<script>
    var progressbar = $("#progressbar"),
        bar         = progressbar.find('.uk-progress-bar'),
        settings    = {

            action: '/console/settings/import', // upload url
            param: 'BinData',
            params: {},
            type: 'json',

            allow : '*.', // allow all file types

            loadstart: function() {
                bar.css("width", "0%").text("0%");
                progressbar.removeClass("uk-hidden");
            },

            progress: function(percent) {
                percent = Math.ceil(percent);
                bar.css("width", percent+"%").text(percent+"%");
            },

            allcomplete: function(response) {

                bar.css("width", "100%").text("100%");

                setTimeout(function(){
                    progressbar.addClass("uk-hidden");
                }, 250);
                
                alert(JSON.stringify(response))
            }
        };

    var select = UIkit.uploadSelect($("#upload-select"), settings),
        drop   = UIkit.uploadDrop($("#upload-drop"), settings);
</script>
{% endblock %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "bucketslist")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# list buckets #}

{% extends "console/layout" %}

{% block breadcrumb %}
<ul class="uk-breadcrumb">
    <li><a href="/console/buckets">buckets</a></li>
</ul>
{% endblock %}

{% block content %}

<div class="uk-width-1-2 uk-container uk-container-center">
{% set res = SearchFiles("buckets", ctx.QueryParam("q"), ctx.QueryParam("page"), 10) %}

<ul class="uk-list uk-list-space">
    <li>
        <div class="uk-float-right">
            <span class="uk-text-muted">{{ (res.PerPage*res.CurrentPage+1)}} to {{ (res.PerPage*(res.CurrentPage+1)+1)}} of {{ res.Total }}</span>
            <div class="uk-button-group">
                {% if res.CurrentPage >= 1 %}
                <a href="?page={{res.CurrentPage-1}}" class="uk-button uk-button-small"><i class="uk-icon-angle-left uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-left"></i></button>
                {% endif %}

                {% if res.HasNext %}
                <a href="?page={{res.NextPage}}" class="uk-button uk-button-small"><i class="uk-icon-angle-right uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-right"></i></button>
                {% endif %}
            </div>    
        </div>
        
    
        <div class="">
            <a class="uk-button uk-button uk-button-small" href="/console/buckets/newbucket" title="add new bucket">New bucket</a>
        </div>
    </li>
{% for file in res.Files %}
    <li>
        <div class="uk-panel">
            <small class="uk-float-right uk-text-muted">{{file.UpdatedAt()|naturaltime}}</small>
            <i class="uk-icon-folder-o uk-text-muted uk-margin-right"></i>
            <a href="/console/buckets/{{ file.ID() }}/files">{{ file.Name() }}</a>
        </div>
    </li>
{% endfor %}

{% if res.HasNext %}
    <li class="uk-margin-top"><a href="?page={{res.NextPage}}">next page</a></li>
{% endif %}
</ul>

</div>

{% endblock %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "fileslist")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# list files #}

{% extends "console/layout" %}

{% block breadcrumb %}
{% set bucket = LoadByID("buckets", ctx.Get("bucket_id")) %}

<ul class="uk-breadcrumb">
    <li><a href="/console/buckets">buckets</a></li>
    <li><span><i class="uk-icon-folder-open-o uk-text-muted"></i> {{ bucket.Name() }}</span></li>
</ul>
{% endblock %}

{% block content %}

<div class="uk-width-1-2 uk-container uk-container-center">
{% set res = SearchFiles(bucket.Name(), ctx.QueryParam("q"), ctx.QueryParam("page"), 10) %}


<ul class="uk-list uk-list-space">
    <li>
        <div class="uk-float-right">
            <span class="uk-text-muted">{{ (res.PerPage*res.CurrentPage+1)}} to {{ (res.PerPage*(res.CurrentPage+1)+1)}} of {{ res.Total }}</span>
            <div class="uk-button-group">
                {% if res.CurrentPage >= 1 %}
                <a href="?page={{res.CurrentPage-1}}" class="uk-button uk-button-small"><i class="uk-icon-angle-left uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-left"></i></button>
                {% endif %}

                {% if res.HasNext %}
                <a href="?page={{res.NextPage}}" class="uk-button uk-button-small"><i class="uk-icon-angle-right uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-right"></i></button>
                {% endif %}
            </div>    
        </div>
        
    
        <div class="">
            <a class="uk-button uk-button uk-button-small" href="/console/buckets/{{ bucket.ID() }}/newfile" title="add new file">New file</a>
        </div>
    </li>
{% for file in res.Files %}
    <li>
        <div class="uk-panel">
            <small class="uk-float-right uk-text-muted">{{file.UpdatedAt()|naturaltime}}</small>
            <i class="uk-icon-file-o uk-text-muted uk-margin-right"></i>
            <a href="/console/buckets/{{ bucket.ID() }}/files/{{ file.ID()}}/view">{{ file.Name() }}</a>
        </div>
    </li>
{% endfor %}

{% if res.HasNext %}
    <li class="uk-margin-top"><a href="?page={{res.NextPage}}">next page</a></li>
{% endif %}
</ul>
</div>
{% endblock %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "viewfile")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# viewfile #}

{% extends "console/layout" %}

{% block breadcrumb %}

{% set bucket = LoadByID("buckets", ctx.Get("bucket_id")) %}
{% set file = LoadByID(bucket.Name(), ctx.Get("file_id")) %}

<ul class="uk-breadcrumb">
    <li><a href="/console/buckets">buckets</a></li>
    <li><a href="/console/buckets/{{ bucket.ID() }}/files"><i class="uk-icon-folder-open-o uk-text-muted"></i> {{ bucket.Name() }}</a></li>
    <li><a class="uk-icon-hover uk-icon uk-icon-plus" href="/console/buckets/{{ bucket.ID() }}/newfile" title="add new file"></a></li>
    <li><span><i class="uk-icon-file-o uk-text-muted"></i> {{ file.Name() }}</span></li>
</ul>

{% endblock %}

{% block content %}

<div id="upload-drop" class="uk-width-1-1 uk-container uk-container-center">

    <div class="uk-panel uk-margin-bottom">
        <dl class="uk-description-list-horizontal uk-text-muted uk-float-right">
            <dt>Bucket ID</dt>
            <dd>{{ bucket.ID() }}</dd>

            <dt>Bucket name</dt>
            <dd>{{ bucket.Name() }}</dd>

            <dt>File ID</dt>
            <dd>{{ file.ID() }}</dd>

            <dt>File name</dt>
            <dd>{{ file.Name() }}</dd>

            <dt>Created</dt>
            <dd><time title="{{file.CreatedAt()}}">{{file.CreatedAt()|naturaltime}}</time></dd>
            
            <dt>Updated</dt>
            <dd><time title="{{file.UpdatedAt()}}">{{file.UpdatedAt()|naturaltime}}</time></dd>

            <dt></dt>
            <dd><a href="#">Delete file</a></dd>
        </dl>
    </div>

    <div class="uk-tab-center">
        <ul class="uk-tab" data-uk-tab="{connect:'#tab-content'}">
            <li class="uk-active"><a href="#">Properties</a></li>
            <li><a href="#">Raw data</a></li>
            <li><a href="#">Struct data</a></li>
        </ul>
    </div>
    <ul id="tab-content" class="uk-switcher uk-margin">
        <li>
            <form class="uk-form" action="/console/buckets/{{ bucket.ID() }}/files/{{ file.ID() }}/properties/put" method="post" enctype="multipart/form-data">
                <input type="hidden" name="FileID" value="{{ file.ID() }}" />
                <input type="hidden" name="BucketID" value="{{ bucket.ID() }}" />

                <div class="uk-form-row">
                    <label class="uk-form-label">Content type</label>
                    <div class="uk-form-controls">
                        <div class="uk-autocomplete uk-form uk-width-1-1" id="select_contenttype-{{file.ID()}}" data-uk-autocomplete>
                            <input type="text" class="uk-width-1-1" name="ContentType" autocomplete="off" value="{{ file.ContentType() }}">
                        </div>
                    </div>
                </div>

                <div class="uk-form-row">
                    <div class="uk-form-controls">
                        <button class="uk-button" name="Mode" value="Manual">Save</button>
                    </div>
                </div>
            </form>
        </li>
        <li>
            {% if file.IsText() %}
            <form class="uk-form" action="/console/buckets/{{ bucket.ID() }}/files/{{ file.ID() }}/textdata/put" method="post" enctype="multipart/form-data">
                <input type="hidden" name="FileID" value="{{ file.ID() }}" />
                <input type="hidden" name="BucketID" value="{{ bucket.ID() }}" />

                <div class="uk-form-row">
                    <div class="uk-form-controls">
                        <input id="rawdata_source-{{file.ID()}}" type="hidden" name="TextData" value="{{file.TextData()|escapejs}}" />
                        <div id="rawdata_view-{{file.ID()}}" class="uk-width-1-1" style="height:400px; border: 1px solid #ddd;"></div>
                    </div>
                </div>

                

                <div class="uk-form-row">
                    <div class="uk-form-controls">
                        <button class="uk-button" name="Mode" value="Manual">Save</button>
                    </div>
                </div>
            </form>
            {% endif %}

            {% if file.IsImage() or file.IsRaw() %}
            
            {% if file.IsImage() %}
            <div class="uk-panel-teaser uk-text-center">
                <img id="image-preview" class="uk-thumbnail" src="{{ file.Name()|urlfile }}" style="max-height: 250px;" alt="" />
            </div>
            {% endif %}
            <div> 
                <div class="uk-placeholder uk-text-center">
                    <i class="uk-icon-cloud-upload uk-icon-medium uk-text-muted uk-margin-small-right"></i> Attach binaries by dropping them here or <a class="uk-form-file">selecting one<input id="upload-select" type="file" name="BinData" /></a>.
                </div>
                <div id="progressbar" class="uk-progress uk-hidden">
                    <div class="uk-progress-bar" style="width: 0%;">0%</div>
                </div>
            </div>
            {% endif %}
        </li>
        <li>
            <form class="uk-form" action="/console/buckets/{{ bucket.ID() }}/files/{{ file.ID() }}/structdata/put" method="post" enctype="multipart/form-data">
                <input type="hidden" name="FileID" value="{{ file.ID() }}" />
                <input type="hidden" name="BucketID" value="{{ bucket.ID() }}" />

                <div class="uk-form-row">
                    <div class="uk-form-controls">
                        <input id="json_source-{{file.ID()}}" type="hidden" name="StructDataJson" value="{{file.MapData()|atojs|escapejs}}" />
                        <div id="json_view-{{file.ID()}}" class="uk-width-1-1" style="height:400px; border: 1px solid #ddd;"></div>
                    </div>
                </div>

                <div class="uk-form-row">
                    <div class="uk-form-controls">
                        <button class="uk-button" name="Mode" value="Manual">Save</button>
                    </div>
                </div>
            </form>
        </li>
    </ul>

</div>
{% endblock %}

{% block scripts %}
<script src="//cdn.jsdelivr.net/uikit/2.26.2/js/components/autocomplete.min.js"></script>
<script>
    var fileId = "{{ ctx.Get("file_id")|escapejs }}";
    var file = {{ file|atojs }};
    
    var progressbar = $("#progressbar"),
        bar         = progressbar.find('.uk-progress-bar'),
        settings    = {

            action: '/console/buckets/{{ bucket.ID() }}/files/{{ file.ID() }}/rawdata/put_via_uploader', // upload url
            param: 'BinData',
            params: {
                FileID: "{{ file.ID() }}",
                BucketID: "{{ bucket.ID() }}", 
            },
            type: 'json',

            allow : '*', // allow all file types

            loadstart: function() {
                bar.css("width", "0%").text("0%");
                progressbar.removeClass("uk-hidden");
            },

            progress: function(percent) {
                percent = Math.ceil(percent);
                bar.css("width", percent+"%").text(percent+"%");
            },

            allcomplete: function(response) {

                bar.css("width", "100%").text("100%");

                setTimeout(function(){
                    progressbar.addClass("uk-hidden");
                }, 250);

                location.reload();
            }
        };

    var select = UIkit.uploadSelect($("#upload-select"), settings),
        drop   = UIkit.uploadDrop($("#upload-drop"), settings);
</script>
<script>
    var autocomplete = UIkit.autocomplete("#select_contenttype-{{file.ID()}}", {
                source: [
                {
                    value: 'image/jpeg'
                }, {
                    value: 'image/pjpeg'
                }, {
                    value: 'image/png'
                }, {
                    value: 'image/vnd.microsoft.icon'
                }, {
                    value: 'image/gif'
                }, {
                    value: 'text/css'
                }, {
                    value: 'text/plain'
                }, {
                    value: 'text/javascript'
                }, {
                    value: 'text/html'
                }, {
                    value: 'text/toml'
                }, {
                    value: 'application/javascript'
                }, {
                    value: 'application/json'
                }, {
                    value: 'application/soap+xml'
                }, {
                    value: 'application/xhtml+xml'
                }, {
                    value: 'text/csv'
                }, {
                    value: 'text/x-jquery-tmpl'
                }, {
                    value: 'text/php'
                }, {
                    value: 'application/x-javascript'
                }],
                minLength: 2,
                delay: 0
            });

    {% comment %}
    json editor
    {% endcomment %}
    
    var jsonEditor;
    var container = document.getElementById("json_view-{{file.ID()}}");
    // https://github.com/josdejong/jsoneditor/blob/master/docs/api.md
    var options = {
        mode: "tree", // Available values: 'tree' (default), 'view', 'form', 'code', 'text'. In 'view' mode, the data and datastructure is read-only. In 'form' mode, only the value can be changed, the datastructure is read-only. Mode 'code' requires the Ace editor to be loaded on the page. Mode 'text' shows the data as plain text.
        modes: ["tree"],
        // modes: ["tree", "code", "form", "view", "text"],
        schema: {
            "title": "root",
            "type": "object",
            "properties": {
                "fullName": {
                    "type": "string"
                },
                "licenses": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "minItems": 2,
                    "uniqueItems": true
                }
            },
            "required": ["fullName", "licenses"]
        },
        indentation: 4,
        ace: ace,

        search: false,
        onChange: function() {
            console.log(jsonEditor.get());
            document.getElementById("json_source-{{file.ID()}}").value = JSON.stringify(jsonEditor.get()); 
        }
    };
    jsonEditor = new JSONEditor(container, options);
    jsonEditor.set({{file.MapData()|atojs}});

    {% comment %}
    raw text editor
    {% endcomment %}

    {% if file.IsText() %}
    var mapContentTypes = JSON.parse(document.getElementById("mapContentTypeToACEType").textContent);

    var textEditor = ace.edit("rawdata_view-{{file.ID()}}");
    textEditor.setFontSize("11pt");
    textEditor.getSession().setMode(mapContentTypes["{{ file.ContentType()|escapejs }}"]);
    textEditor.setValue("{{file.TextData()|escapejs}}");
    textEditor.getSession().on('change', function() {
        document.getElementById("rawdata_source-{{file.ID()}}").value = textEditor.getSession().getValue();
    });
    
    textEditor.$blockScrolling = Infinity;
    {% endif %}
</script>
{% endblock %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "putfile_textdata")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# putfile_textdata #}

{% if ctx.IsPost() %}

{% set bucketID = ctx.FormValue("BucketID") %}
{{ bucketID }}
{% set bucket = LoadByID("buckets", bucketID) %}

{% set fileID = ctx.FormValue("FileID") %}
{{ fileID }}
{% set file = LoadByID(bucket.Name(), fileID) %}

{{ file.SetTextData(ctx.FormValue("TextData")) }}

{{ file.Sync() }}

{# TODO: check error #}

{% set url = "/console/buckets/"|add:bucket.ID()|add:"/files/"|add:file.ID()|add:"/view" %}

{{ ctx.Redirect302(url) }}
{% endif %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "putfile_structdata")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# putfile_structdata #}

{% if ctx.IsPost() %}

{% set bucketID = ctx.FormValue("BucketID") %}
{{ bucketID }}
{% set bucket = LoadByID("buckets", bucketID) %}

{% set fileID = ctx.FormValue("FileID") %}
{{ fileID }}
{% set file = LoadByID(bucket.Name(), fileID) %}

{% set m = M() %}
{{ m.LoadFrom(ctx.FormValue("StructDataJson")) }}

{{ file.SetMapData(m.ToMap()) }}

{{ file.Sync() }}

{# TODO: check error #}

{% set url = "/console/buckets/"|add:bucket.ID()|add:"/files/"|add:file.ID()|add:"/view" %}

{{ ctx.Redirect302(url) }}
{% endif %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "putfile_properties")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# putfile_properties #}

{% if ctx.IsPost() %}

{% set bucketID = ctx.FormValue("BucketID") %}
{{ bucketID }}
{% set bucket = LoadByID("buckets", bucketID) %}

{% set fileID = ctx.FormValue("FileID") %}
{{ fileID }}
{% set file = LoadByID(bucket.Name(), fileID) %}


{{ file.MMeta().Set("ContentType", ctx.FormValue("ContentType")) }}

{{ file.Sync() }}

{# TODO: check error #}

{% set url = "/console/buckets/"|add:bucket.ID()|add:"/files/"|add:file.ID()|add:"/view" %}

{{ ctx.Redirect302(url) }}
{% endif %}`))

		file.Sync()
	}

	// 	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "newfileform")
	// 	if err == dbox.ErrNotFound {
	// 		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

	// 		file.SetContentType("text/html")
	// 		file.RawData().Write([]byte(`{# newfileform #}

	// {% extends "console/layout" %}

	// {% block breadcrumb %}

	// {% set bucket = LoadByID("buckets", ctx.Get("bucket_id")) %}
	// {% set file = LoadByID(bucket.Name(), ctx.Get("file_id")) %}

	// <ul class="uk-breadcrumb">
	//     <li><a href="/console/buckets">buckets</a></li>
	//     <li><a href="/console/buckets/{{ bucket.ID() }}/files"><i class="uk-icon-folder-open-o uk-text-muted"></i> {{ bucket.Name() }}</a></li>
	//     <li><span>new file</span></li>
	// </ul>

	// {% endblock %}

	// {% block content %}

	// <div class="uk-width-1-2 uk-container uk-container-center">

	//     <div class="uk-panel uk-margin-bottom">
	//         <dl class="uk-description-list-horizontal uk-text-muted uk-float-right">
	//             <dt>Bucket ID</dt>
	//             <dd>{{ bucket.ID() }}</dd>

	//             <dt>Bucket name</dt>
	//             <dd>{{ bucket.Name() }}</dd>
	//         </dl>
	//     </div>

	//     <form class="uk-form" action="/console/buckets/{{ bucket.ID() }}/files" method="post" enctype="multipart/form-data">
	//         <input type="hidden" name="BucketID" value="{{ bucket.ID() }}" />

	//         <legend>Main</legend>
	//         <div class="uk-form-row">
	//             <label class="uk-form-label">Name</label>
	//             <div class="uk-form-controls">
	//                 <input type="text" class="uk-width-1-1 uk-form-large" name="Name" autofocus />
	//             </div>
	//         </div>

	//         <div class="uk-form-row">
	//             <label class="uk-form-label">Content type</label>
	//             <div class="uk-form-controls">
	//                 <div class="uk-autocomplete uk-form uk-width-1-1" id="select_contenttype-new_file" data-uk-autocomplete>
	//                     <input type="text" class="uk-width-1-1" name="ContentType" autocomplete="off">
	//                 </div>
	//             </div>
	//         </div>

	//         <legend>Data</legend>

	//         <div class="uk-form-row">
	//             <div class="uk-form-controls">
	//                 <input id="rawdata_source-new_file" type="hidden" name="TextData" value="" />
	//                 <div id="rawdata_view-new_file" class="uk-width-1-1" style="height:400px; border: 1px solid #ddd;"></div>
	//             </div>
	//         </div>

	//         <legend>Struct data</legend>

	//        <div class="uk-form-row">
	//             <div class="uk-form-controls">
	//                 <input id="json_source-new_file" type="hidden" name="StructDataJson" value="" />
	//                 <div id="json_view-new_file" class="uk-width-1-1" style="height:400px; border: 1px solid #ddd;"></div>
	//             </div>
	//         </div>

	//         <div class="uk-form-row">
	//             <div class="uk-form-controls">
	//                 <button class="uk-button" name="Mode" value="Manual">Create</button>
	//             </div>
	//         </div>
	//     </form>

	// </div>
	// {% endblock %}

	// {% block scripts %}
	// <script src="//cdn.jsdelivr.net/uikit/2.26.2/js/components/autocomplete.min.js"></script>
	// <script src="//cdnjs.cloudflare.com/ajax/libs/jsoneditor/5.5.5/jsoneditor.min.js"></script>
	// <script>
	//     var selectedContentType = "";
	//     var autocomplete = UIkit.autocomplete("#select_contenttype-new_file", {
	//                 source: [
	//                 {
	//                     value: 'image/jpeg'
	//                 }, {
	//                     value: 'image/pjpeg'
	//                 }, {
	//                     value: 'image/png'
	//                 }, {
	//                     value: 'image/vnd.microsoft.icon'
	//                 }, {
	//                     value: 'image/gif'
	//                 }, {
	//                     value: 'text/css'
	//                 }, {
	//                     value: 'text/plain'
	//                 }, {
	//                     value: 'text/javascript'
	//                 }, {
	//                     value: 'text/html'
	//                 }, {
	//                     value: 'text/toml'
	//                 }, {
	//                     value: 'application/javascript'
	//                 }, {
	//                     value: 'application/json'
	//                 }, {
	//                     value: 'application/soap+xml'
	//                 }, {
	//                     value: 'application/xhtml+xml'
	//                 }, {
	//                     value: 'text/csv'
	//                 }, {
	//                     value: 'text/x-jquery-tmpl'
	//                 }, {
	//                     value: 'text/php'
	//                 }, {
	//                     value: 'application/x-javascript'
	//                 }],
	//                 minLength: 2,
	//                 delay: 0
	//             });

	//     autocomplete.on("selectitem.uk.autocomplete", function(event, data, acobject) {
	//         selectedContentType = data.value;
	//     }).on("change", function(event, data, acobject) {
	//         selectedContentType = event.target.value;
	//     })

	//     {% comment %}
	//     json editor
	//     {% endcomment %}
	//     var jsonEditor;
	//     var container = document.getElementById("json_view-new_file");
	//     // https://github.com/josdejong/jsoneditor/blob/master/docs/api.md
	//     var options = {
	//         mode: "tree", // Available values: 'tree' (default), 'view', 'form', 'code', 'text'. In 'view' mode, the data and datastructure is read-only. In 'form' mode, only the value can be changed, the datastructure is read-only. Mode 'code' requires the Ace editor to be loaded on the page. Mode 'text' shows the data as plain text.
	//         modes: ["tree", "code", "form", "view", "text"],
	//         schema: {
	//             "title": "root",
	//             "type": "object",
	//             "properties": {
	//                 "fullName": {
	//                     "type": "string"
	//                 },
	//                 "licenses": {
	//                     "type": "array",
	//                     "items": {
	//                         "type": "string"
	//                     },
	//                     "minItems": 2,
	//                     "uniqueItems": true
	//                 }
	//             },
	//             "required": ["fullName", "licenses"]
	//         },
	//         indentation: 4,

	//         search: false,
	//         onChange: function() {
	//             console.log(jsonEditor.get());
	//             document.getElementById("json_source-new_file").value = JSON.stringify(jsonEditor.get());
	//         }
	//     };
	//     jsonEditor = new JSONEditor(container, options);
	//     jsonEditor.set({});

	//     {% comment %}
	//     raw text editor
	//     {% endcomment %}

	//     var textEditor = ace.edit("rawdata_view-new_file");
	//     textEditor.setFontSize("11pt");
	//     textEditor.setValue("");
	//     textEditor.getSession().on('change', function() {
	//         document.getElementById("rawdata_source-new_file").value = textEditor.getSession().getValue();
	//     });

	//     textEditor.$blockScrolling = Infinity;
	// </script>
	// {% endblock %}`))

	// 		file.Sync()
	// 	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "newfileform")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# newfileform #}

{% extends "console/layout" %}

{% block breadcrumb %}

{% set bucket = LoadByID("buckets", ctx.Get("bucket_id")) %}
{% set file = LoadByID(bucket.Name(), ctx.Get("file_id")) %}

<ul class="uk-breadcrumb">
    <li><a href="/console/buckets">buckets</a></li>
    <li><a href="/console/buckets/{{ bucket.ID() }}/files"><i class="uk-icon-folder-open-o uk-text-muted"></i> {{ bucket.Name() }}</a></li>
    <li><span>new file</span></li>
</ul>

{% endblock %}

{% block content %}

<div class="uk-width-1-2 uk-container uk-container-center">

    <div class="uk-panel uk-margin-bottom">
        <dl class="uk-description-list-horizontal uk-text-muted uk-float-right">
            <dt>Bucket ID</dt>
            <dd>{{ bucket.ID() }}</dd>

            <dt>Bucket name</dt>
            <dd>{{ bucket.Name() }}</dd>
        </dl>
    </div>

    <form class="uk-form" action="/console/buckets/{{ bucket.ID() }}/files" method="post" enctype="multipart/form-data">
        <input type="hidden" name="BucketID" value="{{ bucket.ID() }}" />

        <legend>Main</legend>
        <div class="uk-form-row">
            <label class="uk-form-label">Name</label>
            <div class="uk-form-controls">
                <input type="text" class="uk-width-1-1 uk-form-large" name="Name" autofocus />
            </div>
        </div>

        <div class="uk-form-row">
            <label class="uk-form-label">Content type</label>
            <div class="uk-form-controls">
                <div class="uk-autocomplete uk-form uk-width-1-1" id="select_contenttype-new_file" data-uk-autocomplete>
                    <input type="text" class="uk-width-1-1" name="ContentType" autocomplete="off">
                </div>
            </div>
        </div>

        <legend>Data</legend>

        <div class="uk-form-row">
            <div class="uk-form-controls">
                <input id="rawdata_source-new_file" type="hidden" name="TextData" value="" />
                <div id="rawdata_view-new_file" class="uk-width-1-1" style="height:400px; border: 1px solid #ddd;"></div>
            </div>
        </div>

        <legend>Struct data</legend>

       <div class="uk-form-row">
            <div class="uk-form-controls">
                <input id="json_source-new_file" type="hidden" name="StructDataJson" value="" />
                <div id="json_view-new_file" class="uk-width-1-1" style="height:400px; border: 1px solid #ddd;"></div>
            </div>
        </div>


        <div class="uk-form-row">
            <div class="uk-form-controls">
                <button class="uk-button" name="Mode" value="Manual">Create</button>
            </div>
        </div>
    </form>

</div>
{% endblock %}

{% block scripts %}
<script src="//cdn.jsdelivr.net/uikit/2.26.2/js/components/autocomplete.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/jsoneditor/5.5.5/jsoneditor.min.js"></script>
<script>
    var selectedContentType = "";
    var autocomplete = UIkit.autocomplete("#select_contenttype-new_file", {
                source: [
                {
                    value: 'image/jpeg'
                }, {
                    value: 'image/pjpeg'
                }, {
                    value: 'image/png'
                }, {
                    value: 'image/vnd.microsoft.icon'
                }, {
                    value: 'image/gif'
                }, {
                    value: 'text/css'
                }, {
                    value: 'text/plain'
                }, {
                    value: 'text/javascript'
                }, {
                    value: 'text/html'
                }, {
                    value: 'text/toml'
                }, {
                    value: 'application/javascript'
                }, {
                    value: 'application/json'
                }, {
                    value: 'application/soap+xml'
                }, {
                    value: 'application/xhtml+xml'
                }, {
                    value: 'text/csv'
                }, {
                    value: 'text/x-jquery-tmpl'
                }, {
                    value: 'text/php'
                }, {
                    value: 'application/x-javascript'
                }],
                minLength: 2,
                delay: 0
            });

    autocomplete.on("selectitem.uk.autocomplete", function(event, data, acobject) {
        selectedContentType = data.value;
    }).on("change", function(event, data, acobject) {
        selectedContentType = event.target.value;
    })

    {% comment %}
    json editor
    {% endcomment %}
    var jsonEditor;
    var container = document.getElementById("json_view-new_file");
    // https://github.com/josdejong/jsoneditor/blob/master/docs/api.md
    var options = {
        mode: "tree", // Available values: 'tree' (default), 'view', 'form', 'code', 'text'. In 'view' mode, the data and datastructure is read-only. In 'form' mode, only the value can be changed, the datastructure is read-only. Mode 'code' requires the Ace editor to be loaded on the page. Mode 'text' shows the data as plain text.
        modes: ["tree", "code", "form", "view", "text"],
        schema: {
            "title": "root",
            "type": "object",
            "properties": {
                "fullName": {
                    "type": "string"
                },
                "licenses": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "minItems": 2,
                    "uniqueItems": true
                }
            },
            "required": ["fullName", "licenses"]
        },
        indentation: 4,

        search: false,
        onChange: function() {
            console.log(jsonEditor.get());
            document.getElementById("json_source-new_file").value = JSON.stringify(jsonEditor.get()); 
        }
    };
    jsonEditor = new JSONEditor(container, options);
    jsonEditor.set({});

    {% comment %}
    raw text editor
    {% endcomment %}

    var textEditor = ace.edit("rawdata_view-new_file");
    textEditor.setFontSize("11pt");
    textEditor.setValue("");
    textEditor.getSession().on('change', function() {
        document.getElementById("rawdata_source-new_file").value = textEditor.getSession().getValue();
    });
    
    textEditor.$blockScrolling = Infinity;
</script>
{% endblock %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "newbucketform")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# new bucket form #}

{% extends "console/layout" %}

{% block breadcrumb %}

<ul class="uk-breadcrumb">
    <li><a href="/console/buckets">buckets</a></li>
    <li><span>new bucket</span></li>
</ul>

{% endblock %}

{% block content %}

<div class="uk-width-1-2 uk-container uk-container-center">
    <form class="uk-form" action="/console/buckets" method="post" enctype="multipart/form-data">
        <legend>Main</legend>
        <div class="uk-form-row">
            <label class="uk-form-label">Name</label>
            <div class="uk-form-controls">
                <input type="text" class="uk-width-1-1 uk-form-large" name="Name" autofocus />
            </div>
        </div>
        
        <legend>Store types <small>(allowed memory, local or boltdb)</small></legend>
        <div class="uk-form-row">
            <label class="uk-form-label">Meta data files</label>
            <div class="uk-form-controls">
                <div class="uk-autocomplete uk-form uk-width-1-1" data-uk-autocomplete="{delay: 0, minLength: 1, source:[{value:'memory'}, {value:'local'}, {value:'boltdb'}]}">
                    <input type="text" class="uk-width-1-1" name="MetaDataStoreType" autocomplete="off" value="">
                </div>
            </div>
            <p class="uk-form-help-block">
                <input id="SameAsMetaStoreTypeChekbox" type="checkbox" name="SameAsMetaStoreType" />
                <label for="SameAsMetaStoreTypeChekbox" class="uk-form-label">similarly to map data and raw data</label>
            </p>
            <p class="uk-form-help-block">
                <input id="MetaHaveSuffixChekbox" type="checkbox" name="MetaHaveSuffix" />
                <label for="MetaHaveSuffixChekbox" class="uk-form-label">will have the suffix</label>
            </p>
        </div>
        <div class="uk-form-row">
            <label class="uk-form-label">Map data files</label>
            <div class="uk-form-controls">
                <div class="uk-autocomplete uk-form uk-width-1-1" data-uk-autocomplete="{delay: 0, minLength: 1, source:[{value:'memory'}, {value:'local'}, {value:'boltdb'}]}">
                    <input type="text" class="uk-width-1-1" name="MapDataStoreType" autocomplete="off" value="">
                </div>
            </div>
            <p class="uk-form-help-block">
                <input id="MapDataHaveSuffixChekbox" type="checkbox" name="MapDataHaveSuffix" />
                <label for="MapDataHaveSuffixChekbox" class="uk-form-label">will have the suffix</label>
            </p>
        </div>
        <div class="uk-form-row">
            <label class="uk-form-label">Raw data files</label>
            <div class="uk-form-controls">
                <div class="uk-autocomplete uk-form uk-width-1-1" data-uk-autocomplete="{delay: 0, minLength: 1, source:[{value:'memory'}, {value:'local'}, {value:'boltdb'}]}">
                    <input type="text" class="uk-width-1-1" name="RawDataStoreType" autocomplete="off" value="">
                </div>
            </div>
            <p class="uk-form-help-block">
                <input id="RawDataHaveSuffixChekbox" type="checkbox" name="RawDataHaveSuffix" />
                <label for="RawDataHaveSuffixChekbox" class="uk-form-label">will have the suffix</label>
            </p>
        </div>

        <legend>
            Mapping map data files
            <a 
                class="uk-link-muted"
                href="https://www.elastic.co/guide/en/elasticsearch/guide/current/mapping-intro.html" 
                target="_blank">
                <i class="uk-icon uk-icon-question-circle"></i>
            </a>
        </legend>

        <div class="uk-form-row">
            <div class="uk-form-controls">
                <input id="json_source-new_file" type="hidden" name="MappingMapDataFiles" value="" />
                <div id="json_view-new_file" class="uk-width-1-1" style="height:400px; border: 1px solid #ddd;"></div>
            </div>
        </div>


        <div class="uk-form-row">
            <div class="uk-form-controls">
                <button class="uk-button" value="Manual">Create bucket</button>
            </div>
        </div>
    </form>

</div>
{% endblock %}

{% block scripts %}
<script src="//cdn.jsdelivr.net/uikit/2.26.2/js/components/autocomplete.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/jsoneditor/5.5.5/jsoneditor.min.js"></script>
<script>
    $("#SameAsMetaStoreTypeChekbox").on("change", function(e) {
        if ($(e.target).prop("checked")) {
            $("input[name='MapDataStoreType']").prop("disabled", true);
            $("input[name='RawDataStoreType']").prop("disabled", true);
            $("input[name='MapDataStoreType']").val($("input[name='MetaDataStoreType']").prop("value"));
            $("input[name='RawDataStoreType']").val($("input[name='MetaDataStoreType']").prop("value"));
        } else {
            $("input[name='MapDataStoreType']").prop( "disabled", false );
            $("input[name='RawDataStoreType']").prop( "disabled", false );
        }
    })
    {% comment %}
    json editor
    {% endcomment %}
    var jsonEditor;
    var container = document.getElementById("json_view-new_file");
    // https://github.com/josdejong/jsoneditor/blob/master/docs/api.md
    var options = {
        mode: "tree", // Available values: 'tree' (default), 'view', 'form', 'code', 'text'. In 'view' mode, the data and datastructure is read-only. In 'form' mode, only the value can be changed, the datastructure is read-only. Mode 'code' requires the Ace editor to be loaded on the page. Mode 'text' shows the data as plain text.
        modes: ["tree", "code", "form", "view", "text"],
        schema: {
            "title": "root",
            "type": "object",
            // todo: schema mapping bucket
            "properties": {
                "fullName": {
                    "type": "string"
                },
                "licenses": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "minItems": 2,
                    "uniqueItems": true
                }
            },
            "required": ["fullName", "licenses"]
        },
        indentation: 4,

        search: false,
        onChange: function() {
            console.log(jsonEditor.get());
            document.getElementById("json_source-new_file").value = JSON.stringify(jsonEditor.get()); 
        }
    };
    jsonEditor = new JSONEditor(container, options);
    jsonEditor.set({});
</script>
{% endblock %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "newfile")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")
		file.RawData().Write([]byte(`{# newfile #}


{% if ctx.IsPost() %}

{% set bucketID = ctx.FormValue("BucketID") %}
{% set fileName = ctx.FormValue("Name") %}
{% set bucket = LoadByID("buckets", bucketID) %}

{% set file = Load(bucket.Name(), fileName) %}

{{ file.SetContentType(ctx.FormValue("ContentType")) }}
{{ file.SetTextData(ctx.FormValue("TextData")) }}
{{ file.MMapData().LoadFrom(ctx.FormValue("StructDataJson")) }}

{{ file.Sync() }}
{# TODO: check error #}

{% set url = "/console/buckets/"|add:bucket.ID()|add:"/files/"|add:file.ID()|add:"/view" %}

{{ ctx.Redirect302(url) }}
{% endif %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "newbucket")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.RawData().Write([]byte(`{# new bucket #}

{% if ctx.IsPost() %}

{% set opt = ctx.BindFormToMap("Name", "MetaDataStoreType", "SameAsMetaStoreType", "MapDataStoreType", "RawDataStoreType", "MappingMapDataFiles", "MetaHaveSuffix", "MapDataHaveSuffix", "RawDataHaveSuffix" ) %}

{% set res = CreateBucket(opt) %}
{# res is error or bucket #}


{% set url = "/console/buckets/"|add:res.ID()|add:"/files" %}
{{ ctx.Redirect302(url) }}
{% comment %}
{% endcomment %}

{% endif %}`))

		file.Sync()
	}

	file, err = store.LoadOrNewFile(api.ConsoleBucketName, "putfile_rawdata_viauploader")
	if err == dbox.ErrNotFound {
		fmt.Printf("%v: create %q\n", file.Bucket(), file.Name())

		file.RawData().Write([]byte(`{# put rawdata of file via uploader #}

{% if ctx.IsPost() %}

{% set bucketID = ctx.FormValue("BucketID") %}
{{ bucketID }}
{% set bucket = LoadByID("buckets", bucketID) %}

{% set fileID = ctx.FormValue("FileID") %}
{{ fileID }}
{% set file = LoadByID(bucket.Name(), fileID) %}

{% set fileData = ctx.FormFileData("BinData") %}

{{ file.SetRawData(fileData.Data) }}
{{ file.SetContentType(fileData.ContentType) }}
{{ file.MMapData().Set("OrigName", fileData.Name) }}

{{ file.Sync() }}

{# TODO: check error #}

{% set url = "/console/buckets/"|add:bucket.ID()|add:"/files/"|add:file.ID()|add:"/view" %}

{{ ctx.Redirect302(url) }}
{% endif %}`))

		file.Sync()
	}
}
