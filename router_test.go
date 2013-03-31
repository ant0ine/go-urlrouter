package urlrouter

import (
	"net/url"
	"testing"
)

func TestFindRouteAPI(t *testing.T) {

	router := Router{
		Routes: []Route{
			Route{
				PathExp: "/",
				Dest:    "root",
			},
		},
	}

	err := router.Start()
	if err != nil {
		t.Fatal()
	}

	// full url string
	input := "http://example.org/"
	route, params, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}
	if route.Dest != "root" {
		t.Error()
	}
	if len(params) != 0 {
		t.Error()
	}

	// part of the url string
	input = "/"
	route, params, err = router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}
	if route.Dest != "root" {
		t.Error()
	}
	if len(params) != 0 {
		t.Error()
	}

	// url object
	urlObj, err := url.Parse("http://example.org/")
	if err != nil {
		t.Fatal()
	}
	route, params = router.FindRouteFromURL(urlObj)
	if route.Dest != "root" {
		t.Error()
	}
	if len(params) != 0 {
		t.Error()
	}
}

func TestNoRoute(t *testing.T) {

	router := Router{
		Routes: []Route{},
	}

	err := router.Start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/notfound"
	route, params, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}

	if route != nil {
		t.Error("should not be able to find a route")
	}
	if params != nil {
		t.Error("params must be nil too")
	}
}

func TestDuplicatedRoute(t *testing.T) {

	router := Router{
		Routes: []Route{
			Route{
				PathExp: "/",
				Dest:    "root",
			},
			Route{
				PathExp: "/",
				Dest:    "the_same",
			},
		},
	}

	err := router.Start()
	if err == nil {
		t.Error("expected the duplicated route error")
	}
}

func TestRouteOrder(t *testing.T) {

	router := Router{
		Routes: []Route{
			Route{
				PathExp: "/r/:id",
				Dest:    "first",
			},
			Route{
				PathExp: "/r/*rest",
				Dest:    "second",
			},
		},
	}

	err := router.Start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/r/123"
	route, params, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}

	if route.Dest != "first" {
		t.Errorf("both match, expected the first defined, got %s", route.Dest)
	}
	if params["id"] != "123" {
		t.Error()
	}
}

func TestSimpleExample(t *testing.T) {

	router := Router{
		Routes: []Route{
			Route{
				PathExp: "/resources/:id",
				Dest:    "one_resource",
			},
			Route{
				PathExp: "/resources",
				Dest:    "all_resources",
			},
		},
	}

	err := router.Start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/resources/123"
	route, params, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}

	if route.Dest != "one_resource" {
		t.Error()
	}
	if params["id"] != "123" {
		t.Error()
	}
}
