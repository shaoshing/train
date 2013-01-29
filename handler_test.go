package train

import (
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var httpClient = http.Client{}
var httpServer = httptest.NewServer(http.HandlerFunc(Handler))

func get(url string) (body string) {
	response, err := httpClient.Get(httpServer.URL + url)
	if err != nil {
		panic(err)
	}
	b_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	body = string(b_body)

	return
}

func assertGet(t *testing.T, url, body string) {
	assert.Equal(t, body, get(url))
}

func TestHandler(t *testing.T) {
	Config.BundleAssets = true

	assertGet(t, "/assets/static.txt", "static.txt\n")
	assertGet(t, "/assets/images/dummy.png", "dummy\n")
	assertGet(t, "/assets/javascripts/normal.js", "normal.js\n")
	assertGet(t, "/assets/stylesheets/normal.css", "normal.css\n")

	assertGet(t, "/assets/javascripts/require.js", `normal.js

sub/normal.js

sub/require.js

require.js
`)
	assertGet(t, "/assets/stylesheets/require.css", `normal.css

sub/normal.css

sub/require.css

require.css
`)
}
