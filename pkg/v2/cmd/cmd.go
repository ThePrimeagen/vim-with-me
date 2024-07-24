package cmd

import (
	"io"
	"log/slog"
	"os/exec"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type writerFn = func(b []byte) (int, error)

type fnAsWriter struct {
    fn writerFn
}

func (f *fnAsWriter) Write(b []byte) (int, error) {
    return f.fn(b)
}

type Cmder struct {
    Err io.Writer
    Out io.Writer
    In io.Reader
    Name string
    Args []string
    cmd *exec.Cmd
    stdin io.WriteCloser
}

func NewCmder(name string) *Cmder {
    return &Cmder{
        Err: nil,
        Out: nil,
        Name: name,
        Args: []string{},
    }
}

func (c *Cmder) AddArg(name string, value string) *Cmder {
    c.Args = append(c.Args, name, value)
    return c;
}

func (c *Cmder) WithErrFn(fn writerFn) *Cmder {
    c.Err = &fnAsWriter{fn: fn}
    return c;
}

func (c *Cmder) WithErr(writer io.Writer) *Cmder {
    c.Err = writer;
    return c;
}

func (c *Cmder) WithOutFn(fn writerFn) *Cmder {
    c.Out = &fnAsWriter{fn: fn}
    return c;
}

func (c *Cmder) WithOut(writer io.Writer) *Cmder {
    c.Out = writer;
    return c;
}

func (c *Cmder) Close() {
    err := c.cmd.Process.Kill();
    if err != nil {
        slog.Error("cannot close cmder", "err", err)
    }
}

func (c *Cmder) WriteLine(b []byte) error {
    read := 0
    for read < len(b) {
        n, err := c.stdin.Write(b[read:])
        if err != nil {
            return err
        }
        read += n
    }
    if b[len(b) - 1] != '\n' {
        _, _ = c.stdin.Write([]byte{'\n'})
    }

    return nil
}

func (c *Cmder) Run() error {
    assert.Assert(c.Out != nil, "you should never spawn a cmd without at least listening to stdout")
    assert.Assert(c.Name != "", "you need to provide a name for the program to run")

    c.cmd = exec.Command(c.Name, c.Args...)

    stdin, err := c.cmd.StdinPipe()
    if err != nil {
        return err
    }
    c.stdin = stdin

    stdout, err := c.cmd.StdoutPipe()
    if err != nil {
        return err
    }
    stderr, err := c.cmd.StderrPipe()
    if err != nil {
        return err
    }

    err = c.cmd.Start()
    if err != nil {
        return err
    }

    go io.Copy(c.Out, stdout)
    if c.Err != nil {
        go io.Copy(c.Err, stderr)
    }

    return c.cmd.Wait()
}
