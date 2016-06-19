package api

import (
	"github.com/flosch/pongo2"
	"sync"
)

var templateCacheMutex sync.Mutex
var templateCache = make(map[string]*pongo2.Template)

func ClearWidgetTemplatesCache() {

	templateCacheMutex.Lock()
	defer templateCacheMutex.Unlock()

	templateCache = make(map[string]*pongo2.Template)
}

// ExecuteFromCache wraper pongo2 template executor
func ExecuteFromCache(filename string) (*pongo2.Template, error) {
	if tpls.Debug {
		// Recompile on any request
		return tpls.FromFile(filename)
	}
	// Cache the template
	cleanedFilename := tplsLoader.Abs("", filename)

	templateCacheMutex.Lock()
	defer templateCacheMutex.Unlock()

	tpl, has := templateCache[cleanedFilename]

	// Cache miss
	if !has {
		tpl, err := tpls.FromFile(cleanedFilename)
		if err != nil {
			return nil, err
		}
		templateCache[cleanedFilename] = tpl
		return tpl, nil
	}

	// Cache hit
	return tpl, nil
}
