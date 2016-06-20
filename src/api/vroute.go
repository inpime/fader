package api

import "strings"
import "net/url"
import log "github.com/Sirupsen/logrus"
import "fmt"

type matcher interface {
	Match(*Request, *RouteMatch) bool
}

type Request struct {
	URL    *url.URL
	Method string
}

func NewHandlerFromRoute(r Rout) Handler {
	var h Handler
	if r.IsSpecial {
		h = Handler{
			Bucket:         "",
			File:           "",
			SpecialHandler: r.Handler,
		}
	} else {
		h = NewHandlerFromString(r.Handler)
	}

	h.Licenses = r.Licenses
	h.Path = r.Path
	h.Methods = r.Methods

	return h
}

func NewHandlerFromString(str string) Handler {
	hConf := strings.Fields(str)

	if len(hConf) != 2 {
		log.Warningf("not valid handler, route=%v", hConf)
		return Handler{}
	}

	return NewHandler(hConf[0], hConf[1])
}

func NewHandler(bucket, file string) Handler {
	return Handler{
		Bucket: bucket,
		File:   file}
}

type Handler struct {
	Bucket string
	File   string

	SpecialHandler string

	Licenses []string
	Methods  []string
	Path     string
}

func (h Handler) IsEmpty() bool {

	return len(h.Bucket) == 0 && len(h.File) == 0
}

type RouteMatch struct {
	Route   *Route
	Handler Handler
	Vars    map[string]string
}

type Route struct {
	handler Handler

	// // List of matchers.
	matchers []matcher

	// // Manager for the variables from host and path.
	regexp *routeRegexpGroup

	name string

	strictSlash bool

	err error
}

func (r *Route) Match(req *Request, match *RouteMatch) bool {
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

func (r *Route) Path(tpl string) *Route {
	r.addRegexpMatcher(tpl, false, false, false)

	return r
}

// methodMatcher matches the request against HTTP methods.
type methodMatcher []string

func (m methodMatcher) Match(r *Request, match *RouteMatch) bool {
	return matchInArray(m, r.Method)
}

// Methods adds a matcher for HTTP methods.
// It accepts a sequence of one or more methods to be matched, e.g.:
// "GET", "POST", "PUT".
func (r *Route) Methods(methods ...string) *Route {
	for k, v := range methods {
		methods[k] = strings.ToUpper(v)
	}

	return r.addMatcher(methodMatcher(methods))
}

func (r *Route) Handler(v Handler) *Route {
	r.handler = v
	return r
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m matcher) *Route {
	if r.err == nil {
		r.matchers = append(r.matchers, m)
	}
	return r
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

type Router struct {
	NotFoundHandler Handler

	routes []*Route

	namedRoutes map[string]*Route
}

func NewRouter() *Router {

	return &Router{namedRoutes: make(map[string]*Route)}
}

func (r *Router) NewRoute() *Route {
	route := &Route{}
	r.routes = append(r.routes, route)

	return route
}

func (r *Router) Handle(path string, h Handler) *Route {

	return r.NewRoute().Path(path).Handler(h)
}

func (r *Router) Match(req *Request, match *RouteMatch) bool {
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
