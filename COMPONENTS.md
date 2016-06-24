# Components

## Import export

Import\Export app data. 
Группировка файлов для экспорта только определенной части приложения

Location: `settings.main`
Section name: `importexport`

Example:

``` toml
# group files name 'groupname1'
[importexport.groupname1.bucketname1]
    files = ["main"] # only the listed files
[importexport.groupname1.bucketname2]
    all = true # all files bucket

# group files name 'groupname2'
[importexport.groupname2.bucketname1]
    files = ["main"] # only the listed files
[importexport.groupname2.bucketname2]
    all = true # all files bucket
```

### Special handlers

Names:
* `AppImport`
* `AppExport`
  * Если в запросе параметр `sys` имеет значение `true`, выгружаются только системные файлы

## File content

Returns the resource data. Eg. for css, js, images, fonts and etc.

Location: `settings.main`
Section name: `filecontent`

Parameters:
* `bucket` - bucket store files

Example:

``` toml
[filecontent]
bucket="static" # file store
```

### Special handlers

Names:
* `FileContentByName`, args `file`, `file name`
* `FileContentByID`, args `file`, `file ID`

Equivalent `{{ "logo.png"|fc }}` and `{{ URL("FileContentByName", "file", "logo.png") }}`