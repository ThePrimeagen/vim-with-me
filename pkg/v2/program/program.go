package program

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
)

type Program struct {
    path string
    rows int
    cols int
    writer io.Writer
    args []string
    program *os.File
}

func NewProgram(name string) *Program {
    return &Program{
        path: name,
        rows: 80,
        cols: 24,
        writer: nil,
        program: nil,
    }
}

func (a* Program) WithArgs(args []string) *Program {
    a.args = args;
    return a
}

func (a* Program) WithWriter(writer io.Writer) *Program {
    a.writer = writer;
    return a
}

func (a* Program) WithSize(rows, cols int) *Program {
    a.rows = rows;
    a.cols = cols;
    return a
}

func (a* Program) Run(ctx context.Context) error {
    assert.Assert(a.writer != nil, "you must provide a reader before you call run")
    assert.Assert(a.program == nil, "you have already started the program")

	cmd := exec.Command(a.path, a.args...)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

    a.program = ptmx
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, syscall.SIGWINCH)
    go func() {
        for range ch {
            if err := pty.Setsize(os.Stdin, &pty.Winsize{
                X: 0,
                Y: 0,

                Rows: uint16(a.rows),
                Cols: uint16(a.cols),
            }); err != nil {
                slog.Error("unable to resize pty", "err", err)
            }
        }
    }()
    ch <- syscall.SIGWINCH // Initial resize.
    defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	_, err = io.Copy(a.writer, ptmx)
    return err
}

func (a *Program) Close() error {
    err := a.program.Close()
    a.program = nil
	return err
}
