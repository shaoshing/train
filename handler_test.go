package train

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

var httpClient = http.Client{}
var httpServer = httptest.NewServer(http.HandlerFunc(Handler))

func TestHandler(t *testing.T) {
	Config.BundleAssets = true

	assertAsset(t, "/assets/static.txt", "static.txt\n", "text/plain")
	assertAsset(t, "/assets/images/dummy.png", "dummy\n", "image/png")
	assert404(t, "/assets/not/found.js")

	assertAsset(t, "/assets/javascripts/normal.js", "normal.js\n", "application/javascript")
	assertAsset(t, "/assets/stylesheets/normal.css", "normal.css\n", "text/css")

	assertAsset(t, "/assets/javascripts/require.js", `normal.js

sub/normal.js

sub/require.js

require.js
`, "application/javascript")
	assertAsset(t, "/assets/stylesheets/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`, "text/css")
}

func TestBundledAssets(t *testing.T) {
	exec.Command("cp", "-rf", "assets/public", "./").Run()
	defer exec.Command("rm", "-rf", "public").Run()

	assertAsset(t, "/assets/app.js", "app.js\n", "application/javascript")
	assert404(t, "/assets/normal.js")
}

func get(url string) (body, contentType string, status int) {
	response, err := httpClient.Get(httpServer.URL + url)
	if err != nil {
		panic(err)
	}
	b_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	body = string(b_body)
	contentType = response.Header.Get("Content-Type")
	status = response.StatusCode

	return
}

func assertAsset(t *testing.T, url, expectedBody, expectedContentType string) {
	body, contentType, _ := get(url)
	assert.Equal(t, expectedBody, body)
	assert.Equal(t, true, strings.Index(contentType, expectedContentType) != -1)
}

func assert404(t *testing.T, url string) {
	_, _, status := get(url)
	assert.Equal(t, 404, status)
}
