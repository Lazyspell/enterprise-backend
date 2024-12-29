// Pcakge web contains a small web framework contains.
package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// A Handler is a type that handles a http request within our own little mini
// framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// APP is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*http.ServeMux
	shutdown chan os.Signal
}

// NewApp creates an App value that hadles a set of routes for the application
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
	}
}

// HandleFunc sets a hanlder function for a given HTTP method and path pair
// to the application server mux.
func (a *App) HandleFunc(pattern string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Println("error happened")
			return
		}
	}

	a.ServeMux.HandleFunc(pattern, h)
}
