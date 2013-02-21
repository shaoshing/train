package train

import (
	"fmt"
	"github.com/shaoshing/gotest"
	"io/ioutil"
	"testing"
)

func TestSass(t *testing.T) {
	assert.Test = t

	scss, _ := ioutil.ReadFile("assets/stylesheets/app.sass")
	css, e := CompileSASS(scss)
	if e != nil {
		fmt.Println(e)
	}
	assert.Contain("h1", css)
	assert.Contain("h2", css)

	css, e = CompileSASS([]byte("body {color:red "))
	assert.True(e != nil)
	assert.Contain("Could not render sass:", e.Error())
}
