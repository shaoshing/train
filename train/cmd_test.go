package main

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"testing"
)

func assertEqual(t *testing.T, path, content string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, content, string(c))
}

func TestCommand(t *testing.T) {
	copyAssets()
	defer removeAssets()

	assertEqual(t, "public/assets/javascripts/normal.js", "normal.js\n")
	assertEqual(t, "public/assets/javascripts/require.js", `//= require javascripts/normal
//= require javascripts/sub/require
require.js
`)
	assertEqual(t, "public/assets/stylesheets/require.css", `/*
 *= require stylesheets/normal
 *= require stylesheets/sub/require
 */
require.css
`)

	bundleAssets()
	assertEqual(t, "public/assets/javascripts/normal.js", "normal.js\n")
	assertEqual(t, "public/assets/javascripts/require.js", `normal.js

sub/normal.js

sub/require.js

require.js
`)
	assertEqual(t, "public/assets/stylesheets/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`)

	compressAssets()
	assertEqual(t, "public/assets/javascripts/require.js", `normal.js;sub/normal.js;sub/require.js;require.js;`)

}
