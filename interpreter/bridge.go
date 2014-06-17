package interpreter

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
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
	Verbose bool
	AssetsPath string
	SASS    struct {
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
	case ".sass", ".scss", ".coffee":
		content, e := ioutil.ReadFile(filePath)
		if e != nil {
			panic(err)
		}
		result, err = interpreter.Render(strings.Replace(fileExt, ".", "", 1), content)
	default:
		err = errors.New("Unsupported format (" + filePath + "). Valid formats are: sass.")
	}

	return
}

func (this *Interpreter) Render(format string, content []byte) (result string, err error) {
	this.Lock()
	defer this.Unlock()

	this.startRubyInterpreter()

	conn, err := net.Dial("unix", this.SocketName)
	if err != nil {
		panic(err)
	}

	option := getOption()

	conn.Write([]byte(format + "<<" + option + "<<" + string(content)))
	var data bytes.Buffer
	data.ReadFrom(conn)
	conn.Close()

	compiled := strings.Split(data.String(), "<<")

	status := compiled[0]
	result = compiled[1]

	if status == "error" {
		err = errors.New("Could not compile " + format + ":\n" + result)
	}

	return
}

func (this *Interpreter) startRubyInterpreter() {
	if this.IsStarted() {
		return
	}

	_, goFile, _, _ := runtime.Caller(0)
	this.SocketName = generateUniqueSocketName()
	currentPid := strconv.FormatInt(int64(os.Getpid()), 10)
	this.cmd = exec.Command("ruby", path.Dir(goFile)+"/interpreter.rb", this.SocketName, currentPid, Config.AssetsPath)
	waitForStarting := make(chan bool)
	writer := &StdoutCapturer{waitForStarting}
	this.cmd.Stdout = writer
	this.cmd.Stderr = writer
	err := this.cmd.Start()
	if err != nil {
		panic(err)
	}
	<-waitForStarting
}

func (this *Interpreter) IsStarted() bool {
	return this.cmd != nil
}

type StdoutCapturer struct {
	waitForStarting chan bool
}

func (this *StdoutCapturer) Write(p []byte) (int, error) {
	if strings.Contains(string(p), "<<ready") {
		this.waitForStarting <- true
	}

	if Config.Verbose {
		return os.Stdout.Write(p)
	}
	return len(p), nil
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

func generateUniqueSocketName() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	return "/tmp/train.interpreter." + timestamp + ".socket"
}
