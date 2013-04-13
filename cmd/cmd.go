package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/shaoshing/train"
	"os/exec"
)

var helpFlag bool

func main() {
	flag.BoolVar(&helpFlag, "h", false, "")
	flag.Parse()

	command := "bundle"

	args := flag.Args()
	if len(args) == 1 {
		command = args[0]
	}

	if helpFlag {
		showHelp()
		return
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
	fmt.Printf(`usage: train [-h] [command]

OPTIONS
  -h
    Show this help message

COMMANDS
  bundle [default]
    Bundle assets into ./public/assets

  upgrade
    Update the package and Install the train command.

  diagnose
    Trouble shooting for the Pipeline feature.

  help
    Show this help message.

  version
    Show version (current version: %s)
`, train.VERSION)
}
