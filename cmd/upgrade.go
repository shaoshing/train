package main

import (
	"fmt"
)

const CmdBinPath = "$GOPATH/bin/train"

func upgrade() {
	bash("go get -u github.com/huacnlee/train")
	bash("go build -o " + CmdBinPath + " github.com/huacnlee/train/cmd")
	fmt.Println("Installed latest train command into " + CmdBinPath)
}
