package interpreter

import (
	"github.com/shaoshing/gotest"
	"testing"
)

func TestSass(t *testing.T) {
	assert.Test = t

	css, e := Compile("assets/stylesheets/app.sass")
	if e != nil {
		panic(e)
	}
	assert.Contain("h1", css)
	assert.Contain("h2", css)

	css, e = Compile("assets/stylesheets/app.err.sass")
	assert.True(e != nil)
	assert.Contain("Could not compile sass:", e.Error())
}

func TestCoffee(t *testing.T) {
	assert.Test = t

	css, e := Compile("assets/javascripts/app.coffee")
	if e != nil {
		panic(e)
	}
	assert.Contain("square", css)

	css, e = Compile("assets/javascripts/app.err.coffee")
	assert.True(e != nil)
	assert.Contain("Could not compile coffee:", e.Error())
}
