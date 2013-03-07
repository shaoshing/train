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
	"strings"
	"strconv"
	"syscall"
)

// TODO: need refactor
var interpreter *Interpreter

var Config struct {
	SASS struct {
		DebugInfo   bool
		LineNumbers bool
	}
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

type Interpreter struct {
	cmd        *exec.Cmd
	ready      bool
	queue      []chan bool
	socketName string
}

func NewInterpreter() {
	var i Interpreter
	_, goFile, _, _ := runtime.Caller(0)
	i.socketName = "/tmp/train.interpreter.socket"
	i.cmd = exec.Command("ruby", path.Dir(goFile)+"/interpreter.rb")
	i.cmd.Stdout = &StdoutCapturer{&i}
	
	go func() {
		err := i.cmd.Run()
		if err != nil {
			panic(err)
		}
	}()

	// return &i
	interpreter = &i
}

func CloseInterpreter() {
	_, goFile, _, _ := runtime.Caller(0)
	dat, err := ioutil.ReadFile(path.Dir(goFile) + "/interpreter.pid")
	if err != nil { 
	    panic(err)
	}
	
	pid, err := strconv.Atoi(string(dat))
	if err != nil { 
	    panic(err)
	}
	
	err = syscall.Kill(pid, syscall.Signal(9))
	if err != nil {
		panic(err)
	}
	
	interpreter = nil
}

func (this *Interpreter) Render(format string, content []byte) (result string, err error) {
	if !this.ready {
		this.Wait()
	}

	conn, err := net.Dial("unix", this.socketName)
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

func (this *Interpreter) Wait() {
	if this.ready {
		return
	}
	c := make(chan bool)
	this.queue = append(this.queue, c)
	
	<-c
}

func (this *Interpreter) Ready() {
	this.ready = true
	for _, c := range this.queue {
		c <- true
	}
	this.queue = make([]chan bool, 0)
}

type StdoutCapturer struct {
	interpreter *Interpreter
}

func (this *StdoutCapturer) Write(p []byte) (n int, err error) {
	if strings.Contains(string(p), "<<ready") {
		this.interpreter.Ready()
	}
	n, err = os.Stdout.Write(p)
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
