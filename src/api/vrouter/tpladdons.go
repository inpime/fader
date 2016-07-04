package vrouter

import (
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"net/url"
)

func tplContext() {

	// builds the path part of the URL
	pongo2.DefaultSet.Globals["URL"] = func(args ...*pongo2.Value) *pongo2.Value {
		emptyUrl, _ := url.Parse("")

		if len(args) == 0 {
			return pongo2.AsValue(emptyUrl)
		}

		routeName := args[0].String()
		route := AppRouter().Get(routeName)

		if route == nil {
			return pongo2.AsValue(emptyUrl)
		}

		if (len(args)-1)%2 != 0 {
			logrus.WithFields(logrus.Fields{
				"_service": addonName,
			}).Warningf("args expected in multiples of two, want %d", len(args)-1)
			return pongo2.AsValue(emptyUrl)
		}

		stringArgs := []string{}

		for _, arg := range args[1:] {
			stringArgs = append(stringArgs, arg.String())
		}

		_url, err := route.URLPath(stringArgs...)

		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"_service": addonName,
				"args":     stringArgs,
			}).Warning("build url")

			return pongo2.AsValue(emptyUrl)
		}

		return pongo2.AsValue(_url)
	}
}
