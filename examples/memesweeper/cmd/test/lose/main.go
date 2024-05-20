package main

import (
	"context"
	"fmt"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/memesweeper/pkg/memesweeper"
	"github.com/theprimeagen/vim-with-me/pkg/chat"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
)

const ROUND_TIME = 300

func createChat() chan chat.ChatMsg {
	ch := make(chan chat.ChatMsg)
	go func() {
		// ensure we are in a round
		<-time.After(time.Millisecond * 50)

		ch <- chat.ChatMsg{Msg: "2B"}
		<-time.After(time.Millisecond * ROUND_TIME)
	}()

	return ch
}

func main() {
	testies.SetupLogger()

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	state := memesweeper.
		NewMemeSweeperState(3, 3).
		WithDims(3, 3).
		WithSeed(42069)

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
				fmt.Printf("chat: %v\n", ms.Chat(&msg))
			case <-ticker.C:
				fmt.Println("rendering")
				ms.Render(time.Now().UnixMilli() - start)
				fmt.Println(ms.Renderer.Debug())
			}
		}

		ticker.Stop()
	}()

	for !ms.GameOver() {

		fmt.Println("starting round")
		ms.StartRound()
		<-time.After(time.Millisecond * ROUND_TIME)
		fmt.Println("ending round")
		ms.EndRound()
	}

	<-time.After(time.Millisecond * 250)
	fmt.Println("closing")
	cancel()
}
