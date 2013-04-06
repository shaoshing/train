package train

import (
	"github.com/shaoshing/gotest"
	"testing"
)

func init() {
	Config.Verbose = true
}

func TestReadingNormalAssets(t *testing.T) {
	assert.Test = t
	var content string
	var err error

	content, _ = ReadAsset("/assets/javascripts/normal.js")
	assert.Equal("normal.js\n", content)

	content, _ = ReadAsset("/assets/javascripts/sub/normal.js")
	assert.Equal("sub/normal.js\n", content)

	_, err = ReadAsset("/assets/not/exists/normal.js")
	assert.Equal("Asset Not Found: /assets/not/exists/normal.js", err.Error())

	content, _ = ReadAsset("/assets/stylesheets/normal.css")
	assert.Equal("normal.css\n", content)

	content, _ = ReadAsset("/assets/stylesheets/sub/normal.css")
	assert.Equal("sub/normal.css\n", content)

	_, err = ReadAsset("/assets/not/exists/normal.css")
	assert.Equal("Asset Not Found: /assets/not/exists/normal.css", err.Error())

	_, err = ReadAsset("/assets/static.txt")
	assert.Equal("Unsupported Asset: /assets/static.txt", err.Error())
}

func TestReadingSass(t *testing.T) {
	assert.Test = t
	content, err := ReadAsset("/assets/stylesheets/app.css")
	if err != nil {
		panic(err)
	}
	assert.Contain("h1", content)
	assert.Contain("h2", content)

	content, err = ReadAsset("/assets/stylesheets/app2.css")
	if err != nil {
		panic(err)
	}
	assert.Contain("h2", content)
	assert.Contain("h3", content)
}

func TestReadingCoffee(t *testing.T) {
	assert.Test = t
	content, err := ReadAsset("/assets/javascripts/app.js")
	if err != nil {
		panic(err)
	}
	assert.Contain("square", content)
}

func TestRequireDirective(t *testing.T) {
	assert.Test = t
	Config.BundleAssets = true
	defer func() {
		Config.BundleAssets = false
	}()

	content, err := ReadAsset("/assets/javascripts/require2.js")
	if err != nil {
		panic(err)
	}
	assert.Equal(`normal.js

sub/normal.js

`, content)

	content, err = ReadAsset("/assets/stylesheets/require2.css")
	if err != nil {
		panic(err)
	}
	assert.Equal(`normal.css

sub/normal.css

`, content)
}

func TestReadingAssetsWithRequire(t *testing.T) {
	assert.Test = t
	Config.BundleAssets = true
	var content string
	var err error

	content, _ = ReadAsset("/assets/javascripts/require.js")
	assert.Equal(`normal.js

sub/normal.js

sub/require.js

require.js
`, content)

	content, _ = ReadAsset("/assets/stylesheets/require.css")
	assert.Equal(`normal.css

sub/normal.css

sub/require.css

require.css
`, content)

	_, err = ReadAsset("/assets/javascripts/error.js")
	assert.Equal(`Asset Not Found: not/found.js
--- required by /assets/javascripts/error.js`, err.Error())

	_, err = ReadAsset("/assets/javascripts/errors.js")
	assert.Equal(`Asset Not Found: not/found.js
--- required by javascripts/error.js
--- required by /assets/javascripts/errors.js`, err.Error())

	Config.BundleAssets = false
	content, _ = ReadAsset("/assets/stylesheets/require.css")
	assert.Equal(`/*
 *= require stylesheets/normal
 *= require stylesheets/sub/require
 */
require.css
`, content)
}
