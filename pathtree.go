package router

import (
	"errors"
	"log"
)

type node struct {
	nodeType    string
	pathElement string
	route       *route
	children    map[string]*node
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

	if child, found := n.children[el]; found && child.nodeType == "static" && child.pathElement == el {
		log.Default().Println("found static path element", el)
		return child, nil
	}

	log.Default().Println("creating static path element", el)
	newNode := &node{"static", el, nil, make(map[string]*node)}
	n.children[el] = newNode
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
	newNode := &node{"var", el, nil, make(map[string]*node)}
	n.children[el] = newNode
	return newNode, nil
}

func (n *node) getNode(el string) *node {
	if len(n.children) == 1 {
		for _, child := range n.children {
			if child.nodeType == "var" {
				return child
			}
		}
	}

	if child, found := n.children[el]; found {
		return child
	}

	log.Default().Println("could not find node for path element", el)
	return nil
}

type pathTree struct {
	root *node
}

func createPathTree() *pathTree {
	return &pathTree{&node{"root", "", nil, make(map[string]*node)}}
}
