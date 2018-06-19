package minimalrouter

import "strings"

func New() *Router {
	return &Router{children: make(map[string]*Router)}
}

// Represents both the Router and each node
type Router struct {
	children map[string]*Router
	varName  string
	handler  interface{}
}

func (this *Router) Add(method, path string, h interface{}) {
	if path[0] != '/' {
		panic("path must begin with '/' in: " + path)
	} else if path[len(path)-1] == '/' {
		panic("path must not end with '/' in: " + path)
	}

	parts := append([]string{method}, strings.Split(path[1:], "/")...)
	node := this

	for _, p := range parts {
		if p[0] == ':' && len(p) > 2 {
			varName := p[1:]
			p = ":"

			if node.varName == "" {
				node.varName = varName
			} else if node.varName != varName {
				panic("conflict while adding " + path + "; parameter conflicts with :" + node.varName)
			}
		}

		n, ok := node.children[p]
		if !ok {
			node.children[p] = New()
			n = node.children[p]
		}

		node = n
	}

	node.handler = h
}

func (this *Router) Match(method, path string) (interface{}, map[string]string) {
	parts := append([]string{method}, strings.Split(path[1:], "/")...)
	params := make(map[string]string)
	node := this

	for _, p := range parts {
		n, ok := node.children[p]

		if !ok && node.varName != "" {
			n = node.children[":"]
			params[node.varName] = p
		} else if !ok {
			return nil, nil
		}

		node = n
	}

	return node.handler, params
}
