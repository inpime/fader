package main

import (
	"addons/session"
	"api"
	"api/config"
	"encoding/base64"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/inpime/dbox"
	"store"
	"time"
	"utils"
)

func initStroes() {
	db, err := bolt.Open(api.Cfg.Store.BoltDBFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		panic(err)
	}

	// set default buckets store for buckets
	dbox.BucketStore = store.NewBoltDBStore(db, config.BucketsBucketName)
	store.BoltDBClient = db

	// ---------------------
	// Create buckets:
	// * Settings
	// * Static
	// * Users
	// * Console
	// ---------------------

	flagInitViaFixtures := false // to dev a variable

	isNewInstallation := false

	// root bucket buckets
	bucket, err := store.BucketByName(config.BucketsBucketName)
	logrus.Debugf("buckets: meta %#v, %v\n", bucket.MapData(), err)
	bucket.InitInOneStore(dbox.BoltDBStoreType)
	if err == dbox.ErrNotFound {
		logrus.Infof("buckets: create bucket %q\n", bucket.Name())

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		bucket.MMapDataFilesMapping().LoadFromM(store.BucketMapMappingDefault)

		bucket.UpdateMapping()
		bucket.Sync()

		isNewInstallation = true
	}

	if !flagInitViaFixtures {
		if isNewInstallation {
			logrus.Info("The first run. Installation of the console panel...")

			if err := api.AppImportFromLastArchive(); err != nil {
				panic(err)
			}
		} else {
			logrus.Info("Existing application. Initializing settings....")
			if err := api.AppStoresInitFromExistBuckets(); err != nil {
				panic(err)
			}
		}

		return
	}

	// if flagInitViaFixtures == true

	// TODO: anything below removed after careful testing

	// ----------------
	// Bucket: Settings
	// ----------------

	bucket, err = store.BucketByName(config.SettingsBucketName)
	bucket.InitRawDataStore(dbox.BoltDBStoreType, true)  // store key - fs.settings.rawdata
	bucket.InitMetaDataStore(dbox.BoltDBStoreType, true) // store key - boltdb.settings
	bucket.InitMapDataStore(dbox.BoltDBStoreType, true)  // store key - boltdb.settings

	if err == dbox.ErrNotFound {
		logrus.Infof("buckets: create bucket %q\n", bucket.Name())

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	// ----------------
	// Bucket: Static
	// ----------------

	bucket, err = store.BucketByName(config.StaticBucketName)
	bucket.InitRawDataStore(dbox.LocalStoreType, false)
	bucket.InitMetaDataStore(dbox.BoltDBStoreType, true)
	bucket.InitMapDataStore(dbox.BoltDBStoreType, true)

	if err == dbox.ErrNotFound {
		logrus.Infof("buckets: create bucket %q\n", bucket.Name())

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	// ----------------
	// Bucket: Users
	// ----------------

	bucket, err = store.BucketByName(config.UsersBucketName)
	bucket.InitInOneStore(dbox.BoltDBStoreType)

	if err == dbox.ErrNotFound {
		logrus.Infof("buckets: create bucket %q\n", bucket.Name())

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)
		bucket.MMapDataFilesMapping().
			Set("licenses", utils.Map().
				Set("type", "string").
				Set("index", "not_analyzed")).
			Set("fullname", utils.Map().
				Set("type", "string")).
			Set("pwd", utils.Map().
				Set("type", "string").
				Set("index", "not_analyzed"))

		bucket.UpdateMapping()
		bucket.Sync()
	}

	// ----------------
	// Bucket: Console
	// ----------------

	bucket, err = store.BucketByName(config.ConsoleBucketName)
	bucket.InitRawDataStore(dbox.LocalStoreType, false)
	bucket.InitMetaDataStore(dbox.BoltDBStoreType, true)
	bucket.InitMapDataStore(dbox.BoltDBStoreType, true)

	if err == dbox.ErrNotFound {
		logrus.Infof("buckets: create bucket %q\n", bucket.Name())

		bucket.MMetaDataFilesMapping().LoadFromM(store.FileMetaMappingDefault)

		bucket.UpdateMapping()
		bucket.Sync()
	}

	// ---------------------
	// Files
	// * Settings
	//   * console.route
	//   * filecontent.route
	//   * main
	// * Static
	//   * logo.png
	// * Users
	//   * console
	// * Console
	//   * login.form
	//   * login
	//   * logout
	//
	//   * layout
	//   * dashboard
	//   * list.buckets
	//   * list.files
	//   * file.view
	//   * file.new.form
	//   * bucket.new.form
	//
	//   * file.put.textdata
	//   * file.put.rawdata.via_uploader
	//   * file.put.structdata
	//   * file.put.properties
	//   * file.new
	//   * bucket.new
	// ---------------------

	// ------------------------
	// File settings
	// ------------------------

	// console.route
	file, err := store.LoadOrNewFile(config.SettingsBucketName, "console.route")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/toml")
		file.RawData().Write([]byte(`# settings@console.route
# with routes for 'importexport'
        
# -----------------
# console routes
# -----------------

# -----------------
# Sessions
# -----------------

[[routs]]
name = "ConsoleLogout"
path = "/console/sessions/logout"
handler = "console logout"
methods = ["post"]
licenses = ["user", "admin"]

[[routs]]
name = "ConsoleLoginForm"
path = "/console/sessions/login"
handler = "console login.form"
methods = ["get"]
licenses = ["guest"]

[[routs]]
name = "ConsoleLogin"
path = "/console/sessions/login"
handler = "console login"
methods = ["post"]
licenses = ["guest"]

# -----------------
# Console
# -----------------

[[routs]]
name = "ConsoleDashboard"
path = "/console"
handler = "console dashboard"
methods = ["get"]
licenses = ["admin"]

[[routs]]
name = "ConsoleListBuckets"
path = "/console/buckets"
handler = "console list.buckets"
methods = ["get"]
licenses = ["admin"]

[[routs]]
name = "ConsoleListBucketFiles"
path = "/console/buckets/{bucket_id}/files"
handler = "console list.files"
methods = ["get"]
licenses = ["admin"]

[[routs]]
name = "ConsoleViewFile"
path = "/console/buckets/{bucket_id}/files/{file_id}/view"
handler = "console file.view"
methods = ["get"]
licenses = ["admin"]

[[routs]]
name = "ConsoleNewFileForm"
path = "/console/buckets/{bucket_id}/newfile"
handler = "console file.new.form"
methods = ["get"]
licenses = ["admin"]

[[routs]]
name = "ConsoleNewBucketForm"
path = "/console/buckets/newbucket"
handler = "console bucket.new.form"
methods = ["get"]
licenses = ["admin"]

# -----------------
# Manager
# -----------------

[[routs]]
name = "ConsoleMakeUpdateFileTextData"
path = "/console/buckets/{bucket_id}/files/{file_id}/textdata/put"
handler = "console file.put.textdata"
methods = ["post"]
licenses = ["admin"]

[[routs]]
name = "ConsoleMakeUpdateFileRawDataViaUploader"
path = "/console/buckets/{bucket_id}/files/{file_id}/rawdata/put_via_uploader"
handler = "console file.put.rawdata.via_uploader"
methods = ["post"]
licenses = ["admin"]

[[routs]]
name = "ConsoleMakeUpdateFileStructData"
path = "/console/buckets/{bucket_id}/files/{file_id}/structdata/put"
handler = "console file.put.structdata"
methods = ["post"]
licenses = ["admin"]

[[routs]]
name = "ConsoleMakeUpdateFilePropertiest"
path = "/console/buckets/{bucket_id}/files/{file_id}/properties/put"
handler = "console file.put.properties"
methods = ["post"]
licenses = ["admin"]

[[routs]]
name = "ConsoleCreateFile"
path = "/console/buckets/{bucket_id}/files"
handler = "console file.new"
methods = ["post"]
licenses = ["admin"]

[[routs]]
name = "ConsoleCreateBucket"
path = "/console/buckets"
handler = "console bucket.new"
methods = ["post"]
licenses = ["admin"]

# -------
# routes for "import export" component
# -------

[[routs]]
name = "AppExport" # ImportExportImportRouteName
path = "/console/settings/export"
handler = "importexport.export"
special = true
methods = ["get"]
licenses = ["admin"]

[[routs]]
name = "AppImport" # ImportExportExportRouteName
path = "/console/settings/import"
handler = "importexport.import"
special = true
methods = ["post"]
licenses = ["admin"]

`))
		file.Sync()
	}

	// filecontent.route
	file, err = store.LoadOrNewFile(config.SettingsBucketName, "filecontent.route")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/toml")
		file.RawData().Write([]byte(`# settings@filecontent.route

# ------------------------------
# routes for "file content" addon
# ------------------------------

[[routs]]
name = "FileContentByName" # FileContentByNameHandlerName
path = "/fc/{file:[/a-zA-Z0-9._-]+}"
handler = "filecontent.byname"
special = true
methods = ["get"]
licenses = ["guest", "user", "admin"]

[[routs]]
name = "FileContentByID" # FileContentByIDHandlerName
path = "/fci/{file:[0-9a-f]{32}}"
handler = "filecontent.byid"
special = true
methods = ["get"]
licenses = ["guest", "user", "admin"]
`))
		file.Sync()
	}

	// main
	file, err = store.LoadOrNewFile(config.SettingsBucketName, config.MainSettingsFileName)
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/toml")

		file.RawData().Write([]byte(`# settings@main

routs = ["console.route", "filecontent.route"]
pageCaching = false

[filecontent]
bucket="static"

[importexport.sys.settings]
    files = ["console.route", "filecontent.route", "main"]
[importexport.sys.console]
    all = true
[importexport.sys.users]
    files = ["guestuser", "console"]
[importexport.sys.static]
    files = ["logo.png"]
`))
		file.Sync()
	}

	// ------------------------
	// File static
	// ------------------------

	// logo.png
	file, err = store.LoadOrNewFile(config.StaticBucketName, "logo.png")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("image/png")

		imgBytes, _ := base64.StdEncoding.DecodeString(`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`)

		file.RawData().Write(imgBytes)
		file.Sync()
	}

	// ------------------------
	// File user
	// ------------------------

	// guestuser
	file, err = store.LoadOrNewFile(config.UsersBucketName, config.GuestUserFileName)
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		user := session.FileAsUser(file)
		user.MMapData().Set("fullname", "Guest").
			Set("pwd", "console")
		user.SetContentType("text/toml")
		user.AddLicense(session.GuestLicense)

		file.Sync()
	}

	// console
	file, err = store.LoadOrNewFile(config.UsersBucketName, "console")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		user := session.FileAsUser(file)
		user.MMapData().Set("fullname", "Console admin").
			Set("pwd", "console")
		user.SetContentType("text/toml")
		user.AddLicense(session.UserLicense)
		user.AddLicense(session.AdminLicense)

		file.Sync()
	}

	// ------------------------
	// File console
	// ------------------------

	// login.form
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "login.form")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@login.form #}
<html>
    <head>
        <title>Fader console. Login</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta charset="utf-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="stylesheet" href="//cdn.jsdelivr.net/uikit/2.26.2/css/uikit.almost-flat.min.css" />
    </head>
    <body class="uk-height uk-vertical-align">
        <div class="uk-vertical-align-middle uk-width-1-1">
            <div class="uk-container-center uk-width-1-2 uk-panel uk-panel-box" style="top: -50px;">
                <form class="uk-form" action="{{ URL("ConsoleLogin") }}" method="POST">
                    <fieldset>
                        <legend>Fader console</legend>
                        <div class="uk-form-row">
                            <label class="uk-form-label" for="form-login">Login</label>
                            <div class="uk-form-controls">
                                <input class="uk-width-1-1" name="email" type="text" id="form-login" placeholder="">
                            </div>
                        </div>
                        <div class="uk-form-row">
                            <label class="uk-form-label" for="form-passwd">Password</label>
                            <div class="uk-form-controls">
                                <input class="uk-width-1-1" name="password" type="password" id="form-passwd" placeholder="">
                            </div>
                        </div>
                        <div class="uk-form-row">
                            <div class="uk-form-controls">
                                <button class="uk-button">Login</button>
                            </div>
                        </div>
                    </fieldset>
                </form>  
            </div>
        </div>
    </body>
