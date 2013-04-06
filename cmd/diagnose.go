package main

import (
	"fmt"
)

func diagnose() {
	var err error
	allGood := true

	fmt.Println("== Diagnosing")
	_, err = bash("ruby --version")
	if err != nil {
		fmt.Println("-- SASS and CoffeeScript are disabled because ruby is not installed.")
		fmt.Println("   (visit http://www.ruby-lang.org/en/downloads/ for installation instructions)")
		return
	}

	_, err = bash("gem which sass")
	if err != nil {
		fmt.Println("-- SASS is disabled because the required gem is not found.")
		fmt.Println("   (install it if you wish to use SASS: gem install sass)")
		allGood = false
	}

	_, err = bash("gem which coffee-script")
	if err != nil {
		fmt.Println("-- CoffeeScript is disabled because the required gem is not found.")
		fmt.Println("   (install it if you wish to use CoffeeScript : gem install coffee-script)")
		allGood = false
	}

	if allGood {
		fmt.Println("-- Great, your environment is perfect for Train.")
	}
}
