package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/pkg/v2/controller"
)

type Next struct {
	next []string
	idx  int
}

func (n *Next) Next() string {
    out := n.next[n.idx]
    n.idx++
    return out
}

type Send struct {
    received []string
}

func (s *Send) SendKey(k string) {
    s.received = append(s.received, k)
}

func TestController(t *testing.T) {
	input := make(chan time.Time)
	play := make(chan time.Time)
    next := &Next{}
    send := &Send{}

    cont := controller.
        NewController(next, send).
        WithInputTimer(input).
        WithPlayTimer(play)

    next.next = []string{
        "w",
        "a",
        "s",
        "d",
        "da",
    }

    ctx := context.Background()
    go cont.Start(ctx)

    play <- time.Now()
    require.Equal(t, 0, len(send.received))

    input <- time.Now()
    play <- time.Now()
    play <- time.Now()
    play <- time.Now()
    require.Equal(t, []string{
        "w", "w", "w",
    }, send.received)
}
