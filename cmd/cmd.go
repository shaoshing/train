package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/huacnlee/train"
	"os/exec"
)

var helpFlag bool
var sourcePath string
var outPath string

func main() {
	flag.BoolVar(&helpFlag, "h", false, "")
	flag.StringVar(&sourcePath, "source", "./assets", "")
	flag.StringVar(&outPath, "out", "./public", "")
	flag.Parse()

	command := "bundle"

	args := flag.Args()
	if len(args) >= 1 {
		command = args[0]
	}

	if helpFlag {
		showHelp()
		return
	}

	switch command {
	case "bundle":
		Bundle(sourcePath, outPath)
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

  --source
    Assets source path, default: ./assets
    example: $ train --source app/assets bundle

  --out
    Assets output path, default: ./public
    example: $ train --out /tmp/public bundle

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
