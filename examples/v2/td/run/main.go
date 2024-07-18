package main

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/theprimeagen/vim-with-me/examples/v2/td"
	"github.com/theprimeagen/vim-with-me/pkg/assert"
	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/v2/chat"
)

func main() {

    testies.SetupLogger()

	ctx := context.Background()

	twitchChat, err := chat.NewTwitchChat(ctx)
	assert.NoError(err, "twitch cannot initialize")
	chtAgg := chat.
		NewChatAggregator().
		WithFilter(td.TDFilter(24, 80));

	go chtAgg.Pipe(twitchChat)

    interval := time.NewTicker(time.Second * 15)

    // STOP GAP
    for range interval.C {
        occurrences := chtAgg.ResetWithAll()
        slices.SortFunc(occurrences, func(a, b *chat.Occurrence) int {
            return b.Count - a.Count
        });

        one := false
        two := false
        for _, v := range occurrences {
            coord, err := td.FromString(v.Msg)
            if err != nil {
                continue
            }

            if !one && coord.Team == 1 {
                fmt.Printf("%s\n", coord.String())
                one = true
            }
            if !two && coord.Team == 2 {
                fmt.Printf("%s\n", coord.String())
                two = true
            }

            if one && two {
                break
            }
        }
    }
}

