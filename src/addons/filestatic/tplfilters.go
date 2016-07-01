package filestatic

import (
	"api/config"
	"fmt"
	"github.com/flosch/pongo2"
)

// filterFileContentByNameURL
func filterUrlFileByName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// TODO: get the URL based on the name route (after the routs will have the names)
	route := config.Router.Get(ByNameHandlerName)

	if route == nil {
		reason := fmt.Sprintf("not found route %q", ByNameHandlerName)
		return nil, &pongo2.Error{ErrorMsg: reason}
	}

	_url, err := route.URLPath("file", in.String())

	if err != nil {
		reason := fmt.Sprintf("error build url by %q", in.String())
		return nil, &pongo2.Error{ErrorMsg: reason}
	}

	return pongo2.AsValue(_url.String()), nil
}
