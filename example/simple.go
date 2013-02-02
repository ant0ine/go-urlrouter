package main

import (
	"fmt"
	"github.com/ant0ine/go-urlrouter"
)

func main() {

	router := urlrouter.Router{
		Routes: []*urlrouter.Route{
			&urlrouter.Route{
				PathExp: "/resources/:id",
				Dest:    "one_resource",
			},
			&urlrouter.Route{
				PathExp: "/resources",
				Dest:    "all_resources",
			},
		},
	}

	err := router.Prepare()
	if err != nil {
		panic(err)
	}

	input := "http://example.org/resources/123.json"
	route := router.FindRoute(input)
	fmt.Print(route.Dest)
}
