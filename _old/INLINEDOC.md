
# Template enviroment

## Functions

| Name | Description |
| ---- | ----------- |
|   Context    |  |
|   `ctx.IsPut() bool`    | The current request is PUT method  |
|   `ctx.IsPost() bool`   | The current request is POST method  |
|   `ctx.IsDelete() bool` | The current request is DELETE method  |
|   `ctx.IsGet() bool`    | The current request is GET method  |
|   `ctx.BindFormToMap(filenames ...string) Map`    |   Returns the form field values for the provided names. Return object of the type [Map](#Map)|
|   `ctx.FormFileData() FileData`    |  Returns the form file for the provided names. Return object of the type [FileData](#FileData)    |
|   `ctx.Redirect302(url string) *Context`    |   |
|   `ctx.Session() Session`    | Return current session. Return object of the type [Session](#Session) |
|   `ctx.CurrentRoute() Route`    | Return current route.   Return object of the type [Route](#Route)  |
|   `ctx.CurrentRouteIs(name string) bool`    | Сheck the name of the current route name.  |
|   `ctx.CSRFToken() string`    |   CSRF token value.   |
|   `ctx.CSRFFieldName() string`    |   CSRF field name.    |
|   `ctx.CSRFField() string`    |   CSRF form field (`<input type='hidden' value='value' name='name' />`) |
|   Router    |   |
|   `URL(name string, args... string) URL`    | Build URL by name route. Return object of the type [URL](#URL). |
|   Basic    |   |
|   `NewUUID() string`    | Generate UUID v4. |
|   `Config(name string) Map`    | Get config by name secion. Return object of the type [Map](#Map) |
|   `DeleteFile(bucketID, fileID string) error`    | Delete file by ID. |
|   `NewFile(bucketName string) File`    | New file from bucket name. Return object of the type [File](#File) |
|   `LoadByID(bucketName, fileID string) File|error`    | Load file by ID. Return object of the type [File](#File) |
|   `Load(bucketName, fileName string) File|error`    | Load file by name. Return object of the type [File](#File) |
|   `URLQuery(url URL, args ...string) URL`    | |
|   `M() Map`    |     |
|   `A() ArrayString`    |     |
|   `AIface() ArrayIface`    |     |
|   `Validator() Validator`    |     |
|   `SendEmail(to, subject, template string, context Map) string|error`    |     |
|   Search    |   |
|   `SearchFiles(bucketName, query string, page, perpage int) SearchResult`    | Search for files. Return object of the type [SearchResult](#SearchResult) |

### Filters

| Name | Description |
| ---- | ----------- |
|   Basic    |  |
|   `is_error`    |  |
|   `clear`    |  |
|   `logf`    |  |
|   `atojs`    |  |
|   `split`    |  |
|   `btos`    |  |
|   `stob`    |  |
|   `append`    |  |
|   File static    |  |
|   `fc`    | Get url by file name. |
|   `filecontenturl`    | It is an alias for `fc`. |
|   `urlfile`    | It is an alias for `fc`. |

### Tags

| Name | Description |
| ---- | ----------- |
|   Basic    |  |

## Object types

### File

| Name | Description |
| ---- | ----------- |
|   `ID() string`    | File ID |
|   `Bucket() string`    | Bucket name |
|   `Name() string`    | File name |
|   `SetName(string) File`    | Set file name |
|   `Sync() error`  | Save file |
|   `Delete() error`  | Felete file |
|   `MMeta() Map`  | Meta data of file |
|   `MMapData() Map`  | Structured data of file |
|   `TextData() string`  | Return string data of file |
|   `SetTextData(string) File`  |  |
|   `SetRawData([]byte) File`  |  |
|   `ContentType() string`  | Return content type of file |
|   `SetContentType(string) File`  |  |
|   `IsImage() bool`  |  |
|   `IsText() bool`  |  |
|   `IsRaw() bool`  |  |

### FieldData

### Session

### Map

### ArrayString

### ArrayIface

### Validator

### Route

### URL

It is an alias for `http.URL`. [See more](https://golang.org/pkg/net/url/#URL)

### SearchResult

Interface

``` go
type SearchResultIface interface {
	GetFiles() []*File

	GetTotal() int64
	GetHasNext() bool
	GetNextPage() int
	IsError() bool
	Error() string
}
```

## Addons

### Mailer

### Validator

### ImportExport

### FileStatic