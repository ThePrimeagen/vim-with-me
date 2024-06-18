package main

import (
	"context"
	"fmt"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/chat"
)

func main() {

    testies.SetupLogger()

	ctx := context.Background()

	//doom create controller
	twitchChat, err := chat.NewTwitchChat(ctx)
	assert.NoError(err, "twitch cannot initialize")
	chtAgg := chat.
		NewChatAggregator().
		WithFilter(td.TDFilter(10, 10)).
		WithAfterMap(td.TDAfterMap)

	go chtAgg.Pipe(twitchChat)

    interval := time.NewTicker(time.Second * 3)

    for range interval.C {
        out := chtAgg.Reset()
        if len(out.Msg) > 0 {
            fmt.Printf("%s\n", out.Msg)
        }
    }
}

