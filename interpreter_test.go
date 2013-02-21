package train

import (
	"fmt"
	"github.com/shaoshing/gotest"
	"testing"
)

func TestSass(t *testing.T) {
	assert.Test = t

	css, e := CompileSASS("assets/stylesheets/app.sass")
	if e != nil {
		fmt.Println(e)
	}
	assert.Contain("h1", css)
	assert.Contain("h2", css)

	css, e = CompileSASS("assets/stylesheets/app.err.sass")
	assert.True(e != nil)
	assert.Contain("Could not render sass:", e.Error())
}
