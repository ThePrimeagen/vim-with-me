package doom_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
	"github.com/theprimeagen/vim-with-me/pkg/v2/controller"
)

type Send struct {
	received []string
	last     string
	ch       chan struct{}
}

func (s *Send) SendKey(k string) {
	s.received = append(s.received, k)
	s.last = k
    if s.ch != nil {
        s.ch <- struct{}{}
    }
}

type Next struct {
	next []string
	idx  int
}

func (n *Next) Next() string {
	if len(n.next) == n.idx {
		return ""
	}

	out := n.next[n.idx]
	n.idx++
	return out
}

func TestDoomController(t *testing.T) {
	timeBetween := time.Millisecond * 250
	send := &Send{}
	dc := doom.
		NewDoomController(send).
		WithTimeBetweenUse(timeBetween)

    dc.Play()

	<-time.After(timeBetween)

	dc.SendKey("w")
	require.Equal(t, []string{
		"w",
	}, send.received)

	dc.SendKey("a")
	require.Equal(t, []string{
		"w", "a",
	}, send.received)

	dc.SendKey("e")
	require.Equal(t, []string{
		"w", "a", "e",
	}, send.received)

	dc.SendKey("e")
	dc.SendKey("e")
	dc.SendKey("e")
	require.Equal(t, []string{
		"w", "a", "e", "e", "e", "e",
	}, send.received)

	<-time.After(timeBetween / 2)
	dc.SendKey("e")
	dc.SendKey("e")
	dc.SendKey("e")
	require.Equal(t, []string{
		"w", "a", "e", "e", "e", "e", "e", "e", "e",
	}, send.received)

	<-time.After(timeBetween / 2)
	dc.SendKey("e")
	dc.SendKey("e")
	require.Equal(t, []string{
		"w", "a", "e", "e", "e", "e", "e", "e", "e", "e", "e",
	}, send.received)
}

func TestWithController(t *testing.T) {
	input := make(chan time.Time)
	play := make(chan time.Time)
	next := &Next{}
	send := &Send{ch: make(chan struct{})}

	doomCtrl := doom.
		NewDoomController(send).
		WithTimeBetweenUse(time.Millisecond * 0)

	ctrl := controller.
		NewController(next, doomCtrl).
		WithInputTimer(input).
		WithPlayTimer(play)

	go ctrl.Start(context.Background())

	doomCtrl.Play()

	next.next = expectedPass
	next.idx = 0

    for i := range len(expectedPass) {
        input <- time.Now()
        play <- time.Now()
        for range len(expectedPass[i]) {
            <-send.ch
        }
    }
    require.Equal(t, []string{
        "w", "a", "s", "d", "f", "e",

        "w", "a",
        "w", "d",
        "w", "f",
        "w", "e",

        "s", "a",
        "s", "d",
        "s", "f",
        "s", "e",

        "f", "w",
        "f", "a",
        "f", "s",
        "f", "d",
    }, send.received)

}

func mf(str string) bool {
    return doom.DoomFilterFn(doom.DoomChatMapFn(str))
}

func TestWithMapAndFilter(t *testing.T) {
    for _, str := range expectedPass {
        require.Equal(t, true, mf(str), str)
    }
}

var expectedPass = []string{
    "w",
    "a",
    "s",
    "d",
    "f",
    "e",

    "wa",
    "wd",
    "wf",
    "we",

    "sa",
    "sd",
    "sf",
    "se",

    "fw",
    "fa",
    "fs",
    "fd",
}

