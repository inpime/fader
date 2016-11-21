package router

import (
	"errors"
	"fmt"
	"interfaces"
	"net/url"
	"strings"
	"sync"
)

type Route struct {
	sync.Mutex

	// Parent where the route was registered (a Router).
	parent parentRoute

	handler interfaces.RequestHandler

	// // List of matchers.
	matchers []interfaces.RouteMatcher

	// // Manager for the variables from host and path.
	regexp *routeRegexpGroup

	// The name used to build URLs.
	name string

	// If true, this route never matches: it is only used to build URLs.
	buildOnly bool

	// If true, when the path pattern is "/path/", accessing "/path" will
	// redirect to the former and vice versa.
	strictSlash bool

	// If true, when the path pattern is "/path//to", accessing "/path//to"
	// will not redirect
	skipClean bool

	err error

	buildVarsFunc BuildVarsFunc
}

func (r *Route) Options() interfaces.RequestHandler {
	return r.handler
}

// Name -----------------------------------------------------------------------

// Name sets the name for the route, used to build URLs.
// If the name was registered already it will be overwritten.
func (r *Route) Name(name string) interfaces.Route {
	if r.name != "" {
		r.err = fmt.Errorf("mux: route already has name %q, can't set %q",
			r.name, name)
	}
	if r.err == nil {
		r.name = name
		r.getNamedRoutes()[name] = r
	}
	return r
}

// GetName returns the name for the route, if any.
func (r *Route) GetName() string {
	return r.name
}

// BuildVarsFunc --------------------------------------------------------------

// BuildVarsFunc is the function signature used by custom build variable
// functions (which can modify route variables before a route's URL is built).
type BuildVarsFunc func(map[string]string) map[string]string

// BuildVarsFunc adds a custom function to be used to modify build variables
// before a route's URL is built.
func (r *Route) BuildVarsFunc(f BuildVarsFunc) *Route {
	r.buildVarsFunc = f
	return r
}

func (r *Route) Match(req interfaces.RequestParams, match *interfaces.RouteMatch) bool {
	if r.buildOnly || r.err != nil {
		return false
	}

	// Match everything.
	for _, m := range r.matchers {
		if matched := m.Match(req, match); !matched {
			return false
		}
	}

	if match.Route == nil {
		match.Route = r
	}
	if match.Handler.IsEmpty() {
		match.Handler = r.handler
	}
	if match.Vars == nil {
		match.Vars = make(map[string]string)
	}

	// Set variables.
	if r.regexp != nil {
		r.regexp.setMatch(req, match, r)
	}
	return true
}

func (r *Route) Path(tpl string) interfaces.Route {
	r.addRegexpMatcher(tpl, false, false, false)

	return r
}

// methodMatcher matches the request against HTTP methods.
type methodMatcher []string

func (m methodMatcher) Match(r interfaces.RequestParams, match *interfaces.RouteMatch) bool {
	return matchInArray(m, r.Method)
}

// Methods adds a matcher for HTTP methods.
// It accepts a sequence of one or more methods to be matched, e.g.:
// "GET", "POST", "PUT".
func (r *Route) Methods(methods ...string) interfaces.Route {
	for k, v := range methods {
		methods[k] = strings.ToUpper(v)
	}

	return r.addMatcher(methodMatcher(methods))
}

func (r *Route) Handler(v interfaces.RequestHandler) interfaces.Route {
	r.handler = v
	return r
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m interfaces.RouteMatcher) interfaces.Route {
	if r.err == nil {
		r.matchers = append(r.matchers, m)
	}
	return r
}

// ----------------------------------------------------------------------------
// URL building
// ----------------------------------------------------------------------------

