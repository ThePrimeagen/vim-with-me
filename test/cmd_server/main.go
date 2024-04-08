package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/commands"
	"github.com/theprimeagen/vim-with-me/pkg/tcp"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

func render(win *window.Window) {
    bytes, err := os.ReadFile("lua/vim-with-me/integration/theprimeagen")
    str := string(bytes)
    str = strings.ReplaceAll(str, "\n", "")
    str = strings.ReplaceAll(str, "*", " ")

    if err != nil {
        log.Fatalf("Error reading file: %s", err)
    }

    _ = win.SetWindow(str)
}

func partialRender(win *window.Window, row, col int, text string) {
    for i := 0; i < len(text); i++ {
        err := win.Set(row, col+i, rune(text[i]))
        if err != nil {
            log.Fatalf("Error setting partial render: %s", err)
        }
    }
}

func main() {
    server, err := testies.CreateServerFromArgs()
    if err != nil {
        log.Fatalf("Error creating server: %s", err)
    }

    win := window.NewWindow(24, 80)

    fmt.Printf("test\n")
    for {
        fmt.Printf("waiting on from socket \n")
        cmd := <-server.FromSockets
        fmt.Printf("command received %+v\n", cmd)
        switch cmd.Command {
        case "open":
            out := window.OpenCommand(win)
            server.Send(out)
        case "to_int":
            i, err := strconv.Atoi(cmd.Data)
            if err != nil {
                server.Send(&tcp.TCPCommand{
                    Command: "e",
                    Data: "Error: " + err.Error(),
                })
            } else {
                out := tcp.ToTCPInt(i)
                server.Send(&tcp.TCPCommand{
                    Command: "from_int",
                    Data: out,
                })
            }
        case "render":
            render(win)
            str := win.Render()
            out := commands.Render(str)
            server.Send(out)
        case "partial":
            data := strings.Split(cmd.Data, ":")
            row, _ := strconv.Atoi(data[0])
            col, _ := strconv.Atoi(data[1])
            partialRender(win, row, col, "theprimeagen")
            renders := win.PartialRender()
            fmt.Printf("partial render %d\n", len(renders))
            server.Send(commands.PartialRender(renders))
        }
    }
}
