// Efficient URL routing using a Trie data structure.
//
// This Package implements a URL Router, but instead of using the usual
// "evaluate all the routes and return the first regexp that matches" strategy,
// it uses a Trie data structure to perform the routing. This is more efficient,
// and scales better for a large number of routes.
// It supports the :param and *splat placeholders in the route strings.
//
// Example:
//  router := urlrouter.NewRouter()
//	err := router.AddRoutes([]urlrouter.Route{
//			urlrouter.Route{
//				PathExp: "/resources/:id",
//				Dest:    "one_resource",
//			},
//			urlrouter.Route{
//				PathExp: "/resources",
//				Dest:    "all_resources",
//			},
//		},
//	})
//
//  if err != nil {
//      panic(err)
//  }
//	err = router.Start()
//	if err != nil {
//		panic(err)
//	}
//
//	input := "http://example.org/resources/123"
//	route, params, err := router.FindRoute(input)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Print(route.Dest)  // one_resource
//	fmt.Print(params["id"])  // 123
//
// (Blog Post: http://blog.ant0ine.com/typepad/2013/02/better-url-routing-golang-1.html)
package urlrouter

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-urlrouter/trie"
	"net/http"
	"net/url"
	"strings"
)

// TODO
// support for #param placeholder ?
// replace map[string]string by a PathParams object for more flexibility, and support for multiple param with the same names, eg: /users/:id/blogs/:id ?

type Route struct {
	// A string defining the route, like "/resource/:id.json".
	// Placeholders supported are:
	// :param that matches any char to the first '/' or '.'
	// *splat that matches everything to the end of the string
	// PathExp string
	// Can be anything useful to point to the code to run for this route.
	Dest       interface{}
	HttpMethod string
	Path       string
}

type Router struct {
	// list of Routes, the order matters, if multiple Routes match, the first defined will be used.
	Routes                 []Route
	disableTrieCompression bool
	// By default, support for the method OPTIONS is added to all the routes
	// and add a 405 Method Not Allowed response if the method is not existing
	// for the path.  By setting this option to false, the routes are not added.
	disableDefaultOptions  bool
	index                  map[*Route]int
	trie                   *trie.Trie
	routesIndex            map[string]bool
}

func (self *Router) AddRoutes(routes []Route) error {
	for _, route := range routes {
		error := self.AddRoute(route)
		if error != nil {
			return error
		}
	}
	return nil
}

// Add a new route to the router.  If there's already a router with for this
// path and HTTP method, an error is returned.
func (self *Router) AddRoute(route Route) error {
	PathExp := route.HttpMethod + route.Path
	if self.routesIndex[PathExp] == true {
		return errors.New(fmt.Sprintf("Duplicated PathExp: %s", PathExp))
	}
	self.Routes = append(self.Routes, route)
	self.routesIndex[PathExp] = true
	return nil
}

// Returns a new router.
func NewRouter() Router {
	return Router{
		Routes:      []Route{},
		routesIndex: map[string]bool{},
	}
}

// This validates the Routes and prepares the Trie data structure.
/// It must be called once the Routes are defined and before trying to find Routes.
func (self *Router) Start() error {

	self.trie = trie.New()
	self.index = map[*Route]int{}

	for i, _ := range self.Routes {
		// pointer to the Route
		route := &self.Routes[i]
		// index
		self.index[route] = i
		// insert in the Trie
		err := self.trie.AddRoute(route.Path, route.HttpMethod, route)
		if err != nil {
			return err
		}
	}

	if self.disableDefaultOptions == false {
		self.AddNotAllowed()
	}

	if self.disableTrieCompression == false {
		self.trie.Compress()
	}

	// TODO validation of the PathExp ? start with a /
	// TODO url encoding

	return nil
}

func (self *Router) addAllowedRoutes(path string, allowed []string) {
	optionsFunc := func(w http.ResponseWriter, req *http.Request, params map[string]string) {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		w.WriteHeader(204)
		r := []byte{}
		w.Write(r)
	}
	route := Route{
		Path:       path,
		HttpMethod: "OPTIONS",
		Dest:       optionsFunc,
	}
	self.trie.AddRoute(route.Path, route.HttpMethod, &route)

}

func (self *Router) addDisallowedRoutes(path string, allowed []string, unallowed []string) {
	unallowedFunc := func(w http.ResponseWriter, req *http.Request, params map[string]string) {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(405)
		r := []byte{}
		w.Write(r)
	}
	for _, m := range unallowed {
		route := Route{
			Path:       path,
			HttpMethod: m,
			Dest:       unallowedFunc,
		}
		self.trie.AddRoute(route.Path, route.HttpMethod, &route)
	}
}

func (self *Router) addOptions(path string, methods map[string]bool) {
	allowed := []string{}
	unallowed := []string{}

	for _, dm := range trie.HttpDefaultMethods {
		if methods[dm] == false {
			unallowed = append(unallowed, dm)
		} else {
			allowed = append(allowed, dm)
		}
	}
	self.addAllowedRoutes(path, allowed)
	self.addDisallowedRoutes(path, allowed, unallowed)
}

func (self *Router) AddNotAllowed() {
	for path, methods := range self.trie.AllRoutes() {
		// I would like to also add HEAD here, but this means a few changes:
		// - functions should not write the response, but returns the status, headers and body
		// - the main handler write the response
		// - the HEAD function will call the function associated with GET and drop the body,
		//   and write only the HEADERS + status (as per the RFC).
		// Since this is a bigger change, I'm not going to apply that patch for now.
		if methods["OPTIONS"] == false {
			self.addOptions(path, methods)
		}
	}
}

// Return the first matching Route and the corresponding parameters for a given URL object.
func (self *Router) FindRouteFromURL(urlObj *url.URL, method string) (*Route, map[string]string) {

	// lookup the routes in the Trie
	// TODO verify url encoding
	matches := self.trie.FindRoutes(urlObj.Path, method)

	// only return the first Route that matches
	minIndex := -1
	matchesByIndex := map[int]*trie.Match{}

	for _, match := range matches {
		route := match.Route.(*Route)
		routeIndex := self.index[route]
		matchesByIndex[routeIndex] = match
		if minIndex == -1 || routeIndex < minIndex {
			minIndex = routeIndex
		}
	}

	if minIndex == -1 {
		// no route found
		return nil, nil
	}

	// and the corresponding params
	match := matchesByIndex[minIndex]

	return match.Route.(*Route), match.Params
}

// Parse the url string (complete or just the path) and return the first matching Route and the corresponding parameters.
func (self *Router) FindRoute(urlStr string, method string) (*Route, map[string]string, error) {

	// parse the url
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}

	route, params := self.FindRouteFromURL(urlObj, method)
	return route, params, nil
}
