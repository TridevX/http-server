package http_server

import (
	"net/http"

	"github.com/tridevx/http-server/router"
)

// App represents an Express-like HTTP server.
type App struct {
	Middleware []func(http.Handler) http.Handler
	Router     *router.Router // Add a router to the App
}

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

// GetHandler returns the router as an http.Handler
func (app *App) GetHandler() http.Handler {
	var handler http.Handler = http.DefaultServeMux

	// Apply middleware in reverse order.
	for i := len(app.Middleware) - 1; i >= 0; i-- {
		handler = app.Middleware[i](handler)
	}

	// Add the router as the final handler
	handler = app.Router

	return handler
}

func (app *App) AttachRouter(customRouter *router.Router) {
	app.Router = customRouter
}

// Listen starts the HTTP server on the specified address and port.
func (app *App) Listen(addr string) error {

	return http.ListenAndServe(addr, app.GetHandler())
}
