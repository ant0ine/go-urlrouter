// Efficient URL routing using a Trie data structure.
//
// TODO description
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
//	input := "http://example.org/resources/123.json"
//	route, err := router.FindRouteFromString(input)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Print(route.Dest)
//
package urlrouter

import (
	"github.com/ant0ine/go-urlrouter/trie"
	"net/url"
)

// TODO explain why Dest as an interface{} is nice
type Route struct {
	PathExp string
	Dest    interface{}
}

//
type Router struct {
	Routes []*Route
	index  map[*Route]int
	trie   *trie.Trie
}

//
func (self *Router) Prepare() error {

	self.trie = trie.New()

	self.index = map[*Route]int{}

	for i, route := range self.Routes {
		// index
		self.index[route] = i
		// insert in the Trie
		self.trie.AddRoute(route.PathExp, route)
	}

	// TODO route.PathExp should be unique ?
	// TODO validation of the PathExp ? start with a /
	// TODO url encoding
	// TODO compress the Trie (when supported)

	return nil
}

//
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

	return self.Routes[min_index]
}

//
func (self *Router) FindRouteFromString(url_str string) (*Route, error) {

	// parse the url
	url_obj, err := url.Parse(url_str)
	if err != nil {
		return nil, err
	}

	return self.FindRouteFromURL(url_obj), nil
}
