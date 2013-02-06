package urlrouter

import (
	"net/url"
	"testing"
)

func BenchmarkSimple(b *testing.B) {

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

	router.Start()

	url_obj, _ := url.Parse("http://example.org/resources/123")
        for i := 0; i < b.N; i++ {
	        router.FindRouteFromURL(url_obj)
        }
}
