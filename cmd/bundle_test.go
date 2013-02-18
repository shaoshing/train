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

	compressAssets()
	assertEqual("public/assets/javascripts/require.js", `normal.js;sub/normal.js;sub/require.js;require.js;`)
	assertEqual("public/assets/javascripts/require-min.js", `Please
Do
Not
Compresee
Me
`)

}
