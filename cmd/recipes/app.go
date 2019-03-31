package recipes

import (
	"log"
	"net/http"
	"os"

	"github.com/peterjasc/recipes/cmd/config"
)

// App holds the mux used to add the handler and server for internal usage
type App struct {
	Mux    *http.ServeMux
	server *http.Server
}

// NewApp setups up http handler to be used by the caller and http server
// on port specified by PORT environment variable or config.HTTPServerDefaultPort
func NewApp() (*App, error) {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":" + config.HTTPServerDefaultPort
	}

	mux := http.NewServeMux()

	app := &App{
		Mux:    mux,
		server: newHTTPServer(mux, addr),
	}

	return app, nil
}

func newHTTPServer(handler http.Handler, addr string) *http.Server {
	httpServer := &http.Server{
		Handler:      handler,
		Addr:         addr,
		ReadTimeout:  config.HTTPServerReadTimeout,
		WriteTimeout: config.HTTPServerWriteTimeout,
	}

	return httpServer
}

// ListenAndServe calls the namesake function for the internal http server object
func (app *App) ListenAndServe() {
	if err := app.server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
