package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

func main() {
	testies.SetupLogger()
	server, err := testies.CreateServerFromArgs()
	if err != nil {
		slog.Error("could not start server", "error", err)
		os.Exit(1)
	}

	defer server.Close()

	commander := commands.NewCommander()
	server.WelcomeMessage(tcp.MakeWelcome(commander.ToCommands()))

	log.Printf("starting server\n")

	go server.Start()

	for {
		cmd := <-server.FromSockets
		server.Send(cmd.Command)
	}
}
