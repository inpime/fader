package standard

import (
	"encoding/json"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
)

func (Extension) initTplFilters() {

	// pongo2.RegisterFilter("fc", filterFileContentByNameURL)
	// pongo2.RegisterFilter("filecontenturl", filterFileContentByNameURL) // alias fc
	// pongo2.RegisterFilter("urlfile", filterFileContentByNameURL) // OLD

	pongo2.RegisterFilter("is_error", filterIsError)
	pongo2.RegisterFilter("clear", filterClear)
	pongo2.RegisterFilter("muted", filterClear) // alias clear
	pongo2.RegisterFilter("logf", filterLogf)
	pongo2.RegisterFilter("atojs", tagAnyObjectToJS)
// 	pongo2.RegisterFilter("split", filterSplit)
	pongo2.RegisterFilter("btos", filterBytesToString)
	pongo2.RegisterFilter("stob", filterStringToBytes)
	pongo2.RegisterFilter("append", filterAppend)
}

// filterAppend append
func filterAppend(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

	return pongo2.AsValue(append(in.Interface().([]interface{}), param.Interface())), nil
}

// filterSplit split
func filterSplit(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	str := in.String()
	sep := param.String()
	if len(sep) == 0 {
		return pongo2.AsValue(strings.Fields(str)), nil
	}

	return pongo2.AsValue(strings.Split(str, sep)), nil
}

// tagAnyObjectToJS atojs
func tagAnyObjectToJS(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	jsonByte, err := json.Marshal(in.Interface())

	if err != nil {
		logrus.WithField("_api", addonName).WithError(err).Warningf("error marshaling %T to json", in.Interface())
	}

	return pongo2.AsSafeValue(string(jsonByte)), nil
}

// filterClear clear
func filterClear(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

	return pongo2.AsValue(nil), nil
}

// filterIsError is_error
func filterIsError(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	_, ok := in.Interface().(error)

	return pongo2.AsValue(ok), nil
}

// filterLogf logf
func filterLogf(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	format := param.String()
	logrus.WithFields(logrus.Fields{
		"_api": addonName + ".custom",
	}).WithField("target", "web").Infof(format, in.Interface())

	return pongo2.AsValue(nil), nil
}

// filterStringToBytes stob
func filterStringToBytes(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	v, ok := in.Interface().(string)
	if !ok {
		return pongo2.AsValue([]byte{}), nil
	}
	return pongo2.AsValue([]byte(v)), nil
}

// filterBytesToString btos
func filterBytesToString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	v, ok := in.Interface().([]byte)
	if !ok {
		return pongo2.AsValue(nil), nil
	}
	return pongo2.AsValue(string(v)), nil
}
