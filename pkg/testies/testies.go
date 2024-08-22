package testies

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

type DebugFile struct {
    fh io.WriteCloser
}
type NumberedDebugFile struct {
    n uint8
    df *DebugFile
}
type LineWriter interface {
    WriteLine(b []byte) error
    WriteStrLine(b string) error
}

func NewDebugFile(name string) (*DebugFile, error) {
    f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
    if err != nil {
        return nil, err
    }

    return &DebugFile{
        fh: f,
    }, nil
}

func (d *DebugFile) AsNumberedDebugFile(n uint8) *NumberedDebugFile {
    return &NumberedDebugFile{
        n: n,
        df: d,
    }
}

func (d *NumberedDebugFile) WriteStrLine(b string) error {
    return d.df.WriteStrLine(b)
}

func (d *NumberedDebugFile) WriteLine(b []byte) error {
    _, err := d.df.fh.Write([]byte{d.n, ':', ' '})
    if err != nil {
        return err
    }

    return d.df.WriteLine(b)
}

func (d *NumberedDebugFile) Close() error {
    return d.Close()
}

func EmptyDebugFile() *DebugFile {
    return &DebugFile{fh: nil}
}

func (d *DebugFile) WriteStrLine(b string) error {
    return d.WriteLine([]byte(b))
}

func (d *DebugFile) WriteLine(b []byte) error {
    if (d.fh == nil) {
        return nil
    }

    read := 0
    for read < len(b) {
        n, err := d.fh.Write(b[read:])
        if err != nil {
            return err
        }
        read += n
    }

    if b[len(b) - 1] != '\n' {
        _, _ = d.fh.Write([]byte{'\n'})
    }

    return nil
}

func (d *DebugFile) Close() {
    if (d.fh != nil) {
        d.fh.Close()
    }
}

func CreateServerFromArgs() (*tcp.TCP, error) {
	var port uint
	flag.UintVar(&port, "port", 0, "Port to listen on")
	flag.Parse()

	if port == 0 {
		return nil, fmt.Errorf("You need to provide a port")
	}

	slog.Info("starting server", "port", port)

	server, err := tcp.NewTCPServer(uint16(port))
	if err != nil {
		return nil, fmt.Errorf("Error creating server: %w", err)
	}

	return server, nil
}

func SetupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	switch strings.ToLower(os.Getenv("LEVEL")) {
	case "error", "e":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "debug", "d":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info", "i":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}
