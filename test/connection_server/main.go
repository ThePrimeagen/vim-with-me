package main

import (
	"fmt"
	"log"
	"os"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp2"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

func read_conn(conn tcp2.Connection) {
	for {
		log.Printf("Reading(%d)...\n", conn.Id)
        cmd, err := conn.Next()

        if err != nil {
            log.Printf("error with: %+v\n", err)
            break
        }

		log.Printf("got command %+v", cmd)
        err = conn.Write(cmd)

        if err != nil {
            log.Printf("error sending command back: %+v\n", err)
			break
        }
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server, err := testies.CreateServerFromArgs()

    if err != nil {
        log.Fatal("errror could not start", err)
    }

    defer server.Close()

    commander := commands.NewCommander()
    server.WelcomeMessage(commander.ToCommands())

	log.Printf("starting server\n")
	fmt.Printf("starting server from fmt!\n")

    go server.Start()

    for {
        cmd := <- server.FromSockets
        server.Send(cmd.Command)
    }
}
