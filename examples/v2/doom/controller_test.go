package doom_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theprimeagen/vim-with-me/examples/v2/doom"
)

type Send struct {
    received []string
}

func (s *Send) SendKey(k string) {
    s.received = append(s.received, k)
}

func TestDoomController(t *testing.T) {
    timeBetween := time.Millisecond * 250
    send := &Send{}
    dc := doom.
        NewDoomController(send).
        WithTimeBetweenUse(timeBetween)
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
        "w", "a", "e",
    }, send.received)

    <-time.After(timeBetween / 2)
    dc.SendKey("e")
    dc.SendKey("e")
    dc.SendKey("e")
    require.Equal(t, []string{
        "w", "a", "e",
    }, send.received)
    <-time.After(timeBetween / 2)
    dc.SendKey("e")
    dc.SendKey("e")
    require.Equal(t, []string{
        "w", "a", "e", "e",
    }, send.received)
}


