package config

import (
	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	// "io"
)

type FnRegTagTpl func(string, pongo2.TagParser) error
type FnRegFilterTpl func(string, pongo2.FilterFunction) error
type FnRegHandlerTpl func(string, echo.HandlerFunc)
type FnRegMiddlewareTpl func(string, echo.MiddlewareFunc) error
type Extensions map[string]FaderExtension

var extensions = make(Extensions)

type FaderExtension interface {
	Name() string
	Version() string

	// SetLogger(io.Writer)
	Middlewares() []echo.MiddlewareFunc
	RegEchoHandlers(func(string, func(echo.Context) error))

	InjectTplAddons() error

	Destroy()

	// BuilderSearchQuery(string, interface{}) error
}

func AddExtension(ext FaderExtension) {
	extensions[ext.Name()] = ext
}

// ListOfExtensions list of Fader extensions
func ListOfExtensions() Extensions {
	return extensions
}

// empty extenstion

type EmptyFaderExtension struct {
}

func (*EmptyFaderExtension) Name() string {
	return "example"
}

func (*EmptyFaderExtension) RegEchoHandlers(fnReg func(string, func(echo.Context) error)) {
	fnReg("example_spcial_handler", func(ctx echo.Context) error {
		return nil
	})
}

func (*EmptyFaderExtension) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (*EmptyFaderExtension) InjectTplAddons() error {
	// pongo2.Context

	// pongo2.Tags

	// pongo2.Filters

	return nil
}
