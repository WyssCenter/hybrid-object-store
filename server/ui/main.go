package main

import (
	"net/http"

	"github.com/gigantum/hoss-ui/renderer"
	"github.com/rs/cors"
)

func main() {
	webStaticDir := "ui/build/"
	configDir := "/opt/config/"

	renderConfig := &renderer.RenderConfig{
		StaticDir: webStaticDir,
		ConfigDir: configDir,
	}

	r := renderer.New(renderConfig)

	// Allow CORS for frontend devs.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: false,
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"*"},
		Debug:            false,
	})

	http.ListenAndServe(":8080", c.Handler(r))
}
