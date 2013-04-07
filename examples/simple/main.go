package main

import (
	"fmt"
	"github.com/ant0ine/go-urlrouter"
)

func main() {

	router := urlrouter.NewRouter()

	err := router.AddRoutes([]urlrouter.Route{
		urlrouter.Route{
			Path: "/resources/:id",
			Dest: "one_resource",
		},
		urlrouter.Route{
			Path: "/resources",
			Dest: "all_resources",
		},
	})

	if err != nil {
		panic(err)
	}

	err = router.Start()
	if err != nil {
		panic(err)
	}

	input := "http://example.org/resources/123"
	route, params, err := router.FindRoute(input, "GET")
	if err != nil {
		panic(err)
	}
	fmt.Println(route.Dest)   // one_resource
	fmt.Println(params["id"]) // 123
}