</html>`))
		file.Sync()
	}

	// login
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "login")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@login #}
{% set session = ctx.Session() %}

{% if ctx.IsPost() %}
{% set res = session.Signin(ctx.FormValue("email"), ctx.FormValue("password")) %}

{% if res|is_error %}
{{ session.AddFlash(res.Error()) }}
{{ ctx.Redirect302("/?error")}}
{% else %}
{{ ctx.Redirect302(URL("ConsoleDashboard").String())}}
{% endif %}

{% endif %}
`))
		file.Sync()
	}

	// logout
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "logout")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@logout #}
{% set session = ctx.Session() %}

{% if ctx.IsPost() %}
{{ session.Logout() }}
{{ ctx.Redirect302("/")}} # TODO: to change the URL redirect
{% endif %}
`))
		file.Sync()
	}

	// layout
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "layout")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@layout #}
<html>
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
</html>
`))
		file.Sync()
	}

	// navbar
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "navbar")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@navbar #}
<nav class="uk-navbar">

	<a class="uk-navbar-brand" href="/console">Fader</a>
	<ul class="uk-navbar-nav">
		<li class="{# uk-active #}"><a href="{{ URL("ConsoleListBuckets") }}">Buckets</a></li>
	</ul>
	
	
	<div class="uk-navbar-content"><a href="/">Go to site</a></div>

	<div class="uk-navbar-content uk-navbar-flip">
		<div class="uk-navbar-content uk-text-muted">{{ ctx.CurrentUser().Name() }}</div>
		<ul class="uk-navbar-nav">
			<li class="uk-navbar-content">
				<form class="uk-form uk-margin-remove uk-display-inline-block" action="{{ URL("ConsoleLogout") }}" method="post">
					<button class="uk-button uk-button-link"><a href="#">Logout</a></button>
				</form>
			</li>
		</ul>
	</div>

