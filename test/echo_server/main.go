package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/tcp"
)

func main() {
	var port uint
	flag.UintVar(&port, "port", 0, "Port to listen on")
	flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	//logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if port == 0 {
		logger.Error("You need to provide a port")
		os.Exit(1)
	}

	logger.Info("starting server", "port", port)

	server, err := tcp.NewTCPServer(uint16(port))

	if err != nil {
		logger.Error("could not start server", "error", err)
		os.Exit(1)
	}

	defer server.Close()

	logger.Info("server started and waiting for command")
	cmd := <-server.FromSockets
	logger.Info("received command", "command", cmd, "debug", server.Debug())

	server.Send(&tcp.TCPCommand{
		Command: cmd.Command,
		Data:    cmd.Data,
	})

	time.Sleep(1 * time.Second)
}
