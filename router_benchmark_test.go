package urlrouter

import (
	"fmt"
	"net/url"
	"regexp"
	"testing"
)

func routes() []Route {
	// simulate the routes of a real but reasonable app.
	// 6 + 10 * (5 + 2) + 1 = 77 routes
	route_paths := []string{
		"/",
		"/signin",
		"/signout",
		"/profile",
		"/settings",
		"/upload/*file",
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 5; j++ {
			route_paths = append(route_paths, fmt.Sprintf("/resource%d/:id/property%d", i, j))
		}
		route_paths = append(route_paths, fmt.Sprintf("/resource%d/:id", i))
		route_paths = append(route_paths, fmt.Sprintf("/resource%d", i))
	}
	route_paths = append(route_paths, "/*")

	routes := []Route{}
	for _, path := range route_paths {
		routes = append(routes, Route{PathExp: path, Dest: path})
	}
	return routes
}

func request_urls() []*url.URL {
	// simulate a few requests
	url_strs := []string{
		"http://example.org/",
		"http://example.org/resource9/123",
		"http://example.org/resource9/123/property1",
		"http://example.org/doesnotexist",
	}
	url_objs := []*url.URL{}
	for _, url_str := range url_strs {
		url_obj, _ := url.Parse(url_str)
		url_objs = append(url_objs, url_obj)
	}
	return url_objs
}

func BenchmarkNoCompression(b *testing.B) {

	b.StopTimer()

	router := Router{
		Routes: routes(),
		disable_trie_compression: true,
	}
	router.Start()
	url_objs := request_urls()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, url_obj := range url_objs {
			router.FindRouteFromURL(url_obj)
		}
	}
}

func BenchmarkCompression(b *testing.B) {

	b.StopTimer()

	router := Router{
		Routes: routes(),
	}
	router.Start()
	url_objs := request_urls()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, url_obj := range url_objs {
			router.FindRouteFromURL(url_obj)
		}
	}
}

func BenchmarkRegExpLoop(b *testing.B) {
	// reference benchmark using the usual RegExps + Loop strategy

	b.StopTimer()

	routes := routes()
	url_objs := request_urls()

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
		reg_str := r2.ReplaceAllString(route.PathExp, "([^/\\.]+)")
		reg_str = r1.ReplaceAllString(reg_str, "(.+)")
		reg_str = "^" + reg_str + "$"

		// compile it
		reg, err := regexp.Compile(reg_str)
		if err != nil {
			panic(err)
		}

		route_regexps = append(route_regexps, *reg)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// do it for a few urls
		for _, url_obj := range url_objs {
			// stop at the first route that matches
			for index, reg := range route_regexps {
				if reg.FindAllString(url_obj.Path, 1) != nil {
					_ = routes[index]
					break
				}
			}
		}
	}
}
