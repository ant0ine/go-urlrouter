package trie

import (
	"testing"
)

func TestSimpleExample(t *testing.T) {

	trie := New()

	trie.AddRoute("/r/1", "resource1")
	trie.AddRoute("/r/2", "resource2")
	trie.AddRoute("/r/:id", "resources")
	trie.AddRoute("/", "root")

	routes := trie.FindRoutes("/r/1")

	if len(routes) != 3 {
		t.Error()
	}
}
