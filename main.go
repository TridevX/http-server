package http_server

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/tridevx/http-server/router"
)

// App represents an Express-like HTTP server.
type App struct {
	Middleware []func(http.Handler) http.Handler
	Router     *router.Router // Add a router to the App
}

type Handler func(r *http.Request) (statusCode int, data map[string]interface{})

// NewApp creates a new instance of the App.
func HttpServer() *App {
	return &App{
		Middleware: make([]func(http.Handler) http.Handler, 0),
		Router:     router.NewRouter(), // Initialize the router here
	}
}

// Use adds middleware to the application.
func (app *App) Use(middleware func(http.Handler) http.Handler) {
	app.Middleware = append(app.Middleware, middleware)
}

func (app *App) getHandler(method, path string) http.Handler {
	for _, route := range app.Router.Routes {
		re := regexp.MustCompile(route.Pattern)
		if route.Method == method && re.MatchString(path) {
			return route.Handler
		}
	}
	return http.NotFoundHandler()
}

// GetHandler returns the router as an http.Handler
func (app *App) GetHandler() http.Handler {
	var handler http.Handler = http.DefaultServeMux

	// Apply middleware in reverse order.
	for i := len(app.Middleware) - 1; i >= 0; i-- {
		handler = app.Middleware[i](handler)
	}

	// Add the router as the final handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.getHandler(r.Method, r.URL.Path).ServeHTTP(w, r)
	})
}

func (app *App) AttachRouter(customRouter *router.Router) {
	app.Router = customRouter
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	statusCode, data := h(r)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Listen starts the HTTP server on the specified address and port.
func (app *App) Listen(addr string) error {
	return http.ListenAndServe(addr, app.GetHandler())
}
