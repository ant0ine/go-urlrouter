// Efficient URL routing using a Trie data structure.
//
// This Package implements a URL Router, but instead of using the usual
// "evaluate all the routes and return the first rexexp that matches" strategy,
// it uses a Trie data structure to perform the routing. This is more efficient,
// and scales better for a large number of routes.
// It supports the usual :param and *splat route syntax.
//
// Example:
//	router := urlrouter.Router{
//		Routes: []*urlrouter.Route{
//			&urlrouter.Route{
//				PathExp: "/resources/:id",
//				Dest:    "one_resource",
//			},
//			&urlrouter.Route{
//				PathExp: "/resources",
//				Dest:    "all_resources",
//			},
//		},
//	}
//
//	err := router.Prepare()
//	if err != nil {
//		panic(err)
//	}
//
//	input := "http://example.org/resources/123"
//	route, err := router.FindRoute(input)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Print(route.Dest)
//
package urlrouter

import (
	"errors"
	"github.com/ant0ine/go-urlrouter/trie"
	"net/url"
)

// TODO
// benchmarks

type Route struct {
	// a string like "/resource/:id.json"
	PathExp string
	// can be anything useful to point to the code to run for this route.
	Dest interface{}
}

type Router struct {
	Routes []*Route
	index  map[*Route]int
	trie   *trie.Trie
}

// This validates the Routes and prepares the Trie data structure.
// It must be called once the Routes are defined and before trying to find Routes.
func (self *Router) Prepare() error {

	self.trie = trie.New()
	self.index = map[*Route]int{}
	unique := map[string]bool{}

	for i, route := range self.Routes {
		// unique
		if unique[route.PathExp] == true {
			return errors.New("duplicated PathExp")
		}
		unique[route.PathExp] = true
		// index
		self.index[route] = i
		// insert in the Trie
		err := self.trie.AddRoute(route.PathExp, route)
		if err != nil {
			return err
		}
	}

	// TODO validation of the PathExp ? start with a /
	// TODO url encoding
	// TODO compress the Trie (when supported)

	return nil
}

// Return the first matching Route for the given URL object.
func (self *Router) FindRouteFromURL(url_obj *url.URL) *Route {

	// lookup the routes in the Trie
	// TODO verify url encoding
	routes := self.trie.FindRoutes(url_obj.Path)

	// only return the first Route that matches
	min_index := -1
	for _, r := range routes {
		route := r.(*Route)
		i := self.index[route]
		if min_index == -1 || i < min_index {
			min_index = i
		}
	}

	if min_index == -1 {
		// no route found
		return nil
	}

	return self.Routes[min_index]
}

// Parse the url string (complete or just the path) and call FindRouteFromURL
func (self *Router) FindRoute(url_str string) (*Route, error) {

	// parse the url
	url_obj, err := url.Parse(url_str)
	if err != nil {
		return nil, err
	}

	return self.FindRouteFromURL(url_obj), nil
}
