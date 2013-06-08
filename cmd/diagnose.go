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
		fmt.Println("   (install it if you wish to use SASS: gem install sass)")
		allGood = false
	} else {
		_, err = interpreter.Compile(assetsPath + "/stylesheets/font.sass")
		if err != nil {
			fmt.Println("-- SASS is disabled because error raised while compiling. Error:")
			fmt.Printf("%s\n", err.Error())
			fmt.Println("(this might related to your Ruby Environment; try re-installing Ruby)")
			allGood = false
		} else {
			sassVersion, _ := bash(`ruby -e "require 'sass'; puts Sass::VERSION"`)
			supportSourceMap, _ := bash(`ruby -e 'require "sass"; e = Sass::Engine.new ""; puts(e.respond_to?(:render_with_sourcemap) ? "yes" : "no")'`)
			fmt.Printf("-- SASS [supported] version: %s, sourcemap: %s\n", strings.Trim(sassVersion, "\n"), strings.Trim(supportSourceMap, "\n"))
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
			fmt.Println("(this might related to your Ruby Environment; try re-installing Ruby)")
			allGood = false
		} else {
			coffeeVersion, _ := bash(`ruby -e "require 'coffee-script'; puts CoffeeScript.version"`)
			fmt.Printf("-- Coffee [supported] version: %s\n", strings.Trim(coffeeVersion, "\n"))
		}
	}

	if allGood {
		fmt.Println("-- Great, your environment seems perfect for Train.")
	} else {
		fmt.Println("-- (Please create an issue at github.com/shaoshing/train/issues if you need help)")
	}

	return allGood
}
