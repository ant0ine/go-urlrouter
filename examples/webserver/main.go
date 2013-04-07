package main

import (
	"fmt"
	"github.com/ant0ine/go-urlrouter"
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Hello %s", params["name"])
}

func Bonjour(w http.ResponseWriter, req *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Bonjour %s", params["name"])
}

func main() {

	router := urlrouter.NewRouter()

	error := router.AddRoutes(
		[]urlrouter.Route{
			urlrouter.Route{
				Path:       "/hello/:name",
				Dest:       Hello,
				HttpMethod: "GET",
			},
			urlrouter.Route{
				Path:       "/bonjour/:name",
				Dest:       Bonjour,
				HttpMethod: "GET",
			},
		},
	)

	if error != nil {
		panic(error)
	}

	router.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		route, params := router.FindRouteFromURL(r.URL, r.Method)
		handler := route.Dest.(func(http.ResponseWriter, *http.Request, map[string]string))
		handler(w, r, params)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
