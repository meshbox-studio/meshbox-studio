package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/meshbox-studio/meshbox-studio/assets"
	"github.com/meshbox-studio/meshbox-studio/ui/pages"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/templui/templui/utils"
)

func main() {
	initDotEnv()

	mux := http.NewServeMux()
	setupAssetsRoutes(mux)
	setupPageRoutes(mux)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func setupPageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/projects", http.StatusMovedPermanently)
	})
	mux.Handle("GET /create", templ.Handler(pages.Create()))
	mux.Handle("GET /projects", templ.Handler(pages.Projects()))
	mux.Handle("GET /archive", templ.Handler(pages.Archive()))
	mux.Handle("GET /trash", templ.Handler(pages.Trash()))
}

func initDotEnv() {
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Error loading .env file:", err)
	}
}

func setupAssetsRoutes(mux *http.ServeMux) {
	isDevelopment := os.Getenv("GO_ENV") == "development"

	if isDevelopment {
		fmt.Println("Server is running in development mode")
	}

	// Your app assets (CSS, fonts, images, ...)
	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))

	// templUI embedded component scripts
	utils.SetupScriptRoutes(mux, isDevelopment)
}
