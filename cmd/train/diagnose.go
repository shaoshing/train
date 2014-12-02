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

	fmt.Println("== Diagnosing\n")

	var rubyVersion string
	rubyVersion, err = bash(`ruby -e "puts RUBY_VERSION"`)
	if err != nil {
		fmt.Println("-- SASS and CoffeeScript are disabled because ruby is not installed.")
		fmt.Println("   (visit https://github.com/sstephenson/rbenv/#installation for installation instructions)")
		return false
	}

	if !strings.Contains(rubyVersion, "1.9") {
		fmt.Printf("-- Train requires Ruby version to be 1.9.x; you have %s", rubyVersion)
		fmt.Println("   (Please install required Ruby version if you wish to use SASS or CoffeeScript)")
		return false
	}

	allGood := true
	_, filename, _, _ := runtime.Caller(1)
	assetsPath := path.Dir(filename) + "/assets"

	_, err = bash("gem which sass")
	if err != nil {
		fmt.Println("-- SASS is disabled because the required gem is not found.")
		fmt.Println("   (install it if you wish to use SASS: gem install sass)\n")
		allGood = false
	} else {
		_, err = interpreter.Compile(assetsPath + "/stylesheets/font.sass")
		if err != nil {
			fmt.Println("-- SASS is disabled because error raised while compiling. Error:")
			fmt.Printf("%s\n", err.Error())
			fmt.Println("(this might related to your Ruby Environment; try re-installing Ruby)\n")
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
			fmt.Println("-- CoffeeScript is disabled because error raised while compiling. Error: ")
			fmt.Printf("%s\n", err.Error())
			fmt.Println("(this might related to your Ruby Environment; try re-installing Ruby)\n")
			allGood = false
		}
	}

	if allGood {
		fmt.Println("-- Great, your environment seems perfect for Train.")
	} else {
		fmt.Println("-- (Please create an issue at github.com/shaoshing/train/issues if you need help)")
	}

	return allGood
}
