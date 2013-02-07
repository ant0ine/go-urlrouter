// Special Trie implementation for the URLRouter.
//
// This Trie implementation is designed to support strings that includes
// :param and *splat parameters. Strings that are commonly used for URL
// routing. You probably don't need to use this package directly.
//
// Example:
//  TBD
package trie

import (
	"errors"
)

func split_param(remaining string) (string, string) {
	i := 0
	for len(remaining) > i && remaining[i] != '/' && remaining[i] != '.' {
		i++
	}
	return remaining[:i], remaining[i:]
}

// A Node of the Trie
type Node struct {
	Route          interface{}
	Children       map[string]*Node
	ChildrenKeyLen int
	ParamChild     *Node
	ParamName      string
	SplatChild     *Node
	SplatName      string
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
		var name string
		name, remaining = split_param(remaining)
		if self.ParamChild == nil {
			self.ParamChild = &Node{}
			self.ParamName = name
		}
		next_node = self.ParamChild
	} else if token[0] == '*' {
		// *splat case
		name := remaining
		remaining = ""
		if self.SplatChild == nil {
			self.SplatChild = &Node{}
			self.SplatName = name
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

// utility for the Node.find_route recursive method
type pstack struct {
	params []map[string]string
}

func (self *pstack) add(name, value string) {
	self.params = append(
		self.params,
		map[string]string{name: value},
	)
}

func (self *pstack) pop() {
	self.params = self.params[:len(self.params)-1]
}

func (self *pstack) as_map() map[string]string {
	// assume that all param of a route have unique names
	r := map[string]string{}
	for _, param := range self.params {
		for key, value := range param {
			r[key] = value
		}
	}
	return r
}

type Match struct {
	// same Route as in Node
	Route interface{}
	// map of params matched for this result
	Params map[string]string
}

func (self *Node) find_routes(path string, stack *pstack) []*Match {

	matches := []*Match{}

	// route found !
	if self.Route != nil && path == "" {
		matches = append(
			matches,
			&Match{
				Route:  self.Route,
				Params: stack.as_map(),
			},
		)
	}

	if len(path) == 0 {
		return matches
	}

	// *splat branch
	if self.SplatChild != nil {
		stack.add(self.SplatName, path)
		matches = append(
			matches,
			self.SplatChild.find_routes("", stack)...,
		)
		stack.pop()
	}

	// :param branch
	if self.ParamChild != nil {
		value, remaining := split_param(path)
		stack.add(self.ParamName, value)
		matches = append(
			matches,
			self.ParamChild.find_routes(remaining, stack)...,
		)
		stack.pop()
	}

	// main branch
	length := self.ChildrenKeyLen
	if len(path) < length {
		return matches
	}
	token := path[0:length]
	remaining := path[length:]
	if self.Children[token] != nil {
		matches = append(
			matches,
			self.Children[token].find_routes(remaining, stack)...,
		)
	}

	return matches
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
func (self *Trie) FindRoutes(path string) []*Match {
	return self.Root.find_routes(path, &pstack{})
}

// Reduce the size of the tree, must be done after the last AddRoute.
func (self *Trie) Compress() {
	self.Root.compress()
}
