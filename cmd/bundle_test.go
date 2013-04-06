package main

import (
	"github.com/shaoshing/gotest"
	"github.com/shaoshing/train"
	"io/ioutil"
	"testing"
)

func init() {
	train.Config.Verbose = true
}

func assertEqual(path, content string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	assert.Equal(content, string(c))
}

func TestCommand(t *testing.T) {
	assert.Test = t

	assert.TrueM(prepareEnv(), "Unable to prepare env for cmd tests")

	removeAssets()

	copyAssets()

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

	assertEqual("public/assets/stylesheets/scss.css", `h2 {
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

	fingerPrintAssets()
	assertEqual("public/assets/stylesheets/font.css", `h1{color:green}`) // should keep original assets
	train.LoadManifestInfo()
	assertEqual("public"+train.ManifestInfo["/assets/stylesheets/font.css"], `h1{color:green}`)

	removeAssets()
}