</nav>
`))
		file.Sync()
	}

	// footer
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "footer")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@footer #}
<footer style="padding: 120px 0;">
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
			<a href="http://inpime.com"><img src="{{ "logo.png"|fc }}" width="90" height="30" title="Fader is a product inpime.com" alt="Fader is a product inpime.com"></a>
		</div>

	</div>
</footer>
`))
		file.Sync()
	}

	// dashboard
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "dashboard")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@dashboard #}
{% extends "console/layout" %}

{% block breadcrumb %}
<ul class="uk-breadcrumb">
    <li><a href="{{ URL("ConsoleListBuckets") }}">buckets</a></li>
</ul>
{% endblock %}

{% block content %}
<div class="uk-width-1-2 uk-container uk-container-center">
<h1>Dashboard</h1>
<div class="uk-panel uk-panel-box">
    <div class="uk-panel-title">Import/export settings</div>
    <div class="uk-margin"> 
        <div id="upload-drop" class="uk-placeholder uk-text-center">
            <i class="uk-icon-cloud-upload uk-icon-medium uk-text-muted uk-margin-small-right"></i> Attach archive by dropping them here or <a class="uk-form-file">selecting one<input id="upload-select" type="file" name="BinData" /></a> for import settings.
        </div>
        <div id="progressbar" class="uk-progress uk-hidden">
            <div class="uk-progress-bar" style="width: 0%;">0%</div>
        </div>
    </div>
    <div class="uk-scrollable-box">
        <ul class="uk-list">
            <li><a href="{{ URL("AppExport") }}"><i class="uk-icon-download"></i> all data</a></li>
        {% for groupName in ListGroupsImportExport() %}
            <li><a href="{{ URLQuery(URL("AppExport"), "group", groupName) }}"><i class="uk-icon-download"></i> export '{{ groupName }}'</a></li>
        {% endfor %}
        <ul>
    </div>
