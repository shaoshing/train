package train

import (
	"fmt"
	"net/http"
)

func Run() {
	setupFileServer()

    var server func(w http.ResponseWriter, r *http.Request)
	if Config.BundleAssets {
        fmt.Println("[Production]Serving assets from ./public/assets:\n")
		server = servePublicAssets
	} else {
        fmt.Println("[Development]Serving assets from ./assets:\n")
		server = serveAssets
	}
    
	// TODO: support custom ServeMux
    http.Handle(Config.AssetsUrl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        server(w, r)
    }))
}

func Stop() {
	stopConnectInterpreter()
}