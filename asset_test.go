package train

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestReadingNormalAssets(t *testing.T) {
	var content string
	var err error

	content, _ = ReadAsset("/assets/javascripts/normal.js")
	assert.Equal(t, "normal.js\n", content)

	content, _ = ReadAsset("/assets/javascripts/sub/normal.js")
	assert.Equal(t, "sub/normal.js\n", content)

	_, err = ReadAsset("/assets/not/exists/normal.js")
	assert.Equal(t, true, err != nil)

	content, _ = ReadAsset("/assets/stylesheets/normal.css")
	assert.Equal(t, "normal.css\n", content)

	content, _ = ReadAsset("/assets/stylesheets/sub/normal.css")
	assert.Equal(t, "sub/normal.css\n", content)

	_, err = ReadAsset("/assets/not/exists/normal.css")
	assert.Equal(t, true, err != nil)

	_, err = ReadAsset("/assets/static.txt")
	assert.Equal(t, true, err != nil)
}

func TestReadingAssetsWithRequire(t *testing.T) {
	Config.BundleAssets = true
	var content string

	content, _ = ReadAsset("/assets/javascripts/require.js")
	assert.Equal(t, `normal.js

sub/normal.js

sub/require.js

require.js
`, content)

	content, _ = ReadAsset("/assets/stylesheets/require.css")
	assert.Equal(t, `normal.css

sub/normal.css

sub/require.css

require.css
`, content)

	Config.BundleAssets = false
	content, _ = ReadAsset("/assets/stylesheets/require.css")
	assert.Equal(t, `/*
 *= require stylesheets/normal
 *= require stylesheets/sub/require
 */
require.css
`, content)
}
