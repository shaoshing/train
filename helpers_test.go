package train

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestIncludeTags(t *testing.T) {
	Config.AssetsPath = "test"
	Config.BundleAssets = false

	helpers := Helpers{}

	assert.Equal(t, `<script src="/assets/javascripts/normal.js"></script>`, string(helpers.JavascriptIncludeTag("normal")))
	assert.Equal(t, `<script src="/assets/javascripts/normal.js"></script>
<script src="/assets/javascripts/sub/normal.js"></script>
<script src="/assets/javascripts/sub/require.js"></script>
<script src="/assets/javascripts/require.js"></script>`, string(helpers.JavascriptIncludeTag("require")))

	assert.Equal(t, `<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css">`, string(helpers.StylesheetIncludeTag("normal")))
	assert.Equal(t, `<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/sub/normal.css">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/sub/require.css">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/require.css">`, string(helpers.StylesheetIncludeTag("require")))

	Config.BundleAssets = true
	assert.Equal(t, `<script src="/assets/javascripts/require.js"></script>`, string(helpers.JavascriptIncludeTag("require")))
	assert.Equal(t, `<link type="text/css" rel="stylesheet" href="/assets/stylesheets/require.css">`, string(helpers.StylesheetIncludeTag("require")))
}
