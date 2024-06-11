package testies

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

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
