package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/shaoshing/train"
	"os/exec"
)

func main() {
	flag.Parse()

	command := "bundle"

	args := flag.Args()
	if len(args) == 1 {
		command = args[0]
	}

	switch command {
	case "bundle":
		bundle()
	case "upgrade":
		upgrade()
	case "diagnose":
		diagnose()
	case "version":
		fmt.Println("Train version", train.VERSION)
	case "help":
		showHelp()
	default:
		showHelp()
	}
}

func bash(bash string) (out string, err error) {
	cmd := exec.Command("sh", "-c", bash)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err = cmd.Run()
	out = buf.String()
	return
}

func showHelp() {
	fmt.Printf(`usage: train [command]

Commands:

   bundle: bundle assets into ./public/assets [run by default]
  upgrade: install the latest qortex command.
 diagnose: trouble shooting.
  version: %s
`, train.VERSION)
}
