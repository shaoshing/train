package interpreter

import (
	"bytes"
	"errors"
	"os/exec"
	"path"
	"sync"
)

type Interpreter struct {
	cmd        *exec.Cmd
	SocketName string
	sync.Mutex
}

var (
	interpreter Interpreter
)

type config struct {
	Verbose    bool
	AssetsPath string
	SASS       struct {
		DebugInfo   bool
		LineNumbers bool
	}
}

var Config = config{
	// AssetsPath for the SASS files. By default it will look for SASS files under
	// the assets/stylesheets folder.
	AssetsPath: "assets",
}

func Compile(filePath string) (result string, err error) {
	fileExt := path.Ext(filePath)
	switch fileExt {
	case ".sass", ".scss":
		fileDir := path.Dir(filePath)
		cmd := exec.Command("sassc", "-t", "nested", "-I", fileDir, filePath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
		result = out.String()
	case ".coffee":
		cmd := exec.Command("coffee", "-p", filePath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
		result = out.String()
	default:
		err = errors.New("Unsupported format (" + filePath + "). Valid formats are: sass.")
	}

	return
}

func getOption() string {
	if Config.SASS.LineNumbers {
		return "line_numbers"
	}
	if Config.SASS.DebugInfo {
		return "debug_info"
	}
	return ""
}
