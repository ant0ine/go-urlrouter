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
//	route := router.FindRoute(input)
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

	// TODO validation of the PathExp ? start with a /
	// TODO compress the Trie (when supported)

	return nil
}

func (self *Router) FindRoute(url_str string) *Route {
	// TODO provide another method that takes a url object ?

	// parse the url
	url_obj, err := url.Parse(url_str)
	if err != nil {
		panic(err) // XXX
	}

	// lookup the routes in the Trie
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
