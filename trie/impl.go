// Special Trie implementation for the URLRouter.
//
// This Trie implementation is designed to support strings that includes
// :param and *splat parameters. Strings that are commonly used for URL
// routing. You probably don't need to use this package directly.
//
// Example:
//  TBD
package trie

// TODO
// support compression
// support *splat
// API
// tests
// benchmarks
// remove the dependency on json

import (
	"encoding/json"
	"fmt"
)

// A Node of the Trie
type Node struct {
	Route    interface{}
	Children map[string]*Node
	//children_key_length int
}

// TODO rename these token methods
func get_find_token(full string, length int) (string, string) {
	return full[0:length], full[length:]
}

func get_insert_token(chars string) (string, string) {

	// TODO test oversize
	token := chars[0:1]
	remaining := chars[1:]

	if token[0] == ':' {
		// this is a route :param
		for len(remaining) > 0 && remaining[0] != '/' {
			remaining = remaining[1:]
		}
		return ":PARAM", remaining
	}

	return token, remaining
}

func get_param_token(remaining string) (string, string) {
	for len(remaining) > 0 && remaining[0] != '/' {
		remaining = remaining[1:]
	}
	return ":PARAM", remaining
}

func (self *Node) add_route(chars string, route interface{}) {

	if len(chars) == 0 {
		self.Route = route
		return
	}

	if self.Children == nil {
		self.Children = map[string]*Node{}
	}

	// ask for 1 char token during the insert
	token, remaining := get_insert_token(chars)

	next_node := self.Children[token]
	if next_node == nil {
		next_node = &Node{}
		self.Children[token] = next_node
	}
	next_node.add_route(remaining, route)
}

func (self *Node) find_routes(path string) []interface{} {

	routes := []interface{}{}

	if self.Route != nil {
		routes = append(routes, self.Route)
	}

	if len(path) < 1 {
		return routes
	}

	// :param branch
	if self.Children[":PARAM"] != nil {
		_, remaining := get_param_token(path)
		routes = append(routes, self.Children[":PARAM"].find_routes(remaining)...)
	}

	// main branch
	token, remaining := get_find_token(path, 1)
	next := self.Children[token]
	if next != nil {
		routes = append(routes, next.find_routes(remaining)...)
	}

	return routes
}

func (self *Node) print_json() {
	bytes, err := json.MarshalIndent(self, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", bytes)
}

// The Trie
type Trie struct {
	Root *Node
}

// Instanciate a Trie with an empty Node as the root.
func New() *Trie {
	return &Trie{
		Root: &Node{},
	}
}

// Insert the route in the Trie following or creating the Nodes corresponding to the path.
func (self *Trie) AddRoute(path string, route interface{}) {
	self.Root.add_route(path, route)
}

// Given a path, return all the matchin routes.
func (self *Trie) FindRoutes(path string) []interface{} {
	return self.Root.find_routes(path)
}
