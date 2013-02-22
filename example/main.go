package main

import (
	"fmt"
	"github.com/shaoshing/train"
	"net/http"
)

func main() {
	train.Config.BundleAssets = true

	http.Handle(train.Config.AssetsUrl, http.HandlerFunc(train.Handler))
	http.HandleFunc("/", example)

	fmt.Println("Listening to localhost:8000")
	http.ListenAndServe(":8080", nil)
}

func example(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, `
Examples:<br/>

<a href="/assets/javascripts/normal.js">JS</a><br/>
<a href="/assets/javascripts/require.js">JS (with require directive)</a><br/>

<br/>
<a href="/assets/stylesheets/normal.css">CSS</a><br/>
<a href="/assets/stylesheets/require.css">CSS (with require directive)</a><br/>
<a href="/assets/stylesheets/app.css">CSS (compiled from SASS)</a><br/>
<a href="/assets/stylesheets/app.err.css">CSS (compiled from SAA, but failed)</a><br/>
`)
}
