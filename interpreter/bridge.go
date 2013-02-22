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
)

type Interpreter struct {
	cmd        *exec.Cmd
	ready      bool
	queue      []chan bool
	socketName string
}

type StdoutCapturer struct {
	interpreter *Interpreter
}

func NewInterpreter() *Interpreter {
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

	return &i
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

func (this *Interpreter) Render(content []byte) (result string, err error) {
	if !this.ready {
		this.Wait()
	}

	conn, err := net.Dial("unix", this.socketName)
	if err != nil {
		panic(err)
	}

	conn.Write(content)
	var data bytes.Buffer
	data.ReadFrom(conn)
	conn.Close()

	if strings.Contains(data.String(), "<<error") {
		err = errors.New("Could not compile SASS:" + strings.Replace(data.String(), "<<error", "", 1))
	} else {
		result = data.String()
	}

	return
}

func (this *StdoutCapturer) Write(p []byte) (n int, err error) {
	if strings.Contains(string(p), "<<ready") {
		this.interpreter.Ready()
	}
	n, err = os.Stdout.Write(p)
	return
}

var sass *Interpreter

func init() {
	sass = NewInterpreter()
}

func CompileSASS(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return sass.Render(content)
}
