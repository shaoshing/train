package main

import (
	"fmt"
	"github.com/huacnlee/train/interpreter"
	"path"
	"runtime"
)

func diagnose() bool {
	var err error

	fmt.Println("== Diagnosing\n")

	allGood := true
	_, filename, _, _ := runtime.Caller(1)
	assetsPath := path.Dir(filename) + "/assets"

	_, err = bash("node-sass --version")
	if err != nil {
		fmt.Println("-- SASS is disabled because the required npm is not found.")
		fmt.Println("   (install it if you wish to use SASS: npm install -g node-sass)\n")
		allGood = false
	} else {
		_, err = interpreter.Compile(assetsPath + "/stylesheets/font.sass")
		if err != nil {
			fmt.Println("-- SASS is disabled because error raised while compiling. Error:")
			fmt.Printf("%s\n", err.Error())
			allGood = false
		}
	}

	_, err = bash("coffee -v")
	if err != nil {
		fmt.Println("-- CoffeeScript is disabled because the required npm is not found.")
		fmt.Println("   (install it if you wish to use CoffeeScript : npm install -g coffee-script)")
		allGood = false
	} else {
		_, err = interpreter.Compile(assetsPath + "/javascripts/app.coffee")
		if err != nil {
			fmt.Println("-- CoffeeScript is disabled because error raised while compiling. Error: ")
			fmt.Printf("%s\n", err.Error())
			allGood = false
		}
	}

	if allGood {
		fmt.Println("-- Great, your environment seems perfect for Train.")
	} else {
		fmt.Println("-- (Please create an issue at github.com/huacnlee/train/issues if you need help)")
	}

	return allGood
}
