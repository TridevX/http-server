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

func (r *Router) Resource(pattern string, controller interface{}, methodName string) {
	key := fmt.Sprintf("%s-%s", pattern, methodName)
	r.routes[key] = func(w http.ResponseWriter, req *http.Request) {
		callControllerMethod(controller, methodName, w, req)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := fmt.Sprintf("%s-%s", req.URL.Path, req.Method)
	if handler, ok := r.routes[key]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) addRoute(method, pattern string, handler http.HandlerFunc) {
	r.routes[fmt.Sprintf("%s-%s", pattern, method)] = handler
}

func callControllerMethod(controller interface{}, methodName string, w http.ResponseWriter, req *http.Request) {
	switch c := controller.(type) {
	case func(http.ResponseWriter, *http.Request):
		c(w, req)
	case http.HandlerFunc:
		c(w, req)
	default:
		fmt.Fprintf(w, "Invalid controller method: %s", methodName)
	}
}