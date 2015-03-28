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
		catCmd := exec.Command("cat", filePath)
		opts := []string{"--output-style", "nested", "--indent-type", "space", "--indent-width", "2", "--linefeed", "lf"}

		if Config.SASS.LineNumbers {
			opts = append(opts, "--source-comments")
		}
		if Config.SASS.DebugInfo {
			opts = append(opts, "--source-comments")
		}
		if fileExt == ".sass" {
			opts = append(opts, "--indented-syntax")
		}
		opts = append(opts, "--include-path", fileDir)
		cmd := exec.Command("node-sass", opts...)

		out, e := pipeExecCommand(catCmd, cmd)

		result = strings.TrimSpace(string(out))
		if e != nil {
			err = errors.New("Could not compile sass: 'cat " + filePath + " | node-sass" +
				strings.Join(opts, " ") + "' failed: " + e.Error())

		}
	case ".coffee":
		out, e := exec.Command("coffee", "-p", filePath).CombinedOutput()
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

func pipeExecCommand(cmds ...*exec.Cmd) ([]byte, error) {
	for i, cmd := range cmds[:len(cmds)-1] {
		out, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		cmd.Start()
		cmds[i+1].Stdin = out
	}

	ret, err := cmds[len(cmds)-1].Output()
	if err != nil {
		return nil, err
	}

	return ret, nil
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