// URL builds a URL for the route.
//
// It accepts a sequence of key/value pairs for the route variables. For
// example, given this route:
//
//     r := mux.NewRouter()
//     r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler).
//       Name("article")
//
// ...a URL for it can be built using:
//
//     url, err := r.Get("article").URL("category", "technology", "id", "42")
//
// ...which will return an url.URL with the following path:
//
//     "/articles/technology/42"
//
// This also works for host variables:
//
//     r := mux.NewRouter()
//     r.Host("{subdomain}.domain.com").
//       HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler).
//       Name("article")
//
//     // url.String() will be "http://news.domain.com/articles/technology/42"
//     url, err := r.Get("article").URL("subdomain", "news",
//                                      "category", "technology",
//                                      "id", "42")
//
// All variables defined in the route are required, and their values must
// conform to the corresponding patterns.
func (r *Route) URL(pairs ...string) (*url.URL, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.regexp == nil {
		return nil, errors.New("mux: route doesn't have a host or path")
	}
	values, err := r.prepareVars(pairs...)
	if err != nil {
		return nil, err
	}
	var scheme, host, path string
	if r.regexp.host != nil {
		// Set a default scheme.
		scheme = "http"
		if host, err = r.regexp.host.url(values); err != nil {
			return nil, err
		}
	}
	if r.regexp.path != nil {
		if path, err = r.regexp.path.url(values); err != nil {
			return nil, err
		}
	}
	return &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}, nil
}

// URLHost builds the host part of the URL for a route. See Route.URL().
//
// The route must have a host defined.
func (r *Route) URLHost(pairs ...string) (*url.URL, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.regexp == nil || r.regexp.host == nil {
		return nil, errors.New("mux: route doesn't have a host")
	}
	values, err := r.prepareVars(pairs...)
	if err != nil {
		return nil, err
	}
	host, err := r.regexp.host.url(values)
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Scheme: "http",
		Host:   host,
	}, nil
}

// URLPath builds the path part of the URL for a route. See Route.URL().
//
// The route must have a path defined.
func (r *Route) URLPath(pairs ...string) (*url.URL, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.regexp == nil || r.regexp.path == nil {
		return nil, errors.New("mux: route doesn't have a path")
	}
	values, err := r.prepareVars(pairs...)
	if err != nil {
		return nil, err
	}
	path, err := r.regexp.path.url(values)
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Path: path,
	}, nil
}

// prepareVars converts the route variable pairs into a map. If the route has a
// BuildVarsFunc, it is invoked.
func (r *Route) prepareVars(pairs ...string) (map[string]string, error) {
	m, err := mapFromPairsToString(pairs...)
	if err != nil {
		return nil, err
	}
	return r.buildVars(m), nil
}

func (r *Route) buildVars(m map[string]string) map[string]string {
	if r.parent != nil {
		m = r.parent.buildVars(m)
	}
	if r.buildVarsFunc != nil {
		m = r.buildVarsFunc(m)
	}
	return m
}

// addRegexpMatcher adds a host or path matcher and builder to a route.
func (r *Route) addRegexpMatcher(tpl string, matchHost, matchPrefix, matchQuery bool) error {
	if r.err != nil {
		return r.err
	}
	r.regexp = new(routeRegexpGroup)
	// r.regexp = r.getRegexpGroup()
	if !matchHost && !matchQuery {
		if len(tpl) == 0 || tpl[0] != '/' {
			return fmt.Errorf("mux: path must start with a slash, got %q", tpl)
		}
		if r.regexp.path != nil {
			tpl = strings.TrimRight(r.regexp.path.template, "/") + tpl
		}
	}
	rr, err := newRouteRegexp(tpl, matchHost, matchPrefix, matchQuery, r.strictSlash)
	if err != nil {
		return err
	}
	for _, q := range r.regexp.queries {
		if err = uniqueVars(rr.varsN, q.varsN); err != nil {
			return err
		}
	}
	if matchHost {
		if r.regexp.path != nil {
			if err = uniqueVars(rr.varsN, r.regexp.path.varsN); err != nil {
				return err
			}
		}
		r.regexp.host = rr
	} else {
		if r.regexp.host != nil {
			if err = uniqueVars(rr.varsN, r.regexp.host.varsN); err != nil {
				return err
			}
		}
		if matchQuery {
			r.regexp.queries = append(r.regexp.queries, rr)
		} else {
			r.regexp.path = rr
		}
	}
	r.addMatcher(rr)
	return nil
}

// ------------------
// Router
// ------------------

