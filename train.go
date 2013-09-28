package train

import (
	"fmt"
	"net/http"
)

var server func(w http.ResponseWriter, r *http.Request)

func SetFileServer() {
	setupFileServer()

	if IsInProduction() {
		fmt.Println("[Production] Serving assets from ./public/assets")
		server = servePublicAssets
	} else {
		fmt.Println("[Development] Serving assets from ./assets")
		server = serveAssets
	}
}

// setting serveMux to nil will use http.DefaultServeMux instead.
func ConfigureHttpHandler(serveMux *http.ServeMux) {
	if serveMux == nil {
		serveMux = http.DefaultServeMux
	}

	SetFileServer()

	serveMux.Handle(Config.AssetsUrl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeRequest(w, r)
	}))
}

// use train.ServeRequest(w, r) to server a request manually.
func ServeRequest(w http.ResponseWriter, r *http.Request) {
	server(w, r)
}