</div>

</div>
{% endblock %}

{% block scripts %}
<script>
    var progressbar = $("#progressbar"),
        bar         = progressbar.find('.uk-progress-bar'),
        settings    = {

            action: '{{ URL("AppImport") }}', // upload url
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
                
                alert("TODO: to check the status of a import")
            }
        };

    var select = UIkit.uploadSelect($("#upload-select"), settings),
        drop   = UIkit.uploadDrop($("#upload-drop"), settings);
</script>
{% endblock %}
`))
		file.Sync()
	}

	// list.buckets
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "list.buckets")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@list.buckets #}

{% extends "console/layout" %}

{% block breadcrumb %}
<ul class="uk-breadcrumb">
    <li><a href="{{ URL("ConsoleListBuckets") }}">buckets</a></li>
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
                
                <a href="{{ URLQuery(URL("ConsoleListBucketFiles", "bucket_id", file.ID()), "page", res.CurrentPage-1) }}" class="uk-button uk-button-small"><i class="uk-icon-angle-left uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-left"></i></button>
                {% endif %}

                {% if res.HasNext %}
                <a href="{{ URLQuery(URL("ConsoleListBucketFiles", "bucket_id", file.ID()), "page", res.NextPage) }}" class="uk-button uk-button-small"><i class="uk-icon-angle-right uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-right"></i></button>
                {% endif %}
            </div>    
        </div>
        
    
        <div class="">
            <a class="uk-button uk-button uk-button-small" href="{{ URL("ConsoleNewBucketForm")}}" title="add new bucket">New bucket</a>
        </div>
    </li>
{% for file in res.Files %}
    <li>
        <div class="uk-panel">
            <small class="uk-float-right uk-text-muted">{{file.UpdatedAt()|naturaltime}}</small>
            <i class="uk-icon-folder-o uk-text-muted uk-margin-right"></i>
            <a href="{{ URL("ConsoleListBucketFiles", "bucket_id", file.ID()) }}">{{ file.Name() }}</a>
        </div>
    </li>
{% endfor %}

{% if res.HasNext %}
    <li class="uk-margin-top"><a href="{{ URLQuery(URL("ConsoleListBucketFiles", "bucket_id", file.ID()), "page", res.NextPage ) }}">next page</a></li>
{% endif %}
</ul>

</div>

{% endblock %}
`))
		file.Sync()
	}

	// list.files
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "list.files")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@list.files #}

