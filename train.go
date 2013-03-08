package train

import (
	"fmt"
	"github.com/shaoshing/train/interpreter"
	"net/http"
)

// setting serveMux to nil will use http.DefaultServeMux instead.
func Run(serveMux *http.ServeMux) {
	setupFileServer()

	if serveMux == nil {
		serveMux = http.DefaultServeMux
	}

	var server func(w http.ResponseWriter, r *http.Request)
	if Config.BundleAssets {
		fmt.Println("[Production]Serving assets from ./public/assets:\n")
		server = servePublicAssets
	} else {
		fmt.Println("[Development]Serving assets from ./assets:\n")
		server = serveAssets
	}

	serveMux.Handle(Config.AssetsUrl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server(w, r)
	}))
}

func Stop() {
	interpreter.CloseInterpreter()
}