type Router struct {
	sync.Mutex

	NotFoundHandler interfaces.RequestHandler

	parent parentRoute

	routes []*Route

	// If true, when the path pattern is "/path/", accessing "/path" will
	// redirect to the former and vice versa.
	strictSlash bool

	// If true, when the path pattern is "/path//to", accessing "/path//to"
	// will not redirect
	skipClean bool

	namedRoutes map[string]*Route
}

// Get returns a route registered with the given name.
func (r *Router) Get(name string) interfaces.Route {
	r.Lock()
	defer r.Unlock()

	f, exists := r.getNamedRoutes()[name]

	if !exists {
		return nil
	}

	return f
}

// GetRoute returns a route registered with the given name. This method
// was renamed to Get() and remains here for backwards compatibility.
func (r *Router) GetRoute(name string) *Route {
	r.Lock()
	defer r.Unlock()

	f, exists := r.getNamedRoutes()[name]

	if !exists {
		return nil
	}

	return f
}

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}

func (r *Router) clear() {
	r.routes = []*Route{}
	r.namedRoutes = make(map[string]*Route)
}

func (r *Router) RefreshRoutes(routs []interfaces.RequestHandler) {
	r.Lock()
	defer r.Unlock()

	r.clear()

	for _, h := range routs {
		r.Handle(h.Path, h).
			Methods(h.AllowedMethods...).
			Name(h.Name)
	}

}

func (r *Router) NewRoute() *Route {

	route := &Route{parent: r, strictSlash: r.strictSlash, skipClean: r.skipClean}
	r.routes = append(r.routes, route)

	return route
}

func (r *Router) Handle(path string, h interfaces.RequestHandler) interfaces.Route {

	return r.NewRoute().Path(path).Handler(h)
}

func (r Router) Match(req interfaces.RequestParams, match *interfaces.RouteMatch) bool {
	for _, route := range r.routes {
		if route.Match(req, match) {
			return true
		}
	}

	if !r.NotFoundHandler.IsEmpty() {
		match.Handler = r.NotFoundHandler
		return true
	}

	return false
}

// ----------------------------------------------------------------------------
// parentRoute
// ----------------------------------------------------------------------------

// parentRoute allows routes to know about parent host and path definitions.
type parentRoute interface {
	getNamedRoutes() map[string]*Route
	getRegexpGroup() *routeRegexpGroup
	buildVars(map[string]string) map[string]string
}

// ----------------------------------------------------------------------------
// parentRoute
// ----------------------------------------------------------------------------

// getNamedRoutes returns the map where named routes are registered.
func (r *Router) getNamedRoutes() map[string]*Route {
	if r.namedRoutes == nil {
		if r.parent != nil {
			r.namedRoutes = r.parent.getNamedRoutes()
		} else {
			r.namedRoutes = make(map[string]*Route)
		}
	}
	return r.namedRoutes
}

// getRegexpGroup returns regexp definitions from the parent route, if any.
func (r *Router) getRegexpGroup() *routeRegexpGroup {
	if r.parent != nil {
		return r.parent.getRegexpGroup()
	}
	return nil
}

func (r *Router) buildVars(m map[string]string) map[string]string {
	if r.parent != nil {
		m = r.parent.buildVars(m)
	}
	return m
}

// getNamedRoutes returns the map where named routes are registered.
func (r *Route) getNamedRoutes() map[string]*Route {
	r.Lock()
	defer r.Unlock()

	if r.parent == nil {
		// During tests router is not always set.
		r.parent = NewRouter()
	}
	return r.parent.getNamedRoutes()
}

// getRegexpGroup returns regexp definitions from this route.
func (r *Route) getRegexpGroup() *routeRegexpGroup {
	if r.regexp == nil {
		if r.parent == nil {
			// During tests router is not always set.
			r.parent = NewRouter()
		}
		regexp := r.parent.getRegexpGroup()
		if regexp == nil {
			r.regexp = new(routeRegexpGroup)
		} else {
			// Copy.
			r.regexp = &routeRegexpGroup{
				host:    regexp.host,
				path:    regexp.path,
				queries: regexp.queries,
			}
		}
	}
	return r.regexp
}