{% extends "console/layout" %}

{% block breadcrumb %}
{% set bucket = LoadByID("buckets", ctx.Get("bucket_id")) %}

<ul class="uk-breadcrumb">
    <li><a href="{{ URL("ConsoleListBuckets") }}">buckets</a></li>
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
                <a href="{{ URLQuery(URL("ConsoleListBucketFiles", "bucket_id", file.ID()), "page", res.CurrentPage-1) }}" class="uk-button uk-button-small"><i class="uk-icon-angle-left uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-left"></i></button>
                {% endif %}

                {% if res.HasNext %}
                <a href="{{ URLQuery(URL("ConsoleListBucketFiles", "bucket_id", file.ID()), "page", res.NextPage) }}" class="uk-button uk-button-small"><i class="uk-icon-angle-right uk-margin-small-top"></i></a>
                {% else %}
                <button href="#" class="uk-button uk-button-small uk-text-muted" disabled><i class="uk-icon-angle-right"></i></button>
                {% endif %}
            </div>    
        </div>
        
    
        <div class="">
            <a class="uk-button uk-button uk-button-small" href="{{ URL("ConsoleNewFileForm", "bucket_id", bucket.ID()) }}" title="add new file">New file</a>
        </div>
    </li>
{% for file in res.Files %}
    <li>
        <div class="uk-panel">
            <small class="uk-float-right uk-text-muted">{{file.UpdatedAt()|naturaltime}}</small>
            <i class="uk-icon-file-o uk-text-muted uk-margin-right"></i>
            <a href="{{ URL("ConsoleViewFile", "bucket_id", bucket.ID(), "file_id", file.ID()) }}">{{ file.Name() }}</a>
        </div>
    </li>
{% endfor %}

