package urlrouter

import (
	"net/url"
	"regexp"
	"testing"
)

func routes() []Route {
	route_paths := []string{
		"/",
		"/signin",
		"/signout",
		"/profile",
		"/settings",
		"/upload/*file",
		"/apps/:id/property1",
		"/apps/:id/property2",
		"/apps/:id/property3",
		"/apps/:id/property4",
		"/apps/:id/property5",
		"/apps/:id",
		"/apps",
		"/users/:id/property1",
		"/users/:id/property2",
		"/users/:id/property3",
		"/users/:id/property4",
		"/users/:id/property5",
		"/users/:id",
		"/users",
		"/resources/:id/property1",
		"/resources/:id/property2",
		"/resources/:id/property3",
		"/resources/:id/property4",
		"/resources/:id/property5",
		"/resources/:id",
		"/resources",
		"/*",
	}
	routes := []Route{}
	for _, path := range route_paths {
		routes = append(routes, Route{PathExp: path, Dest: path})
	}
	return routes
}

func BenchmarkNoCompression(b *testing.B) {

	b.StopTimer()

	router := Router{
		Routes: routes(),
		disable_trie_compression: true,
	}
	router.Start()
	url_obj, _ := url.Parse("http://example.org/resources/123")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		router.FindRouteFromURL(url_obj)
	}
}

func BenchmarkCompression(b *testing.B) {

	b.StopTimer()

	router := Router{
		Routes: routes(),
	}
	router.Start()
	url_obj, _ := url.Parse("http://example.org/resources/123")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		router.FindRouteFromURL(url_obj)
	}
}

func BenchmarkRegExpLoop(b *testing.B) {
	// reference benchmark using the usual RegExps + Loop strategy

	b.StopTimer()

	routes := routes()

	// build the route regexps
	r1, err := regexp.Compile(":[^/\\.]*")
	if err != nil {
		panic(err)
	}
	r2, err := regexp.Compile("\\*.*")
	if err != nil {
		panic(err)
	}
	route_regexps := []regexp.Regexp{}
	for _, route := range routes {

		// generate the regexp string
		reg_str := r2.ReplaceAllString(route.PathExp, "[^/\\.]+")
		reg_str = r1.ReplaceAllString(reg_str, ".+")
		reg_str = "^" + reg_str + "$"

		// compile it
		reg, err := regexp.Compile(reg_str)
		if err != nil {
			panic(err)
		}

		route_regexps = append(route_regexps, *reg)
	}

	// url to route
	url_obj, _ := url.Parse("http://example.org/resources/123")

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for index, reg := range route_regexps {
			if reg.MatchString(url_obj.Path) {
				_ = routes[index]
				break
			}
		}
	}
}
