package urlrouter

import (
	"testing"
	//"fmt"
)

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

	input := "http://example.org/resources/123.json"
	route, err := router.FindRouteFromString(input)
	if err != nil {
		t.Fatal()
	}

	if route.Dest != "one_resource" {
		t.Error()
	}
}
