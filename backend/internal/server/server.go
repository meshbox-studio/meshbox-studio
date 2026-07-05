package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/meshbox-studio/meshbox-studio/internal/activity"
	"github.com/meshbox-studio/meshbox-studio/internal/auth"
	"github.com/meshbox-studio/meshbox-studio/internal/projects"
	"github.com/meshbox-studio/meshbox-studio/internal/stats"
	"github.com/meshbox-studio/meshbox-studio/internal/webui"
)

func NewAPI(isDevelopment bool, dataDir string) (huma.API, *http.ServeMux) {
	mux := http.NewServeMux()

	mux.Handle("/", webui.NewHandler(isDevelopment))

	config := huma.DefaultConfig("Meshbox Studio API", stats.AppVersion())
	config.OpenAPIPath = "/api/openapi"
	config.DocsPath = "/api/docs"
	config.SchemasPath = "/api/schemas"

	api := humago.New(mux, config)

	stats.RegisterRoutes(api)
	auth.RegisterRoutes(api)
	activity.RegisterRoutes(api)
	projects.RegisterRoutes(api, mux, dataDir)

	return api, mux
}