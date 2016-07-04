package filestatic

import (
	"api/vrouter"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
)

// filterFileContentByNameURL
func filterUrlFileByName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// TODO: get the URL based on the name route (after the routs will have the names)
	route := vrouter.AppRouter().Get(ByNameRouteName)

	if route == nil {
		logrus.WithError(fmt.Errorf("not found route")).WithFields(logrus.Fields{
			"_service":  addonName,
			"routename": ByNameRouteName,
		}).Warning("not found route")

		return pongo2.AsValue(""), nil
	}

	_url, err := route.URLPath("file", in.String())

	if err != nil {
		logrus.WithError(fmt.Errorf("error build url by file")).WithFields(logrus.Fields{
			"_service":      addonName,
			"args_filename": in.String(),
		}).Warning("error build url")

		return pongo2.AsValue(""), nil
	}

	return pongo2.AsValue(_url.String()), nil
}
