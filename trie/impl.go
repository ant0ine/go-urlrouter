// Special Trie implementation for the URLRouter.
//
// This Trie implementation is designed to support strings that includes
// :param and *splat parameters. Strings that are commonly used for URL
// routing. You probably don't need to use this package directly.
//
// Example:
//  TBD
package trie

// TODO support compression

import (
	"errors"
)

// A Node of the Trie
type Node struct {
	Route      interface{}
	Children   map[string]*Node
	ParamChild *Node
	SplatChild *Node
	//children_key_length int
}

func get_param_remaining(remaining string) string {
	for len(remaining) > 0 && remaining[0] != '/' && remaining[0] != '.' {
		remaining = remaining[1:]
	}
	return remaining
}

func (self *Node) add_route(path string, route interface{}) error {

	if len(path) == 0 {
		// end of the path, set the Route
		if self.Route != nil {
			return errors.New("Node.Route already set, duplicated path")
		}
		self.Route = route
		return nil
	}

	token := path[0:1]
	remaining := path[1:]
	var next_node *Node

	if token[0] == ':' {
		// :param case
		remaining = get_param_remaining(remaining)
		if self.ParamChild == nil {
			self.ParamChild = &Node{}
		}
		next_node = self.ParamChild
	} else if token[0] == '*' {
		// *splat case
		remaining = ""
		if self.SplatChild == nil {
			self.SplatChild = &Node{}
		}
		next_node = self.SplatChild
	} else {
		// general case
		if self.Children == nil {
			self.Children = map[string]*Node{}
		}
		if self.Children[token] == nil {
			self.Children[token] = &Node{}
		}
		next_node = self.Children[token]
	}

	return next_node.add_route(remaining, route)
}

func (self *Node) find_routes(path string) []interface{} {

	routes := []interface{}{}

	if self.Route != nil {
		routes = append(routes, self.Route)
	}

	if len(path) == 0 {
		return routes
	}

	// *splat branch
	if self.SplatChild != nil {
		routes = append(
			routes,
			self.SplatChild.find_routes("")...,
		)
	}

	// :param branch
	if self.ParamChild != nil {
		remaining := get_param_remaining(path)
		routes = append(
			routes,
			self.ParamChild.find_routes(remaining)...,
		)
	}

	// main branch
	length := 1
	token := path[0:length]
	remaining := path[length:]
	if self.Children[token] != nil {
		routes = append(
			routes,
			self.Children[token].find_routes(remaining)...,
		)
	}

	return routes
}

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
func (self *Trie) AddRoute(path string, route interface{}) error {
	return self.Root.add_route(path, route)
}

// Given a path, return all the matchin routes.
func (self *Trie) FindRoutes(path string) []interface{} {
	return self.Root.find_routes(path)
}
