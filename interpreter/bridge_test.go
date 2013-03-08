package interpreter

import (
	"github.com/shaoshing/gotest"
	"testing"
)

var closingCount = 2
func closeTestInterpreter() {
	if closingCount == 0 {
		CloseInterpreter()
	}
}

func TestSass(t *testing.T) {
	closingCount--
	defer closeTestInterpreter()

	assert.Test = t

	css, e := Compile("assets/stylesheets/app.sass")
	if e != nil {
		panic(e)
	}
	assert.Contain("h1", css)
	assert.Contain("h2", css)

	css, e = Compile("assets/stylesheets/app2.scss")
	if e != nil {
		panic(e)
	}
	assert.Contain("h2", css)
	assert.Contain("h3", css)

	css, e = Compile("assets/stylesheets/app.err.sass")
	assert.True(e != nil)
	assert.Contain("Could not compile sass:", e.Error())

	Config.SASS.DebugInfo = true
	css, e = Compile("assets/stylesheets/app.sass")
	assert.Contain("-sass-debug-info", css)

	Config.SASS.LineNumbers = true
	css, e = Compile("assets/stylesheets/app.sass")
	assert.Contain("line 1", css)
}

func TestCoffee(t *testing.T) {
	closingCount--
	defer closeTestInterpreter()
	
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
