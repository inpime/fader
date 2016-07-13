# Official and system addons

Required and system  addons of the Fader
* session
* search
* filestatic
* importexport
* standard

## Session

Name `fader.addons.session`

## Search

Name `fader.addons.search`

### Template context functions

| Function name  | Signature | Description |
|---|---|---|
|`SearchFiles`|  |   |

## Filtestatic

Name `fader.addons.filestatic`

### Handlers

| Handler name  | Signature | Description |
|---|---|---|
| `fader.addons.filestatic.byname` |  |    |
| `fader.addons.filestatic.byid` |  |    | 

### Template filters

| Filter name  | Signature | Description |
|---|---|---|
|`fc`|   |   |
|`filecontenturl`|   |   |
|`urlfile`|   |   |

## Importexport

Name `fader.addons.importexport`

### Handlers

| Handler name  | Signature | Description |
|---|---|---|
| `fader.addons.importexport.import` |  |    |
| `fader.addons.importexport.export` |  |    | 

### Template context functions

| Function name  | Signature | Description |
|---|---|---|
|`ListGroupsImportExport`|  |   |

## Standard

Name `fader.addons.standard`

Вспомогательне инструменты. TODO: Дозаполнить из pongo2

### Template context functions

| Function name  | Signature | Description |
|---|---|---|
|`NewUUID`|   |   |
|`SectionAppConfig`|   |   |
|`DeleteFile`|   |   |
|`NewFile`|   |   |
|`LoadByID`|   |   |
|`Load`|   |   |
|`URLQuery`|   |   |
|`URL`|   |   |
|`M`|   |   |
|`A`|   |   |
|`AIface`|   |   |
|`CreateBucket`|   |   |


### Template filters

| Filter name  | Signature | Description |
|---|---|---|
|`is_error`|   |   |
|`clear`|   |   |
|`logf`|   |   |
|`atojs`|   |   |
|`split`|   |   |
|`btos`|   |   |
|`stob`|   |   |
|`append`|   |   |

### Template tags

| Tag name  | Signature | Description |
|---|---|---|