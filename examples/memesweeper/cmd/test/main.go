package main

import (
	"context"
	"fmt"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/memesweeper/pkg/memesweeper"
	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

func createChat() chan chat.ChatMsg {
    ch := make(chan chat.ChatMsg)
    go func() {
        <-time.After(time.Millisecond * 150)
        ch <- chat.ChatMsg{Msg: "B2"}

        <-time.After(time.Millisecond * 100)
        ch <- chat.ChatMsg{Msg: "B3"}
        ch <- chat.ChatMsg{Msg: "B3"}
    }()

    return ch
}

func main() {
    testies.SetupLogger()

    ctx, cancel := context.WithCancel(context.Background())
    _ = cancel

    state := memesweeper.NewMemeSweeperState(10, 5).WithDims(5, 10)
    ms := memesweeper.NewMemeSweeper(state)
    ch := createChat()

    go func() {
        ticker := time.NewTicker(time.Millisecond * 100)
        outer:
        for {
            start := time.Now().UnixMilli()
            select {
            case <-ctx.Done():
                fmt.Println("done")
                break outer
            case msg := <-ch:
                fmt.Printf("msg received \"%s\"\n", msg.Msg)
                ms.Chat(&msg)
            case <-ticker.C:
                fmt.Println("rendering")
                ms.Render(time.Now().UnixMilli() - start)
                ms.Renderer.Debug()
            }
        }

        ticker.Stop()
    }()

    fmt.Println("starting round")
    ms.StartRound()
    <-time.After(time.Millisecond * 350)
    fmt.Println("ending round")
    ms.EndRound()

    <-time.After(time.Millisecond * 250)
    fmt.Println("closing")
    cancel()
}


