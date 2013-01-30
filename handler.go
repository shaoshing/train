package train

import (
	"io"
	"net/http"
	"path"
	"strings"
)

var contentTypes = map[string]string{
	".js":  "application/javascript",
	".css": "text/css",
}

var assetServer *http.Handler

func Handler(w http.ResponseWriter, r *http.Request) {
	setupFileServer()

	url := r.URL.Path
	ext := path.Ext(url)

	switch ext {
	case ".js", ".css":
		w.Header().Set("Content-Type", contentTypes[ext])
		content := ReadAsset(url)
		reader := strings.NewReader(content)
		io.Copy(w, reader)
	default:
		(*assetServer).ServeHTTP(w, r)
	}
	return
}

func setupFileServer() {
	if assetServer == nil {
		server := http.FileServer(http.Dir(Config.AssetsPath + "/../"))
		assetServer = &server
	}
}
