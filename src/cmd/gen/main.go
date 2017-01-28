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

	faderConsoleV1 := "fader.consolev1"
	bucketFile = interfaces.NewBucket()
	bucketFile.BucketID = faderConsoleBucketID
	bucketFile.BucketName = faderConsoleV1
	err = bucketManager.CreateBucket(bucketFile)
	log.Printf("create bucket %q", bucketFile.BucketName)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// main.toml
	////////////////////////////////////////////////////////////////////////////

	createFile(
		settingBucketID,
		api.MainConfigFileName, // filename
		``, // lua
		"text/toml",
		`
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
`,
	)

	////////////////////////////////////////////////////////////////////////////
	// CONSOLE v1
	////////////////////////////////////////////////////////////////////////////

	createFile(
		settingBucketID,
		"fader.console.v1.toml", // filename
		``, // lua
		"text/toml",
		`
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
`,
	)

	// routing
	createFile(
		settingBucketID,
		"fader.console.v1.routing.toml", // filename
		``, // lua
		"text/toml",
		`
[[routing.routs]]
name = "cv1_Dashboard"
path = "/console/v1/dashboard"
bucket = "`+faderConsoleV1+`"
file = "dashboard.html"
licenses = ["guest"]
methods = ["get"]

[[routing.routs]]
name = "cv1_ListBuckets"
path = "/console/v1/buckets"
bucket = "`+faderConsoleV1+`"
file = "list_buckets.html"
licenses = ["guest"]
methods = ["get"]

[[routing.routs]]
name = "cv1_ListFiles"
path = "/console/v1/buckets/{bucket_id}/files"
bucket = "`+faderConsoleV1+`"
file = "list_files.html"
licenses = ["guest"]
methods = ["get"]

[[routing.routs]]
name = "cv1_EditFile"
path = "/console/v1/files/{file_id}/edit"
bucket = "`+faderConsoleV1+`"
file = "edit_file.html"
licenses = ["guest"]
methods = ["get"]

[[routing.routs]]
name = "cv1_EditFile_UpdateByCmd"
path = "/console/v1/files/{file_id}/cmd_{cmd}"
bucket = "`+faderConsoleV1+`"
file = "edit_file.do.html"
licenses = ["guest"]
methods = ["post"]
`,
	)

	// createFile(
	// 	faderConsoleBucketID,
	// 	"", // filename
	// 	``, // lua
	// 	"text/toml",
	// 	``,
	// )

	// layout
	createFile(
		faderConsoleBucketID,
		"_layout.html", // filename
		``,             // lua
		"text/html",
		`
<html>
    <head>
        <title>{% block title %}{% endblock %}</title>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta charset="utf-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />

		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/uikit/2.27.2/css/uikit.almost-flat.min.css">
		<link rel="stylesheet" href="//cdn.jsdelivr.net/uikit/2.27.2/css/components/autocomplete.almost-flat.min.css" />
		
		<script id="ACETypes" type="application/json">{
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
	</head>
	<body>
		<div class="uk-container uk-container-center uk-margin-top uk-margin-bottom">

			<nav class="uk-navbar uk-margin-bottom">
				<a class="uk-navbar-brand uk-hidden-small" 
					href="{{ Route("cv1_Dashboard").URLPath() }}">FD-Cv1</a>
				<ul class="uk-navbar-nav uk-hidden-small">
					{% include "`+faderConsoleV1+`/_navbar.html" %}
				</ul>
				<a href="#offcanvas" class="uk-navbar-toggle uk-visible-small" data-uk-offcanvas></a>
				<div class="uk-navbar-brand uk-navbar-center uk-visible-small">FD-Cv1</div>
			</nav>
		</div>

		<div class="uk-container uk-container-center uk-margin-top uk-margin-bottom">
			{% block content %}{% endblock %}
		</div>


	<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.2.4/jquery.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/uikit/2.27.2/js/uikit.min.js"></script>
	<script src="//cdnjs.cloudflare.com/ajax/libs/clipboard.js/1.5.15/clipboard.min.js"></script>
	<script>
	$(document).ready(
        function(){
            var c = new Clipboard('[data-clipboard-text]');
            var c = new Clipboard('[data-clipboard-target]');
        }
    )
	</script>
	{% block scripts %}{% endblock %}

	<div id="offcanvas" class="uk-offcanvas">
        <div class="uk-offcanvas-bar">
            <ul class="uk-nav uk-nav-offcanvas">
                {% include "`+faderConsoleV1+`/_navbar.html" %}
            </ul>
        </div>
    </div>

	</body>
</html>
`,
	)

	// navbar

	createFile(
		faderConsoleBucketID,
		"_navbar.html", // filename
		``,             // lua
		"text/html",
		`
<li>
	<a href="{{ Route("cv1_ListBuckets").URLPath() }}">Buckets</a>
</li>
`,
	)

	// dashboard
	createFile(
		faderConsoleBucketID,
		"dashboard.html", // filename
		``,               // lua
		"text/html",
		`
{% extends "`+faderConsoleV1+`/_layout.html" %}

{% block content %}
<h1>Content</h1>
{% endblock %}
`,
	)

	// list buckets
	createFile(
		faderConsoleBucketID,
		"list_buckets.html", // filename
		``,                  // lua
		"text/html",
		`
{% extends "`+faderConsoleV1+`/_layout.html" %}

{% block content %}
<div class="uk-grid">
	<div class="uk-width-1-1">
		<ul class="uk-breadcrumb">
			<li><a href="{{ Route("cv1_Dashboard").URLPath() }}" class="uk-icon-home"></a></li>
			<li><a href="{{ Route("cv1_ListBuckets").URLPath() }}">Buckets</a></li>
		</ul>
	</div>
	<div class="uk-width-1-1">
		<h1 class="uk-heading">Buckets</h1>
		<ul class="uk-list uk-list-striped">
		{% for b in ListBuckets() %}
			<li><a href="{{ Route("cv1_ListFiles").URLPath("bucket_id", b.BucketID.String()) }}">{{ b.BucketName }}</a></li>
		{% endfor %}
		</ul>
	</div>
</div>

{% endblock %}
`,
	)

	// list files
	createFile(
		faderConsoleBucketID,
		"list_files.html", // filename
		`
local basic = require("basic")
ctx():Set("CurrentBucketName", "TODO: basket name")
`, // lua
		"text/html",
		`
{% extends "`+faderConsoleV1+`/_layout.html" %}

{% block content %}
<div class="uk-grid">
	<div class="uk-width-1-1">
		<ul class="uk-breadcrumb">
			<li><a href="{{ Route("cv1_Dashboard").URLPath() }}" class="uk-icon-home"></a></li>
			<li><a href="{{ Route("cv1_ListBuckets").URLPath() }}">Buckets</a></li>
			<li><a href="{{ Route("cv1_ListFiles").URLPath() }}">TODO: current bucket</a></li>
		</ul>
	</div>
	<div class="uk-width-1-1">
		<h1 class="uk-heading">Files by '{{ ctx.Get("CurrentBucketName") }}'</h1>
		<ul class="uk-list uk-list-striped">
		{% for f in ListFilesByBucketID(ctx.Get("bucket_id") ) %}
			<li><a href="{{ Route("cv1_EditFile").URLPath("file_id", f.FileID.String()) }}">{{ f.FileName }}</a></li>
		{% endfor %}
		</ul>
	</div>
</div>

{% endblock %}
`,
	)

	// edit file
	createFile(
		faderConsoleBucketID,
		"edit_file.html", // filename
		`
local std = require("basic")

file = std.FindFile(ctx():Get("file_id"))
bucketOfFile = std.FindBucket(file:BucketID())

ctx():Set("BucketName", bucketOfFile:BucketName())
ctx():Set("File", file)
`, // lua
		"text/html",
		`
{% extends "`+faderConsoleV1+`/_layout.html" %}

{% block content %}
{% set file = ctx.Get("File") %}
{% set bucketID = file.BucketID %}
{% set bucketName = ctx.Get("BucketName") %}

<div class="uk-grid">
	<div class="uk-width-1-1 uk-margin-bottom">
		<ul class="uk-breadcrumb">
			<li><a href="{{ Route("cv1_Dashboard").URLPath() }}" class="uk-icon-home"></a></li>
			<li><a href="{{ Route("cv1_ListBuckets").URLPath() }}">Buckets</a></li>
			<li><a href="{{ Route("cv1_ListFiles").URLPath("bucket_id", bucketID) }}">{{ bucketName }}</a></li>
			<li>
				<span>
					{{ file.FileName }}
					<div class="uk-button-group">
						<button class="uk-button uk-button-mini" data-clipboard-text="{{ file.FileName }}">
							<i class="uk-icon-clipboard"></i>
						</button>

						<div data-uk-dropdown="{mode:'click'}">
							<button class="uk-button uk-button-mini" style="vertical-align: inherit;"><i class="uk-icon-caret-down"></i></button>

							<!-- This is the dropdown -->
							<div class="uk-dropdown uk-dropdown-small">
								<ul class="uk-nav uk-nav-dropdown">
									<li>
										<a 
											class="uk-dropdown-close"
											data-clipboard-text='{{ bucketName }}/{{ file.FileName }}'>{{ bucketName }}/{{ file.FileName }}</a>
									</li>
									<li>
										<a 
											class="uk-dropdown-close"
											data-clipboard-text="&#123;&#37; extends &quot;{{ bucketName }}/{{ file.FileName }}&quot; &#37;&#125;">
											&#123;&#37; extend &quot; ... &quot; &#37;&#125;
										</a>
									</li>
									<li>
										<a 
											class="uk-dropdown-close"
											data-clipboard-text="&#123;&#37; include &quot;{{ bucketName }}/{{ file.FileName }}&quot; &#37;&#125;">
											&#123;&#37; include &quot; ... &quot; &#37;&#125;
										</a>
									</li>
								</ul>
							</div>

						</div>
					</div>
				</span>
			</li>
		</ul>
	</div>
	<div class="uk-width-1-1">
		<div class="uk-width-1-1">
			<form 
				action="{{ Route("cv1_EditFile_UpdateByCmd").URLPath("file_id", file.FileID.String(), "cmd", "update_name") }}" 
				class="uk-form uk-form-stacked" 
				method="POST">
				
				<fieldset data-uk-margin>
					<input type="text" 
							id="form-filename" 
							placeholder="File name"
							class="uk-width-1-3"
							name="FileName"
							value="{{ file.FileName }}">
					<button class="uk-button"><i class="uk-icon-save"></i></button>
				</fieldset>
			</form>
		</div>
		
		<div class="uk-width-medium-1-1 uk-row-first">

			<ul class="uk-tab" data-uk-tab="{connect:'#file-content'}">
				<li><a href="#">Raw data</a></li>
				<li><a href="#">Properties</a></li>
				<li><a href="#">Lua</a></li>
				<li><a href="#">Structural data</a></li>
			</ul>

			<ul id="file-content" class="uk-switcher uk-margin">
				<li>
					{# Raw data #}
					<div class="uk-width-1-1">
						<form class="uk-form uk-form-stacked">
							<div class="uk-form-row">
								<label class="uk-form-label">Content type</label>
								<div class="uk-form-controls">
									<div 
									class="uk-autocomplete uk-form uk-width-1-1" 
									id="file-{{ file.FileID.String() }}-choose-content-type" 
									data-uk-autocomplete>
										<input 
										type="text" 
										class="uk-width-1-1" 
										name="ContentType" 
										autocomplete="off" 
										value="{{ file.ContentType }}">
									</div>
								</div>
							</div>

							<div class="uk-form-row">
								<label class="uk-form-label">Content</label>
								<div class="uk-form-controls">
									<input 
										id="file-{{ file.FileID.String() }}-raw-data" 
										type="hidden" 
										name="TextData" 
										value="{{ file.RawData|btos }}" />
									<div 
										id="file-{{file.FileID.String()}}-raw-data-editor" 
										class="uk-width-1-1 uk-scrollable-box" 
										style="height:500px; border: 1px solid #ddd;"></div>
								</div>
							</div>

							<div class="uk-form-row">
								<div class="uk-form-controls uk-float-right">
									<button class="uk-button"><i class="uk-icon-save"></i> Save</button>
								</div>
							</div>
						</form>
					</div>
				</li>
				<li>
					{# Properties #}
					<div class="uk-width-1-1">
						<form class="uk-form uk-form-stacked">
							<div class="uk-form-row">
								<label class="uk-form-label" for="form-filename">Name</label>
								<div class="uk-form-controls">
									<input type="text" 
										id="form-filename" 
										placeholder="File name"
										class="uk-width-1-1">
								</div>
							</div>

							<div class="uk-form-row">
								<div class="uk-form-controls uk-float-right">
									<button class="uk-button"><i class="uk-icon-save"></i> Save</button>
								</div>
							</div>
						</form>
					</div>
				</li>
				<li>
					{# Lua #}
					<div class="uk-width-1-1">
						<form class="uk-form uk-form-stacked">
							<div class="uk-form-row">
								<label class="uk-form-label" for="form-contenttype">Content type</label>
								<div class="uk-form-controls">
									<input type="text" 
										id="form-contenttype" 
										placeholder="File name"
										class="uk-width-1-1"
										>
								</div>
							</div>

							<div class="uk-form-row">
								<div class="uk-form-controls uk-float-right">
									<button class="uk-button"><i class="uk-icon-save"></i> Save</button>
								</div>
							</div>
						</form>
					</div>
				</li>
				<li>
					{# Structural data #}
					<div class="uk-width-1-1">
						Structural data
					</div>
				</li>
			</ul>

		</div>
	</div>
</div>

{# end content block #}
{% endblock %}

{% block scripts %}
<script src="//cdn.jsdelivr.net/uikit/2.27.2/js/components/autocomplete.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/ace.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/mode-toml.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/mode-django.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/mode-html.js"></script>

<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/mode-lua.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/mode-luapage.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/snippets/lua.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/snippets/luapage.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.2.6/worker-lua.js"></script>

<script>
{# Autocomplete content type #}
var autocomplete = UIkit.autocomplete("#file-{{ file.FileID.String() }}-choose-content-type", {
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

{% if true %}
{# editor raw data #}
var rawTypes = JSON.parse(document.getElementById("ACETypes").textContent);

var rawDataEditor = ace.edit("file-{{file.FileID.String()}}-raw-data-editor");
rawDataEditor.setFontSize("11pt");
rawDataEditor.getSession().setMode(rawTypes["{{ file.ContentType|escapejs }}"]);
rawDataEditor.setValue(document.getElementById("file-{{ file.FileID.String() }}-raw-data").value);
rawDataEditor.getSession().on('change', function() {
	document.getElementById("file-{{ file.FileID.String() }}-raw-data").value = rawDataEditor.getSession().getValue();
});
rawDataEditor.$blockScrolling = Infinity;
{% endif %}

</script>

{# end scripts block #}
{% endblock %}
`,
	)

	// edit_file.do.html
	createFile(
		faderConsoleBucketID,
		"edit_file.do.html", // filename
		`
local std = require("basic")

fileID = ctx():Get("file_id")
newName = ctx():FormValue("FileName")

file = std.FindFile(fileID)

if file == nil then
	ctx():NoContext(404)
else
	file:SetFileName(newName)
	std.UpdateFileFrom(file, std.PrimaryNamesData)
	
	goTo = ctx():Route("cv1_EditFile"):URLPath("file_id", file:FileID())
	print("redirect to", goTo)
	
	ctx():Redirect(goTo)
end 
`, // lua
		"text/lua",
		``,
	)

}

func createFile(
	bucketID uuid.UUID,
	fileName,
	luaScript,
	contentType,
	rawData string,
) {
	file := interfaces.NewFile()
	file.FileID = uuid.NewV4()
	file.BucketID = bucketID
	file.FileName = fileName
	file.LuaScript = []byte(luaScript)
	file.ContentType = contentType
	file.RawData = []byte(rawData)

	err := fileManager.CreateFile(file)
	log.Printf("create file %q", file.FileName)
	if err != nil {
		panic(err)
	}
}
