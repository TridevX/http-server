package router

import (
	"fmt"
	"net/http"
)

type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Get(pattern string, handler http.HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *Router) Post(pattern string, handler http.HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

func (r *Router) Put(pattern string, handler http.HandlerFunc) {
	r.addRoute("PUT", pattern, handler)
}

func (r *Router) Patch(pattern string, handler http.HandlerFunc) {
	r.addRoute("PATCH", pattern, handler)
}

func (r *Router) Delete(pattern string, handler http.HandlerFunc) {
	r.addRoute("DELETE", pattern, handler)
}

func (r *Router) Resource(pattern string, controller string) {
	r.Get(pattern, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "GET %s - %s\n", pattern, controller)
	})

	r.Post(pattern, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "POST %s - %s\n", pattern, controller)
	})

	r.Put(pattern, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "PUT %s - %s\n", pattern, controller)
	})

	r.Patch(pattern, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "PATCH %s - %s\n", pattern, controller)
	})

	r.Delete(pattern, func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "DELETE %s - %s\n", pattern, controller)
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Method == req.Method && route.Pattern == req.URL.Path {
			route.Handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func (r *Router) addRoute(method, pattern string, handler http.HandlerFunc) {
	r.routes = append(r.routes, Route{
		Pattern: pattern,
		Method:  method,
		Handler: handler,
	})
}