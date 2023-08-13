package framework

import (
	stdPath "path"
	"strings"
)

type HandlerFunc func(Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Router struct {
	routes map[string]HandlerFunc
	app    *Application
}

type Group struct {
	router     Router
	path       string
	middleware []MiddlewareFunc
}

func (g *Group) Add(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	path = stdPath.Join(g.path, path)
	middleware = append(g.middleware, middleware...)
	g.router.Add(path, handler, middleware...)
}

func (g *Group) Group(path string, middleware ...MiddlewareFunc) *Group {
	return &Group{
		router:     g.router,
		path:       stdPath.Join(g.path, path),
		middleware: append(g.middleware, middleware...),
	}
}

// NewRouter returns a new Router instance.
func NewRouter(app *Application) *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
		app:    app,
	}
}

func (r Router) Group(path string, middleware ...MiddlewareFunc) *Group {
	return &Group{
		router:     r,
		path:       path,
		middleware: middleware,
	}
}

func (r *Router) Add(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	path = cleanPath(path)
	for _, m := range middleware {
		handler = m(handler)
	}
	r.routes[path] = handler
}

func (r *Router) Find(path string) (HandlerFunc, bool) {
	path = cleanPath(path)
	handle, ok := r.routes[path]
	if !ok {
		path = strings.TrimLeft(path, "/")
		handle, ok = r.routes[path]
	}
	return handle, ok
}
