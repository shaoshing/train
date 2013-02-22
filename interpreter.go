package train

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const SassSocketName = "/tmp/train.sass.socket"

type Interpreter struct {
	cmd        *exec.Cmd
	ready      bool
	queue      []chan bool
	socketName string
}

type StdoutCapturer struct {
	interpreter *Interpreter
}

func NewInterpreter(file, socketName string) *Interpreter {
	var i Interpreter

	_, goFile, _, _ := runtime.Caller(0)
	i.socketName = socketName
	i.cmd = exec.Command("ruby", path.Dir(goFile)+"/interpreters/"+file+".rb")
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
		err = errors.New("Could not render sass:" + data.String())
	} else {
		result = data.String()
	}

	return
}

func (this *StdoutCapturer) Write(p []byte) (n int, err error) {
	if strings.Contains(string(p), "<<ready") {
		this.interpreter.Ready()
	}
	return
}

var sass *Interpreter

func init() {
	sass = NewInterpreter("sass", SassSocketName)
}

func CompileSASS(path string) (string, error) {
	content, _ := ioutil.ReadFile(path)
	return sass.Render(content)
}
