package main

import (
	"fmt"
	"github.com/shaoshing/train"
	"html/template"
	"net/http"
)

func main() {
	train.ConfigureHttpHandler(nil)

	http.HandleFunc("/", example)
	http.HandleFunc("/toggle_bundle_assets", toggle_bundle_assets)

	fmt.Println("Listening to localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}

func example(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tpl := template.New("example")
	tpl.Funcs(train.HelperFuncs)
	tpl.ParseFiles("example/index.html")
	tpl.ExecuteTemplate(w, "example", struct{ BundleAssets bool }{train.Config.BundleAssets})
}

func toggle_bundle_assets(w http.ResponseWriter, r *http.Request) {
	train.Config.BundleAssets = !train.Config.BundleAssets
	http.Redirect(w, r, "/", 302)
}
