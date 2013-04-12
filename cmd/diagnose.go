package main

import (
	"fmt"
	"github.com/shaoshing/train/interpreter"
	"path"
	"runtime"
	"strings"
)

func diagnose() bool {
	var err error

	fmt.Println("== Diagnosing")

	var rubyVersion string
	rubyVersion, err = bash(`ruby -e "puts RUBY_VERSION"`)
	if err != nil {
		fmt.Println("-- SASS and CoffeeScript are disabled because ruby is not installed.")
		fmt.Println("   (visit http://www.ruby-lang.org/en/downloads/ for installation instructions)")
		return false
	}

	if !strings.Contains(rubyVersion, "1.9") {
		fmt.Printf("-- Train requires Ruby version to be 1.9.x; you have %s", rubyVersion)
		fmt.Println("   (Please install required Ruby version if you wish to use SASS or CoffeeScript)")
		return false
	}

	allGood := true
	interpreter.Config.Verbose = true
	_, filename, _, _ := runtime.Caller(1)
	assetsPath := path.Dir(filename) + "/assets"

	_, err = bash("gem which sass")
	if err != nil {
		fmt.Println("-- SASS is disabled because the required gem is not found.")
		fmt.Println("   (install it if you wish to use SASS: gem install sass)")
		allGood = false
	} else {
		_, err = interpreter.Compile(assetsPath + "/stylesheets/app.sass")
		if err != nil {
			fmt.Println("-- Could not compile SASS:")
			fmt.Printf("   %s\n", err.Error())
			allGood = false
		}
	}

	_, err = bash("gem which coffee-script")
	if err != nil {
		fmt.Println("-- CoffeeScript is disabled because the required gem is not found.")
		fmt.Println("   (install it if you wish to use CoffeeScript : gem install coffee-script)")
		allGood = false
	} else {
		_, err = interpreter.Compile(assetsPath + "/javascripts/app.coffee")
		if err != nil {
			fmt.Println("-- Could not compile CoffeeScript:")
			fmt.Printf("   %s\n", err.Error())
			allGood = false
		}
	}

	if allGood {
		fmt.Println("-- Great, your environment is perfect for Train.")
	} else {
		fmt.Println("-- (If you need help, please create an issue at github.com/shaoshing/train/issues)")
	}

	return allGood
}