{% if res.HasNext %}
    <li class="uk-margin-top"><a href="{{ URLQuery(URL("ConsoleListBucketFiles", "bucket_id", file.ID()), "page", res.NextPage) }}">next page</a></li>
{% endif %}
</ul>
</div>
{% endblock %}
`))
		file.Sync()
	}

	// file.view
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "file.view")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@file.view #}

{% extends "console/layout" %}

{% block breadcrumb %}

{% set bucket = LoadByID("buckets", ctx.Get("bucket_id")) %}
{% set file = LoadByID(bucket.Name(), ctx.Get("file_id")) %}

<ul class="uk-breadcrumb">
    <li><a href="{{ URL("ConsoleListBuckets") }}">buckets</a></li>
    <li><a href="{{ URL("ConsoleListBucketFiles", "bucket_id", bucket.ID()) }}"><i class="uk-icon-folder-open-o uk-text-muted"></i> {{ bucket.Name() }}</a></li>
    <li><a class="uk-icon-hover uk-icon uk-icon-plus" href="{{ URL("ConsoleNewFileForm", "bucket_id", bucket.ID()) }}" title="add new file"></a></li>
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
            <form class="uk-form" action="{{ URL("ConsoleMakeUpdateFilePropertiest", "bucket_id", bucket.ID(), "file_id", file.ID()) }}" method="post" enctype="multipart/form-data">
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
            <form class="uk-form" action="{{ URL("ConsoleMakeUpdateFileTextData", "bucket_id", bucket.ID(), "file_id", file.ID()) }}" method="post" enctype="multipart/form-data">
                <input type="hidden" name="FileID" value="{{ file.ID() }}" />
                <input type="hidden" name="BucketID" value="{{ bucket.ID() }}" />

                <div class="uk-form-row">
                    <div class="uk-form-controls">
                        <input id="rawdata_source-{{file.ID()}}" type="hidden" name="TextData" value="{{file.TextData()}}" />
                        <div id="rawdata_view-{{file.ID()}}" class="uk-width-1-1 uk-scrollable-box" style="height:500px; border: 1px solid #ddd;"></div>
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
                <img id="image-preview" class="uk-thumbnail" src="{{ file.Name()|fc }}" style="max-height: 250px;" alt="" />
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
            <form class="uk-form" action="{{ URL("ConsoleMakeUpdateFileStructData", "bucket_id", bucket.ID(), "file_id", file.ID()) }}" method="post" enctype="multipart/form-data">
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

            action: '{{ URL("ConsoleMakeUpdateFileRawDataViaUploader", "bucket_id", bucket.ID(), "file_id", file.ID()) }}', // upload url
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
    textEditor.setValue(document.getElementById("rawdata_source-{{ file.ID() }}").value);
    textEditor.getSession().on('change', function() {
        document.getElementById("rawdata_source-{{file.ID()}}").value = textEditor.getSession().getValue();
    });
    
    textEditor.$blockScrolling = Infinity;
    {% endif %}
</script>
{% endblock %}
`))
		file.Sync()
	}

	// file.new.form
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "file.new.form")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@file.new.form #}

