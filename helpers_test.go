package train

import (
	"github.com/shaoshing/gotest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestHelpers(t *testing.T) {
	assert.Test = t
	Config.BundleAssets = false
	now := time.Now()
	stamp := strconv.FormatInt(now.Unix(), 10)
	updateAssetTimes(now)

	assert.Equal(`<script src="/assets/javascripts/normal.js?`+stamp+`"></script>`, string(JavascriptTag("normal")))
	assert.Equal(`<script src="/assets/javascripts/normal.js?`+stamp+`"></script>
<script src="/assets/javascripts/sub/normal.js?`+stamp+`"></script>
<script src="/assets/javascripts/sub/require.js?`+stamp+`"></script>
<script src="/assets/javascripts/require.js?`+stamp+`"></script>`, string(JavascriptTag("require")))

	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css?`+stamp+`">`, string(StylesheetTag("normal")))
	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/sub/normal.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/sub/require.css?`+stamp+`">
<link type="text/css" rel="stylesheet" href="/assets/stylesheets/require.css?`+stamp+`">`, string(StylesheetTag("require")))

	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/stylesheets/normal.css?`+stamp+`" media="print">`, string(StylesheetTagWithParam("normal", `media="print"`)))

	assert.Equal(`<script src="/assets/javascripts/app.js?`+stamp+`"></script>`, string(JavascriptTag("app")))
	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/stylesheets/app.css?`+stamp+`">`, string(StylesheetTag("app")))

	Config.Mode = PRODUCTION_MODE
	defer func() {
		Config.Mode = DEVELOPMENT_MODE
	}()
	ManifestInfo = FpAssets{
		"/assets/javascripts/require.js":  "/assets/javascripts/require-fingerprintinghash.js",
		"/assets/stylesheets/require.css": "/assets/stylesheets/require-fingerprintinghash.css",
	}

	assert.Equal(`<script src="/assets/javascripts/require-fingerprintinghash.js"></script>`, string(JavascriptTag("require")))
	assert.Equal(`<link type="text/css" rel="stylesheet" href="/assets/stylesheets/require-fingerprintinghash.css">`, string(StylesheetTag("require")))

	ManifestInfo = FpAssets{}
}

func updateAssetTimes(t time.Time) {
	filepath.Walk(Config.AssetsPath, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			os.Chtimes(filePath, t, t)
		}
		return nil
	})
}
