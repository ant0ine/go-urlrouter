package trie

import (
	"testing"
)

func TestPathInsert(t *testing.T) {

	trie := New()
	if trie.Root == nil {
		t.Error()
	}

	trie.AddRoute("/", "1")
	if trie.Root.Children["/"] == nil {
		t.Error()
	}

	trie.AddRoute("/r", "2")
	if trie.Root.Children["/"].Children["r"] == nil {
		t.Error()
	}

	trie.AddRoute("/r/", "3")
	if trie.Root.Children["/"].Children["r"].Children["/"] == nil {
		t.Error()
	}
}

func TestTrieCompression(t *testing.T) {

	trie := New()
	trie.AddRoute("/abc", "3")
	trie.AddRoute("/adc", "3")

	// before compression
	if trie.Root.Children["/"].Children["a"].Children["b"].Children["c"] == nil {
		t.Error()
	}
	if trie.Root.Children["/"].Children["a"].Children["d"].Children["c"] == nil {
		t.Error()
	}

	trie.Compress()

	// after compression
	if trie.Root.Children["/abc"] == nil {
		t.Errorf("%+v", trie.Root)
	}
	if trie.Root.Children["/adc"] == nil {
		t.Errorf("%+v", trie.Root)
	}

}
func TestParamInsert(t *testing.T) {
	trie := New()

	trie.AddRoute("/:id/", "")
	if trie.Root.Children["/"].ParamChild.Children["/"] == nil {
		t.Error()
	}

	trie.AddRoute("/:id/:property.:format", "")
	if trie.Root.Children["/"].ParamChild.Children["/"].ParamChild.Children["."].ParamChild == nil {
		t.Error()
	}
}

func TestSplatInsert(t *testing.T) {
	trie := New()
	trie.AddRoute("/*splat", "")
	if trie.Root.Children["/"].SplatChild == nil {
		t.Error()
	}
}

func TestDupeInsert(t *testing.T) {
	trie := New()
	trie.AddRoute("/", "1")
	err := trie.AddRoute("/", "2")
	if err == nil {
		t.Error()
	}
	if trie.Root.Children["/"].Route != "1" {
		t.Error()
	}
}

func is_in(test string, sources []interface{}) bool {
	for _, source := range sources {
		if source.(string) == test {
			return true
		}
	}
	return false
}

func TestFindRoute(t *testing.T) {

	trie := New()

	trie.AddRoute("/", "root")
	trie.AddRoute("/r/:id", "resource")
	trie.AddRoute("/r/:id/property", "property")
	trie.AddRoute("/r/:id/property*", "property_format")

	trie.Compress()

	routes := trie.FindRoutes("/")
	if len(routes) != 1 {
		t.Errorf("expected one route, got %d", len(routes))
	}
	if !is_in("root", routes) {
		t.Error("expected 'root'")
	}

	routes = trie.FindRoutes("/notfound")
	if len(routes) != 0 {
		t.Errorf("expected zero route, got %d", len(routes))
	}

	routes = trie.FindRoutes("/r/1")
	if len(routes) != 1 {
		t.Errorf("expected one route, got %d", len(routes))
	}
	if !is_in("resource", routes) {
		t.Errorf("expected 'resource', got %+v", routes)
	}

	routes = trie.FindRoutes("/r/1/property")
	if len(routes) != 1 {
		t.Errorf("expected one route, got %d", len(routes))
	}
	if !is_in("property", routes) {
		t.Error("expected 'property'")
	}

	routes = trie.FindRoutes("/r/1/property.json")
	if len(routes) != 1 {
		t.Errorf("expected one route, got %d", len(routes))
	}
	if !is_in("property_format", routes) {
		t.Error("expected 'property_format'")
	}
}

func TestFindRouteMultipleMatches(t *testing.T) {

	trie := New()

	trie.AddRoute("/r/1", "resource1")
	trie.AddRoute("/r/2", "resource2")
	trie.AddRoute("/r/:id", "resource_generic")
	trie.AddRoute("/s/*rest", "special_all")
	trie.AddRoute("/s/:param", "special_generic")
	trie.AddRoute("/", "root")

	trie.Compress()

	routes := trie.FindRoutes("/r/1")
	if len(routes) != 2 {
		t.Errorf("expected two routes, got %d", len(routes))
	}
	if !is_in("resource_generic", routes) {
		t.Error()
	}
	if !is_in("resource1", routes) {
		t.Error()
	}

	routes = trie.FindRoutes("/s/1")
	if len(routes) != 2 {
		t.Errorf("expected two routes, got %d", len(routes))
	}
	if !is_in("special_all", routes) {
		t.Error()
	}
	if !is_in("special_generic", routes) {
		t.Error()
	}
}
