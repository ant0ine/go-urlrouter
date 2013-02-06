// Special Trie implementation for the URLRouter.
//
// This Trie implementation is designed to support strings that includes
// :param and *splat parameters. Strings that are commonly used for URL
// routing. You probably don't need to use this package directly.
//
// Example:
//  TBD
package trie

// TODO support param map for matched routes

import (
	"errors"
)

// A Node of the Trie
type Node struct {
	Route          interface{}
	Children       map[string]*Node
	ChildrenKeyLen int
	ParamChild     *Node
	SplatChild     *Node
}

func get_param_remaining(remaining string) string {
	i := 0
	for len(remaining) > i && remaining[i] != '/' && remaining[i] != '.' {
		i++
	}
	return remaining[i:]
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
			self.ChildrenKeyLen = 1
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

	if self.Route != nil && path == "" {
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
	length := self.ChildrenKeyLen
	if len(path) < length {
		return routes
	}
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

func (self *Node) compress() {
	// *splat branch
	if self.SplatChild != nil {
		self.SplatChild.compress()
	}
	// :param branch
	if self.ParamChild != nil {
		self.ParamChild.compress()
	}
	// main branch
	if len(self.Children) == 0 {
		return
	}
	// compressable ?
	can_compress := true
	for _, node := range self.Children {
		if node.Route != nil || node.SplatChild != nil || node.ParamChild != nil {
			can_compress = false
		}
	}
	// compress
	if can_compress {
		merged := map[string]*Node{}
		for key, node := range self.Children {
			for gd_key, gd_node := range node.Children {
				merged_key := key + gd_key
				merged[merged_key] = gd_node
			}
		}
		self.Children = merged
		self.ChildrenKeyLen++
		self.compress()
		// continue
	} else {
		for _, node := range self.Children {
			node.compress()
		}
	}
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

// Reduce the size of the tree, must be done after the last AddRoute.
func (self *Trie) Compress() {
	self.Root.compress()
}
