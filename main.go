package http_server

import (
	"net/http"

	"github.com/tridevx/http-server/router"
)

// App represents an Express-like HTTP server.
type App struct {
	middleware []func(http.Handler) http.Handler
	router     *router.Router // Add a router to the App
}

// NewApp creates a new instance of the App.
func HttpServer() *App {
	return &App{
		middleware: make([]func(http.Handler) http.Handler, 0),
		router:     router.NewRouter(), // Initialize the router here
	}
}

// Use adds middleware to the application.
func (app *App) Use(middleware func(http.Handler) http.Handler) {
	app.middleware = append(app.middleware, middleware)
}

// ServeHTTP implements the http.Handler interface to handle HTTP requests.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.DefaultServeMux

	// Apply middleware in reverse order.
	for i := len(app.middleware) - 1; i >= 0; i-- {
		handler = app.middleware[i](handler)
	}

	handler.ServeHTTP(w, r)
}

func (app *App) AttachRouter(customRouter *router.Router) {
	app.router = customRouter
}

// Listen starts the HTTP server on the specified address and port.
func (app *App) Listen(addr string) error {
	return http.ListenAndServe(addr, app)
}
