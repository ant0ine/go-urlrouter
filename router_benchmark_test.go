package urlrouter

import (
	"net/url"
	"testing"
)

func BenchmarkNoCompression(b *testing.B) {

	b.StopTimer()

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

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		router.FindRouteFromURL(url_obj)
	}
}
