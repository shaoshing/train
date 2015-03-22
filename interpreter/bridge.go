package interpreter

import (
	"errors"
	"os/exec"
	"path"
	"strings"
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
		opts := []string{"-t", "nested"}

		if Config.SASS.LineNumbers {
			opts = append(opts, "--line-numbers")
		}
		if Config.SASS.DebugInfo {
			opts = append(opts, "--line-comments")
		}
		opts = append(opts, "-I", fileDir, filePath)
		out, e := exec.Command("sassc", opts...).Output()
		result = string(out)
		if e != nil {
			err = errors.New("Could not compile sass: 'sassc " +
				strings.Join(opts, " ") + "' failed: " + e.Error())
		}
	case ".coffee":
		out, e := exec.Command("coffee", "-p", filePath).Output()
		result = string(out)
		if e != nil {
			err = errors.New("Could not compile coffee: 'coffee -p " +
				filePath + "' failed: " + e.Error())
		}
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
