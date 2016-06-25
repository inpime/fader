package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
)

var pongo2InitAddonsOnce sync.Once

func pongo2InitAddons() {
	pongo2InitAddonsOnce.Do(func() {
		// tags

		pongo2.ReplaceTag("ssi", tagSSI)

		// TODO: static file with check access (ACL)
		// Загружать props файла и проверять от сессии доступ к этому файлу

		// filters

		// file content
		pongo2.RegisterFilter("fc", filterFileContentByNameURL)
		pongo2.RegisterFilter("filecontenturl", filterFileContentByNameURL) // alias fc
		pongo2.RegisterFilter("urlfile", filterFileContentByNameURL)        // OLD

		pongo2.RegisterFilter("is_error", filterIsError)
		pongo2.RegisterFilter("clear", filterClear)
		pongo2.RegisterFilter("logf", filterLogf)
		pongo2.RegisterFilter("atojs", tagAnyObjectToJS)
		pongo2.RegisterFilter("split", filterSplit)
		pongo2.RegisterFilter("btos", filterBytesToString)
		pongo2.RegisterFilter("stob", filterStringToBytes)
		pongo2.RegisterFilter("append", filterAppend)
	})
}

func filterAppend(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

	return pongo2.AsValue(append(in.Interface().([]interface{}), param.Interface())), nil
}

// ------
// filter static file
// ------

func filterFileContentByNameURL(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// TODO: get the URL based on the name route (after the routs will have the names)
	route := router.Get(FileContentByNameRouteName)

	if route == nil {
		reason := fmt.Sprintf("not found route %q", FileContentByNameRouteName)
		return nil, &pongo2.Error{ErrorMsg: reason}
	}

	_url, err := route.URLPath("file", in.String())

	if err != nil {
		reason := fmt.Sprintf("error build url by %q", in.String())
		return nil, &pongo2.Error{ErrorMsg: reason}
	}

	return pongo2.AsValue(_url.String()), nil
}

// ------
// filter string to bytes
// ------

func filterLogf(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	format := param.String()
	logrus.WithField("target", "web").Infof(format, in.Interface())

	return pongo2.AsValue(nil), nil
}

// ------
// filter string to bytes
// ------

func filterStringToBytes(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	v, ok := in.Interface().(string)
	if !ok {
		return pongo2.AsValue([]byte{}), nil
	}
	return pongo2.AsValue([]byte(v)), nil
}

// ------
// filter bytes to string
// ------

func filterBytesToString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	v, ok := in.Interface().([]byte)
	if !ok {
		return pongo2.AsValue(nil), nil
	}
	return pongo2.AsValue(string(v)), nil
}

// ------
// filter split
// ------

func filterSplit(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	str := in.String()
	sep := param.String()
	if len(sep) == 0 {
		return pongo2.AsValue(strings.Fields(str)), nil
	}

	return pongo2.AsValue(strings.Split(str, sep)), nil
}

// ------
// filter any object to json\js
// ------

func tagAnyObjectToJS(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	jsonByte, err := json.Marshal(in.Interface())

	if err != nil {
		logrus.WithError(err).Warningf("error marshaling %T to json", in.Interface())
	}

	return pongo2.AsSafeValue(string(jsonByte)), nil
}

// ------
// filter is_error
// ------

func filterIsError(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	_, ok := in.Interface().(error)

	return pongo2.AsValue(ok), nil
}

// ------
// filter clear
// ------

func filterClear(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

	return pongo2.AsValue(nil), nil
}

// // ------
// // tag static
// // ------

// type tagStaticNode struct {
// 	filename string
// }

// func (node *tagStaticNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
// 	writer.WriteString("/statics/" + node.filename)

// 	return nil
// }

// func tagStatic(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {

// 	if fileToken := arguments.MatchType(pongo2.TokenString); fileToken != nil {
// 		return &tagStaticNode{fileToken.Val}, nil
// 	} else if fileToken := arguments.MatchType(pongo2.TokenIdentifier); fileToken != nil {

// 		logrus.Debugf("static: %#v", fileToken)
// 		return &tagStaticNode{fileToken.Val}, nil
// 	} else {

// 		return nil, arguments.Error("First argument must be a string.", nil)
// 	}

// 	if arguments.Remaining() > 0 {
// 		return nil, arguments.Error("Malformed SSI-tag argument.", nil)
// 	}

// 	return nil, nil
// }

// ---------
// tag ssi
// ---------

type tagSSINode struct {
	filename string
	content  string
	template *pongo2.Template
}

func (node *tagSSINode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	if node.template != nil {
		// Execute the template within the current context
		includeCtx := make(pongo2.Context)
		includeCtx.Update(ctx.Public)
		includeCtx.Update(ctx.Private)
		content, err := node.template.Execute(includeCtx)
		if err != nil {
			return err.(*pongo2.Error)
		}
		writer.WriteString(content)
	} else {
		// Just print out the content
		writer.WriteString(node.content)
	}
	return nil
}

func tagSSI(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	SSINode := &tagSSINode{}

	if fileToken := arguments.MatchType(pongo2.TokenString); fileToken != nil {
		SSINode.filename = fileToken.Val

		if arguments.Match(pongo2.TokenIdentifier, "parsed") != nil {
			// parsed
			temporaryTpl, err := tpls.FromFile(fileToken.Val)
			if err != nil {
				return nil, err.(*pongo2.Error)
			}
			SSINode.template = temporaryTpl
		} else {
			// plaintext

			fileReader, err := tplsLoader.Get(fileToken.Val)

			if err != nil {

				logrus.WithError(err).Warningf("pongo2: tag `ssi`, load file by name %q", fileToken.Val)
				SSINode.content = ""
				return SSINode, nil
				// return nil, (&pongo2.Error{
				// 	Sender:   "tag:ssi",
				// 	ErrorMsg: err.Error(),
				// })
			}

			buf, err := ioutil.ReadAll(fileReader)

			if err != nil {
				return nil, (&pongo2.Error{
					Sender:   "tag:ssi",
					ErrorMsg: err.Error(),
				})
			}

			SSINode.content = string(buf)
		}
	} else {
		return nil, arguments.Error("First argument must be a string.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed SSI-tag argument.", nil)
	}

	return SSINode, nil
}
