package train

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestReadingNormalAssets(t *testing.T) {
	var content string

	content = ReadAsset("/assets/javascripts/normal.js")
	assert.Equal(t, "normal.js\n", content)

	content = ReadAsset("/assets/javascripts/sub/normal.js")
	assert.Equal(t, "sub/normal.js\n", content)

	content = ReadAsset("/assets/not/exists/normal.js")
	assert.Equal(t, "", content)

	content = ReadAsset("/assets/stylesheets/normal.css")
	assert.Equal(t, "normal.css\n", content)

	content = ReadAsset("/assets/stylesheets/sub/normal.css")
	assert.Equal(t, "sub/normal.css\n", content)

	content = ReadAsset("/assets/not/exists/normal.css")
	assert.Equal(t, "", content)
}

func TestReadingAssetsWithRequire(t *testing.T) {
	Config.BundleAssets = true
	var content string

	content = ReadAsset("/assets/javascripts/require.js")
	assert.Equal(t, `normal.js

sub/normal.js

sub/require.js

require.js
`, content)

	content = ReadAsset("/assets/stylesheets/require.css")
	assert.Equal(t, `normal.css

sub/normal.css

sub/require.css

require.css
`, content)

	Config.BundleAssets = false
	content = ReadAsset("/assets/stylesheets/require.css")
	assert.Equal(t, `/*
 *= require stylesheets/normal
 *= require stylesheets/sub/require
 */
require.css
`, content)
}

func TestReadingStaticAssets(t *testing.T) {
	var content string

	content = ReadAsset("/assets/static.txt")
	assert.Equal(t, "static.txt\n", content)

	content = ReadAsset("/assets/images/dummy.png")
	assert.Equal(t, "dummy\n", content)
}
