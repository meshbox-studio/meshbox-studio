package webui

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func NewHandler(isDevelopment bool) http.Handler {
	if isDevelopment {
		return newDevelopmentHandler()
	}
	return newProductionHandler()
}

func newProductionHandler() http.Handler {
	staticFS, err := fs.Sub(distFS, "dist/app")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "frontend assets are missing: run `npm --prefix frontend run build`", http.StatusServiceUnavailable)
		})
	}

	fileServer := http.FileServer(http.FS(staticFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := normalizePath(r.URL.Path)
		if requestPath == "" {
			requestPath = "index.html"
		}

		if shouldServeIndex(requestPath) {
			serveEmbeddedIndex(w, r, staticFS)
			return
		}

		if fileExists(staticFS, requestPath) {
			fileServer.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(requestPath, "api/") {
			http.NotFound(w, r)
			return
		}

		if path.Ext(requestPath) != "" {
			http.NotFound(w, r)
			return
		}

		serveEmbeddedIndex(w, r, staticFS)
	})
}

func newDevelopmentHandler() http.Handler {
	localDist := http.Dir("./internal/webui/dist/app")
	fileServer := http.FileServer(localDist)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := normalizePath(r.URL.Path)
		if requestPath == "" {
			requestPath = "index.html"
		}

		if shouldServeIndex(requestPath) {
			serveDiskIndex(w, r, localDist)
			return
		}

		if diskFileExists(localDist, requestPath) {
			w.Header().Set("Cache-Control", "no-store")
			fileServer.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(requestPath, "api/") {
			http.NotFound(w, r)
			return
		}

		if path.Ext(requestPath) != "" {
			http.NotFound(w, r)
			return
		}

		serveDiskIndex(w, r, localDist)
	})
}

func normalizePath(rawPath string) string {
	trimmed := strings.TrimPrefix(path.Clean(rawPath), "/")
	if trimmed == "." {
		return ""
	}
	return trimmed
}

func shouldServeIndex(requestPath string) bool {
	return requestPath == "index.html"
}

func serveEmbeddedIndex(w http.ResponseWriter, r *http.Request, staticFS fs.FS) {
	w.Header().Set("Cache-Control", "no-store")
	index, err := staticFS.Open("index.html")
	if err != nil {
		http.Error(w, "frontend assets are missing: run `npm --prefix frontend run build`", http.StatusServiceUnavailable)
		return
	}
	defer index.Close()

	stat, err := index.Stat()
	if err != nil {
		http.Error(w, "failed to read frontend metadata", http.StatusInternalServerError)
		return
	}

	content, err := fs.ReadFile(staticFS, "index.html")
	if err != nil {
		http.Error(w, "failed to read frontend index", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, "index.html", fileModTime(stat), bytes.NewReader(content))
}

func serveDiskIndex(w http.ResponseWriter, r *http.Request, localDist http.Dir) {
	w.Header().Set("Cache-Control", "no-store")
	indexPath := string(localDist) + "/index.html"
	if _, err := os.Stat(indexPath); err != nil {
		http.Error(w, "frontend assets are missing: run `npm --prefix frontend run build`", http.StatusServiceUnavailable)
		return
	}
	http.ServeFile(w, r, indexPath)
}

func fileExists(staticFS fs.FS, name string) bool {
	f, err := staticFS.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return !stat.IsDir()
}

func diskFileExists(localDist http.Dir, name string) bool {
	fullPath := string(localDist) + "/" + name
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func fileModTime(info fs.FileInfo) time.Time {
	if info == nil {
		return time.Time{}
	}
	return info.ModTime()
}
