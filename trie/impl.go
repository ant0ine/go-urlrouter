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

func splitParam(remaining string) (string, string) {
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

func (self *Node) addRoute(path string, route interface{}) error {

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
	var nextNode *Node

	if token[0] == ':' {
		// :param case
		var name string
		name, remaining = splitParam(remaining)
		if self.ParamChild == nil {
			self.ParamChild = &Node{}
			self.ParamName = name
		}
		nextNode = self.ParamChild
	} else if token[0] == '*' {
		// *splat case
		name := remaining
		remaining = ""
		if self.SplatChild == nil {
			self.SplatChild = &Node{}
			self.SplatName = name
		}
		nextNode = self.SplatChild
	} else {
		// general case
		if self.Children == nil {
			self.Children = map[string]*Node{}
			self.ChildrenKeyLen = 1
		}
		if self.Children[token] == nil {
			self.Children[token] = &Node{}
		}
		nextNode = self.Children[token]
	}

	return nextNode.addRoute(remaining, route)
}

// utility for the Node.findRoute recursive method
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

func (self *pstack) asMap() map[string]string {
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

func (self *Node) findRoutes(path string, stack *pstack) []*Match {

	matches := []*Match{}

	// route found !
	if self.Route != nil && path == "" {
		matches = append(
			matches,
			&Match{
				Route:  self.Route,
				Params: stack.asMap(),
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
			self.SplatChild.findRoutes("", stack)...,
		)
		stack.pop()
	}

	// :param branch
	if self.ParamChild != nil {
		value, remaining := splitParam(path)
		stack.add(self.ParamName, value)
		matches = append(
			matches,
			self.ParamChild.findRoutes(remaining, stack)...,
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
			self.Children[token].findRoutes(remaining, stack)...,
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
	canCompress := true
	for _, node := range self.Children {
		if node.Route != nil || node.SplatChild != nil || node.ParamChild != nil {
			canCompress = false
		}
	}
	// compress
	if canCompress {
		merged := map[string]*Node{}
		for key, node := range self.Children {
			for gdKey, gdNode := range node.Children {
				mergedKey := key + gdKey
				merged[mergedKey] = gdNode
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
	return self.Root.addRoute(path, route)
}

// Given a path, return all the matchin routes.
func (self *Trie) FindRoutes(path string) []*Match {
	return self.Root.findRoutes(path, &pstack{})
}

// Reduce the size of the tree, must be done after the last AddRoute.
func (self *Trie) Compress() {
	self.Root.compress()
}
