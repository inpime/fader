package interfaces

import "net/url"

type RouteMatch struct {
	// Route matched route
	Route   Route
	Handler RequestHandler
	Vars    map[string]string
}

type RequestParams struct {
	URL    *url.URL
	Method string
}

type RouteMatcher interface {
	Match(RequestParams, *RouteMatch) bool
}

type Route interface {
	// RouteMatcher

	Options() RequestHandler

	GetName() string
	Name(string) Route
	Path(string) Route
	Methods(...string) Route
	Handler(RequestHandler) Route

	URLPath(pairs ...string) (*url.URL, error)
}

type Router interface {
	// RouteMatcher

	Get(string) Route
	Handle(path string, h RequestHandler) Route
}

type RequestHandler struct {
	Name string `toml:"name"`
	Path string `toml:"path"`

	Bucket string `toml:"bucket"`
	File   string `toml:"file"`

	LuaScript     string `toml:"lua"`
	LuaArgsScript string `toml:"lua_args"`

	AllowedLicenses []string `toml:"licenses"`
	AllowedMethods  []string `toml:"methods"`

	CSRF bool `toml:"csrf"`
}

func (h RequestHandler) IsEmpty() bool {
	return len(h.Bucket) == 0 && len(h.File) == 0
}
