package interpreter

import (
	"github.com/shaoshing/gotest"
	"testing"
	"time"
)

func init() {
	Config.Verbose = true
}

func TestSass(t *testing.T) {
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

func TestConcurrency(t *testing.T) {
	assert.Test = t
	concurrency := 10
	compileChan := make(chan bool)
	for i := 0; i < concurrency; i++ {
		go func() {
			Compile("assets/stylesheets/app.sass")
			compileChan <- true
		}()
	}

	completeChan := make(chan bool)
	go func() {
		for i := 0; i < concurrency; i++ {
			<-compileChan
		}
		completeChan <- true
	}()

	success := true
	select {
	case <-completeChan:
		success = true
	case <-time.After(10 * time.Second):
		success = false
	}

	assert.True(success)
}
