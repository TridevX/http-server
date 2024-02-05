package router

import (
	"fmt"
	"net/http"
	"reflect"
)

type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

type Router struct {
	Routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Get(pattern string, controller interface{}, methodName string) {
	r.addRoute("GET", pattern, controller, methodName)
}

func (r *Router) Post(pattern string, controller interface{}, methodName string) {
	r.addRoute("POST", pattern, controller, methodName)
}

func (r *Router) Put(pattern string, controller interface{}, methodName string) {
	r.addRoute("PUT", pattern, controller, methodName)
}

func (r *Router) Patch(pattern string, controller interface{}, methodName string) {
	r.addRoute("PATCH", pattern, controller, methodName)
}

func (r *Router) Delete(pattern string, controller interface{}, methodName string) {
	r.addRoute("DELETE", pattern, controller, methodName)
}

func (r *Router) Resource(pattern string, controller interface{}) {
	r.Get(pattern, controller, "index")
	r.Post(pattern, controller, "store")
	r.Get(fmt.Sprintf("%s/:id", pattern), controller, "show")
	r.Put(fmt.Sprintf("%s/:id", pattern), controller, "update")
	r.Delete(fmt.Sprintf("%s/:id", pattern), controller, "destroy")
}

func (r *Router) addRoute(method, pattern string, controller interface{}, methodName string) {
	handlerFunc := getHandlerFromController(controller, methodName)
	r.Routes = append(r.Routes, Route{
		Pattern: pattern,
		Method:  method,
		Handler: handlerFunc,
	})
}

func getHandlerFromController(controller interface{}, methodName string) http.HandlerFunc {
	controllerValue := reflect.ValueOf(controller)
	method := controllerValue.MethodByName(methodName)

	if !method.IsValid() || method.Kind() != reflect.Func {
		return func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "Invalid controller method", http.StatusInternalServerError)
		}
	}

	// Ensure the method has the correct signature
	methodType := method.Type()
	if methodType.NumIn() != 2 || methodType.In(0) != reflect.TypeOf((*http.ResponseWriter)(nil)).Elem() || methodType.In(1) != reflect.TypeOf((*http.Request)(nil)).Elem() {
		return func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "Invalid controller method signature", http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, req *http.Request) {
		method.Call([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(req)})
	}
}
