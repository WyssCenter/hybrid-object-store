package renderer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

// RenderConfig stores variables pulled from the environment
type RenderConfig struct {
	StaticDir string
	ConfigDir string
}

// RenderServer is the heart of this package. It routes requests to functions and does logging.
type RenderServer struct {
	router    *httprouter.Router // HTTP request routing
	staticDir string             // Location of general file resources (js, images)
	configDir string
}

// ServeHTTP fulfill's RenderServer's obligation to the Handler interface.
func (a *RenderServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *RenderServer) singleFile(filename string) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			http.Error(w, "404 - Not found.", http.StatusNotFound)
		} else {
			fmt.Fprint(w, string(fileBytes))
		}
	}
}

// New creates a root handler for the server.
func New(conf *RenderConfig) http.Handler {
	router := httprouter.New()

	app := &RenderServer{
		router,
		conf.StaticDir,
		conf.ConfigDir,
	}

	router.GET("/", app.singleFile(filepath.Join(conf.StaticDir, "index.html")))
	router.GET("/favicon.png", app.singleFile(filepath.Join(conf.ConfigDir, "favicon.png")))
	router.GET("/logo.svg", app.singleFile(filepath.Join(conf.ConfigDir, "logo.svg")))
	router.GET("/config.json", app.singleFile(filepath.Join(conf.ConfigDir, "config.json")))
	router.GET("/manifest.json", app.singleFile(filepath.Join(conf.StaticDir, "manifest.json")))
	router.ServeFiles("/static/css/*filepath", http.Dir(filepath.Join(conf.StaticDir, "static/css")))
	router.ServeFiles("/static/js/*filepath", http.Dir(filepath.Join(conf.StaticDir, "static/js")))
	router.ServeFiles("/static/media/*filepath", http.Dir(filepath.Join(conf.StaticDir, "static/media")))

	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fileBytes, err := ioutil.ReadFile(filepath.Join(conf.StaticDir, "index.html"))
		if err != nil {
			http.Error(rw, "404 - Not found.", http.StatusNotFound)
		} else {
			fmt.Fprint(rw, string(fileBytes))
		}
	})

	return app
}
