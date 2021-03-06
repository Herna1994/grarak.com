package shell

import (
	"os/exec"
	"io"
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"time"
)

type ShellErr string

func (e ShellErr) Error() string {
	return string(e)
}

type Shell struct {
	cmd    *exec.Cmd
	writer io.WriteCloser
	reader io.ReadCloser

	exit, running bool
}

func NewShell() *Shell {
	cmd := exec.Command("sh")

	writer, err := cmd.StdinPipe()
	if err != nil {
		panic("Can't open shell")
	}

	reader, err := cmd.StdoutPipe()
	if err != nil {
		panic("Can't open shell")
	}

	cmd.Start()
	return &Shell{cmd, writer, reader,
		false, false}
}

func (sh *Shell) Run(cmd []byte) ([]byte, int, error) {
	if sh.exit {
		return nil, -1, ShellErr("Shell is killed")
	}
	sh.running = true
	callback := "/shellcallback/"
	_, err := io.WriteString(sh.writer,
		fmt.Sprintf("%s\necho $?%s\n", string(cmd), callback))
	if err != nil {
		return nil, -1, err
	}
	sh.cmd.Run()

	read := bufio.NewReader(sh.reader)
	output := make([]string, 0)
	status := -1
	for {
		buf, _, err := read.ReadLine()
		if err != nil {
			return nil, -1, ShellErr("Something went horribly wrong")
		}
		if strings.Contains(string(buf), callback) {
			output := strings.Replace(string(buf), callback, "", 1)
			status, err = strconv.Atoi(output)
			if err != nil {
				return nil, -1, ShellErr("Something went horribly wrong")
			}
			break
		}
		output = append(output, string(buf))
	}
	sh.running = false
	return []byte(strings.Join(output, "\n")), status, nil
}

func (sh *Shell) kill() {
	sh.cmd.Process.Kill()
	sh.writer.Close()
	sh.reader.Close()
}

func (sh *Shell) Exit() {
	sh.exit = true
	for sh.running {
		time.Sleep(time.Second / 3)
	}
	if sh.writer != nil {
		io.WriteString(sh.writer, "exit\n")

		if sh.cmd != nil {
			sh.cmd.Process.Wait()
		}
		sh.kill()
	}
}
