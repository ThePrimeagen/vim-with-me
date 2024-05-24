package program

import (
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"golang.org/x/sys/unix"
)

type Program struct {
	*os.File
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	path   string
	rows   int
	cols   int
	writer io.Writer
	args   []string
}

func NewProgram(path string) *Program {
	return &Program{
		path:   path,
		rows:   80,
		cols:   24,
		writer: nil,
		cmd:    nil,
		File:   nil,
	}
}

func (a *Program) SendKey(key string) {
    for _, k := range key {
        a.Write([]byte{byte(k)})
        <-time.After(time.Millisecond * 40)
    }
}

func (a *Program) WithArgs(args []string) *Program {
	a.args = args
	return a
}

func (a *Program) WithWriter(writer io.Writer) *Program {
	if a.writer != nil {
		a.writer = io.MultiWriter(a.writer, writer)
	} else {
		a.writer = writer
	}
	return a
}

func (a *Program) WithSize(rows, cols int) *Program {
	a.rows = rows
	a.cols = cols
	return a
}

func echoOff(f *os.File) {
	fd := int(f.Fd())
	//      const ioctlReadTermios = unix.TIOCGETA // OSX.
	const ioctlReadTermios = unix.TCGETS // Linux
	//      const ioctlWriterTermios =  unix.TIOCSETA // OSX.
	const ioctlWriteTermios = unix.TCSETS // Linux

	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		panic(err)
	}

	newState := *termios
	newState.Lflag &^= unix.ECHO
	newState.Lflag |= unix.ICANON | unix.ISIG
	newState.Iflag |= unix.ICRNL
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		panic(err)
	}
}

func (a *Program) Run(ctx context.Context) error {
	assert.Assert(a.writer != nil, "you must provide a reader before you call run")
	assert.Assert(a.File == nil, "you have already started the program")

	cmd := exec.Command(a.path, a.args...)

	a.cmd = cmd

	ptmx, err := pty.Start(cmd)
	echoOff(ptmx)

	if err != nil {
		return err
	}

	a.File = ptmx

	_, err = io.Copy(a.writer, ptmx)
	return err
}

func (a *Program) Close() error {
	err := a.File.Close()
	a.File = nil
	return err
}
