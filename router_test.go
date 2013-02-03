package urlrouter

import (
	"net/url"
	"testing"
)

func TestFindRouteAPI(t *testing.T) {

	router := Router{
		Routes: []*Route{
			&Route{
				PathExp: "/",
				Dest:    "root",
			},
		},
	}

	err := router.Prepare()
	if err != nil {
		t.Fatal()
	}

	// full url string
	input := "http://example.org/"
	route, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}
	if route.Dest != "root" {
		t.Error()
	}

	// part of the url string
	input = "/"
	route, err = router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}
	if route.Dest != "root" {
		t.Error()
	}

	// url object
	url_obj, err := url.Parse("http://example.org/")
	if err != nil {
		t.Fatal()
	}
	route = router.FindRouteFromURL(url_obj)
	if route.Dest != "root" {
		t.Error()
	}
}

func TestNoRoute(t *testing.T) {

	router := Router{
		Routes: []*Route{},
	}

	err := router.Prepare()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/notfound"
	route, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}

	if route != nil {
		t.Error("should not be able to find a route")
	}
}

func TestDuplicatedRoute(t *testing.T) {

	router := Router{
		Routes: []*Route{
			&Route{
				PathExp: "/",
				Dest:    "root",
			},
			&Route{
				PathExp: "/",
				Dest:    "the_same",
			},
		},
	}

	err := router.Prepare()
	if err == nil {
		t.Error("expected the duplicated route error")
	}
}

func TestRouteOrder(t *testing.T) {

	router := Router{
		Routes: []*Route{
			&Route{
				PathExp: "/r/:id",
				Dest:    "first",
			},
			&Route{
				PathExp: "/r/*rest",
				Dest:    "second",
			},
		},
	}

	err := router.Prepare()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/r/123"
	route, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}

	if route.Dest != "first" {
		t.Error("both match, expected the first defined")
	}
}

func TestSimpleExample(t *testing.T) {

	router := Router{
		Routes: []*Route{
			&Route{
				PathExp: "/resources/:id",
				Dest:    "one_resource",
			},
			&Route{
				PathExp: "/resources",
				Dest:    "all_resources",
			},
		},
	}

	err := router.Prepare()
	if err != nil {
		t.Fatal()
	}

	input := "http://example.org/resources/123"
	route, err := router.FindRoute(input)
	if err != nil {
		t.Fatal()
	}

	if route.Dest != "one_resource" {
		t.Error()
	}
}
