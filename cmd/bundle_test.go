package main

import (
	"github.com/shaoshing/gotest"
	"io/ioutil"
	"testing"
)

func assertEqual(path, content string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	assert.Equal(content, string(c))
}

func TestCommand(t *testing.T) {
	assert.Test = t
	copyAssets()
	defer removeAssets()

	assertEqual("public/assets/javascripts/normal.js", "normal.js\n")
	assertEqual("public/assets/javascripts/require.js", `//= require javascripts/normal
//= require javascripts/sub/require
require.js
`)
	assertEqual("public/assets/stylesheets/require.css", `/*
 *= require stylesheets/normal
 *= require stylesheets/sub/require
 */
require.css
`)

	bundleAssets()
	assertEqual("public/assets/javascripts/normal.js", "normal.js\n")
	assertEqual("public/assets/javascripts/require.js", `normal.js

sub/normal.js

sub/require.js

require.js
`)
	assertEqual("public/assets/stylesheets/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`)

	assertEqual("public/assets/stylesheets/font.css", `h1 {
  color: green; }
`)
	assertEqual("public/assets/stylesheets/app.css", `h1 {
  color: green; }

h2 {
  color: green; }
`)
	assertEqual("public/assets/javascripts/app.js", `(function() {
  var a;

  a = 12;

}).call(this);
`)

	compressAssets()
	assertEqual("public/assets/javascripts/require.js", `normal.js;sub/normal.js;sub/require.js;require.js;`)
	assertEqual("public/assets/javascripts/require-min.js", `Please
Do
Not
Compresee
Me
`)
	assertEqual("public/assets/stylesheets/font.css", `h1{color:green}`)
	assertEqual("public/assets/javascripts/app.js", `(function(){var b;b=12}).call(this);`)
}
