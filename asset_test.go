package train

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestReadAsset(t *testing.T) {
	Config.AssetsPath = "test"
	var content string

	content = ReadAsset("/assets/normal.js")
	assert.Equal(t, "test/javascripts/normal.js\n", content)

	content = ReadAsset("/assets/sub/normal.js")
	assert.Equal(t, "test/javascripts/sub/normal.js\n", content)

	content = ReadAsset("/assets/not/exists/normal.js")
	assert.Equal(t, "", content)

	content = ReadAsset("/assets/normal.css")
	assert.Equal(t, "test/stylesheets/normal.css\n", content)

	content = ReadAsset("/assets/sub/normal.css")
	assert.Equal(t, "test/stylesheets/sub/normal.css\n", content)

	content = ReadAsset("/assets/not/exists/normal.css")
	assert.Equal(t, "", content)
}
