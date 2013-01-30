package train

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var httpClient = http.Client{}
var httpServer = httptest.NewServer(http.HandlerFunc(Handler))

func get(url string) (body, contentType string) {
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

	return
}

func assertGet(t *testing.T, url, expectedBody, expectedContentType string) {
	body, contentType := get(url)
	assert.Equal(t, expectedBody, body)
	assert.Equal(t, true, strings.Index(contentType, expectedContentType) != -1)
}

func TestHandler(t *testing.T) {
	Config.BundleAssets = true

	assertGet(t, "/assets/static.txt", "static.txt\n", "text/plain")
	assertGet(t, "/assets/images/dummy.png", "dummy\n", "image/png")
	assertGet(t, "/assets/javascripts/normal.js", "normal.js\n", "application/javascript")
	assertGet(t, "/assets/stylesheets/normal.css", "normal.css\n", "text/css")

	assertGet(t, "/assets/javascripts/require.js", `normal.js

sub/normal.js

sub/require.js

require.js
`, "application/javascript")
	assertGet(t, "/assets/stylesheets/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`, "text/css")
}
