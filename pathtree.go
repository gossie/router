package router

import (
	"errors"
	"log"
)

const (
	NodeTypeRoot = iota
	NodeTypeStatic
	NodeTypeVar
)

type node struct {
	nodeType    int
	pathElement string
	route       *route
	children    []*node
	middleware  []Middleware
}

func contains(nodes []*node, el string) (*node, bool) {
	for _, n := range nodes {
		if n.pathElement == el {
			return n, true
		}
	}
	return nil, false
}

func (n *node) createOrGetStaticChild(el string) (*node, error) {
	if n.children == nil {
		n.children = make([]*node, 0, 1)
	}

	foundVariable := false
	for _, child := range n.children {
		if child.nodeType == NodeTypeVar {
			foundVariable = true
		}
	}

	if foundVariable {
		return nil, errors.New("a static path element cannot be added, if there is already a path variable at that position")
	}

	pathElement := el
	if child, found := contains(n.children, pathElement); found && child.nodeType == NodeTypeStatic && child.pathElement == pathElement {
		log.Default().Println("found static path element", pathElement)
		return child, nil
	}

	log.Default().Println("creating static path element", pathElement)
	newNode := &node{nodeType: NodeTypeStatic, pathElement: pathElement}
	n.children = append(n.children, newNode)
	return newNode, nil
}

func (n *node) createOrGetVarChild(el string) (*node, error) {
	if n.children == nil {
		n.children = make([]*node, 0, 1)
	}

	if child, found := contains(n.children, el); found && child.nodeType == NodeTypeVar && child.pathElement == el {
		log.Default().Println("found variable path element", el)
		return child, nil
	}

	if len(n.children) != 0 {
		return nil, errors.New("a path variable cannot be added as a path element that is already present")
	}

	log.Default().Println("creating variable path element", el)
	newNode := &node{nodeType: NodeTypeVar, pathElement: el}
	n.children = append(n.children, newNode)
	return newNode, nil
}

func (n *node) childNode(el string) *node {
	if len(n.children) == 1 && n.children[0].nodeType == NodeTypeVar {
		return n.children[0]
	}

	if child, found := contains(n.children, el); found {
		return child
	}

	log.Default().Println("could not find node for path element", el)
	return nil
}

type pathTree struct {
	root *node
}

func newPathTree() *pathTree {
	return &pathTree{&node{nodeType: NodeTypeRoot}}
}
