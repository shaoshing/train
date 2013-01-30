package train

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

var assetServer *http.Handler
var publicAssetServer *http.Handler

func Handler(w http.ResponseWriter, r *http.Request) {
	setupFileServer()

	if hasPublicAssets() {
		servePublicAssets(w, r)
	} else {
		serveAssets(w, r)
	}
	return
}

func hasPublicAssets() bool {
	_, err := os.Stat("public" + Config.AssetsUrl)
	return err == nil
}

func servePublicAssets(w http.ResponseWriter, r *http.Request) {
	(*publicAssetServer).ServeHTTP(w, r)
}

var contentTypes = map[string]string{
	".js":  "application/javascript",
	".css": "text/css",
}

func serveAssets(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	ext := path.Ext(url)

	switch ext {
	case ".js", ".css":
		content, err := ReadAsset(url)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			io.Copy(w, strings.NewReader(err.Error()))
		} else {
			w.Header().Set("Content-Type", contentTypes[ext])
			io.Copy(w, strings.NewReader(content))
		}
	default:
		(*assetServer).ServeHTTP(w, r)
	}
}

func setupFileServer() {
	if assetServer == nil {
		server := http.FileServer(http.Dir(Config.AssetsPath + "/.."))
		assetServer = &server
	}
	if publicAssetServer == nil {
		server := http.FileServer(http.Dir("public"))
		publicAssetServer = &server
	}
}
