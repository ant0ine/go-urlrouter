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

func TestParamInsert(t *testing.T) {
	trie := New()

	trie.AddRoute("/:id/", "")
	if trie.Root.Children["/"].ParamChild.Children["/"] == nil {
		t.Error()
	}

	trie.AddRoute("/:id/:property.json", "")
	if trie.Root.Children["/"].ParamChild.Children["/"].ParamChild.Children["."] == nil {
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

func TestSimpleExample(t *testing.T) {

	trie := New()

	trie.AddRoute("/r/1", "resource1")
	trie.AddRoute("/r/2", "resource2")
	trie.AddRoute("/r/:id", "resources")
	trie.AddRoute("/s/*rest", "resources")
	trie.AddRoute("/", "root")

	routes := trie.FindRoutes("/r/1")
	if len(routes) != 3 {
		t.Error()
	}

	routes = trie.FindRoutes("/s/1")
	if len(routes) != 2 {
		t.Error()
	}

	routes = trie.FindRoutes("/t")
	if len(routes) != 1 {
		t.Error()
	}
}
