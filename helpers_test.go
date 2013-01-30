package train

import (
	"github.com/bmizerany/assert"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestHelpers(t *testing.T) {
	Config.BundleAssets = false
	now := time.Now()
	stamp := strconv.FormatInt(now.Unix(), 10)
	updateAssetTimes(now)

	helpers := Helpers{}

	assert.Equal(t, `<script src="/assets/javascripts/normal.js?`+stamp+`"></script>`, string(helpers.JavascriptTag("normal")))
	assert.Equal(t, `<script src="/assets/javascripts/normal.js?`+stamp+`"></script>
<script src="/assets/javascripts/sub/normal.js?`+stamp+`"></script>
<script src="/assets/javascripts/sub/require.js?`+stamp+`"></script>
<script src="/assets/javascripts/require.js?`+stamp+`"></script>`, string(helpers.JavascriptTag("require")))

	assert.Equal(t, `<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css?`+stamp+`">`, string(helpers.StylesheetTag("normal")))
	assert.Equal(t, `<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/sub/normal.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/sub/require.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/require.css?`+stamp+`">`, string(helpers.StylesheetTag("require")))

	assert.Equal(t, `<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css?`+stamp+`" media="print">`, string(helpers.StylesheetTagWithParam("normal", `media="print"`)))

	Config.BundleAssets = true
	assert.Equal(t, `<script src="/assets/javascripts/require.js?`+stamp+`"></script>`, string(helpers.JavascriptTag("require")))
	assert.Equal(t, `<link type="text/css" rel="stylesheet" href="/assets/stylesheets/require.css?`+stamp+`">`, string(helpers.StylesheetTag("require")))

	Config.BundleAssets = false
}

func updateAssetTimes(t time.Time) {
	filepath.Walk(Config.AssetsPath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			os.Chtimes(filePath, t, t)
		}
		return nil
	})
}
