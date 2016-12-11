package urlrouter

import (
	"net/url"
	"testing"
)

var HttpMethods = [4]string{"GET", "PUT", "POST", "DELETE"}

func TestFindRouteAPI(t *testing.T) {

	router := NewRouter()

	err := router.AddRoute(Route{
		Path: "/",
		Dest: "root",
	})

	if err != nil {
		t.Fatal()
	}

	router.Start()

	// full url string
	input := "http://example.org/"
	for _, method := range HttpMethods {
		route, params, err := router.FindRoute(input, method)
		if err != nil {
			t.Fatal()
		}
		if route.Dest != "root" {
			t.Error()
		}
		if len(params) != 0 {
			t.Error()
		}
	}

	// part of the url string
	input = "/"
	route, params, err := router.FindRoute(input, "GET")
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
	route, params = router.FindRouteFromURL(urlObj, "GET")
	if route.Dest != "root" {
		t.Error()
	}
	if len(params) != 0 {
		t.Error()
	}
}

func TestNoRoute(t *testing.T) {

	router := NewRouter()
	err := router.Start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/notfound"
	route, params, err := router.FindRoute(input, "GET")
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

	router := NewRouter()

	err := router.AddRoutes([]Route{
		Route{
			Path: "/",
			Dest: "root",
		},
		Route{
			Path: "/",
			Dest: "the_same",
		},
	})

	if err == nil {
		t.Error("expected the duplicated route error")
	}
}

func TestRouteOrder(t *testing.T) {

	router := NewRouter()
	err := router.AddRoutes([]Route{
		Route{
			Path: "/r/:id",
			Dest: "first",
		},
		Route{
			Path: "/r/*rest",
			Dest: "second",
		},
	})

	if err != nil {
		t.Fatal()
	}

	err = router.Start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/r/123"
	route, params, err := router.FindRoute(input, "GET")
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

	router := NewRouter()
	err := router.AddRoutes([]Route{
		Route{
			Path: "/resources/:id",
			Dest: "one_resource",
		},
		Route{
			Path: "/resources",
			Dest: "all_resources",
		},
	})

	if err != nil {
		t.Fatal()
	}

	err = router.Start()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/resources/123"
	route, params, err := router.FindRoute(input, "GET")
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
