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
	if trie.Root.Children["/"].ParamName != "id" {
		t.Error()
	}

	trie.AddRoute("/:id/:property.:format", "")
	if trie.Root.Children["/"].ParamChild.Children["/"].ParamChild.Children["."].ParamChild == nil {
		t.Error()
	}
	if trie.Root.Children["/"].ParamName != "id" {
		t.Error()
	}
	if trie.Root.Children["/"].ParamChild.Children["/"].ParamName != "property" {
		t.Error()
	}
	if trie.Root.Children["/"].ParamChild.Children["/"].ParamChild.Children["."].ParamName != "format" {
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

func is_in_matches(test string, matches []*Match) bool {
	for _, match := range matches {
		if match.Route.(string) == test {
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
	trie.AddRoute("/r/:id/property.*format", "property_format")

	trie.Compress()

	matches := trie.FindRoutes("/")
	if len(matches) != 1 {
		t.Errorf("expected one route, got %d", len(matches))
	}
	if !is_in_matches("root", matches) {
		t.Error("expected 'root'")
	}

	matches = trie.FindRoutes("/notfound")
	if len(matches) != 0 {
		t.Errorf("expected zero route, got %d", len(matches))
	}

	matches = trie.FindRoutes("/r/1")
	if len(matches) != 1 {
		t.Errorf("expected one route, got %d", len(matches))
	}
	if !is_in_matches("resource", matches) {
		t.Errorf("expected 'resource', got %+v", matches)
	}
	if matches[0].Params["id"] != "1" {
		t.Error()
	}

	matches = trie.FindRoutes("/r/1/property")
	if len(matches) != 1 {
		t.Errorf("expected one route, got %d", len(matches))
	}
	if !is_in_matches("property", matches) {
		t.Error("expected 'property'")
	}
	if matches[0].Params["id"] != "1" {
		t.Error()
	}

	matches = trie.FindRoutes("/r/1/property.json")
	if len(matches) != 1 {
		t.Errorf("expected one route, got %d", len(matches))
	}
	if !is_in_matches("property_format", matches) {
		t.Error("expected 'property_format'")
	}
	if matches[0].Params["id"] != "1" {
		t.Error()
	}
	if matches[0].Params["format"] != "json" {
		t.Error()
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

	matches := trie.FindRoutes("/r/1")
	if len(matches) != 2 {
		t.Errorf("expected two matches, got %d", len(matches))
	}
	if !is_in_matches("resource_generic", matches) {
		t.Error()
	}
	if !is_in_matches("resource1", matches) {
		t.Error()
	}

	matches = trie.FindRoutes("/s/1")
	if len(matches) != 2 {
		t.Errorf("expected two matches, got %d", len(matches))
	}
	if !is_in_matches("special_all", matches) {
		t.Error()
	}
	if !is_in_matches("special_generic", matches) {
		t.Error()
	}
}
