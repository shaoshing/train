package train

import (
	"github.com/shaoshing/gotest"
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
	assert.Test = t
	Config.BundleAssets = true

	assertAsset("/assets/static.txt", "static.txt\n", "text/plain")
	assertAsset("/assets/images/dummy.png", "dummy\n", "image/png")
	assert404("/assets/not/found.js")

	assertAsset("/assets/javascripts/normal.js", "normal.js\n", "application/javascript")
	assertAsset("/assets/stylesheets/normal.css", "normal.css\n", "text/css")

	assertAsset("/assets/javascripts/require.js", `normal.js

sub/normal.js

sub/require.js

require.js
`, "application/javascript")
	assertAsset("/assets/stylesheets/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`, "text/css")

	assertAsset("/assets/stylesheets/app.css", `h1 {
  color: green; }

h2 {
  color: green; }
`, "text/css")

	body, _, status := get("/assets/stylesheets/app.err.css")
	assert.Contain("Could not compile sass", body)
	assert.Equal(500, status)
}

func TestBundledAssets(t *testing.T) {
	assert.Test = t
	exec.Command("cp", "-rf", "assets/public", "./").Run()
	defer exec.Command("rm", "-rf", "public").Run()

	assertAsset("/assets/app.js", "app.js\n", "application/javascript")
	assert404("/assets/normal.js")
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

func assertAsset(url, expectedBody, expectedContentType string) {
	body, contentType, _ := get(url)
	assert.Equal(expectedBody, body)
	assert.Equal(true, strings.Index(contentType, expectedContentType) != -1)
}

func assert404(url string) {
	_, _, status := get(url)
	assert.Equal(404, status)
}
