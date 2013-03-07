package train

import (
	"fmt"
	"net/http"
)

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
    
	// TODO: support using custom ServeMux
    http.Handle(Config.AssetsUrl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        server(w, r)
    }))
}

func Stop() {
	stopConnectInterpreter()
}