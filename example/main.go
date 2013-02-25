package main

import (
	"fmt"
	"github.com/shaoshing/train"
	"html/template"
	"net/http"
)

func main() {
	train.Config.BundleAssets = true

	http.Handle(train.Config.AssetsUrl, http.HandlerFunc(train.Handler))
	http.HandleFunc("/", example)

	fmt.Println("Listening to localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func example(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tpl := template.New("example")

	tpl.Funcs(train.HelperFuncs)
	tpl.Parse(`
{{define "example"}}
	Examples:<br/>

	{{javascript_tag "normal"}}
	{{stylesheet_tag "normal"}}

	<a href="/assets/javascripts/normal.js">JS 1</a><br/>
	<a href="/assets/javascripts/require.js">JS 2</a> (with require directive)<br/>
	<a href="/assets/javascripts/app.js">JS 3</a> (compiled from CoffeeScript <a href="/assets/javascripts/app.coffee">source</a>)<br/>
	<a href="/assets/javascripts/app.err.js">JS 4</a> (compiled from CoffeeScript, but failed <a href="/assets/javascripts/app.err.coffee">source</a>)<br/>

	<br/>
	<a href="/assets/stylesheets/normal.css">CSS 1</a><br/>
	<a href="/assets/stylesheets/require.css">CSS 2</a> (with require directive)<br/>
	<a href="/assets/stylesheets/app.css">CSS 3</a> (compiled from SASS <a href="/assets/stylesheets/app.sass">source</a>)<br/>
	<a href="/assets/stylesheets/app2.css">CSS 4</a> (compiled from SCSS <a href="/assets/stylesheets/app2.scss">source</a>)<br/>
	<a href="/assets/stylesheets/app.err.css">CSS 5</a> (compiled from SASS, but failed <a href="/assets/stylesheets/app.err.sass">source</a>)<br/>
{{end}}
`)

	tpl.ExecuteTemplate(w, "example", nil)
}
