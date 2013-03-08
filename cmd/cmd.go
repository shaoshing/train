package main

import (
	"bytes"
	"flag"
	"fmt"
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
	case "help":
		fmt.Printf(`Available commands:

   bundle: bundle assets into ./public/assets [default]
  upgrade: install the latest qortex command.
     help: show this help info.
`)
	default:
		panic(fmt.Sprintf("Unknown command %s", command))
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
