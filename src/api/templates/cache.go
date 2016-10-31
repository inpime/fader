package templates

import (
	"strings"
	"sync"

	"github.com/flosch/pongo2"
)

var templateCacheMutex sync.Mutex
var templateCache = make(map[string]*pongo2.Template)

func ClearMemCache() {

	templateCacheMutex.Lock()
	defer templateCacheMutex.Unlock()

	templateCache = make(map[string]*pongo2.Template)
}

// ExecuteFromMemCache wraper pongo2 template executor
func ExecuteFromMemCache(filename string) (*pongo2.Template, error) {
	if pongo2.DefaultSet.Debug {
		// Recompile on any request
		return pongo2.DefaultSet.FromFile(filename)
	}
	// Cache the template
	cleanedFilename := strings.TrimSpace(filename)

	templateCacheMutex.Lock()
	defer templateCacheMutex.Unlock()

	tpl, has := templateCache[cleanedFilename]

	// Cache miss
	if !has {
		tpl, err := pongo2.DefaultSet.FromFile(cleanedFilename)
		if err != nil {
			return nil, err
		}
		templateCache[cleanedFilename] = tpl
		return tpl, nil
	}

	// Cache hit
	return tpl, nil
}
