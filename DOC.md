# Philosophy

# Context Functions

* `ctx.FormFileData(filename) *FileData` - file data from the form

# Page Variables

* `ctx` - the context of the current request, [see more details](#link to description of context)

# Page Functions

* `Load(bucketname, filename)` - load the file by name (bucket and filename)
* `LoadByID(bucketname, fileid)` - load the file by id (bucket and file id)
* `SearchFiles(bucketname, querystring, page, perpage)` - search files in bucket
* `M()` - helpful function, [see more details](#link to description of helpful functions)
* `A()` - helpful function, [see more details](#link to description of helpful functions)

# Buckets

* settings - contains files for application settings and components
* users - users of the application
* pages - page of the web site
* console - web-ui management application

## Settings

* main - general settings
* routing - router 

### Routs

Schema (toml)

```
[[routs]]
path = string
handler = string
special = boolean
methods = array string
licenses = array string
```

#### Special Handlers

* `usercontent` - returns the file contents

