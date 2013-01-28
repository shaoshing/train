package train

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestReadingNormalAssets(t *testing.T) {
	Config.AssetsPath = "test"
	var content string

	content = ReadAsset("/assets/normal.js")
	assert.Equal(t, "normal.js\n\n", content)

	content = ReadAsset("/assets/sub/normal.js")
	assert.Equal(t, "sub/normal.js\n\n", content)

	content = ReadAsset("/assets/not/exists/normal.js")
	assert.Equal(t, "", content)

	content = ReadAsset("/assets/normal.css")
	assert.Equal(t, "normal.css\n\n", content)

	content = ReadAsset("/assets/sub/normal.css")
	assert.Equal(t, "sub/normal.css\n\n", content)

	content = ReadAsset("/assets/not/exists/normal.css")
	assert.Equal(t, "", content)
}

func TestReadingAssetsWithRequire(t *testing.T) {
	Config.AssetsPath = "test"
	var content string

	content = ReadAsset("/assets/require.js")
	assert.Equal(t, `normal.js

sub/normal.js

sub/require.js

require.js

`, content)

	content = ReadAsset("/assets/require.css")
	assert.Equal(t, `normal.css

sub/normal.css

sub/require.css

require.css

`, content)
}

func TestReadingStaticAssets(t *testing.T) {
	Config.AssetsPath = "test"
	var content string

	content = ReadAsset("/assets/static.txt")
	assert.Equal(t, "static.txt\n", content)

	content = ReadAsset("/assets/dummy.png")
	assert.Equal(t, "dummy\n", content)
}