{% extends "console/layout" %}

{% block breadcrumb %}

{% set bucket = LoadByID("buckets", ctx.Get("bucket_id")) %}
{% set file = LoadByID(bucket.Name(), ctx.Get("file_id")) %}

<ul class="uk-breadcrumb">
    <li><a href="{{ URL("ConsoleListBuckets") }}">buckets</a></li>
    <li><a href="{{ URL("ConsoleListBucketFiles", "bucket_id", bucket.ID()) }}"><i class="uk-icon-folder-open-o uk-text-muted"></i> {{ bucket.Name() }}</a></li>
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

    <form class="uk-form" action="{{ URL("ConsoleCreateFile", "bucket_id", bucket.ID()) }}" method="post" enctype="multipart/form-data">
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
{% endblock %}
`))
		file.Sync()
	}

	// bucket.new.form
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "bucket.new.form")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@bucket.new.form #}

{% extends "console/layout" %}

{% block breadcrumb %}

<ul class="uk-breadcrumb">
    <li><a href="{{ URL("ConsoleListBuckets") }}">buckets</a></li>
    <li><span>new bucket</span></li>
</ul>

{% endblock %}

{% block content %}

<div class="uk-width-1-2 uk-container uk-container-center">
    <form class="uk-form" action="{{ URL("ConsoleCreateBucket") }}" method="post" enctype="multipart/form-data">
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
{% endblock %}
`))
		file.Sync()
	}

	// file.put.textdata
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "file.put.textdata")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@file.put.textdata #}

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

{{ ctx.Redirect302(URL("ConsoleViewFile", "bucket_id", bucket.ID(), "file_id", file.ID()).String()) }}
{% endif %}
`))
		file.Sync()
	}

	// file.put.rawdata.via_uploader
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "file.put.rawdata.via_uploader")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@file.put.rawdata.via_uploader #}

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

{{ ctx.Redirect302(URL("ConsoleViewFile", "bucket_id", bucket.ID(), "file_id", file.ID()).String()) }}
{% endif %}
`))
		file.Sync()
	}

	// file.put.structdata
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "file.put.structdata")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@file.put.structdata #}

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

{{ ctx.Redirect302(URL("ConsoleViewFile", "bucket_id", bucket.ID(), "file_id", file.ID()).String()) }}
{% endif %}
`))
		file.Sync()
	}

	// file.put.properties
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "file.put.properties")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@file.put.properties #}

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

{{ ctx.Redirect302(URL("ConsoleViewFile", "bucket_id", bucket.ID(), "file_id", file.ID()).String()) }}
{% endif %}
`))
		file.Sync()
	}

	// file.new
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "file.new")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@file.new #}

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

{{ ctx.Redirect302(URL("ConsoleViewFile", "bucket_id", bucket.ID(), "file_id", file.ID()).String()) }}
{% endif %}
`))
		file.Sync()
	}

	// bucket.new
	file, err = store.LoadOrNewFile(config.ConsoleBucketName, "bucket.new")
	if err == dbox.ErrNotFound {
		logrus.Infof("%v: create %q\n", file.Bucket(), file.Name())

		file.SetContentType("text/html")

		file.RawData().Write([]byte(`{# console@bucket.new #}
{% if ctx.IsPost() %}

{% set opt = ctx.BindFormToMap("Name", "MetaDataStoreType", "SameAsMetaStoreType", "MapDataStoreType", "RawDataStoreType", "MappingMapDataFiles", "MetaHaveSuffix", "MapDataHaveSuffix", "RawDataHaveSuffix" ) %}

{% set res = CreateBucket(opt) %}
{# TODO: check error #}

{{ ctx.Redirect302(URL("ConsoleListBucketFiles", "bucket_id", res.ID()).String()) }}
{% comment %}
{% endcomment %}

{% endif %}
`))
		file.Sync()
	}
}
