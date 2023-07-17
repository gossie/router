package router

import (
	"errors"
	"log"
	"strings"
)

type node struct {
	nodeType            string
	pathElement         string
	route               *route
	children            map[string]*node
	middlewareFunctions []Middleware
}

func (n *node) createOrGetStaticChild(el string) (*node, error) {
	foundVariable := false
	for _, child := range n.children {
		if child.nodeType == "var" {
			foundVariable = true
		}
	}

	if foundVariable {
		return nil, errors.New("a static path element cannot be added, if there is already a path variable at that position")
	}

	pathElement := strings.ToLower(el)
	if child, found := n.children[pathElement]; found && child.nodeType == "static" && child.pathElement == pathElement {
		log.Default().Println("found static path element", pathElement)
		return child, nil
	}

	log.Default().Println("creating static path element", pathElement)
	newNode := &node{"static", pathElement, nil, make(map[string]*node), make([]Middleware, 0)}
	n.children[pathElement] = newNode
	return newNode, nil
}

func (n *node) createOrGetVarChild(el string) (*node, error) {
	if child, found := n.children[el]; found && child.nodeType == "var" && child.pathElement == el {
		log.Default().Println("found variable path element", el)
		return child, nil
	}

	if len(n.children) != 0 {
		return nil, errors.New("a path variable cannot be added as a path element that is already present")
	}

	log.Default().Println("creating variable path element", el)
	newNode := &node{"var", el, nil, make(map[string]*node), make([]Middleware, 0)}
	n.children[el] = newNode
	return newNode, nil
}

func (n *node) childNode(el string) *node {
	if len(n.children) == 1 {
		for _, child := range n.children {
			if child.nodeType == "var" {
				return child
			}
		}
	}

	if child, found := n.children[strings.ToLower(el)]; found {
		return child
	}

	log.Default().Println("could not find node for path element", el)
	return nil
}

type pathTree struct {
	root *node
}

func newPathTree() *pathTree {
	return &pathTree{&node{"root", "", nil, make(map[string]*node), make([]Middleware, 0)}}
}
