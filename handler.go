package train

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
    "fmt"
)

var assetServer *http.Handler
var publicAssetServer *http.Handler

// func Handler(w http.ResponseWriter, r *http.Request) {
//     setupFileServer()
// 
//     if IsInProduction() {
//         fmt.Println("In Production\n")
//         servePublicAssets(w, r)
//     } else {
//         serveAssets(w, r)
//     }
//     // return
// }

func Run() {
	setupFileServer()

    var server func(w http.ResponseWriter, r *http.Request)
	if Config.BundleAssets {
        fmt.Println("Mode[Production]")
		server = servePublicAssets
	} else {
        fmt.Println("Mode[Development]")
		server = serveAssets
	}
    
    http.Handle(Config.AssetsUrl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        server(w, r)
    }))
}

func IsInProduction() bool {
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
			if strings.Contains(err.Error(), "Could not compile") {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}

			io.Copy(w, strings.NewReader(err.Error()))
			log.Printf("Failed to deliver asset\nGET %s\n-----------------------\n%s\n", url, err.Error())
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
