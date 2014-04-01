package yolo

import (
	"github.com/robfig/pathtree"
	"net/http"
)

// Add routes to router here
type router struct {
	routes *pathtree.Node
}

type Handler interface{}

func newRouter() *router {
	return &router{
		pathtree.New(),
	}
}

func (r *router) addRoute(route string, handlers []Handler) {
	r.routes.Add(route, handlers)
}

// findRoute takes an http.Request and returns the Handlers to call along with
// any values from the url (if any).
func (r *router) findRoute(req *http.Request) ([]Handler, map[string]interface{}) {
	node, expanded := r.routes.Find("/" + req.Method + req.URL.Path)

	if node == nil {
		return nil, nil
	}

	var params = make(map[string]interface{})
	for i, name := range node.Wildcards {
		params[name] = expanded[i]
	}

	return node.Value.([]Handler), params
}
