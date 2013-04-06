package train

import (
	"fmt"
	"net/http"
)

// setting serveMux to nil will use http.DefaultServeMux instead.
func ConfigureHttpHandler(serveMux *http.ServeMux) {
	setupFileServer()

	if serveMux == nil {
		serveMux = http.DefaultServeMux
	}

	var server func(w http.ResponseWriter, r *http.Request)
	if IsInProduction() {
		fmt.Println("[Production] Serving assets from ./public/assets")
		server = servePublicAssets
	} else {
		fmt.Println("[Development] Serving assets from ./assets")
		server = serveAssets
	}

	serveMux.Handle(Config.AssetsUrl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server(w, r)
	}))
}
